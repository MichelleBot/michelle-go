package ai

import (
	"fmt"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

type AIResponse struct {
	Response string `json:"response"`
}

func init() {
	core.Use(&core.Command{
		Name:        "gemini",
		Aliases:     []string{"gm"},
		Description: "Chat dengan Gemini AI",
		Usage:       "gemini <message>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Limit:       core.PerUserLimit(8, time.Minute),
		Handler: func(ptz *core.Ptz) error {
			text := ptz.RawArgs

			if text == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !gemini <message> atau reply pesan")
				}
				text = replyText
			}
			ptz.React("🤖")
			defer ptz.Unreact()

			cookie := ptz.Bot.Config.SiputzX.GeminiCookie
			if cookie == "" {
				return ptz.ReplyText("❌ Gemini cookie tidak diset di .env")
			}

			ctx, cancel := ptz.ContextWithTimeout(60 * time.Second)
			defer cancel()

			params := map[string]string{
				"text":   text,
				"cookie": cookie,
			}

			data, err := api.Request[AIResponse](ctx, ptz.Bot.API, "/api/ai/gemini", params)
			if err != nil {
				ptz.Bot.Log.Errorf("Gemini error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan saat memproses permintaan.")
			}

			res := fmt.Sprintf("🤖 *Gemini AI*\n\n%s", data.Response)
			return ptz.ReplyText(res)
		},
	})
}
