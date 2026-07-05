package sticker

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "sticker",
		Aliases:     []string{"s", "stiker"},
		Description: "Buat sticker dari image atau video (reply atau kirim langsung)",
		Usage:       "sticker (image/video)",
		Category:    "sticker",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			input := serialize.GetInputMedia(ptz.Message, "image", "video")
			if input == nil {
				return ptz.ReplyText("Kirim atau reply image/video yang ingin dijadikan sticker")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			data, err := serialize.DownloadMedia(ptz.Bot.Client, input.Message)
			if err != nil {
				return ptz.ReplyText("❌ Gagal download: " + err.Error())
			}

			mime := serialize.GetMediaMIME(input.Message)
			ext := serialize.GetMediaExtFromMIME(mime)
			meta := serialize.StickerMetadata{
				PackName:   ptz.Bot.Config.StickerPackName,
				Author:     ptz.Bot.Config.StickerAuthor,
				Categories: []string{""},
			}

			if input.MsgType == "video" {
				webp, err := serialize.ToAnimatedWebpExif(data, ext, true, meta)
				if err != nil {
					return ptz.ReplyText("❌ Gagal convert video ke sticker: " + err.Error())
				}
				return ptz.ReplySticker(webp, "image/webp", true)
			}

			webp, err := serialize.ToStaticWebpExif(data, ext, meta)
			if err != nil {
				return ptz.ReplyText("❌ Gagal convert image ke sticker: " + err.Error())
			}
			return ptz.ReplySticker(webp, "image/webp", false)
		},
	})
}
