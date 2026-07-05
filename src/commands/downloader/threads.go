package downloader

import (
	"context"
	"encoding/json"
	"regexp"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

var threadsRegex = regexp.MustCompile(`^https?://(www\.)?(threads\.net|threads\.com|instagram\.com)/`)

type threadsMediaURL struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
	Ext  string `json:"ext"`
}

type threadsMeta struct {
	Title     string `json:"title"`
	Source    string `json:"source"`
	Shortcode string `json:"shortcode"`
	LikeCount int64  `json:"like_count"`
	TakenAt   int64  `json:"taken_at"`
}

type threadsData struct {
	URL   []threadsMediaURL `json:"url"`
	Meta  threadsMeta       `json:"meta"`
	Thumb string            `json:"thumb"`
}

type threadsResp struct {
	Status bool        `json:"status"`
	Data   threadsData `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "threads",
		Aliases:     []string{"threadsdl"},
		Description: "Download media Threads",
		Usage:       "threads <url>",
		Category:    "downloader",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*threads* — Download media Threads\n\nUsage: .threads <url>")
			}

			url := ptz.Args[0]
			if !threadsRegex.MatchString(url) {
				return ptz.ReplyText("❌ URL Threads tidak valid.")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/d/ummy", map[string]string{"url": url})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result threadsResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			if !result.Status || len(result.Data.URL) == 0 {
				return ptz.ReplyText("❌ Media tidak ditemukan.")
			}

			mediaURL := result.Data.URL[0].URL
			mediaType := result.Data.URL[0].Type
			likes := serialize.NumFmt64(result.Data.Meta.LikeCount)

			caption := "🧵 *Threads*\n\n" +
				"📝 " + result.Data.Meta.Title + "\n" +
				"❤️ " + likes

			mediaData, err := serialize.Fetch(mediaURL)
			if err != nil {
				return ptz.ReplyText("❌ Gagal download: " + err.Error())
			}

			if mediaType == "mp4" {
				return ptz.ReplyVideo(mediaData, "video/mp4", caption)
			}
			return ptz.ReplyImage(mediaData, "image/jpeg", caption)
		},
	})
}
