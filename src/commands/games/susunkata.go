package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type susunKataResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Tipe    string `json:"tipe"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "susunkata",
		Aliases:     []string{"sk"},
		Description: "Game susun kata",
		Usage:       "susunkata",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "susunkata",
				Endpoint: "/api/games/susunkata",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res susunKataResp
					err = json.Unmarshal(raw, &res)
					soal = fmt.Sprintf("%s\nPetunjuk: %s", res.Data.Soal, res.Data.Tipe)
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🧩 *Susun Kata*\n\n%s", soal)
				},
			})
		},
	})
}
