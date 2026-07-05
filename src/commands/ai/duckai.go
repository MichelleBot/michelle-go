package ai

import (
	"fmt"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "duckai",
		Aliases:     []string{"dai", "ai"},
		Description: "Chat dengan DuckAI",
		Usage:       "duckai <message>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			message := ptz.RawArgs

			if message == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !duckai <message> atau reply pesan")
				}
				message = replyText
			}

			ptz.React("🤖")
			defer ptz.Unreact()

			ctx, cancel := ptz.ContextWithTimeout(60 * time.Second)
			defer cancel()

			params := map[string]string{
				"message": message,
			}

			type DuckAIData struct {
				Message string `json:"message"`
				Model   string `json:"model"`
			}

			data, err := api.Request[DuckAIData](ctx, ptz.Bot.API, "/api/ai/duckai", params)
			if err != nil {
				ptz.Bot.Log.Errorf("DuckAI error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan saat memproses permintaan.")
			}

			res := fmt.Sprintf("🤖 *DuckAI*\n\n%s\n\n_Model: %s_", data.Message, data.Model)
			return ptz.ReplyText(res)
		},
	})
}
