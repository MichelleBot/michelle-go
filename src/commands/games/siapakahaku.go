package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type siapakahAkuResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "siapakahaku",
		Aliases:     []string{"sa"},
		Description: "Game siapakah aku",
		Usage:       "siapakahaku",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "siapakahaku",
				Endpoint: "/api/games/siapakahaku",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res siapakahAkuResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("👤 *Siapakah Aku*\n\n%s", soal)
				},
			})
		},
	})
}
