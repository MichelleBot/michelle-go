package group

import (
	"go.mau.fi/whatsmeow/types"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "tagall",
		Aliases:     []string{"everyone", "all"},
		Description: "Tag semua member group",
		Usage:       "tagall [pesan]",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if err := ptz.LoadGroupInfo(); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			msg := ptz.RawArgs
			if msg == "" {
				msg = "📢 Perhatian semua member!"
			}
			var jids []types.JID
			var sb strings.Builder
			sb.WriteString(msg + "\n\n")
			for _, p := range ptz.GroupInfo.Participants {
				jids = append(jids, p.JID)
				sb.WriteString("@" + p.JID.User + " ")
			}
			return serialize.SendTextReplyMention(ptz.Bot.Client, ptz.Chat, sb.String(), jids, ptz.Message, ptz.Info)
		},
	})
}
