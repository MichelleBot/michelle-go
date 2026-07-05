package group

import (
	"go.mau.fi/whatsmeow/types"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "tagadmin",
		Aliases:     []string{"admins"},
		Description: "Tag semua admin group",
		Usage:       "tagadmin [pesan]",
		Category:    "group",
		GroupOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if err := ptz.LoadGroupInfo(); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			msg := ptz.RawArgs
			if msg == "" {
				msg = "📢 Admin dibutuhkan!"
			}
			var jids []types.JID
			var sb strings.Builder
			sb.WriteString(msg + "\n\n")
			for _, p := range ptz.GroupInfo.Participants {
				if p.IsAdmin || p.IsSuperAdmin {
					jids = append(jids, p.JID)
					sb.WriteString("@" + p.JID.User + " ")
				}
			}
			if len(jids) == 0 {
				return ptz.ReplyText("Tidak ada admin ditemukan.")
			}
			return serialize.SendTextReplyMention(ptz.Bot.Client, ptz.Chat, sb.String(), jids, ptz.Message, ptz.Info)
		},
	})
}
