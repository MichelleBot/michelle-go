package stalk

import (
	"fmt"
	"sptzx/src/serialize"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "stalkthreads",
		Aliases:     []string{"threadsstalk"},
		Description: "Stalk profil Threads",
		Usage:       "stalkthreads <username>",
		Category:    "stalk",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*stalkthreads* — Stalk profil Threads\n\nUsage: .stalkthreads <username>")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			username := strings.Join(ptz.Args, "")
			username = strings.TrimPrefix(username, "@")

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type ThreadsData struct {
				ID               string   `json:"id"`
				Username         string   `json:"username"`
				Name             string   `json:"name"`
				Bio              string   `json:"bio"`
				ProfilePicture   string   `json:"profile_picture"`
				HdProfilePicture string   `json:"hd_profile_picture"`
				IsVerified       bool     `json:"is_verified"`
				Followers        float64  `json:"followers"`
				Links            []string `json:"links"`
			}

			data, err := api.Request[ThreadsData](ctx, ptz.Bot.API, "/api/stalk/threads", map[string]string{"q": username})
			if err != nil {
				ptz.Bot.Log.Errorf("Threads stalk error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			verified := ""
			if data.IsVerified {
				verified = " ✅"
			}

			caption := fmt.Sprintf("🧵 *Threads*\n\n"+
				"👤 %s (@%s)%s\n",
				data.Name, data.Username, verified)

			if data.Bio != "" {
				caption += "📝 " + data.Bio + "\n"
			}

			caption += fmt.Sprintf("\n👥 *%s* Followers", serialize.NumFmt(data.Followers))

			if len(data.Links) > 0 {
				caption += "\n🔗 " + data.Links[0]
			}

			avatarURL := data.HdProfilePicture
			if avatarURL == "" {
				avatarURL = data.ProfilePicture
			}

			imgData, err := serialize.Fetch(avatarURL)
			if err != nil {
				return ptz.ReplyText(caption)
			}
			return ptz.ReplyImage(imgData, "image/jpeg", caption)
		},
	})
}
