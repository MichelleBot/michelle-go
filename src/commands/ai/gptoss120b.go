package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

type GPTOSSResp struct {
	Status bool `json:"status"`
	Data   struct {
		Response string `json:"response"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "gptoss120b",
		Aliases:     []string{"gpt120b", "gpto"},
		Description: "Chat dengan GPT OSS 120B",
		Usage:       "gptoss120b <message>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			prompt := ptz.RawArgs

			if prompt == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !gptoss120b <message> atau reply pesan")
				}
				prompt = replyText
			}

			ptz.React("🤖")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{"prompt": prompt}

			raw, err := client.Get(context.Background(), "/api/ai/gptoss120b", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result GPTOSSResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Parse error")
			}

			if !result.Status {
				return ptz.ReplyText("❌ Failed")
			}

			text := fmt.Sprintf("🤖 *GPT OSS 120B*\n\n%s", result.Data.Response)
			return ptz.ReplyText(text)
		},
	})
}
