package downloader

import (
	"regexp"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

var instagramRegex = regexp.MustCompile(`^https?://(www\.)?instagram\.com/(reel|p|tv)/`)

func init() {
	core.Use(&core.Command{
		Name:        "igdl",
		Aliases:     []string{"instagram"},
		Description: "Download video/foto Instagram",
		Usage:       "igdl <url>",
		Category:    "downloader",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*igdl* — Download media Instagram\n\nUsage: .igdl <url>")
			}

			url := ptz.Args[0]
			if !instagramRegex.MatchString(url) {
				return ptz.ReplyText("❌ URL Instagram tidak valid.")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			ctx, cancel := ptz.ContextWithTimeout(45 * time.Second)
			defer cancel()

			type IGMediaURL struct {
				URL  string `json:"url"`
				Name string `json:"name"`
				Type string `json:"type"`
				Ext  string `json:"ext"`
			}

			type IGVideoInfo struct {
				URL  []IGMediaURL `json:"url"`
				Meta struct {
					Title        string `json:"title"`
					Username     string `json:"username"`
					LikeCount    int64  `json:"like_count"`
					CommentCount int64  `json:"comment_count"`
				} `json:"meta"`
			}

			type IGData struct {
				Data []IGVideoInfo `json:"data"`
			}

			data, err := api.Request[[]IGData](ctx, ptz.Bot.API, "/api/d/savefrom", map[string]string{"url": url})
			if err != nil {
				ptz.Bot.Log.Errorf("Instagram error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan saat mengambil media Instagram.")
			}

			if len(data) == 0 || len(data[0].Data) == 0 {
				return ptz.ReplyText("❌ Media tidak ditemukan.")
			}

			mediaItem := data[0].Data[0]
			if len(mediaItem.URL) == 0 {
				return ptz.ReplyText("❌ Tidak ada URL download.")
			}

			mediaURL := mediaItem.URL[0].URL
			mediaType := mediaItem.URL[0].Type

			caption := "📸 *Instagram*\n\n"
			if mediaItem.Meta.Username != "" {
				caption += "👤 @" + mediaItem.Meta.Username + "\n"
			}
			if mediaItem.Meta.Title != "" {
				caption += "📝 " + mediaItem.Meta.Title + "\n"
			}
			caption += "❤️ " + serialize.NumFmt64(mediaItem.Meta.LikeCount) + "  💬 " + serialize.NumFmt64(mediaItem.Meta.CommentCount)

			mediaData, err := serialize.Fetch(mediaURL)
			if err != nil {
				ptz.Bot.Log.Errorf("Instagram download error: %v", err)
				return ptz.ReplyText("❌ Gagal mendownload media.")
			}

			switch mediaType {
			case "mp4", "video_mp4", "video":
				return ptz.ReplyVideo(mediaData, "video/mp4", caption)
			default:
				return ptz.ReplyImage(mediaData, "image/jpeg", caption)
			}
		},
	})
}
