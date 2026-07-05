package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakTebakanResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebaktebakan",
		Aliases:     []string{"ttb"},
		Description: "Game tebak tebakan",
		Usage:       "tebaktebakan",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebaktebakan",
				Endpoint: "/api/games/tebaktebakan",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakTebakanResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🤔 *Tebak Tebakan*\n\n%s", soal)
				},
			})
		},
	})
}
