package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

type GLMResp struct {
	Status bool `json:"status"`
	Data   struct {
		Response string `json:"response"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "glm47flash",
		Aliases:     []string{"glm", "glmflash"},
		Description: "Chat dengan GLM 4.7 Flash",
		Usage:       "glm47flash <message>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			prompt := ptz.RawArgs

			if prompt == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !glm47flash <message> atau reply pesan")
				}
				prompt = replyText
			}

			ptz.React("🤖")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{"prompt": prompt}

			raw, err := client.Get(context.Background(), "/api/ai/glm47flash", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result GLMResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Parse error")
			}

			if !result.Status {
				return ptz.ReplyText("❌ Failed")
			}

			text := fmt.Sprintf("🤖 *GLM 4.7 Flash*\n\n%s", result.Data.Response)
			return ptz.ReplyText(text)
		},
	})
}
