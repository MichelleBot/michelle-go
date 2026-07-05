package random

import (
	"fmt"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

type QuoteAnimeData struct {
	Link     string `json:"link"`
	Gambar   string `json:"gambar"`
	Karakter string `json:"karakter"`
	Anime    string `json:"anime"`
	Episode  string `json:"episode"`
	Quotes   string `json:"quotes"`
}

func init() {
	core.Use(&core.Command{
		Name:        "quotesanime",
		Aliases:     []string{"quote", "qa"},
		Description: "Random anime quote dari siputzx",
		Usage:       "quotesanime",
		Category:    "random",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("✨")
			defer ptz.Unreact()

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			data, err := api.Request[[]QuoteAnimeData](ctx, ptz.Bot.API, "/api/r/quotesanime", nil)
			if err != nil {
				ptz.Bot.Log.Errorf("QuoteAnime error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server.")
			}

			if len(data) == 0 {
				return ptz.ReplyText("❌ Tidak ada data ditemukan.")
			}

			q := data[0]
			caption := fmt.Sprintf("💬 *%s*\n\n\"%s\"\n\n👤 %s\n📺 %s\n📍 %s\n🔗 %s", q.Anime, q.Quotes, q.Karakter, q.Anime, q.Episode, q.Link)

			if q.Gambar != "" {
				imgData, err := serialize.Fetch(q.Gambar)
				if err == nil {
					return ptz.ReplyImage(imgData, "image/jpeg", caption)
				}
			}

			return ptz.ReplyText(caption)
		},
	})
}
