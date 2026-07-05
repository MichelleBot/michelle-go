package primbon

import (
	"context"
	"encoding/json"
	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "potensipenyakit",
		Aliases:     []string{"penyakit", "cekhesalat"},
		Description: "Potensi penyakit berdasarkan tanggal lahir",
		Usage:       "potensipenyakit <tgl> <bln> <thn>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) < 3 {
				return ptz.ReplyText("*potensipenyakit* — Potensi penyakit\n\nUsage: .potensipenyakit <tgl> <bln> <thn>\nContoh: .potensipenyakit 12 5 1998")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/cek_potensi_penyakit", map[string]string{
				"tgl": ptz.Args[0],
				"bln": ptz.Args[1],
				"thn": ptz.Args[2],
			})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var resp struct {
				Status bool `json:"status"`
				Data   struct {
					Analisa string `json:"analisa"`
					Sektor  string `json:"sektor"`
					Elemen  string `json:"elemen"`
					Catatan string `json:"catatan"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			d := resp.Data
			return ptz.ReplyText("🏥 *Potensi Penyakit*\n\n" +
				"🔬 *Analisa:*\n" + d.Analisa + "\n\n" +
				"⚡ *Elemen:*\n" + d.Elemen)
		},
	})
}
