package commands

import (
	"fmt"
	"time"

	"michelle/system/core"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:       []string{"runtime"},
		Category:    "miscs",
		Handler: func(ptz *core.Ptz) error {
			uptime := time.Since(utils.StartTime)
			return ptz.ReplyText(fmt.Sprintf("*Aktif selama : [ %s ]*", utils.FmtUptime(uptime)))
		},
	})
}
