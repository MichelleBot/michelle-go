package converter

import (
	"fmt"
	"math/rand"
	"net/url"

	"michelle/system/core"
	"michelle/system/serialize"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"fakeml"},
		UsageHint: "text",
		Category:  "converter",
		Handler:   runFakeML,
	})
}

func runFakeML(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText("🚩 Masukkan nama pengguna.")
	}
	text := ptz.RawArgs
	if len(text) > 15 {
		return ptz.ReplyText("🚩 Maksimal 15 karakter.")
	}

	ptz.React("🕒")

	// Get profile picture URL
	pic, err := serialize.GetProfilePicture(ptz.Bot.Client, ptz.Sender)
	avatarURL := ""
	if err == nil && pic != nil {
		avatarURL = pic.URL
	}

	border := rand.Intn(16) + 1
	apiUrl := fmt.Sprintf("https://satriacanvas.vercel.app/fake-ml?usr=%s&rank=imo&border=%d&lobby_type=indo&avatar=%s", 
        url.QueryEscape(text), border, url.QueryEscape(avatarURL))

	data, err := utils.FetchAsBuffer(apiUrl)
	if err != nil {
		return ptz.ReplyText("❌ Gagal membuat gambar: " + err.Error())
	}

	return ptz.ReplyImage(data, "image/jpeg", "")
}
