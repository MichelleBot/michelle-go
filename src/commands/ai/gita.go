package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

type GitaResp struct {
	Status bool   `json:"status"`
	Data   string `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "gita",
		Aliases:     []string{"gitagpt"},
		Description: "Tanya GitaGPT",
		Usage:       "gita <question>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			question := ptz.RawArgs

			if question == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !gita <question> atau reply pesan")
				}
				question = replyText
			}
			ptz.React("🤖")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{"q": question}

			raw, err := client.Get(context.Background(), "/api/ai/gita", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result GitaResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Parse error")
			}

			if !result.Status {
				return ptz.ReplyText("❌ Failed")
			}

			text := fmt.Sprintf("🤖 *GitaGPT*\n\n%s", result.Data)
			return ptz.ReplyText(text)
		},
	})
}
