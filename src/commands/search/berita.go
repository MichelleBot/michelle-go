package search

import (
	"context"
	"fmt"
	"strings"

	"sptzx/src/api"
	"sptzx/src/core"
)

var beritaSources = map[string]string{
	"antara":    "/api/search/berita/antara",
	"cnbc":      "/api/search/berita/cnbc",
	"cnn":       "/api/search/berita/cnn",
	"jpnn":      "/api/search/berita/jpnn",
	"kumparan":  "/api/search/berita/kumparan",
	"merdeka":   "/api/search/berita/merdeka",
	"okezone":   "/api/search/berita/okezone",
	"republika": "/api/search/berita/republika",
	"sindo":     "/api/search/berita/sindonews",
	"tempo":     "/api/search/berita/tempo",
	"terang":    "/api/search/berita/terang",
	"tribun":    "/api/search/berita/tribun",
}

func formatBeritaHelp() string {
	res := "*Berita* — Cari berita terbaru\n\nUsage: .berita <sumber>\n\nSumber tersedia:\n"
	for name := range beritaSources {
		res += fmt.Sprintf("- %s\n", name)
	}
	return res
}

func init() {
	core.Use(&core.Command{
		Name:        "berita",
		Aliases:     []string{"news"},
		Description: "Mencari berita terbaru dari berbagai sumber",
		Usage:       "berita <sumber>",
		Category:    "search",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText(formatBeritaHelp())
			}

			sumber := strings.ToLower(ptz.Args[0])
			endpoint, ok := beritaSources[sumber]
			if !ok {
				return ptz.ReplyText("❌ Sumber berita tidak valid.\n\n" + formatBeritaHelp())
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			res, err := api.Request[[]struct {
				Title string `json:"title"`
				Link  string `json:"link"`
			}](context.Background(), client, endpoint, nil)

			if err != nil {
				return ptz.ReplyText("❌ Gagal mengambil berita: " + err.Error())
			}

			if len(res) == 0 {
				return ptz.ReplyText("❌ Tidak ada berita terbaru.")
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("📰 *Berita Terbaru dari %s*\n\n", strings.ToUpper(sumber)))

			for i, item := range res {
				if i >= 10 {
					break
				}
				sb.WriteString(fmt.Sprintf("%d. *%s*\n", i+1, item.Title))
				sb.WriteString(fmt.Sprintf("🔗 %s\n\n", item.Link))
			}

			return ptz.ReplyText(sb.String())
		},
	})
}
