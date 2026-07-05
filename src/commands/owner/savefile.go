package owner

import (
	"os"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "savefile",
		Description: "Simpan file media dari quoted message ke server",
		Usage:       "savefile [filename]",
		Category:    "owner",
		OwnerOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			quoted := serialize.GetQuotedMessage(ptz.Message)
			if quoted == nil {
				return ptz.ReplyText("Reply pesan media yang ingin disimpan.")
			}
			if !serialize.IsMediaType(quoted) {
				return ptz.ReplyText("❌ Bukan pesan media.")
			}
			data, err := serialize.DownloadMedia(ptz.Bot.Client, quoted)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			filename := serialize.GetMediaFilename(quoted)
			if len(ptz.Args) > 0 {
				filename = strings.Join(ptz.Args, " ")
			}
			if err := os.MkdirAll("./downloads", 0755); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			path := "./downloads/" + filename
			if err := os.WriteFile(path, data, 0644); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Tersimpan:\n📁 " + path)
		},
	})
}
