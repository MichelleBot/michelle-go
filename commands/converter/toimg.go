package converter

import (
	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"toimg"},
		UsageHint: "reply sticker",
		Category:  "converter",
		Quota:     core.PerUserQuota(1),
		Handler:   runToImg,
	})
}

func runToImg(ptz *core.Ptz) error {
	if ptz.Message.ExtendedTextMessage == nil || ptz.Message.ExtendedTextMessage.ContextInfo == nil || ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage == nil {
		return ptz.ReplyText("🚩 Balas stiker yang ingin Anda ubah menjadi gambar.")
	}

	quotedMsg := ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage
	if quotedMsg.StickerMessage == nil {
		return ptz.ReplyText("🚩 Harap balas pesan stiker.")
	}

	ptz.React("🕒")

	data, err := serialize.DownloadMedia(ptz.Bot.Client, quotedMsg)
	if err != nil {
		return ptz.ReplyText("❌ Gagal mendownload stiker: " + err.Error())
	}

	imgData, err := serialize.ToJPEG(data, ".webp")
	if err != nil {
		return ptz.ReplyText("❌ Gagal mengubah stiker ke gambar: " + err.Error())
	}

	return ptz.ReplyImage(imgData, "image/jpeg", "")
}
