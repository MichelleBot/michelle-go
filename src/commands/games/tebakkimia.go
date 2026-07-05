package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakKimiaResp struct {
	Data struct {
		Unsur   string `json:"unsur"`
		Lambang string `json:"lambang"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakkimia",
		Aliases:     []string{"tki"},
		Description: "Game tebak kimia",
		Usage:       "tebakkimia",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakkimia",
				Endpoint: "/api/games/tebakkimia",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakKimiaResp
					err = json.Unmarshal(raw, &res)
					soal = fmt.Sprintf("Apa lambang kimia dari unsur *%s*?", res.Data.Unsur)
					jawaban = res.Data.Lambang
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🧪 *Tebak Kimia*\n\n%s", soal)
				},
			})
		},
	})
}
