package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type karakterFreefireResp struct {
	Data struct {
		Name   string `json:"name"`
		Gambar string `json:"gambar"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "karakterfreefire",
		Aliases:     []string{"kff"},
		Description: "Game tebak karakter FreeFire",
		Usage:       "karakterfreefire",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "karakterfreefire",
				Endpoint: "/api/games/karakter-freefire",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res karakterFreefireResp
					err = json.Unmarshal(raw, &res)
					soal = "Tebak nama karakter FreeFire dari gambar tersebut."
					jawaban = res.Data.Name
					imageURL = res.Data.Gambar
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🔫 *Tebak Karakter FF*\n\n%s", soal)
				},
			})
		},
	})
}
