package owner

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "join",
		Description: "Bot join group via invite link",
		Usage:       "join <link>",
		Category:    "owner",
		OwnerOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Usage: .join <link>")
			}
			link := strings.TrimPrefix(strings.TrimSpace(ptz.Args[0]), "https://chat.whatsapp.com/")
			jid, err := serialize.JoinGroupWithLink(ptz.Bot.Client, link)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Berhasil join ke " + jid.String())
		},
	})
}
