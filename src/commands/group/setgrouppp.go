package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "setgrouppp",
		Aliases:     []string{"setgpp"},
		Description: "Ubah foto profil group",
		Usage:       "setgrouppp (kirim/reply image)",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
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
			if _, err := serialize.SetGroupPhoto(ptz.Bot.Client, ptz.Chat, jpeg); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Foto group berhasil diubah.")
		},
	})
}
