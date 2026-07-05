package primbon

import (
	"context"
	"encoding/json"
	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "usahabisnis",
		Aliases:     []string{"bisnis", "sifatbisnis"},
		Description: "Sifat usaha bisnis berdasarkan tanggal lahir",
		Usage:       "usahabisnis <tgl> <bln> <thn>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) < 3 {
				return ptz.ReplyText("*usahabisnis* — Sifat usaha bisnis\n\nUsage: .usahabisnis <tgl> <bln> <thn>\nContoh: .usahabisnis 1 1 2000")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/sifat_usaha_bisnis", map[string]string{
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
					HariLahir string `json:"hari_lahir"`
					Usaha     string `json:"usaha"`
					Catatan   string `json:"catatan"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			d := resp.Data
			return ptz.ReplyText("💼 *Sifat Usaha Bisnis*\n\n" +
				"📅 " + d.HariLahir + "\n\n" +
				d.Usaha)
		},
	})
}
