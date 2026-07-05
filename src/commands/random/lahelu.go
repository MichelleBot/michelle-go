package random

import (
	"fmt"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

type LaheluData struct {
	Title        string `json:"title"`
	TotalUpvotes int    `json:"totalUpvotes"`
	CreateTime   string `json:"createTime"`
	UserInfo     struct {
		Username string `json:"username"`
	} `json:"userInfo"`
	Content []struct {
		Type  int    `json:"type"`
		Value string `json:"value"`
	} `json:"content"`
	Media          string `json:"media"`
	MediaType      int    `json:"mediaType"`
	MediaThumbnail string `json:"mediaThumbnail"`
}

func init() {
	core.Use(&core.Command{
		Name:        "lahelu",
		Aliases:     []string{"lhl"},
		Description: "Random post dari Lahelu",
		Usage:       "lahelu",
		Category:    "random",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("😂")
			defer ptz.Unreact()

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			data, err := api.Request[[]LaheluData](ctx, ptz.Bot.API, "/api/r/lahelu", nil)
			if err != nil {
				ptz.Bot.Log.Errorf("Lahelu error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server.")
			}

			if len(data) == 0 {
				return ptz.ReplyText("❌ Tidak ada data ditemukan.")
			}

			p := data[0]
			caption := fmt.Sprintf("😂 *%s*\n👤 %s\n⬆️ %d\n📅 %s", p.Title, p.UserInfo.Username, p.TotalUpvotes, p.CreateTime)

			if p.Media != "" {
				media, err := serialize.Fetch(p.Media)
				if err == nil {
					if p.MediaType == 1 {
						return ptz.ReplyVideo(media, "video/mp4", caption)
					} else if p.MediaType == 2 {
						return ptz.ReplyImage(media, "image/webp", caption)
					}
				}
			}

			return ptz.ReplyText(caption)
		},
	})
}
