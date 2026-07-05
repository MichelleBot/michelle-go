package tools

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "subdomain",
		Aliases:     []string{"subenum", "subscan"},
		Description: "Enumerate subdomain dari sebuah domain",
		Usage:       "subdomain <domain>",
		Category:    "tools",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*subdomain* — Enum subdomain\n\nUsage: .subdomain <domain>\nContoh: .subdomain google.com")
			}

			domain := ptz.Args[0]

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			client.SetLogger(ptz.Bot.Log)
			raw, err := client.GetRaw(ctx, "/api/tools/subdomains", map[string]string{
				"domain": domain,
			})
			if err != nil {
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			var resp struct {
				Status bool     `json:"status"`
				Data   []string `json:"data"`
			}

			if err := json.Unmarshal(raw, &resp); err != nil || !resp.Status {
				if err != nil {
					ptz.Bot.Log.Errorf("Subdomain unmarshal error: %v", err)
				}
				return ptz.ReplyText("❌ Gagal scan subdomain.")
			}

			seen := map[string]bool{}
			var unique []string
			for _, entry := range resp.Data {
				for _, line := range strings.Split(entry, "\n") {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "*") || seen[line] {
						continue
					}
					seen[line] = true
					unique = append(unique, line)
				}
			}

			if len(unique) == 0 {
				return ptz.ReplyText("❌ Tidak ada subdomain ditemukan.")
			}

			limit := 30
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("🌐 *Subdomain — %s*\n\n", domain))
			sb.WriteString(fmt.Sprintf("📊 Ditemukan *%d* subdomain\n\n", len(unique)))
			for i, s := range unique {
				if i >= limit {
					break
				}
				sb.WriteString(fmt.Sprintf("• %s\n", s))
			}
			if len(unique) > limit {
				sb.WriteString(fmt.Sprintf("\n_...dan %d lainnya_", len(unique)-limit))
			}
			return ptz.ReplyText(sb.String())
		},
	})
}
