package stalk

import (
	"fmt"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "stalkpinterest",
		Aliases:     []string{"pintereststalk"},
		Description: "Stalk profil Pinterest",
		Usage:       "stalkpinterest <username>",
		Category:    "stalk",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*stalkpinterest* — Stalk profil Pinterest\n\nUsage: .stalkpinterest <username>")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			username := strings.Join(ptz.Args, "")
			username = strings.TrimPrefix(username, "@")

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type PinterestData struct {
				ID         string  `json:"id"`
				Username   string  `json:"username"`
				FullName   string  `json:"full_name"`
				Bio        string  `json:"bio"`
				ProfileURL string  `json:"profile_url"`
				CreatedAt  string  `json:"created_at"`
				Website    *string `json:"website"`
				Location   *string `json:"location"`
				Country    *string `json:"country"`
				Image      struct {
					Original string `json:"original"`
					Large    string `json:"large"`
				} `json:"image"`
				Stats struct {
					Pins      float64 `json:"pins"`
					Followers float64 `json:"followers"`
					Following float64 `json:"following"`
					Boards    float64 `json:"boards"`
					Likes     float64 `json:"likes"`
				} `json:"stats"`
				Meta struct {
					FirstName string  `json:"first_name"`
					LastName  string  `json:"last_name"`
					Gender    *string `json:"gender"`
				} `json:"meta"`
			}

			data, err := api.Request[PinterestData](ctx, ptz.Bot.API, "/api/stalk/pinterest", map[string]string{"q": username})
			if err != nil {
				ptz.Bot.Log.Errorf("Pinterest stalk error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			caption := fmt.Sprintf("📌 *Pinterest*\n\n"+
				"👤 %s (@%s)\n",
				data.FullName, data.Username)

			if data.Bio != "" {
				caption += "📝 " + data.Bio + "\n"
			}
			if data.Location != nil && *data.Location != "" {
				caption += "📍 " + *data.Location + "\n"
			}
			if data.Website != nil && *data.Website != "" {
				caption += "🔗 " + *data.Website + "\n"
			}

			caption += fmt.Sprintf("\n"+
				"📌 *%s* Pins\n"+
				"👥 *%s* Followers\n"+
				"➕ *%s* Following\n"+
				"🗂 *%s* Boards",
				serialize.NumFmt(data.Stats.Pins),
				serialize.NumFmt(data.Stats.Followers),
				serialize.NumFmt(data.Stats.Following),
				serialize.NumFmt(data.Stats.Boards),
			)

			if data.CreatedAt != "" {
				caption += "\n📅 Joined: " + data.CreatedAt
			}

			avatarURL := data.Image.Original
			if avatarURL == "" {
				avatarURL = data.Image.Large
			}

			imgData, err := serialize.Fetch(avatarURL)
			if err != nil {
				return ptz.ReplyText(caption)
			}
			return ptz.ReplyImage(imgData, "image/jpeg", caption)
		},
	})
}
