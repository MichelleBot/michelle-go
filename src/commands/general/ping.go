package general

import (
	"fmt"
	"time"

	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "ping",
		Aliases:     []string{"p"},
		Description: "Cek response time bot",
		Usage:       "ping",
		Category:    "general",
		Handler: func(ptz *core.Ptz) error {
			start := time.Now()
			if err := ptz.React("🏓"); err != nil {
				ptz.Bot.Log.Debugf("react error: %v", err)
			}
			elapsed := time.Since(start)
			return ptz.ReplyText(fmt.Sprintf("🏓 *Pong!* `%dms`\n_sptzx v%s by %s_", elapsed.Milliseconds(), core.SptzxVersion, core.SptzxAuthor))
		},
	})
}
