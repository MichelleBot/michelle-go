package downloader

import (
	"fmt"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "savemedia",
		Description: "Simpan media dari pesan reply atau yang dikirim langsung",
		Usage:       "savemedia (media)",
		Category:    "downloader",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			input := serialize.GetInputMedia(ptz.Message)
			if input == nil {
				return ptz.ReplyText("Kirim atau reply pesan media yang ingin disimpan")
			}

			data, err := serialize.DownloadMedia(ptz.Bot.Client, input.Message)
			if err != nil {
				return ptz.ReplyText("Gagal download: " + err.Error())
			}

			filename := serialize.GetMediaFilename(input.Message)
			mime := serialize.GetMediaMIME(input.Message)
			caption := fmt.Sprintf("Saved media\nType: %s\nSize: %d bytes", mime, len(data))

			switch input.MsgType {
			case "image":
				return serialize.SendImageReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
			case "video":
				return serialize.SendVideoReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
			case "audio":
				return serialize.SendAudioReply(ptz.Bot.Client, ptz.Chat, data, mime, false, ptz.Message, ptz.Info)
			case "document":
				return serialize.SendDocumentReply(ptz.Bot.Client, ptz.Chat, data, mime, filename, caption, ptz.Message, ptz.Info)
			case "sticker":
				return serialize.SendStickerReply(ptz.Bot.Client, ptz.Chat, data, mime, false, ptz.Message, ptz.Info)
			}
			return ptz.ReplyText("Tipe media tidak didukung")
		},
	})
}
