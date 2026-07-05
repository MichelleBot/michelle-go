package group

import (
	"go.mau.fi/whatsmeow/types"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "testwelcome",
		Description: "Test welcome message",
		Usage:       "testwelcome",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			s := ptz.Bot.Settings.GetGroupSettings(ptz.Chat)
			if !s.WelcomeEnabled {
				return ptz.ReplyText("Welcome message belum diaktifkan. Gunakan .enablewelcome")
			}
			msg := strings.ReplaceAll(s.WelcomeMessage, "@user", "@"+ptz.Sender.User)
			return serialize.SendTextMention(ptz.Bot.Client, ptz.Chat, msg, []types.JID{ptz.Sender})
		},
	})
}
