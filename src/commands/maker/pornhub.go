package maker

import (
	"context"
	"fmt"
	"strings"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "pornhub",
		Aliases:     []string{"ph"},
		Description: "Generate PornHub style logo",
		Usage:       "pornhub <text1> | <text2>",
		Category:    "maker",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			parts := strings.Split(ptz.RawArgs, "|")
			if len(parts) < 2 {
				return ptz.ReplyText("Format: .pornhub text1 | text2")
			}
			text1 := strings.TrimSpace(parts[0])
			text2 := strings.TrimSpace(parts[1])

			ptz.React("✨")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{
				"url":   "https://en.ephoto360.com/create-pornhub-style-logos-online-free-549.html",
				"text1": text1,
				"text2": text2,
			}

			imgData, err := client.GetRaw(context.Background(), "/api/m/ephoto360", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyImage(imgData, "image/jpeg", fmt.Sprintf("✨ *PornHub Effect*\n_Text: %s %s_", text1, text2))
		},
	})
}
