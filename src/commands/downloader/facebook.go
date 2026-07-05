package downloader

import (
	"context"
	"encoding/json"
	"regexp"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

var facebookRegex = regexp.MustCompile(`^https?://(www\.|m\.)?facebook\.com/|^https?://fb\.watch/`)

type fbMediaURL struct {
	URL     string `json:"url"`
	Ext     string `json:"ext"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Subname string `json:"subname"`
}

type fbMeta struct {
	Title    string `json:"title"`
	Source   string `json:"source"`
	Duration string `json:"duration"`
}

type fbVideoInfo struct {
	URL  []fbMediaURL `json:"url"`
	Meta fbMeta       `json:"meta"`
}

type fbData struct {
	Type string        `json:"type"`
	Data []fbVideoInfo `json:"data"`
}

type fbResp struct {
	Success bool     `json:"success"`
	Data    []fbData `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "facebook",
		Aliases:     []string{"fb", "fbdl"},
		Description: "Download video Facebook",
		Usage:       "facebook <url>",
		Category:    "downloader",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*facebook* — Download video Facebook\n\nUsage: .facebook <url>")
			}

			url := ptz.Args[0]
			if !facebookRegex.MatchString(url) {
				return ptz.ReplyText("❌ URL Facebook tidak valid.")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/d/savefrom", map[string]string{"url": url})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result fbResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			if !result.Success || len(result.Data) == 0 {
				return ptz.ReplyText("❌ Video tidak ditemukan.")
			}

			var videoItem fbVideoInfo
			for _, d := range result.Data {
				if d.Type == "video" && len(d.Data) > 0 {
					videoItem = d.Data[0]
					break
				}
			}

			if len(videoItem.URL) == 0 {
				return ptz.ReplyText("❌ Tidak ada URL download.")
			}

			var videoURL, quality string
			for _, m := range videoItem.URL {
				if m.Subname == "HD" {
					videoURL = m.URL
					quality = "HD"
					break
				}
			}
			if videoURL == "" {
				videoURL = videoItem.URL[0].URL
				quality = videoItem.URL[0].Subname
				if quality == "" {
					quality = "SD"
				}
			}

			title := videoItem.Meta.Title
			duration := videoItem.Meta.Duration

			caption := "📘 *Facebook*\n\n" +
				"📝 " + title + "\n" +
				"⏱ " + duration + "\n" +
				"🎬 " + quality

			data, err := serialize.Fetch(videoURL)
			if err != nil {
				return ptz.ReplyText("❌ Gagal download: " + err.Error())
			}

			return ptz.ReplyVideo(data, "video/mp4", caption)
		},
	})
}
