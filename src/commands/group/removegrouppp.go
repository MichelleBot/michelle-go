package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "removegrouppp",
		Aliases:     []string{"removegpp"},
		Description: "Hapus foto profil group",
		Usage:       "removegrouppp",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if _, err := serialize.RemoveGroupPhoto(ptz.Bot.Client, ptz.Chat); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Foto group berhasil dihapus.")
		},
	})
}
