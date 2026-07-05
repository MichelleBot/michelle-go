package games

import (
	"encoding/json"
	"fmt"

	"sptzx/src/core"
)

type tebakWarnaResp struct {
	Data struct {
		Image   string `json:"image"`
		Correct string `json:"correct"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakwarna",
		Aliases:     []string{"tw"},
		Description: "Game tebak warna (Ishihara)",
		Usage:       "tebakwarna",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakwarna",
				Endpoint: "/api/games/tebakwarna",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res tebakWarnaResp
					err = json.Unmarshal(raw, &res)
					soal = "Tebak angka yang tertera di gambar tersebut."
					jawaban = res.Data.Correct
					imageURL = res.Data.Image
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🎨 *Tebak Warna (Ishihara)*\n\n%s", soal)
				},
			})
		},
	})
}
