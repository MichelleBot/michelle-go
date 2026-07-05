package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type lengkapiKalimatResp struct {
	Data struct {
		Pertanyaan string `json:"pertanyaan"`
		Jawaban    string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "lengkapikalimat",
		Aliases:     []string{"lk"},
		Description: "Game lengkapi kalimat",
		Usage:       "lengkapikalimat",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "lengkapikalimat",
				Endpoint: "/api/games/lengkapikalimat",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res lengkapiKalimatResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Pertanyaan
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("📝 *Lengkapi Kalimat*\n\n%s", soal)
				},
			})
		},
	})
}
