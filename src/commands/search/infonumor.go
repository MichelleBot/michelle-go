package search

import (
	"fmt"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

func init() {
	core.Use(&core.Command{
		Name:        "info",
		Description: "Dapatkan info nomor (nama, status, business)",
		Usage:       "info <nomor>",
		Category:    "search",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Format: .info <nomor>")
			}

			if err := ptz.React("⏳"); err != nil {
				ptz.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ptz.Unreact()

			phone := strings.Join(ptz.Args, " ")
			normalized := serialize.NormalizePhone(phone)
			jid := serialize.PhoneToJID(normalized)

			userInfo, err := serialize.GetUserInfo(ptz.Bot.Client, []types.JID{jid})
			if err != nil {
				return ptz.ReplyText("Error: " + err.Error())
			}

			info, ok := userInfo[jid]
			if !ok {
				return ptz.ReplyText("Info tidak ditemukan")
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Info Nomor: %s\n\n", normalized))

			if info.Status != "" {
				sb.WriteString(fmt.Sprintf("Status: %s\n", info.Status))
			}

			if info.PictureID != "" {
				sb.WriteString("Foto Profil: Ada\n")
			} else {
				sb.WriteString("Foto Profil: Tidak ada\n")
			}

			if len(info.Devices) > 0 {
				sb.WriteString(fmt.Sprintf("Jumlah Device: %d\n", len(info.Devices)))
			}

			return ptz.ReplyText(sb.String())
		},
	})
}
