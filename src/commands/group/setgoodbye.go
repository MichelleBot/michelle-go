package group

import (
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "setgoodbye",
		Description: "Set pesan goodbye member keluar",
		Usage:       "setgoodbye <pesan> (gunakan @user untuk mention)",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if ptz.RawArgs == "" {
				return ptz.ReplyText("Usage: .setgoodbye <pesan>\nGunakan @user untuk mention.")
			}
			s := ptz.Bot.Settings.GetGroupSettings(ptz.Chat)
			s.GoodbyeMessage = ptz.RawArgs
			ptz.Bot.Settings.SetGroupSettings(ptz.Chat, s)
			return ptz.ReplyText("✅ Pesan goodbye diset:\n" + ptz.RawArgs)
		},
	})
}
