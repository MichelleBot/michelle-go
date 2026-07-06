package converter

import (
	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"sticker", "s", "sk", "stiker", "sgif"},
		UsageHint: "reply media",
		Category:  "converter",
		Quota:     core.PerUserQuota(1),
		Handler:   runSticker,
	})
}

func runSticker(ptz *core.Ptz) error {
	input := serialize.GetInputMedia(ptz.Message, "image", "video")
	if input == nil {
		return ptz.ReplyText("❌ Kirim atau reply image/video yang ingin dijadikan sticker")
	}

	ptz.React("🕒")

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
}
