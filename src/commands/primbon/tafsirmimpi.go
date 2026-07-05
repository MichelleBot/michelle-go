package primbon

import (
	"context"
	"encoding/json"
	"fmt"
	"sptzx/src/api"
	"sptzx/src/core"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "tafsirmimpi",
		Aliases:     []string{"mimpi"},
		Description: "Tafsir mimpi menurut primbon",
		Usage:       "tafsirmimpi <kata kunci>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*tafsirmimpi* — Tafsir mimpi\n\nUsage: .tafsirmimpi <kata kunci>\nContoh: .tafsirmimpi ular")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/tafsirmimpi", map[string]string{
				"mimpi": strings.Join(ptz.Args, " "),
			})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var resp struct {
				Status bool `json:"status"`
				Data   struct {
					Keyword string `json:"keyword"`
					Hasil   []struct {
						Mimpi  string `json:"mimpi"`
						Tafsir string `json:"tafsir"`
					} `json:"hasil"`
					Total int `json:"total"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Tidak ditemukan.")
			}

			d := resp.Data
			if len(d.Hasil) == 0 {
				return ptz.ReplyText("❌ Tidak ada tafsir untuk: " + d.Keyword)
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("🌙 *Tafsir Mimpi*\n\n🔍 *%s* — %d hasil\n\n", d.Keyword, d.Total))
			limit := 10
			if len(d.Hasil) < limit {
				limit = len(d.Hasil)
			}
			for i := 0; i < limit; i++ {
				h := d.Hasil[i]
				sb.WriteString(fmt.Sprintf("*%d. %s*\n%s\n\n", i+1, h.Mimpi, h.Tafsir))
			}
			if d.Total > limit {
				sb.WriteString(fmt.Sprintf("_...dan %d lainnya_", d.Total-limit))
			}
			return ptz.ReplyText(sb.String())
		},
	})
}
