package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "setname",
		Description: "Ubah nama group",
		Usage:       "setname <nama>",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Usage: .setname <nama>")
			}
			name := strings.Join(ptz.Args, " ")
			if err := serialize.SetGroupName(ptz.Bot.Client, ptz.Chat, name); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Nama group diubah ke: " + name)
		},
	})
}
