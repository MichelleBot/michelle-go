package search

import (
	"fmt"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "cek",
		Description: "Cek apakah nomor terdaftar di WhatsApp",
		Usage:       "cek <nomor>",
		Category:    "search",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Format: .cek <nomor>")
			}

			if err := ptz.React("⏳"); err != nil {
				ptz.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ptz.Unreact()

			phone := strings.Join(ptz.Args, " ")
			normalized := serialize.NormalizePhone(phone)

			results, err := serialize.IsOnWhatsApp(ptz.Bot.Client, []string{"+" + normalized})
			if err != nil {
				return ptz.ReplyText("Error: " + err.Error())
			}

			if len(results) == 0 {
				return ptz.ReplyText("Nomor tidak ditemukan")
			}

			result := results[0]
			if result.IsIn {
				jid := result.JID.ToNonAD()
				msg := fmt.Sprintf("Nomor %s terdaftar di WhatsApp\nJID: %s", normalized, jid.String())
				return ptz.ReplyText(msg)
			}

			return ptz.ReplyText(fmt.Sprintf("Nomor %s tidak terdaftar di WhatsApp", normalized))
		},
	})
}
