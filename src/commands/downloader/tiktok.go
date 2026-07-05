package downloader

import (
	"fmt"
	"regexp"
	"sptzx/src/serialize"
	"strconv"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

var tiktokRegex = regexp.MustCompile(`^https?://(www\.|vt\.|vm\.)?tiktok\.com/`)

func init() {
	core.Use(&core.Command{
		Name:        "tiktok",
		Aliases:     []string{"tiktokdl", "tkdl"},
		Description: "Download video TikTok",
		Usage:       "tiktok <url>",
		Category:    "downloader",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*tiktok* — Download video TikTok\n\nUsage: .tiktok <url>")
			}

			url := ptz.Args[0]
			if !tiktokRegex.MatchString(url) {
				return ptz.ReplyText("❌ URL TikTok tidak valid.")
			}

			if err := core.EnsureQuotaAvailable(ptz, 1); err != nil {
				return err
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			ctx, cancel := ptz.ContextWithTimeout(45 * time.Second)
			defer cancel()

			type TikTokV2Data struct {
				ItemID            string `json:"itemId"`
				NoWatermarkLink   string `json:"no_watermark_link"`
				NoWatermarkLinkHD string `json:"no_watermark_link_hd"`
				MusicLink         string `json:"music_link"`
				CoverLink         string `json:"cover_link"`
				AuthorNickname    string `json:"author_nickname"`
				AuthorUniqueID    string `json:"author_unique_id"`
				Duration          string `json:"duration"`
				LikeCount         string `json:"like_count"`
				PlayCount         string `json:"play_count"`
				CommentCount      string `json:"comment_count"`
				ShareCount        string `json:"share_count"`
			}

			data, err := api.Request[TikTokV2Data](ctx, ptz.Bot.API, "/api/d/tiktok/v2", map[string]string{"url": url})
			if err != nil {
				ptz.Bot.Log.Errorf("TikTok error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan saat mengambil data video.")
			}

			videoURL := data.NoWatermarkLinkHD
			quality := "HD"
			if videoURL == "" {
				videoURL = data.NoWatermarkLink
				quality = "SD"
			}
			if videoURL == "" {
				return ptz.ReplyText("❌ Tidak ada URL video.")
			}

			author := data.AuthorNickname
			if author == "" {
				author = data.AuthorUniqueID
			}

			durationMs, _ := strconv.ParseInt(data.Duration, 10, 64)
			durationSec := durationMs / 1000
			duration := fmt.Sprintf("%d:%02d", durationSec/60, durationSec%60)

			caption := "🎵 *TikTok*\n\n" +
				"👤 " + author + "\n" +
				"⏱ " + duration + "\n" +
				"❤️ " + data.LikeCount + "  ▶️ " + data.PlayCount + "\n" +
				"💬 " + data.CommentCount + "  🔁 " + data.ShareCount + "\n" +
				"🎬 " + quality

			vData, err := serialize.Fetch(videoURL)
			if err != nil {
				ptz.Bot.Log.Errorf("TikTok download error: %v", err)
				return ptz.ReplyText("❌ Gagal mendownload video.")
			}

			if err := ptz.ReplyVideo(vData, "video/mp4", caption); err != nil {
				return err
			}

			if err := core.ConsumeQuota(ptz, 1); err != nil {
				ptz.Bot.Log.Errorf("tiktok success quota consume error: %v", err)
				return err
			}

			return nil
		},
	})
}
