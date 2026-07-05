package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "setdesc",
		Aliases:     []string{"setdescription"},
		Description: "Ubah deskripsi group",
		Usage:       "setdesc <deskripsi>",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if ptz.RawArgs == "" {
				return ptz.ReplyText("Usage: .setdesc <deskripsi>")
			}
			if err := serialize.SetGroupDescription(ptz.Bot.Client, ptz.Chat, ptz.RawArgs); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Deskripsi group diubah.")
		},
	})
}
