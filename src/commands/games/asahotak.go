package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type asahOtakResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "asahotak",
		Aliases:     []string{"ao"},
		Description: "Game asah otak",
		Usage:       "asahotak",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "asahotak",
				Endpoint: "/api/games/asahotak",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res asahOtakResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🧠 *Asah Otak*\n\n%s", soal)
				},
			})
		},
	})
}
