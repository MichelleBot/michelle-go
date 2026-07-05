package primbon

import (
	"context"
	"encoding/json"
	"sptzx/src/api"
	"sptzx/src/core"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "artinama",
		Aliases:     []string{"namaartik", "nama"},
		Description: "Arti dan karakter berdasarkan nama",
		Usage:       "artinama <nama>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*artinama* — Arti nama\n\nUsage: .artinama <nama>\nContoh: .artinama putu")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/artinama", map[string]string{
				"nama": strings.Join(ptz.Args, " "),
			})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var resp struct {
				Status bool `json:"status"`
				Data   struct {
					Nama string `json:"nama"`
					Arti string `json:"arti"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Nama tidak ditemukan.")
			}

			d := resp.Data
			var sb strings.Builder
			sb.WriteString("📖 *Arti Nama*\n\n")
			sb.WriteString("👤 *" + d.Nama + "*\n\n")
			sb.WriteString(d.Arti)
			return ptz.ReplyText(sb.String())
		},
	})
}
