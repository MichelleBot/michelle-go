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
		Name:        "stalktiktok",
		Aliases:     []string{"tiktokstalk"},
		Description: "Stalk profil TikTok",
		Usage:       "stalktiktok <username>",
		Category:    "stalk",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*stalktiktok* — Stalk profil TikTok\n\nUsage: .stalktiktok <username>")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			username := strings.Join(ptz.Args, "")
			username = strings.TrimPrefix(username, "@")

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type TikTokData struct {
				User struct {
					UniqueID  string `json:"uniqueId"`
					Nickname  string `json:"nickname"`
					Signature string `json:"signature"`
					Verified  bool   `json:"verified"`
					Avatar    string `json:"avatarLarger"`
					BioLink   struct {
						Link string `json:"link"`
					} `json:"bioLink"`
				} `json:"user"`
				Stats struct {
					FollowerCount  float64 `json:"followerCount"`
					FollowingCount float64 `json:"followingCount"`
					HeartCount     float64 `json:"heartCount"`
					VideoCount     float64 `json:"videoCount"`
					DiggCount      float64 `json:"diggCount"`
					FriendCount    float64 `json:"friendCount"`
				} `json:"stats"`
			}

			data, err := api.Request[TikTokData](ctx, ptz.Bot.API, "/api/stalk/tiktok", map[string]string{"username": username})
			if err != nil {
				ptz.Bot.Log.Errorf("TikTok error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			u := data.User
			s := data.Stats

			verified := ""
			if u.Verified {
				verified = " ✅"
			}

			caption := fmt.Sprintf("🎵 *TikTok*\n\n"+
				"👤 %s (@%s)%s\n"+
				"📝 %s\n\n"+
				"👥 *%s* Followers\n"+
				"➕ *%s* Following\n"+
				"❤️ *%s* Likes\n"+
				"🎬 *%s* Videos\n"+
				"👫 *%s* Friends",
				u.Nickname, u.UniqueID, verified,
				u.Signature,
				serialize.NumFmt(s.FollowerCount),
				serialize.NumFmt(s.FollowingCount),
				serialize.NumFmt(s.HeartCount),
				serialize.NumFmt(s.VideoCount),
				serialize.NumFmt(s.FriendCount),
			)

			if u.BioLink.Link != "" {
				caption += "\n🔗 " + u.BioLink.Link
			}

			imgData, err := serialize.Fetch(u.Avatar)
			if err != nil {
				return ptz.ReplyText(caption)
			}
			return ptz.ReplyImage(imgData, "image/jpeg", caption)
		},
	})
}
