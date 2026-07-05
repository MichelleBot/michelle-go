package random

import (
	"fmt"
	"time"

	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "cecan",
		Aliases:     []string{"cc"},
		Description: "Random cecan images",
		Usage:       "cecan [country]",
		Category:    "random",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			countries := []string{"indonesia", "thailand", "vietnam", "china", "japan", "korea"}
			country := "indonesia"

			if len(ptz.Args) > 0 {
				for _, c := range countries {
					if ptz.Args[0] == c {
						country = c
						break
					}
				}
			}

			ptz.React("😍")
			defer ptz.Unreact()

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			endpoint := fmt.Sprintf("/api/r/cecan/%s", country)

			imgData, err := ptz.Bot.API.GetRaw(ctx, endpoint, nil)
			if err != nil {
				ptz.Bot.Log.Errorf("Cecan error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan saat mengambil gambar.")
			}

			caption := fmt.Sprintf("😍 *Cecan %s*", country)
			return ptz.ReplyImage(imgData, "image/jpeg", caption)
		},
	})
}
