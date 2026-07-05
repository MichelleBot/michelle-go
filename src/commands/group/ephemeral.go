package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "ephemeral",
		Aliases:     []string{"disappear"},
		Description: "Set timer pesan menghilang (off/24h/7d/90d)",
		Usage:       "ephemeral <off|24h|7d|90d>",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Usage: .ephemeral <off|24h|7d|90d>")
			}
			val := strings.ToLower(ptz.Args[0])
			var fn func() error
			switch val {
			case "off", "0":
				fn = func() error { return serialize.SetDisappearingOff(ptz.Bot.Client, ptz.Chat) }
			case "24h", "1d":
				fn = func() error { return serialize.SetDisappearing24h(ptz.Bot.Client, ptz.Chat) }
			case "7d", "1w":
				fn = func() error { return serialize.SetDisappearing7d(ptz.Bot.Client, ptz.Chat) }
			case "90d", "3m":
				fn = func() error { return serialize.SetDisappearing90d(ptz.Bot.Client, ptz.Chat) }
			default:
				return ptz.ReplyText("Pilihan: off, 24h, 7d, 90d")
			}
			if err := fn(); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Timer pesan hilang: " + val)
		},
	})
}
