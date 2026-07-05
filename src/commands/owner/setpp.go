package owner

import (
	"context"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"time"
)

func init() {
	core.Use(&core.Command{
		Name:        "setpp",
		Description: "Ubah foto profil bot",
		Usage:       "setpp (kirim/reply image)",
		Category:    "owner",
		OwnerOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			input := serialize.GetInputMedia(ptz.Message, "image")
			if input == nil {
				return ptz.ReplyText("Kirim atau reply image.")
			}
			data, err := serialize.DownloadMedia(ptz.Bot.Client, input.Message)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			mime := serialize.GetMediaMIME(input.Message)
			ext := serialize.GetMediaExtFromMIME(mime)
			jpeg, err := serialize.ToJPEG(data, ext)
			if err != nil {
				jpeg = data
			}
			jid := ptz.Bot.Client.Store.GetJID().ToNonAD()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			id, err := ptz.Bot.Client.SetGroupPhoto(ctx, jid, jpeg)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ PP berhasil diubah! ID: " + id)
		},
	})
}
