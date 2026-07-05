package games

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type GameDB struct {
	db  *sql.DB
	mu  sync.RWMutex
	mem map[string]*Session
	qid map[string]string
}

type Session struct {
	Key        string
	GameType   string
	Answer     string
	QuestionID string
	ChatJID    string
	StarterJID string
	IsGroup    bool
	ClueCount  int
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

var (
	db   *GameDB
	once sync.Once
)

const sessionTTL = 5 * time.Minute

func InitDB(dsn string) error {
	var initErr error
	once.Do(func() {
		sqlDB, err := sql.Open("sqlite3", dsn)
		if err != nil {
			initErr = err
			return
		}
		db = &GameDB{
			db:  sqlDB,
			mem: make(map[string]*Session),
			qid: make(map[string]string),
		}
		if err := db.migrate(); err != nil {
			initErr = err
			return
		}
		if err := db.loadFromDB(); err != nil {
			initErr = err
			return
		}
		go db.runCleanup()
	})
	return initErr
}

func (g *GameDB) migrate() error {
	_, err := g.db.Exec(`
		CREATE TABLE IF NOT EXISTS game_sessions (
			key         TEXT PRIMARY KEY,
			game_type   TEXT NOT NULL,
			answer      TEXT NOT NULL,
			question_id TEXT NOT NULL,
			chat_jid    TEXT NOT NULL,
			starter_jid TEXT NOT NULL,
			is_group    INTEGER NOT NULL DEFAULT 0,
			clue_count  INTEGER NOT NULL DEFAULT 0,
			created_at  INTEGER NOT NULL,
			expires_at  INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS game_disabled (
			chat_jid    TEXT PRIMARY KEY,
			disabled_at INTEGER NOT NULL
		);
	`)
	return err
}

func (g *GameDB) loadFromDB() error {
	rows, err := g.db.Query(
		`SELECT key, game_type, answer, question_id, chat_jid, starter_jid, is_group, clue_count, created_at, expires_at
		 FROM game_sessions WHERE expires_at > ?`, time.Now().Unix(),
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		s := &Session{}
		var ca, ea int64
		var isGroup int
		if err := rows.Scan(&s.Key, &s.GameType, &s.Answer, &s.QuestionID,
			&s.ChatJID, &s.StarterJID, &isGroup, &s.ClueCount, &ca, &ea); err != nil {
			continue
		}
		s.IsGroup = isGroup == 1
		s.CreatedAt = time.Unix(ca, 0)
		s.ExpiresAt = time.Unix(ea, 0)
		g.mem[s.Key] = s
		g.qid[s.QuestionID] = s.Key
	}
	return rows.Err()
}

func SetSession(sess *Session) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if old, ok := db.mem[sess.Key]; ok {
		delete(db.qid, old.QuestionID)
		db.db.Exec(`DELETE FROM game_sessions WHERE key = ?`, sess.Key)
	}
	db.mem[sess.Key] = sess
	db.qid[sess.QuestionID] = sess.Key
	isGroup := 0
	if sess.IsGroup {
		isGroup = 1
	}
	db.db.Exec(
		`INSERT OR REPLACE INTO game_sessions
		 (key, game_type, answer, question_id, chat_jid, starter_jid, is_group, clue_count, created_at, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sess.Key, sess.GameType, sess.Answer, sess.QuestionID,
		sess.ChatJID, sess.StarterJID, isGroup, sess.ClueCount,
		sess.CreatedAt.Unix(), sess.ExpiresAt.Unix(),
	)
}

func UpdateClueCount(sess *Session) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if s, ok := db.mem[sess.Key]; ok {
		s.ClueCount = sess.ClueCount
	}
	db.db.Exec(`UPDATE game_sessions SET clue_count = ? WHERE key = ?`, sess.ClueCount, sess.Key)
}

func MatchByQuestionID(questionID string) (*Session, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	key, ok := db.qid[questionID]
	if !ok {
		return nil, false
	}
	sess, ok := db.mem[key]
	if !ok || time.Now().After(sess.ExpiresAt) {
		return nil, false
	}
	return sess, true
}

func GetActiveSessionByKey(key string) (*Session, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	sess, ok := db.mem[key]
	if !ok || time.Now().After(sess.ExpiresAt) {
		return nil, false
	}
	return sess, true
}

func GetActiveChatSession(chatJID, gameType string) (*Session, bool) {
	return GetActiveSessionByKey(chatSessionKey(chatJID, gameType))
}

func GetActiveUserSession(chatJID, senderJID, gameType string) (*Session, bool) {
	return GetActiveSessionByKey(userSessionKey(chatJID, senderJID, gameType))
}

func HasActiveChatSession(chatJID string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, sess := range db.mem {
		if sess.ChatJID == chatJID && sess.IsGroup && time.Now().Before(sess.ExpiresAt) {
			return true
		}
	}
	return false
}

func GetActiveChatSessionAny(chatJID string) (*Session, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, sess := range db.mem {
		if sess.ChatJID == chatJID && sess.IsGroup && time.Now().Before(sess.ExpiresAt) {
			return sess, true
		}
	}
	return nil, false
}

func GetActiveUserSessionAny(chatJID, senderJID string) (*Session, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, sess := range db.mem {
		if sess.ChatJID == chatJID && sess.StarterJID == senderJID && !sess.IsGroup && time.Now().Before(sess.ExpiresAt) {
			return sess, true
		}
	}
	return nil, false
}

func DeleteSession(sess *Session) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.mem, sess.Key)
	delete(db.qid, sess.QuestionID)
	db.db.Exec(`DELETE FROM game_sessions WHERE key = ?`, sess.Key)
}

func IsGameEnabled(chatJID string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var count int
	db.db.QueryRow(`SELECT COUNT(*) FROM game_disabled WHERE chat_jid = ?`, chatJID).Scan(&count)
	return count == 0
}

func SetGameEnabled(chatJID string, enabled bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if enabled {
		db.db.Exec(`DELETE FROM game_disabled WHERE chat_jid = ?`, chatJID)
	} else {
		db.db.Exec(
			`INSERT OR REPLACE INTO game_disabled (chat_jid, disabled_at) VALUES (?, ?)`,
			chatJID, time.Now().Unix(),
		)
	}
}

func chatSessionKey(chatJID, gameType string) string {
	return chatJID + "|group|" + gameType
}

func userSessionKey(chatJID, senderJID, gameType string) string {
	return chatJID + "|" + senderJID + "|" + gameType
}

func NewSession(chatJID, senderJID, gameType, answer, questionID string, isGroup bool) *Session {
	now := time.Now()
	var key string
	if isGroup {
		key = chatSessionKey(chatJID, gameType)
	} else {
		key = userSessionKey(chatJID, senderJID, gameType)
	}
	return &Session{
		Key:        key,
		GameType:   gameType,
		Answer:     answer,
		QuestionID: questionID,
		ChatJID:    chatJID,
		StarterJID: senderJID,
		IsGroup:    isGroup,
		ClueCount:  0,
		CreatedAt:  now,
		ExpiresAt:  now.Add(sessionTTL),
	}
}

func (g *GameDB) runCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		g.cleanup()
	}
}

func (g *GameDB) cleanup() {
	now := time.Now()
	g.mu.Lock()
	defer g.mu.Unlock()
	for key, sess := range g.mem {
		if now.After(sess.ExpiresAt) {
			delete(g.mem, key)
			delete(g.qid, sess.QuestionID)
		}
	}
	g.db.Exec(`DELETE FROM game_sessions WHERE expires_at <= ?`, now.Unix())
}
