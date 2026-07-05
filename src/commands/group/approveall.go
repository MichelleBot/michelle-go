package group

import (
	"fmt"
	"go.mau.fi/whatsmeow/types"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "approveall",
		Description: "Setujui semua permintaan join",
		Usage:       "approveall",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
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
			var jids []types.JID
			for _, r := range reqs {
				jids = append(jids, r.JID)
			}
			if _, err := serialize.ApproveJoinRequests(ptz.Bot.Client, ptz.Chat, jids); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText(fmt.Sprintf("✅ Berhasil approve %d permintaan.", len(reqs)))
		},
	})
}
