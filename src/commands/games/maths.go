package games

import (
	"encoding/json"
	"fmt"
	"strconv"

	"sptzx/src/core"
)

type mathsResp struct {
	Data struct {
		Str    string `json:"str"`
		Result int64  `json:"result"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "maths",
		Aliases:     []string{"math"},
		Description: "Game matematika",
		Usage:       "maths",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			return playGame(ptz, GameDef{
				Type:     "maths",
				Endpoint: "/api/games/maths",
				Parse: func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error) {
					var res mathsResp
					err = json.Unmarshal(raw, &res)
					soal = res.Data.Str
					jawaban = strconv.FormatInt(res.Data.Result, 10)
					return
				},
				Format: func(soal string) string {
					return fmt.Sprintf("🧮 *Matematika*\n\n%s", soal)
				},
			})
		},
	})
}
