package games

import (
	"encoding/json"
	"sptzx/src/core"
)

type flagResp struct {
	Data struct {
		Name string `json:"name"`
		Img  string `json:"img"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "tebakbendera",
		Aliases:     []string{"tb"},
		Description: "Game tebak bendera negara",
		Usage:       "tebakbendera",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "tebakbendera",
				Endpoint: "/api/games/tebakbendera",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res flagResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Name
					jawaban = res.Data.Name
					imageURL = res.Data.Img
					return
				},
				Format: func(_ string) string {
					return "🚩 *Tebak Bendera*\n\nNegara apakah ini?"
				},
			})
		},
	})
}
