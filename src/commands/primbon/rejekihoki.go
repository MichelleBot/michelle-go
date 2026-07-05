package primbon

import (
	"context"
	"encoding/json"
	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "rejekihoki",
		Aliases:     []string{"rejeki", "weton"},
		Description: "Rejeki dan hoki berdasarkan weton kelahiran",
		Usage:       "rejekihoki <tgl> <bln> <thn>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) < 3 {
				return ptz.ReplyText("*rejekihoki* — Rejeki & hoki weton\n\nUsage: .rejekihoki <tgl> <bln> <thn>\nContoh: .rejekihoki 1 1 2025")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/rejeki_hoki_weton", map[string]string{
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
					Rejeki    string `json:"rejeki"`
					Catatan   string `json:"catatan"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			d := resp.Data
			return ptz.ReplyText("💰 *Rejeki & Hoki Weton*\n\n" +
				"📅 " + d.HariLahir + "\n\n" +
				d.Rejeki)
		},
	})
}
