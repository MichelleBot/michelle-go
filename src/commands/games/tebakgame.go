package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakGameResp struct {
	Data struct {
		Img     string `json:"img"`
		Jawaban string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakgame",
		Aliases:     []string{"tgm"},
		Description: "Game tebak judul game",
		Usage:       "tebakgame",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakgame",
				Endpoint: "/api/games/tebakgame",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakGameResp
					err = json.Unmarshal(raw, &res)
					soal = "Tebak judul game dari gambar tersebut."
					jawaban = res.Data.Jawaban
					imageURL = res.Data.Img
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🎮 *Tebak Game*\n\n%s", soal)
				},
			})
		},
	})
}
