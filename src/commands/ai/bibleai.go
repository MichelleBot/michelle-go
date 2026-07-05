package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

type BibleResp struct {
	Status bool `json:"status"`
	Data   struct {
		Question string `json:"question"`
		Results  struct {
			Answer string `json:"answer"`
		} `json:"results"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "bibleai",
		Aliases:     []string{"bible"},
		Description: "Tanya BibleAI",
		Usage:       "bibleai <question>",
		Category:    "ai",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			question := ptz.RawArgs

			if question == "" {
				replyText := ptz.GetReplyText()
				if replyText == "" {
					return ptz.ReplyText("Format: !bibleai <question> atau reply pesan")
				}
				question = replyText
			}
			ptz.React("🤖")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{"question": question}

			raw, err := client.Get(context.Background(), "/api/ai/bibleai", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result BibleResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Parse error")
			}

			if !result.Status {
				return ptz.ReplyText("❌ Failed")
			}

			text := fmt.Sprintf("🤖 *BibleAI*\n\nQ: %s\n\nA: %s", result.Data.Question, result.Data.Results.Answer)
			return ptz.ReplyText(text)
		},
	})
}
