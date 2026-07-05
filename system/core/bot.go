package core

import (
	"database/sql"
	"michelle/system/config"
	"michelle/system/middleware"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Bot struct {
	Client         *whatsmeow.Client
	Config         *config.Config
	Antispam       *middleware.Antispam
	CommandLimiter *middleware.CommandLimiter
	Settings       *SettingsStore
	Users          *UserStore
	Messages       *MessageStore
	Polls          *PollStore
	BotConfig      *BotSettings
	Registry       *Registry
	Container      *sqlstore.Container
	Log            waLog.Logger
}

func NewBot(cfg *config.Config, container *sqlstore.Container, client *whatsmeow.Client, log waLog.Logger, db *sql.DB) *Bot {

	return &Bot{
		Client:         client,
		Config:         cfg,
		Antispam:       middleware.NewAntispam(cfg.Antispam.MaxMsgPerSecond, cfg.Antispam.MaxMsgPerMinute, cfg.Antispam.BanDurationSecs),
		CommandLimiter: middleware.NewCommandLimiter(),
		Settings:       NewSettingsStore(db, log),
		Users:          NewUserStore(db, log),
		Messages:       NewMessageStore(),
		Polls:          NewPollStore(),
		BotConfig:      NewBotSettings(),
		Container:      container,
		Log:            log,
	}
}

func (b *Bot) GetPrefix() string {
	if len(b.Config.Prefixes) > 0 {
		return b.Config.Prefixes[0]
	}
	return "!"
}
