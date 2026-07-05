package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakLirikResp struct {
	Data struct {
		Soal    string `json:"soal"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebaklirik",
		Aliases:     []string{"tl"},
		Description: "Game tebak lirik lagu",
		Usage:       "tebaklirik",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebaklirik",
				Endpoint: "/api/games/tebaklirik",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakLirikResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Soal
					jawaban = res.Data.Jawaban
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🎵 *Tebak Lirik*\n\n%s", soal)
				},
			})
		},
	})
}
