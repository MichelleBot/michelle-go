package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakGambarResp struct {
	Data struct {
		Img       string `json:"img"`
		Deskripsi string `json:"deskripsi"`
		Jawaban   string `json:"jawaban"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakgambar",
		Aliases:     []string{"tg"},
		Description: "Game tebak gambar",
		Usage:       "tebakgambar",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakgambar",
				Endpoint: "/api/games/tebakgambar",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakGambarResp
					err = json.Unmarshal(raw, &res)
					soal = fmt.Sprintf("Clue: %s", res.Data.Deskripsi)
					jawaban = res.Data.Jawaban
					imageURL = res.Data.Img
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🖼️ *Tebak Gambar*\n\n%s", soal)
				},
			})
		},
	})
}
