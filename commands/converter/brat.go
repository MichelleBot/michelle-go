package converter

import (
	"net/url"

	"michelle/system/core"
	"michelle/system/serialize"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"brat"},
		UsageHint: "text",
		Category:  "converter",
		Quota:     core.PerUserQuota(1),
		Handler:   runBrat,
	})
}

func runBrat(ptz *core.Ptz) error {
	if ptz.RawArgs == "" {
		return ptz.ReplyText(utils.Example(ptz.Prefix, "brat", "michelle"))
	}

	if len(ptz.RawArgs) > 100 {
		return ptz.ReplyText("🚩 Maksimal 100 karakter.")
	}

	ptz.React("🕒")

	apiUrl := "https://brat.siputzx.my.id/image?text=" + url.QueryEscape(ptz.RawArgs)
	data, err := utils.FetchAsBuffer(apiUrl)
	if err != nil {
		return ptz.ReplyText("❌ Gagal membuat stiker: " + err.Error())
	}

	meta := serialize.StickerMetadata{
		PackName:   ptz.Bot.Config.StickerPackName,
		Author:     ptz.Bot.Config.StickerAuthor,
		Categories: []string{""},
	}

	webp, err := serialize.ToStaticWebpExif(data, ".png", meta)
	if err != nil {
		return ptz.ReplyText("❌ Gagal convert ke sticker: " + err.Error())
	}

	return ptz.ReplySticker(webp, "image/webp", false)
}
