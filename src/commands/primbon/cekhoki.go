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
		Name:        "cekhoki",
		Aliases:     []string{"nomorhoki", "hokinom"},
		Description: "Cek energi dan keberuntungan nomor telepon",
		Usage:       "cekhoki <nomor>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*cekhoki* — Cek nomor hoki\n\nUsage: .cekhoki <nomor>\nContoh: .cekhoki 085658939117")
			}

			phone := strings.NewReplacer("+", "", " ", "", "-", "").Replace(strings.Join(ptz.Args, ""))
			if strings.HasPrefix(phone, "0") {
				phone = "62" + phone[1:]
			}

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/primbon/nomorhoki", map[string]string{
				"phoneNumber": phone,
			})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var resp struct {
				Status bool `json:"status"`
				Data   struct {
					Nomor          string `json:"nomor"`
					AnkaBaguaShuzi struct {
						Value float64 `json:"value"`
					} `json:"angka_bagua_shuzi"`
					EnergiPositif struct {
						Total   float64 `json:"total"`
						Details struct {
							Kekayaan   float64 `json:"kekayaan"`
							Kesehatan  float64 `json:"kesehatan"`
							Cinta      float64 `json:"cinta"`
							Kestabilan float64 `json:"kestabilan"`
						} `json:"details"`
					} `json:"energi_positif"`
					EnergiNegatif struct {
						Total   float64 `json:"total"`
						Details struct {
							Perselisihan float64 `json:"perselisihan"`
							Kehilangan   float64 `json:"kehilangan"`
							Malapetaka   float64 `json:"malapetaka"`
							Kehancuran   float64 `json:"kehancuran"`
						} `json:"details"`
					} `json:"energi_negatif"`
					Analisis struct {
						Status bool `json:"status"`
					} `json:"analisis"`
				} `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				return ptz.ReplyText("❌ Gagal analisis nomor.")
			}

			d := resp.Data
			hoki := "❌ Kurang beruntung"
			if d.Analisis.Status {
				hoki = "✅ *HOKI!*"
			}

			var sb strings.Builder
			sb.WriteString("🔮 *Cek Nomor Hoki*\n\n")
			sb.WriteString(fmt.Sprintf("📱 *%s*\n", d.Nomor))
			sb.WriteString(fmt.Sprintf("📊 Bagua Shuzi: *%.2f%%*\n\n", d.AnkaBaguaShuzi.Value))
			sb.WriteString(fmt.Sprintf("✨ *Energi Positif (%.1f%%):*\n", d.EnergiPositif.Total))
			pos := d.EnergiPositif.Details
			sb.WriteString(fmt.Sprintf("💰 Kekayaan: %.0f  ❤️ Cinta: %.0f\n", pos.Kekayaan, pos.Cinta))
			sb.WriteString(fmt.Sprintf("💊 Kesehatan: %.0f  ⚖️ Kestabilan: %.0f\n\n", pos.Kesehatan, pos.Kestabilan))
			sb.WriteString(fmt.Sprintf("⚠️ *Energi Negatif (%.1f%%):*\n", d.EnergiNegatif.Total))
			neg := d.EnergiNegatif.Details
			sb.WriteString(fmt.Sprintf("💢 Perselisihan: %.0f  💸 Kehilangan: %.0f\n", neg.Perselisihan, neg.Kehilangan))
			sb.WriteString(fmt.Sprintf("☄️ Malapetaka: %.0f  💀 Kehancuran: %.0f\n\n", neg.Malapetaka, neg.Kehancuran))
			sb.WriteString("📝 Status: " + hoki)
			return ptz.ReplyText(sb.String())
		},
	})
}
