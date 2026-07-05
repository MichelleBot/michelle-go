package owner

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "leave",
		Description: "Bot keluar dari group",
		Usage:       "leave",
		Category:    "owner",
		OwnerOnly:   true,
		GroupOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if err := serialize.LeaveGroup(ptz.Bot.Client, ptz.Chat); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Berhasil leave group.")
		},
	})
}
