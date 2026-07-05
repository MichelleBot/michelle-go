package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tekaTekiResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tekateki",
		Aliases:     []string{"tt"},
		Description: "Game teka-teki",
		Usage:       "tekateki",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tekateki",
				Endpoint: "/api/games/tekateki",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tekaTekiResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🧩 *Teka-Teki*\n\n%s", soal)
				},
			})
		},
	})
}
