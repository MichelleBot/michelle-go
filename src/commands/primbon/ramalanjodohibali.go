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
		Name:        "ramaljodohibali",
		Aliases:     []string{"jodohibali"},
		Description: "Ramalan jodoh Bali berdasarkan nama & tanggal lahir",
		Usage:       "ramaljodohibali <nama1> <tgl> <bln> <thn> | <nama2> <tgl> <bln> <thn>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			parts := strings.Split(ptz.RawArgs, "|")
			if len(parts) != 2 {
				return ptz.ReplyText("*ramaljodohibali* — Ramalan jodoh Bali\n\nUsage: .ramaljodohibali <nama1> <tgl> <bln> <thn> | <nama2> <tgl> <bln> <thn>\nContoh: .ramaljodohibali putu 16 11 2007 | keyla 1 1 2008")
			}

			p1 := strings.Fields(strings.TrimSpace(parts[0]))
			p2 := strings.Fields(strings.TrimSpace(parts[1]))

			if len(p1) < 4 || len(p2) < 4 {
				return ptz.ReplyText("❌ Format tidak lengkap.")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/ramalanjodohbali", map[string]string{
				"nama1": p1[0], "tgl1": p1[1], "bln1": p1[2], "thn1": p1[3],
				"nama2": p2[0], "tgl2": p2[1], "bln2": p2[2], "thn2": p2[3],
			})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var resp struct {
				Status bool `json:"status"`
				Data   struct {
					NamaAnda struct {
						Nama     string `json:"nama"`
						TglLahir string `json:"tgl_lahir"`
					} `json:"nama_anda"`
					NamaPasangan struct {
						Nama     string `json:"nama"`
						TglLahir string `json:"tgl_lahir"`
					} `json:"nama_pasangan"`
					Result  string `json:"result"`
					Catatan string `json:"catatan"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			d := resp.Data
			var sb strings.Builder
			sb.WriteString("🌺 *Ramalan Jodoh Bali*\n\n")
			sb.WriteString(fmt.Sprintf("👤 *%s* — %s\n", d.NamaAnda.Nama, d.NamaAnda.TglLahir))
			sb.WriteString(fmt.Sprintf("👤 *%s* — %s\n\n", d.NamaPasangan.Nama, d.NamaPasangan.TglLahir))
			sb.WriteString("🔮 *Hasil:*\n" + d.Result)
			return ptz.ReplyText(sb.String())
		},
	})
}
