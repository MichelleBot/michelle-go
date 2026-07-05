package group

import (
	"go.mau.fi/whatsmeow/types"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "testgoodbye",
		Description: "Test goodbye message",
		Usage:       "testgoodbye",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			s := ptz.Bot.Settings.GetGroupSettings(ptz.Chat)
			if !s.GoodbyeEnabled {
				return ptz.ReplyText("Goodbye message belum diaktifkan. Gunakan .enablegoodbye")
			}
			msg := strings.ReplaceAll(s.GoodbyeMessage, "@user", "@"+ptz.Sender.User)
			return serialize.SendTextMention(ptz.Bot.Client, ptz.Chat, msg, []types.JID{ptz.Sender})
		},
	})
}
