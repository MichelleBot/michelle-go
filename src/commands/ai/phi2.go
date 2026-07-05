package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

type Phi2Resp struct {
	Status bool `json:"status"`
	Data   struct {
		Response string `json:"response"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "phi2",
		Aliases:     []string{"phi"},
		Description: "Chat dengan Phi-2 AI",
		Usage:       "phi2 <message>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			prompt := ptz.RawArgs

			if prompt == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !phi2 <message> atau reply pesan")
				}
				prompt = replyText
			}

			ptz.React("🤖")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{
				"prompt": prompt,
			}

			raw, err := client.Get(context.Background(), "/api/ai/phi2", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result Phi2Resp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Parse error")
			}

			if !result.Status {
				return ptz.ReplyText("❌ Failed to get response")
			}

			text := fmt.Sprintf("🤖 *Phi-2 AI*\n\n%s", result.Data.Response)
			return ptz.ReplyText(text)
		},
	})
}
