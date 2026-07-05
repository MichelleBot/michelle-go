package group

import (
	"fmt"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "joinrequest",
		Aliases:     []string{"requests", "listreq"},
		Description: "Lihat daftar permintaan join",
		Usage:       "joinrequest",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			reqs, err := serialize.GetGroupRequestParticipants(ptz.Bot.Client, ptz.Chat)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			if len(reqs) == 0 {
				return ptz.ReplyText("Tidak ada permintaan join.")
			}
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("*Permintaan Join (%d)*\n\n", len(reqs)))
			for _, r := range reqs {
				sb.WriteString(fmt.Sprintf("+%s\n%s\n", r.JID.User, r.RequestedAt.Format("02 Jan 2006 15:04")))
			}
			sb.WriteString("\nGunakan .approveall atau .rejectall")
			return ptz.ReplyText(sb.String())
		},
	})
}
