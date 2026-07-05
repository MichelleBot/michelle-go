package stalk

import (
	"fmt"
	"sptzx/src/serialize"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "stalkyoutube",
		Aliases:     []string{"youtubestalk"},
		Description: "Stalk channel YouTube",
		Usage:       "stalkyoutube <username>",
		Category:    "stalk",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*stalkyoutube* — Stalk channel YouTube\n\nUsage: .stalkyoutube <username>")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			username := strings.Join(ptz.Args, "")
			username = strings.TrimPrefix(username, "@")

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type YoutubeData struct {
				Channel struct {
					Username        string `json:"username"`
					Name            string `json:"name"`
					SubscriberCount string `json:"subscriberCount"`
					VideoCount      string `json:"videoCount"`
					AvatarURL       string `json:"avatarUrl"`
					ChannelURL      string `json:"channelUrl"`
					Description     string `json:"description"`
				} `json:"channel"`
				LatestVideos []struct {
					Title         string `json:"title"`
					VideoURL      string `json:"videoUrl"`
					Duration      string `json:"duration"`
					ViewCount     string `json:"viewCount"`
					PublishedTime string `json:"publishedTime"`
				} `json:"latest_videos"`
			}

			data, err := api.Request[YoutubeData](ctx, ptz.Bot.API, "/api/stalk/youtube", map[string]string{"username": username})
			if err != nil {
				ptz.Bot.Log.Errorf("Youtube stalk error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			ch := data.Channel

			caption := fmt.Sprintf("📺 *YouTube*\n\n"+
				"👤 %s\n"+
				"🔗 %s\n\n"+
				"👥 *%s*\n"+
				"🎬 *%s*\n\n"+
				"📝 %s",
				ch.Username,
				ch.ChannelURL,
				ch.SubscriberCount,
				ch.VideoCount,
				ch.Description,
			)

			if len(data.LatestVideos) > 0 {
				caption += "\n\n🕐 *Video Terbaru:*"
				limit := 3
				if len(data.LatestVideos) < limit {
					limit = len(data.LatestVideos)
				}
				for i := 0; i < limit; i++ {
					v := data.LatestVideos[i]
					caption += fmt.Sprintf("\n• %s (%s) — %s", v.Title, v.Duration, v.ViewCount)
				}
			}

			imgData, err := serialize.Fetch(ch.AvatarURL)
			if err != nil {
				return ptz.ReplyText(caption)
			}
			return ptz.ReplyImage(imgData, "image/jpeg", caption)
		},
	})
}
