package maker

import (
	"context"
	"fmt"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	effects := []struct {
		Name string
		URL  string
		Desc string
	}{
		{"gradient", "https://textpro.me/create-a-gradient-text-shadow-effect-online-1141.html", "Gradient"},
		{"naruto", "https://textpro.me/create-naruto-logo-style-text-effect-online-1125.html", "Naruto"},
		{"glow", "https://textpro.me/create-light-glow-sliced-text-effect-online-1068.html", "Glow"},
		{"glitch", "https://textpro.me/create-impressive-glitch-text-effects-online-1027.html", "Glitch"},
		{"bear", "https://textpro.me/online-black-and-white-bear-mascot-logo-creation-1012.html", "Bear"},
	}

	for _, e := range effects {
		effect := e
		core.Use(&core.Command{
			Name:        effect.Name,
			Description: fmt.Sprintf("Generate %s effect", effect.Desc),
			Usage:       fmt.Sprintf("%s <text>", effect.Name),
			Category:    "maker",
			Quota:       core.PerUserQuota(1),
			Handler: func(ptz *core.Ptz) error {
				if len(ptz.Args) < 1 {
					return ptz.ReplyText(fmt.Sprintf("Format: .%s <text>", effect.Name))
				}
				ptz.React("✨")
				defer ptz.Unreact()

				client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
				params := map[string]string{
					"url":   effect.URL,
					"text1": ptz.Args[0],
				}

				imgData, err := client.GetRaw(context.Background(), "/api/m/textpro", params)
				if err != nil {
					return ptz.ReplyText("❌ " + err.Error())
				}
				return ptz.ReplyImage(imgData, "image/jpeg", fmt.Sprintf("✨ *%s Effect*\n_Text: %s_", effect.Desc, ptz.Args[0]))
			},
		})
	}
}
