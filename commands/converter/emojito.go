package converter

import (
	"time"

	"michelle/system/core"
	"michelle/system/serialize"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"emojito"},
		UsageHint: "emoji",
		Category:  "converter",
		Limit:     core.PerUserLimit(1, 24*time.Hour),
		Handler:   runEmojito,
	})
}

func runEmojito(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText("🚩 Contoh: " + ptz.Prefix + ptz.Command + " 😳")
	}

	ptz.React("🕒")

	json, err := core.Api.Michelle("/emojito", map[string]string{
		"q": ptz.Args[0],
	})

	if err != nil {
		return ptz.ReplyText("❌ Terjadi kesalahan: " + err.Error())
	}

	if status, ok := json["status"].(bool); !ok || !status {
		msg := "Gagal memproses"
		if m, ok := json["msg"].(string); ok {
			msg = m
		}
		return ptz.ReplyText("🚩 " + msg)
	}

	data := json["data"].(map[string]interface{})
	url := data["url"].(string)

	buffer, err := utils.FetchAsBuffer(url)
	if err != nil {
		return ptz.ReplyText("❌ Gagal mendownload stiker: " + err.Error())
	}

	meta := serialize.StickerMetadata{
		PackName:   ptz.Bot.Config.StickerPackName,
		Author:     ptz.Bot.Config.StickerAuthor,
		Categories: []string{ptz.Args[0]},
	}

	webp, err := serialize.AddExifToWebp(buffer, meta)
	if err != nil {
		return ptz.ReplyText("❌ Gagal convert ke sticker: " + err.Error())
	}

	return ptz.ReplySticker(webp, "image/webp", false)
}
