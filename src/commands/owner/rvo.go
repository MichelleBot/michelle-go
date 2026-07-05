package owner

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "rvo",
		Aliases:     []string{"revealvo", "viewonce"},
		Description: "Buka pesan view once (image/video/voice note)",
		Usage:       "rvo (reply pesan view once)",
		Category:    "owner",
		OwnerOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			if ptz.Event.IsViewOnce {
				msg := ptz.Event.Message
				if msg == nil {
					return ptz.ReplyText("Tidak ada pesan view once ditemukan")
				}

				msgType := serialize.GetMessageType(msg)
				if msgType != "image" && msgType != "video" && msgType != "audio" {
					return ptz.ReplyText("Tipe view once tidak didukung: " + msgType)
				}

				if err := ptz.React("⏳"); err != nil {
					ptz.Bot.Log.Debugf("Failed to react: %v", err)
				}
				defer ptz.Unreact()

				data, err := serialize.DownloadMedia(ptz.Bot.Client, msg)
				if err != nil {
					return ptz.ReplyText("Gagal download view once: " + err.Error())
				}
				mime := serialize.GetMediaMIME(msg)
				caption := serialize.GetMediaCaption(msg)

				switch msgType {
				case "image":
					return serialize.SendImageReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
				case "video":
					return serialize.SendVideoReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
				case "audio":
					ptt := msg.AudioMessage != nil && msg.AudioMessage.PTT != nil && *msg.AudioMessage.PTT
					return serialize.SendAudioReply(ptz.Bot.Client, ptz.Chat, data, mime, ptt, ptz.Message, ptz.Info)
				}
				return nil
			}

			quoted := serialize.GetQuotedMessage(ptz.Message)
			if quoted == nil {
				return ptz.ReplyText("Reply pesan view once yang ingin dibuka")
			}
			msgType := serialize.GetMessageType(quoted)
			if msgType != "image" && msgType != "video" && msgType != "audio" {
				return ptz.ReplyText("Pesan yang di-reply bukan image, video, atau voice note")
			}

			if err := ptz.React("⏳"); err != nil {
				ptz.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ptz.Unreact()

			data, err := serialize.DownloadMedia(ptz.Bot.Client, quoted)
			if err != nil {
				return ptz.ReplyText("Gagal download media: " + err.Error())
			}
			mime := serialize.GetMediaMIME(quoted)
			caption := serialize.GetMediaCaption(quoted)

			switch msgType {
			case "image":
				return serialize.SendImageReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
			case "video":
				return serialize.SendVideoReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
			case "audio":
				ptt := quoted.AudioMessage != nil && quoted.AudioMessage.PTT != nil && *quoted.AudioMessage.PTT
				return serialize.SendAudioReply(ptz.Bot.Client, ptz.Chat, data, mime, ptt, ptz.Message, ptz.Info)
			}
			return nil
		},
	})
}
