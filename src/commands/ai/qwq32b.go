package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

type QwQResp struct {
	Status bool `json:"status"`
	Data   struct {
		Response string `json:"response"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "qwq32b",
		Aliases:     []string{"qwq", "reasoning"},
		Description: "Chat dengan QwQ 32B (Reasoning)",
		Usage:       "qwq32b <message>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			prompt := ptz.RawArgs

			if prompt == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !qwq32b <message> atau reply pesan")
				}
				prompt = replyText
			}

			ptz.React("🤖")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{"prompt": prompt}

			raw, err := client.Get(context.Background(), "/api/ai/qwq32b", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result QwQResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Parse error")
			}

			if !result.Status {
				return ptz.ReplyText("❌ Failed")
			}

			text := fmt.Sprintf("🤖 *QwQ 32B*\n\n%s", result.Data.Response)
			return ptz.ReplyText(text)
		},
	})
}
