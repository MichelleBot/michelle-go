package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakLogoResp struct {
	Data struct {
		Data struct {
			Image   string `json:"image"`
			Jawaban string `json:"jawaban"`
		} `json:"data"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebaklogo",
		Aliases:     []string{"tlogo"},
		Description: "Game tebak logo",
		Usage:       "tebaklogo",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebaklogo",
				Endpoint: "/api/games/tebaklogo",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakLogoResp
					err = json.Unmarshal(raw, &res)
					soal = "Tebak logo apakah ini?"
					jawaban = res.Data.Data.Jawaban
					imageURL = res.Data.Data.Image
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🛡️ *Tebak Logo*\n\n%s", soal)
				},
			})
		},
	})
}
