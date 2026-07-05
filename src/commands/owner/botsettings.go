package owner

import (
	"sptzx/src/core"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "botsettings",
		Aliases:     []string{"settings", "bs"},
		Description: "Lihat pengaturan bot saat ini",
		Usage:       "botsettings",
		Category:    "owner",
		OwnerOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			flag := func(v bool) string {
				if v {
					return "✅ ON"
				}
				return "❌ OFF"
			}
			var sb strings.Builder
			sb.WriteString("🤖 *Bot Settings*\n\n")
			sb.WriteString("Self Mode: " + flag(ptz.Bot.BotConfig.GetSelfMode()) + "\n")
			sb.WriteString("Private Only: " + flag(ptz.Bot.BotConfig.GetPrivateOnly()) + "\n")
			sb.WriteString("Group Only: " + flag(ptz.Bot.BotConfig.GetGroupOnly()) + "\n")
			return ptz.ReplyText(sb.String())
		},
	})
}
