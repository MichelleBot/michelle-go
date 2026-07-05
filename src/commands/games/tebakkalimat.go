package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakKalimatResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakkalimat",
		Aliases:     []string{"tkal"},
		Description: "Game tebak kalimat",
		Usage:       "tebakkalimat",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakkalimat",
				Endpoint: "/api/games/tebakkalimat",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakKalimatResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("💬 *Tebak Kalimat*\n\n%s", soal)
				},
			})
		},
	})
}
