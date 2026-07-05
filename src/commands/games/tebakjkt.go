package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakJktResp struct {
	Data struct {
		Gambar  string `json:"gambar"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakjkt",
		Aliases:     []string{"tjk"},
		Description: "Game tebak member JKT48",
		Usage:       "tebakjkt",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakjkt",
				Endpoint: "/api/games/tebakjkt",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakJktResp
					err = json.Unmarshal(raw, &res)
					soal = "Tebak nama member JKT48 dari gambar tersebut."
					jawaban = res.Data.Jawaban
					imageURL = res.Data.Gambar
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🎶 *Tebak Member JKT48*\n\n%s", soal)
				},
			})
		},
	})
}
