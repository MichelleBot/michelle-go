package ai

import (
	"fmt"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "deepseekr1",
		Aliases:     []string{"deepseek", "ds"},
		Description: "Chat dengan DeepSeek R1 AI",
		Usage:       "deepseekr1 <message>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			prompt := ptz.RawArgs

			if prompt == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !deepseekr1 <message> atau reply pesan")
				}
				prompt = replyText
			}

			ptz.React("🤖")
			defer ptz.Unreact()

			ctx, cancel := ptz.ContextWithTimeout(60 * time.Second)
			defer cancel()

			params := map[string]string{
				"prompt": prompt,
			}

			data, err := api.Request[AIResponse](ctx, ptz.Bot.API, "/api/ai/deepseekr1", params)
			if err != nil {
				ptz.Bot.Log.Errorf("DeepSeek error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan saat memproses permintaan.")
			}

			res := fmt.Sprintf("🤖 *DeepSeek R1*\n\n%s", data.Response)
			return ptz.ReplyText(res)
		},
	})
}
