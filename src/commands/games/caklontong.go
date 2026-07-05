package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type cakLontongResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "caklontong",
		Aliases:     []string{"cl"},
		Description: "Game tebak-tebakan ala Cak Lontong",
		Usage:       "caklontong",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "caklontong",
				Endpoint: "/api/games/caklontong",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res cakLontongResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🤣 *Tebak Cak Lontong*\n\n%s", soal)
				},
			})
		},
	})
}
