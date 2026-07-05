package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakKartunResp struct {
	Data struct {
		Name string `json:"name"`
		Img  string `json:"img"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakkartun",
		Aliases:     []string{"tk"},
		Description: "Game tebak kartun",
		Usage:       "tebakkartun",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakkartun",
				Endpoint: "/api/games/tebakkartun",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakKartunResp
					err = json.Unmarshal(raw, &res)
					soal = "Tebak nama kartun dari gambar tersebut."
					jawaban = res.Data.Name
					imageURL = res.Data.Img
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("📺 *Tebak Kartun*\n\n%s", soal)
				},
			})
		},
	})
}
