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
		Name:        "ramaljodoh",
		Aliases:     []string{"jodoh"},
		Description: "Ramalan jodoh Jawa berdasarkan nama & tanggal lahir",
		Usage:       "ramaljodoh <nama1> <tgl> <bln> <thn> | <nama2> <tgl> <bln> <thn>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			parts := strings.Split(ptz.RawArgs, "|")
			if len(parts) != 2 {
				return ptz.ReplyText("*ramaljodoh* — Ramalan jodoh Jawa\n\nUsage: .ramaljodoh <nama1> <tgl> <bln> <thn> | <nama2> <tgl> <bln> <thn>\nContoh: .ramaljodoh putu 16 11 2007 | keyla 1 1 2008")
			}

			p1 := strings.Fields(strings.TrimSpace(parts[0]))
			p2 := strings.Fields(strings.TrimSpace(parts[1]))

			if len(p1) < 4 || len(p2) < 4 {
				return ptz.ReplyText("❌ Format tidak lengkap. Pastikan: nama tgl bln thn")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/ramalanjodoh", map[string]string{
				"nama1": p1[0], "tgl1": p1[1], "bln1": p1[2], "thn1": p1[3],
				"nama2": p2[0], "tgl2": p2[1], "bln2": p2[2], "thn2": p2[3],
			})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var resp struct {
				Status bool `json:"status"`
				Data   struct {
					Result struct {
						OrangPertama struct {
							Nama         string `json:"nama"`
							TanggalLahir string `json:"tanggal_lahir"`
						} `json:"orang_pertama"`
						OrangKedua struct {
							Nama         string `json:"nama"`
							TanggalLahir string `json:"tanggal_lahir"`
						} `json:"orang_kedua"`
						HasilRamalan []string `json:"hasil_ramalan"`
					} `json:"result"`
					Peringatan string `json:"peringatan"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			r := resp.Data.Result
			var sb strings.Builder
			sb.WriteString("💕 *Ramalan Jodoh*\n\n")
			sb.WriteString(fmt.Sprintf("👤 *%s* — %s\n", r.OrangPertama.Nama, r.OrangPertama.TanggalLahir))
			sb.WriteString(fmt.Sprintf("👤 *%s* — %s\n\n", r.OrangKedua.Nama, r.OrangKedua.TanggalLahir))
			sb.WriteString("🔮 *Hasil Ramalan:*\n")
			for i, h := range r.HasilRamalan {
				sb.WriteString(fmt.Sprintf("\n%d. %s", i+1, h))
			}
			if resp.Data.Peringatan != "" {
				sb.WriteString("\n\n⚠️ " + resp.Data.Peringatan)
			}
			return ptz.ReplyText(sb.String())
		},
	})
}
