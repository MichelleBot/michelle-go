package primbon

import (
	"context"
	"encoding/json"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "kecocokannama",
		Aliases:     []string{"cocokannama", "namacocok"},
		Description: "Cek kecocokan nama pasangan",
		Usage:       "kecocokannama <nama1> <nama2>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) < 2 {
				return ptz.ReplyText("*kecocokannama* — Kecocokan nama pasangan\n\nUsage: .kecocokannama <nama1> <nama2>\nContoh: .kecocokannama putu keyla")
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/kecocokan_nama_pasangan", map[string]string{
				"nama1": ptz.Args[0],
				"nama2": ptz.Args[1],
			})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var resp struct {
				Status bool `json:"status"`
				Data   struct {
					NamaAnda     string `json:"nama_anda"`
					NamaPasangan string `json:"nama_pasangan"`
					SisiPositif  string `json:"sisi_positif"`
					SisiNegatif  string `json:"sisi_negatif"`
					Gambar       string `json:"gambar"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Gagal parse response.")
			}

			d := resp.Data
			var sb strings.Builder
			sb.WriteString("💑 *Kecocokan Nama*\n\n")
			sb.WriteString("👤 " + d.NamaAnda + " & " + d.NamaPasangan + "\n\n")
			sb.WriteString("✅ *Sisi Positif:*\n" + d.SisiPositif + "\n\n")
			sb.WriteString("❌ *Sisi Negatif:*\n" + d.SisiNegatif)

			caption := sb.String()

			if d.Gambar != "" {
				img, err := serialize.Fetch(d.Gambar)
				if err == nil {
					return ptz.ReplyImage(img, "image/png", caption)
				}
			}

			return ptz.ReplyText(caption)
		},
	})
}
