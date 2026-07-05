package downloader

import (
	"context"
	"encoding/json"
	"regexp"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

var twitterRegex = regexp.MustCompile(`^https?://(www\.)?(twitter\.com|x\.com)/`)

type twitterMedia struct {
	URL     string `json:"url"`
	Ext     string `json:"ext"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Quality int    `json:"quality"`
	Subname string `json:"subname"`
}

type twitterMeta struct {
	Title  string `json:"title"`
	Source string `json:"source"`
}

type twitterResp struct {
	ID   string         `json:"id"`
	URL  []twitterMedia `json:"url"`
	Meta twitterMeta    `json:"meta"`
	HD   *struct {
		URL string `json:"url"`
	} `json:"hd"`
	Thumb string `json:"thumb"`
}

func init() {
	core.Use(&core.Command{
		Name:        "twitter",
		Aliases:     []string{"x", "twitterdl"},
		Description: "Download video Twitter/X",
		Usage:       "twitter <url>",
		Category:    "downloader",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*twitter* — Download video Twitter/X\n\nUsage: .twitter <url>")
			}

			url := ptz.Args[0]
			if !twitterRegex.MatchString(url) {
				return ptz.ReplyText("❌ URL Twitter/X tidak valid.")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/d/ssstwiter", map[string]string{"url": url})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result twitterResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			if len(result.URL) == 0 {
				return ptz.ReplyText("❌ Video tidak ditemukan.")
			}

			var videoURL string
			var quality string

			if result.HD != nil && result.HD.URL != "" {
				videoURL = result.HD.URL
				quality = "HD"
			} else {
				best := 0
				for _, m := range result.URL {
					if m.Quality > best {
						best = m.Quality
						videoURL = m.URL
						quality = m.Subname + "p"
					}
				}
			}

			caption := "🐦 *Twitter / X*\n\n" +
				"📝 " + result.Meta.Title + "\n" +
				"🎬 " + quality

			data, err := serialize.Fetch(videoURL)
			if err != nil {
				return ptz.ReplyText("❌ Gagal download: " + err.Error())
			}

			return ptz.ReplyVideo(data, "video/mp4", caption)
		},
	})
}
