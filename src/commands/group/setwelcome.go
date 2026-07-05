package group

import (
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "setwelcome",
		Description: "Set pesan welcome member baru",
		Usage:       "setwelcome <pesan> (gunakan @user untuk mention)",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if ptz.RawArgs == "" {
				return ptz.ReplyText("Usage: .setwelcome <pesan>\nGunakan @user untuk mention.")
			}
			s := ptz.Bot.Settings.GetGroupSettings(ptz.Chat)
			s.WelcomeMessage = ptz.RawArgs
			ptz.Bot.Settings.SetGroupSettings(ptz.Chat, s)
			return ptz.ReplyText("✅ Pesan welcome diset:\n" + ptz.RawArgs)
		},
	})
}
