package core

import (
	"database/sql"
	"fmt"
	"math"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

const (
	freeDailyLimit    = 100
	premiumDailyLimit = 500
)

type SettingsStore struct {
	mu       sync.RWMutex
	db       *sql.DB
	log      waLog.Logger
	settings map[types.JID]*GroupSettings
}

type UserProfile struct {
	JID          string
	LimitBalance int
	DailyLimit   int
	DailyUsed    int
	DailyRemain  int
	ExtraLimit   int
	Credit       int
	XP           int
	Interactions int
	PremiumUntil time.Time
}

type UserRank struct {
	JID   string
	Value int
}

type UserStore struct {
	mu  sync.Mutex
	db  *sql.DB
	log waLog.Logger
}

type GroupSettings struct {
	WelcomeEnabled    bool
	WelcomeMessage    string
	GoodbyeEnabled    bool
	GoodbyeMessage    string
	AntispamEnabled   bool
	AntiDeleteEnabled bool
}

type BotSettings struct {
	mu          sync.RWMutex
	SelfMode    bool
	PrivateOnly bool
	GroupOnly   bool
}

func NewSettingsStore(db *sql.DB, log waLog.Logger) *SettingsStore {
	s := &SettingsStore{
		db:       db,
		log:      log,
		settings: make(map[types.JID]*GroupSettings),
	}
	s.migrate()
	s.loadAll()
	return s
}

func NewUserStore(db *sql.DB, log waLog.Logger) *UserStore {
	u := &UserStore{db: db, log: log}
	u.migrate()
	return u
}

func (s *SettingsStore) migrate() {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS group_settings (
		chat_jid         TEXT PRIMARY KEY,
		welcome_enabled  INTEGER NOT NULL DEFAULT 0,
		welcome_message  TEXT NOT NULL DEFAULT '',
		goodbye_enabled  INTEGER NOT NULL DEFAULT 0,
		goodbye_message  TEXT NOT NULL DEFAULT '',
		antispam_enabled INTEGER NOT NULL DEFAULT 0,
		anti_delete_enabled INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil && s.log != nil {
		s.log.Errorf("DB migrate error: %v", err)
	}
	_, _ = s.db.Exec(`ALTER TABLE group_settings ADD COLUMN anti_delete_enabled INTEGER NOT NULL DEFAULT 0`)
}

func (u *UserStore) migrate() {
	_, err := u.db.Exec(`CREATE TABLE IF NOT EXISTS user_profiles (
		jid           TEXT PRIMARY KEY,
		limit_balance INTEGER NOT NULL DEFAULT 0,
		credit        INTEGER NOT NULL DEFAULT 0,
		xp            INTEGER NOT NULL DEFAULT 0,
		interactions  INTEGER NOT NULL DEFAULT 0,
		last_active   INTEGER NOT NULL DEFAULT 0,
		daily_used    INTEGER NOT NULL DEFAULT 0,
		limit_reset_day INTEGER NOT NULL DEFAULT 0,
		premium_until INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil && u.log != nil {
		u.log.Errorf("DB user migrate error: %v", err)
	}

	u.tryAddColumn("ALTER TABLE user_profiles ADD COLUMN xp INTEGER NOT NULL DEFAULT 0")
	u.tryAddColumn("ALTER TABLE user_profiles ADD COLUMN interactions INTEGER NOT NULL DEFAULT 0")
	u.tryAddColumn("ALTER TABLE user_profiles ADD COLUMN last_active INTEGER NOT NULL DEFAULT 0")
	u.tryAddColumn("ALTER TABLE user_profiles ADD COLUMN daily_used INTEGER NOT NULL DEFAULT 0")
	u.tryAddColumn("ALTER TABLE user_profiles ADD COLUMN limit_reset_day INTEGER NOT NULL DEFAULT 0")
}

func (u *UserStore) tryAddColumn(query string) {
	if _, err := u.db.Exec(query); err != nil && u.log != nil {
		u.log.Debugf("DB user alter skip: %v", err)
	}
}

func (u *UserStore) ensureRow(jid string) {
	_, err := u.db.Exec(`INSERT OR IGNORE INTO user_profiles (jid, limit_balance, credit, xp, interactions, last_active, daily_used, limit_reset_day, premium_until) VALUES (?, 0, 0, 0, 0, 0, 0, ?, 0)`, jid, currentUTCResetDay())
	if err != nil && u.log != nil {
		u.log.Errorf("DB user ensure row error: %v", err)
	}
}

func (u *UserStore) GetUserProfile(jid string) *UserProfile {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.ensureRow(jid)

	var extraLimit, credit, xp, interactions, dailyUsed, limitResetDay int
	var premiumUntilUnix int64
	err := u.db.QueryRow(`SELECT limit_balance, credit, xp, interactions, daily_used, limit_reset_day, premium_until FROM user_profiles WHERE jid = ?`, jid).
		Scan(&extraLimit, &credit, &xp, &interactions, &dailyUsed, &limitResetDay, &premiumUntilUnix)
	if err != nil {
		if u.log != nil {
			u.log.Errorf("DB get user profile error: %v", err)
		}
		return &UserProfile{JID: jid, LimitBalance: 0, Credit: 0}
	}

	if limitResetDay != currentUTCResetDay() {
		dailyUsed = 0
		if _, err := u.db.Exec(`UPDATE user_profiles SET daily_used = 0, limit_reset_day = ? WHERE jid = ?`, currentUTCResetDay(), jid); err != nil && u.log != nil {
			u.log.Errorf("DB daily reset error: %v", err)
		}
	}

	dailyLimit := dailyAllowanceFromPremiumUntil(premiumUntilUnix)
	dailyRemain := maxInt(0, dailyLimit-dailyUsed)

	profile := &UserProfile{
		JID:          jid,
		LimitBalance: dailyRemain + extraLimit,
		DailyLimit:   dailyLimit,
		DailyUsed:    dailyUsed,
		DailyRemain:  dailyRemain,
		ExtraLimit:   extraLimit,
		Credit:       credit,
		XP:           xp,
		Interactions: interactions,
	}
	if premiumUntilUnix > 0 {
		profile.PremiumUntil = time.Unix(premiumUntilUnix, 0)
	}
	return profile
}

func (u *UserStore) AddLimit(jid string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount limit harus lebih dari 0")
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	u.ensureRow(jid)

	_, err := u.db.Exec(`UPDATE user_profiles SET limit_balance = limit_balance + ? WHERE jid = ?`, amount, jid)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) AddCredit(jid string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount kredit harus lebih dari 0")
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	u.ensureRow(jid)

	_, err := u.db.Exec(`UPDATE user_profiles SET credit = credit + ? WHERE jid = ?`, amount, jid)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) AddPremiumDays(jid string, days int) error {
	if days <= 0 {
		return fmt.Errorf("hari premium harus lebih dari 0")
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	u.ensureRow(jid)

	var premiumUntilUnix int64
	err := u.db.QueryRow(`SELECT premium_until FROM user_profiles WHERE jid = ?`, jid).Scan(&premiumUntilUnix)
	if err != nil {
		return err
	}

	now := time.Now()
	base := now
	if premiumUntilUnix > now.Unix() {
		base = time.Unix(premiumUntilUnix, 0)
	}
	next := base.Add(time.Duration(days) * 24 * time.Hour)

	_, err = u.db.Exec(`UPDATE user_profiles SET premium_until = ? WHERE jid = ?`, next.Unix(), jid)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) IsPremium(jid string) bool {
	profile := u.GetUserProfile(jid)
	return !profile.PremiumUntil.IsZero() && time.Now().Before(profile.PremiumUntil)
}

func (u *UserStore) ConsumeLimit(jid string, amount int) (bool, error) {
	if amount <= 0 {
		return true, nil
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	u.ensureRow(jid)

	var extraLimit, dailyUsed, limitResetDay int
	var premiumUntilUnix int64
	err := u.db.QueryRow(`SELECT limit_balance, daily_used, limit_reset_day, premium_until FROM user_profiles WHERE jid = ?`, jid).
		Scan(&extraLimit, &dailyUsed, &limitResetDay, &premiumUntilUnix)
	if err != nil {
		return false, err
	}

	if limitResetDay != currentUTCResetDay() {
		dailyUsed = 0
		limitResetDay = currentUTCResetDay()
	}

	dailyLimit := dailyAllowanceFromPremiumUntil(premiumUntilUnix)
	dailyRemain := maxInt(0, dailyLimit-dailyUsed)
	if dailyRemain+extraLimit < amount {
		return false, nil
	}

	consumeDaily := minInt(amount, dailyRemain)
	consumeExtra := amount - consumeDaily
	dailyUsed += consumeDaily
	extraLimit -= consumeExtra

	_, err = u.db.Exec(`UPDATE user_profiles SET limit_balance = ?, daily_used = ?, limit_reset_day = ? WHERE jid = ?`, extraLimit, dailyUsed, limitResetDay, jid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserStore) BuyLimit(jid string, qty int, pricePerUnit int) error {
	if qty <= 0 {
		return fmt.Errorf("jumlah limit harus lebih dari 0")
	}
	if pricePerUnit <= 0 {
		return fmt.Errorf("harga limit tidak valid")
	}

	cost := qty * pricePerUnit

	u.mu.Lock()
	defer u.mu.Unlock()
	u.ensureRow(jid)

	var credit int
	err := u.db.QueryRow(`SELECT credit FROM user_profiles WHERE jid = ?`, jid).Scan(&credit)
	if err != nil {
		return err
	}
	if credit < cost {
		return fmt.Errorf("kredit tidak cukup")
	}

	_, err = u.db.Exec(`UPDATE user_profiles SET credit = credit - ?, limit_balance = limit_balance + ? WHERE jid = ?`, cost, qty, jid)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserStore) TrackInteraction(jid string) (int, int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.ensureRow(jid)

	var interactions int
	var lastActive int64
	err := u.db.QueryRow(`SELECT interactions, last_active FROM user_profiles WHERE jid = ?`, jid).Scan(&interactions, &lastActive)
	if err != nil {
		return 0, 0, err
	}

	now := time.Now()
	delta := now.Sub(time.Unix(lastActive, 0))

	base := 8
	switch {
	case lastActive == 0:
		base = 10
	case delta < 15*time.Second:
		base = 2
	case delta < time.Minute:
		base = 4
	case delta < 5*time.Minute:
		base = 7
	}

	streakBonus := int(math.Min(6, math.Log1p(float64(interactions+1))*1.8))
	xpGain := base + streakBonus
	if xpGain < 1 {
		xpGain = 1
	}

	_, err = u.db.Exec(`UPDATE user_profiles SET interactions = interactions + 1, xp = xp + ?, last_active = ? WHERE jid = ?`, xpGain, now.Unix(), jid)
	if err != nil {
		return 0, 0, err
	}

	return xpGain, interactions + 1, nil
}

func (u *UserStore) TopByXP(limit int) []UserRank {
	if limit <= 0 {
		limit = 5
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	rows, err := u.db.Query(`SELECT jid, xp FROM user_profiles ORDER BY xp DESC, interactions DESC LIMIT ?`, limit)
	if err != nil {
		if u.log != nil {
			u.log.Errorf("DB top xp error: %v", err)
		}
		return nil
	}
	defer rows.Close()

	result := make([]UserRank, 0, limit)
	for rows.Next() {
		var row UserRank
		if err := rows.Scan(&row.JID, &row.Value); err != nil {
			continue
		}
		result = append(result, row)
	}
	return result
}

func (u *UserStore) TopByCredit(limit int) []UserRank {
	if limit <= 0 {
		limit = 5
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	rows, err := u.db.Query(`SELECT jid, credit FROM user_profiles ORDER BY credit DESC, xp DESC LIMIT ?`, limit)
	if err != nil {
		if u.log != nil {
			u.log.Errorf("DB top credit error: %v", err)
		}
		return nil
	}
	defer rows.Close()

	result := make([]UserRank, 0, limit)
	for rows.Next() {
		var row UserRank
		if err := rows.Scan(&row.JID, &row.Value); err != nil {
			continue
		}
		result = append(result, row)
	}
	return result
}

func currentUTCResetDay() int {
	now := time.Now().UTC()
	return now.Year()*10000 + int(now.Month())*100 + now.Day()
}

func dailyAllowanceFromPremiumUntil(premiumUntilUnix int64) int {
	if premiumUntilUnix > time.Now().Unix() {
		return premiumDailyLimit
	}
	return freeDailyLimit
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *SettingsStore) loadAll() {
	rows, err := s.db.Query(`SELECT chat_jid, welcome_enabled, welcome_message,
		goodbye_enabled, goodbye_message, antispam_enabled, anti_delete_enabled FROM group_settings`)
	if err != nil {
		if s.log != nil {
			s.log.Errorf("DB loadAll error: %v", err)
		}
		return
	}
	defer rows.Close()
	for rows.Next() {
		var jidStr string
		var we, ge, ae, ade int
		gs := &GroupSettings{}
		if err := rows.Scan(&jidStr, &we, &gs.WelcomeMessage, &ge, &gs.GoodbyeMessage, &ae, &ade); err != nil {
			if s.log != nil {
				s.log.Errorf("DB scan error: %v", err)
			}
			continue
		}
		gs.WelcomeEnabled = we == 1
		gs.GoodbyeEnabled = ge == 1
		gs.AntispamEnabled = ae == 1
		gs.AntiDeleteEnabled = ade == 1
		jid, err := types.ParseJID(jidStr)
		if err != nil {
			continue
		}
		s.settings[jid] = gs
	}
	if err := rows.Err(); err != nil && s.log != nil {
		s.log.Errorf("DB rows error: %v", err)
	}
}

func (s *SettingsStore) GetGroupSettings(jid types.JID) *GroupSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if gs, ok := s.settings[jid]; ok {
		return gs
	}
	return &GroupSettings{
		WelcomeEnabled:    false,
		WelcomeMessage:    "Selamat datang @user! Semoga betah di grup ini.",
		GoodbyeEnabled:    false,
		GoodbyeMessage:    "Selamat tinggal @user!",
		AntispamEnabled:   false,
		AntiDeleteEnabled: false,
	}
}

func (s *SettingsStore) SetGroupSettings(jid types.JID, gs *GroupSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings[jid] = gs
	we, ge, ae, ade := 0, 0, 0, 0
	if gs.WelcomeEnabled {
		we = 1
	}
	if gs.GoodbyeEnabled {
		ge = 1
	}
	if gs.AntispamEnabled {
		ae = 1
	}
	if gs.AntiDeleteEnabled {
		ade = 1
	}
	_, err := s.db.Exec(`INSERT OR REPLACE INTO group_settings
		(chat_jid, welcome_enabled, welcome_message, goodbye_enabled, goodbye_message, antispam_enabled, anti_delete_enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		jid.String(), we, gs.WelcomeMessage, ge, gs.GoodbyeMessage, ae, ade,
	)
	if err != nil && s.log != nil {
		s.log.Errorf("DB save error: %v", err)
	}
}

func NewBotSettings() *BotSettings {
	return &BotSettings{}
}

func (bs *BotSettings) GetSelfMode() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.SelfMode
}

func (bs *BotSettings) SetSelfMode(v bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.SelfMode = v
}

func (bs *BotSettings) GetPrivateOnly() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.PrivateOnly
}

func (bs *BotSettings) SetPrivateOnly(v bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.PrivateOnly = v
}

func (bs *BotSettings) GetGroupOnly() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.GroupOnly
}

func (bs *BotSettings) SetGroupOnly(v bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.GroupOnly = v
}
