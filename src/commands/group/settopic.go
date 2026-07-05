package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "settopic",
		Description: "Ubah topic/deskripsi group",
		Usage:       "settopic <topic>",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Usage: .settopic <topic>")
			}
			topic := strings.Join(ptz.Args, " ")
			if err := serialize.SetGroupTopic(ptz.Bot.Client, ptz.Chat, topic); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Topic group diubah.")
		},
	})
}
