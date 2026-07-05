package tools

import (
	"fmt"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "kodepos",
		Aliases:     []string{"postalcode", "zip"},
		Description: "Cari kode pos berdasarkan nama desa/kelurahan",
		Usage:       "kodepos <nama desa>",
		Category:    "tools",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*kodepos* — Cari kode pos\n\nUsage: .kodepos <nama desa>\nContoh: .kodepos pasiran jaya")
			}

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type KodeposItem struct {
				Kodepos   string `json:"kodepos"`
				Desa      string `json:"desa"`
				Kecamatan string `json:"kecamatan"`
				Kota      string `json:"kota"`
				Provinsi  string `json:"provinsi"`
			}

			data, err := api.Request[[]KodeposItem](ctx, ptz.Bot.API, "/api/tools/kodepos", map[string]string{
				"form": strings.Join(ptz.Args, " "),
			})
			if err != nil {
				ptz.Bot.Log.Errorf("Kodepos error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			if len(data) == 0 {
				return ptz.ReplyText("❌ Kode pos tidak ditemukan.")
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("📮 *Kode Pos*\n\n"))
			for i, d := range data {
				if i >= 10 {
					sb.WriteString(fmt.Sprintf("\n_...dan %d hasil lainnya_", len(data)-10))
					break
				}
				sb.WriteString(fmt.Sprintf("*%d.* `%s`\n", i+1, d.Kodepos))
				sb.WriteString(fmt.Sprintf("   🏘 %s, %s\n", d.Desa, d.Kecamatan))
				sb.WriteString(fmt.Sprintf("   🏙 %s, %s\n\n", d.Kota, d.Provinsi))
			}
			return ptz.ReplyText(strings.TrimSpace(sb.String()))
		},
	})
}
