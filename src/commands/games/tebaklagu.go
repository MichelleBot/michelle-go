package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type laguResp struct {
	Data struct {
		Lagu  string `json:"lagu"`
		Judul string `json:"judul"`
		Artis string `json:"artis"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebaklagu",
		Aliases:     []string{"tl"},
		Description: "Game tebak judul lagu",
		Usage:       "tebaklagu",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebaklagu",
				Endpoint: "/api/games/tebaklagu",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res laguResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Artis
					jawaban = res.Data.Judul
					audioURL = res.Data.Lagu
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🎵 *Tebak Lagu*\n\nArtis: *%s*\n\nSebutkan judul lagunya!", soal)
				},
			})
		},
	})
}
