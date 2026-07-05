package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Owners          []string
	Prefixes        []string
	SessionDB       string
	LoginMethod     string
	PairingPhone    string
	LogLevel        string
	Timezone        string
	StickerPackName string
	StickerAuthor   string
	Antispam        AntispamConfig
	SiputzX         SiputzXConfig
}

type AntispamConfig struct {
	MaxMsgPerSecond int
	MaxMsgPerMinute int
	BanDurationSecs int
}

type SiputzXConfig struct {
	Enabled      bool
	BaseURL      string
	GeminiCookie string
}

func Load() *Config {
	godotenv.Load()

	owners := strings.Split(os.Getenv("BOT_OWNERS"), ",")
	filtered := make([]string, 0, len(owners))
	for _, o := range owners {
		o = strings.TrimSpace(o)
		if o != "" {
			filtered = append(filtered, o)
		}
	}

	prefixStr := os.Getenv("BOT_PREFIX")
	if prefixStr == "" {
		prefixStr = "!,.,"
	}
	prefixes := strings.Split(prefixStr, ",")
	cleanPrefixes := make([]string, 0, len(prefixes))
	for _, p := range prefixes {
		p = strings.TrimSpace(p)
		if p != "" {
			cleanPrefixes = append(cleanPrefixes, p)
		}
	}
	if len(cleanPrefixes) == 0 {
		cleanPrefixes = []string{"!", ".", "/"}
	}

	sessionDB := os.Getenv("SESSION_DB")
	if sessionDB == "" {
		sessionDB = "file:session.db?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000"
	}

	loginMethod := os.Getenv("LOGIN_METHOD")
	if loginMethod == "" {
		loginMethod = "qr"
	}
	loginMethod = strings.ToLower(strings.TrimSpace(loginMethod))
	if loginMethod == "pair" {
		loginMethod = "paircode"
	}

	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	stickerPackName := os.Getenv("STICKER_PACK_NAME")
	if stickerPackName == "" {
		stickerPackName = "WhatsApp Bot"
	}

	stickerAuthor := os.Getenv("STICKER_AUTHOR")
	if stickerAuthor == "" {
		stickerAuthor = "Siputzx"
	}

	pairingPhone := os.Getenv("PAIRING_PHONE")
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	siputzxEnabled := os.Getenv("SIPUTZX_ENABLED")
	siputzxBaseURL := os.Getenv("SIPUTZX_BASE_URL")
	if siputzxBaseURL == "" {
		siputzxBaseURL = "https://api.siputzx.my.id"
	}
	geminiCookie := os.Getenv("GEMINI_COOKIE")

	maxMsgPerSecond := parseEnvInt("ANTISPAM_MAX_PER_SECOND", 3)
	maxMsgPerMinute := parseEnvInt("ANTISPAM_MAX_PER_MINUTE", 20)
	banDurationSecs := parseEnvInt("ANTISPAM_BAN_DURATION_SECS", 30)

	return &Config{
		Owners:          filtered,
		Prefixes:        cleanPrefixes,
		SessionDB:       sessionDB,
		LoginMethod:     loginMethod,
		PairingPhone:    pairingPhone,
		LogLevel:        logLevel,
		Timezone:        timezone,
		StickerPackName: stickerPackName,
		StickerAuthor:   stickerAuthor,
		Antispam: AntispamConfig{
			MaxMsgPerSecond: maxMsgPerSecond,
			MaxMsgPerMinute: maxMsgPerMinute,
			BanDurationSecs: banDurationSecs,
		},
		SiputzX: SiputzXConfig{
			Enabled:      siputzxEnabled == "true",
			BaseURL:      siputzxBaseURL,
			GeminiCookie: geminiCookie,
		},
	}
}

func parseEnvInt(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(v)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func (c *Config) Validate() error {
	if len(c.Owners) == 0 {
		return fmt.Errorf("BOT_OWNERS tidak boleh kosong di .env")
	}
	if c.SessionDB == "" {
		return fmt.Errorf("SESSION_DB tidak diset di .env")
	}
	if c.LoginMethod != "qr" && c.LoginMethod != "paircode" {
		return fmt.Errorf("LOGIN_METHOD harus qr atau paircode")
	}
	if c.Antispam.MaxMsgPerSecond <= 0 || c.Antispam.MaxMsgPerMinute <= 0 || c.Antispam.BanDurationSecs <= 0 {
		return fmt.Errorf("konfigurasi antispam harus lebih dari 0")
	}
	return nil
}
