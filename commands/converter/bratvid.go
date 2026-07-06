package converter

import (
	"net/url"
	"time"

	"michelle/system/core"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"bratvid"},
		UsageHint: "text",
		Category:  "converter",
		Quota:     core.PerUserQuota(1),
		Limit:     core.PerUserLimit(3, time.Minute),
		Handler:   runBratVid,
	})
}

func runBratVid(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText("🚩 Masukkan teks.")
	}

	text := ptz.RawArgs
	if len(text) > 100 {
		return ptz.ReplyText("🚩 Maksimal 100 karakter.")
	}

	ptz.React("🕒")

	// Use URL encoding for the text
	apiUrl := "https://brat.siputzx.my.id/mp4?text=" + url.QueryEscape(text)
	data, err := utils.FetchAsBuffer(apiUrl)
	if err != nil {
		return ptz.ReplyText("❌ Gagal membuat video: " + err.Error())
	}
    
	return ptz.ReplySticker(data, "image/webp", true)
}
