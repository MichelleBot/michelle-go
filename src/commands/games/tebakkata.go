package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type kataResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakkata",
		Aliases:     []string{"tk"},
		Description: "Game tebak kata",
		Usage:       "tebakkata",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakkata",
				Endpoint: "/api/games/tebakkata",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res kataResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🔤 *Tebak Kata*\n\n%s", soal)
				},
			})
		},
	})
}
