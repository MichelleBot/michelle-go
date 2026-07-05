package downloader

import (
	"context"
	"encoding/json"
	"regexp"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

var soundcloudRegex = regexp.MustCompile(`^https?://(www\.|m\.)?soundcloud\.com/`)

type scMediaURL struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
	Ext  string `json:"ext"`
}

type scMeta struct {
	Source   string `json:"source"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

type scInfo struct {
	ID   string       `json:"id"`
	URL  []scMediaURL `json:"url"`
	Meta scMeta       `json:"meta"`
}

type scData struct {
	Type string `json:"type"`
	Data scInfo `json:"data"`
}

type scResp struct {
	Success bool     `json:"success"`
	Data    []scData `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "soundcloud",
		Aliases:     []string{"sc", "scdl"},
		Description: "Download audio SoundCloud",
		Usage:       "soundcloud <url>",
		Category:    "downloader",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*soundcloud* — Download audio SoundCloud\n\nUsage: .soundcloud <url>")
			}

			url := ptz.Args[0]
			if !soundcloudRegex.MatchString(url) {
				return ptz.ReplyText("❌ URL SoundCloud tidak valid.")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/d/savefrom", map[string]string{"url": url})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result scResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			if !result.Success || len(result.Data) == 0 {
				return ptz.ReplyText("❌ Audio tidak ditemukan.")
			}

			var audioInfo scInfo
			for _, d := range result.Data {
				if d.Type == "audio" {
					audioInfo = d.Data
					break
				}
			}

			if len(audioInfo.URL) == 0 {
				return ptz.ReplyText("❌ Tidak ada URL download.")
			}

			audioURL := audioInfo.URL[0].URL

			audioData, err := serialize.Fetch(audioURL)
			if err != nil {
				return ptz.ReplyText("❌ Gagal download: " + err.Error())
			}

			return ptz.ReplyAudio(audioData, "audio/mpeg")
		},
	})
}
