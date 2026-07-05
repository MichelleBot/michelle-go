package sticker

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

type bratResponse struct {
	URL    string `json:"url"`
	Params struct {
		Animated bool `json:"animated"`
	} `json:"params"`
}

func init() {
	core.Use(&core.Command{
		Name:        "brat",
		Aliases:     []string{"bratsticker"},
		Description: "Buat sticker brat dari teks",
		Usage:       "brat <teks>  |  brat a <teks> (animated)",
		Category:    "sticker",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Masukkan teks.\nContoh: .brat halo\nAnimated: .brat a halo dunia")
			}

			animated := false
			text := ptz.RawArgs
			if ptz.Args[0] == "a" {
				animated = true
				if len(ptz.Args) < 2 {
					return ptz.ReplyText("Masukkan teks setelah 'a'.\nContoh: .brat a halo dunia")
				}
				text = ptz.RawArgs[2:]
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			apiURL := "https://bratgenerator.siputzx.my.id/?q=" + url.QueryEscape(text)
			if animated {
				apiURL += "&a=1"
			}

			apiBody, err := serialize.Fetch(apiURL)
			if err != nil {
				return ptz.ReplyText("❌ Gagal request API: " + err.Error())
			}

			var data bratResponse
			if err := json.Unmarshal(apiBody, &data); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}
			if data.URL == "" {
				return ptz.ReplyText(fmt.Sprintf("❌ URL sticker kosong."))
			}

			imgData, err := serialize.Fetch(data.URL)
			if err != nil {
				return ptz.ReplyText("❌ Gagal download sticker.")
			}

			if data.Params.Animated {
				trimmed, err := serialize.ToAnimatedWebpExif(imgData, ".webp", true, serialize.StickerMetadata{PackName: ptz.Bot.Config.StickerPackName, Author: ptz.Bot.Config.StickerAuthor, Categories: []string{""}})
				if err != nil {
					return ptz.ReplyText("❌ " + err.Error())
				}
				return ptz.ReplySticker(trimmed, "image/webp", true)
			}

			static, err := serialize.ToStaticWebpExif(imgData, ".webp", serialize.StickerMetadata{PackName: ptz.Bot.Config.StickerPackName, Author: ptz.Bot.Config.StickerAuthor, Categories: []string{""}})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplySticker(static, "image/webp", false)
		},
	})
}
