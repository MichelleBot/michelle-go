package downloader

import (
	"context"
	"encoding/json"
	"regexp"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

var douyinRegex = regexp.MustCompile(`^https?://(www\.|v\.)?douyin\.com/`)

type douyinDownload struct {
	Quality string `json:"quality"`
	URL     string `json:"url"`
}

type douyinData struct {
	Title     string           `json:"title"`
	Thumbnail string           `json:"thumbnail"`
	Downloads []douyinDownload `json:"downloads"`
}

type douyinResp struct {
	Status bool       `json:"status"`
	Data   douyinData `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "douyin",
		Aliases:     []string{"dy", "douyindl"},
		Description: "Download video Douyin",
		Usage:       "douyin <url>",
		Category:    "downloader",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*douyin* — Download video Douyin\n\nUsage: .douyin <url>")
			}

			url := ptz.Args[0]
			if !douyinRegex.MatchString(url) {
				return ptz.ReplyText("❌ URL Douyin tidak valid.")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/d/douyin", map[string]string{"url": url})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result douyinResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			if !result.Status || len(result.Data.Downloads) == 0 {
				return ptz.ReplyText("❌ Video tidak ditemukan.")
			}

			videoURL := result.Data.Downloads[0].URL

			caption := "🎬 *Douyin*\n\n" +
				"📝 " + result.Data.Title

			data, err := serialize.Fetch(videoURL)
			if err != nil {
				return ptz.ReplyText("❌ Gagal download: " + err.Error())
			}

			return ptz.ReplyVideo(data, "video/mp4", caption)
		},
	})
}
