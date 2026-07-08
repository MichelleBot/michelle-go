package core

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sql.DB
}

func NewDB(db *sql.DB) (*DB, error) {
	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}
	// Apply optimizations
	pragmas := []string{
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA cache_size=-2000;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, err
		}
	}
	// Configure connection pool
	db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	
	d := &DB{Conn: db}
	err := d.migrate()
	return d, err
}

func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (jid TEXT PRIMARY KEY, lid TEXT, name TEXT, banned BOOLEAN, ban_temporary INTEGER, ban_times INTEGER, point INTEGER, balance INTEGER, pocket INTEGER, deposito INTEGER, guard INTEGER, last_claim INTEGER, last_rob INTEGER, premium BOOLEAN, expired INTEGER, last_notified INTEGER, last_seen INTEGER, last_jadibot INTEGER, hit INTEGER, warning INTEGER, attempt INTEGER, code TEXT, code_expire INTEGER, email TEXT, verified BOOLEAN, taken BOOLEAN, partner TEXT, saving_history TEXT, player_data TEXT)`,
		`CREATE TABLE IF NOT EXISTS groups (jid TEXT PRIMARY KEY, name TEXT, activity INTEGER, adzan BOOLEAN, antibot BOOLEAN, antiporn BOOLEAN, antidelete BOOLEAN, antilink BOOLEAN, antiphishing BOOLEAN, antitagsw BOOLEAN, antivirtex BOOLEAN, antiforward BOOLEAN, antisticker BOOLEAN, adminonly BOOLEAN, captcha BOOLEAN, filter BOOLEAN, game BOOLEAN, mysterybox BOOLEAN, left BOOLEAN, localonly BOOLEAN, list_data TEXT, mute BOOLEAN, autosticker BOOLEAN, restrict BOOLEAN, member_data TEXT, text_left TEXT, text_welcome TEXT, welcome BOOLEAN, expired INTEGER, last_notified INTEGER, blocked TEXT, blacklist TEXT, stay BOOLEAN, open_at TEXT, close_at TEXT)`,
		`CREATE TABLE IF NOT EXISTS chats (jid TEXT PRIMARY KEY, chat INTEGER, last_seen INTEGER, last_reply INTEGER)`,
		`CREATE TABLE IF NOT EXISTS command_stats (cmd TEXT PRIMARY KEY, total_hit INTEGER, today_hit INTEGER, last_hit INTEGER)`,
		`CREATE TABLE IF NOT EXISTS jadibot_sessions (phone TEXT PRIMARY KEY, owner_jid TEXT, created_at INTEGER)`,
	}

	for _, q := range queries {
		if _, err := db.Conn.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
