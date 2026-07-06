package converter

import (
	"net/url"

	"michelle/system/core"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"fakeff"},
		UsageHint: "text",
		Category:  "converter",
		Quota:     core.PerUserQuota(1),
		Handler:   runFakeFF,
	})
}

func runFakeFF(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText("🚩 Masukkan teks.")
	}

	text := ptz.RawArgs
	if len(text) > 20 {
		return ptz.ReplyText("🚩 Maksimal 20 karakter.")
	}

	ptz.React("🕒")

	apiUrl := "https://satriacanvas.vercel.app/fake-ff?usr=" + url.QueryEscape(text)
	data, err := utils.FetchAsBuffer(apiUrl)
	if err != nil {
		return ptz.ReplyText("❌ Gagal membuat gambar: " + err.Error())
	}

	return ptz.ReplyImage(data, "image/jpeg", "")
}
