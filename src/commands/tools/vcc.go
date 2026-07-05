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
		Name:        "vcc",
		Aliases:     []string{"genvcc", "virtualcard"},
		Description: "Generate virtual credit card",
		Usage:       "vcc [type] [jumlah]",
		Category:    "tools",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			validTypes := map[string]string{
				"visa":       "Visa",
				"mastercard": "MasterCard",
				"amex":       "Amex",
				"cup":        "CUP",
				"jcb":        "JCB",
				"diners":     "Diners",
				"rupay":      "RuPay",
			}

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*vcc* — Generate virtual credit card\n\nUsage: .vcc [type] [jumlah]\nContoh: .vcc visa 3\n\nTipe: Visa, MasterCard, Amex, CUP, JCB, Diners, RuPay\nMaks: 5 kartu")
			}

			cardType := "Visa"
			count := "1"

			if len(ptz.Args) >= 1 {
				if v, ok := validTypes[strings.ToLower(ptz.Args[0])]; ok {
					cardType = v
				}
			}
			if len(ptz.Args) >= 2 {
				count = ptz.Args[1]
			}

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type VCCData []struct {
				CardNumber     string `json:"cardNumber"`
				ExpirationDate string `json:"expirationDate"`
				CardholderName string `json:"cardholderName"`
				CVV            string `json:"cvv"`
			}

			data, err := api.Request[VCCData](ctx, ptz.Bot.API, "/api/tools/vcc-generator", map[string]string{
				"type":  cardType,
				"count": count,
			})
			if err != nil {
				ptz.Bot.Log.Errorf("VCC error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("💳 *Virtual Credit Card (%s)*\n\n", cardType))
			for i, c := range data {
				sb.WriteString(fmt.Sprintf("*%d.*\n", i+1))
				sb.WriteString(fmt.Sprintf("🔢 `%s`\n", c.CardNumber))
				sb.WriteString(fmt.Sprintf("👤 %s\n", c.CardholderName))
				sb.WriteString(fmt.Sprintf("📅 %s  🔒 %s\n\n", c.ExpirationDate, c.CVV))
			}
			sb.WriteString("_⚠️ Hanya untuk keperluan testing_")
			return ptz.ReplyText(sb.String())
		},
	})
}
