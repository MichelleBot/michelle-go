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
		Name:        "stalkgithub",
		Aliases:     []string{"githubstalk"},
		Description: "Stalk profil GitHub",
		Usage:       "stalkgithub <username>",
		Category:    "stalk",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*stalkgithub* — Stalk profil GitHub\n\nUsage: .stalkgithub <username>")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			username := strings.Join(ptz.Args, "")

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type GithubData struct {
				Username   string  `json:"username"`
				Nickname   string  `json:"nickname"`
				Bio        *string `json:"bio"`
				ProfilePic string  `json:"profile_pic"`
				URL        string  `json:"url"`
				Type       string  `json:"type"`
				Company    string  `json:"company"`
				Blog       string  `json:"blog"`
				Location   string  `json:"location"`
				Email      *string `json:"email"`
				PublicRepo int     `json:"public_repo"`
				PublicGist int     `json:"public_gists"`
				Followers  int     `json:"followers"`
				Following  int     `json:"following"`
				CreatedAt  string  `json:"created_at"`
				UpdatedAt  string  `json:"updated_at"`
			}

			data, err := api.Request[GithubData](ctx, ptz.Bot.API, "/api/stalk/github", map[string]string{"user": username})
			if err != nil {
				ptz.Bot.Log.Errorf("GitHub stalk error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			caption := fmt.Sprintf("🐙 *GitHub*\n\n"+
				"👤 %s (@%s)\n",
				data.Nickname, data.Username)

			if data.Bio != nil && *data.Bio != "" {
				caption += "📝 " + *data.Bio + "\n"
			}
			caption += "\n"
			if data.Company != "" {
				caption += "🏢 " + data.Company + "\n"
			}
			if data.Location != "" {
				caption += "📍 " + data.Location + "\n"
			}
			if data.Blog != "" {
				caption += "🔗 " + data.Blog + "\n"
			}
			if data.Email != nil && *data.Email != "" {
				caption += "✉️ " + *data.Email + "\n"
			}

			caption += fmt.Sprintf("\n"+
				"📦 *%d* Repos  •  🌟 *%d* Gists\n"+
				"👥 *%d* Followers  •  ➕ *%d* Following",
				data.PublicRepo, data.PublicGist,
				data.Followers, data.Following,
			)

			if len(data.CreatedAt) >= 10 {
				caption += "\n📅 Joined: " + data.CreatedAt[:10]
			}
			caption += "\n🌐 " + data.URL

			imgData, err := serialize.Fetch(data.ProfilePic)
			if err != nil {
				return ptz.ReplyText(caption)
			}
			return ptz.ReplyImage(imgData, "image/jpeg", caption)
		},
	})
}
