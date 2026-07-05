package maker

import (
	"context"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "blackpink",
		Aliases:     []string{"bp"},
		Description: "Generate Blackpink style logo",
		Usage:       "blackpink <text>",
		Category:    "maker",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) < 1 {
				return ptz.ReplyText("Format: .blackpink <text>")
			}
			ptz.React("✨")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			params := map[string]string{
				"url":  "https://photooxy.com/create-blackpink-style-logo-effects-online-for-free-417.html",
				"text": ptz.Args[0],
			}

			imgData, err := client.GetRaw(context.Background(), "/api/m/photooxy", params)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyImage(imgData, "image/jpeg", fmt.Sprintf("✨ *Blackpink Effect*\n_Text: %s_", ptz.Args[0]))
		},
	})
}
