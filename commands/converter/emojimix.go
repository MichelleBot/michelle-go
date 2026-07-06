package converter

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"michelle/system/core"
	"michelle/system/serialize"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"emojimix", "mix", "emomix"},
		UsageHint: "emoji + emoji",
		Category:  "converter",
		Quota:     core.PerUserQuota(1),
		Handler:   runEmojimix,
	})
}

func runEmojimix(ptz *core.Ptz) error {
	if ptz.RawArgs == "" {
		return ptz.ReplyText(fmt.Sprintf("• *Example* : %semojimix 😳+😩", ptz.Bot.GetPrefix()))
	}

	parts := strings.Split(ptz.RawArgs, "+")
	if len(parts) < 2 {
		return ptz.ReplyText("🚩 Berikan 2 emoji untuk dimix (contoh: 😳+😩).")
	}

	ptz.React("🕒")

	emo1 := strings.TrimSpace(parts[0])
	emo2 := strings.TrimSpace(parts[1])

	apiUrl := fmt.Sprintf("https://tenor.googleapis.com/v2/featured?key=AIzaSyAyimkuYQYF_FXVALexPuGQctUWRURdCYQ&contentfilter=high&media_filter=png_transparent&component=proactive&collection=emoji_kitchen_v5&q=%s_%s",
		url.QueryEscape(emo1), url.QueryEscape(emo2))

	res, err := utils.FetchAsJSON(apiUrl)
	if err != nil {
		return ptz.ReplyText("❌ Gagal mengambil data: " + err.Error())
	}
	
	results, ok := res["results"].([]interface{})
	if !ok || len(results) == 0 {
		return ptz.ReplyText("🚩 Emoji tidak bisa dimix.")
	}

	for _, result := range results {
		r := result.(map[string]interface{})
		url := r["url"].(string)
		data, err := utils.FetchAsBuffer(url)
		if err != nil {
			continue
		}

		meta := serialize.StickerMetadata{
			PackName:   ptz.Bot.Config.StickerPackName,
			Author:     ptz.Bot.Config.StickerAuthor,
			Categories: []string{emo1, emo2},
		}

		webp, err := serialize.ToStaticWebpExif(data, ".png", meta)
		if err != nil {
			continue
		}

		ptz.ReplySticker(webp, "image/webp", false)
		time.Sleep(1500 * time.Millisecond)
	}

	return nil
}
