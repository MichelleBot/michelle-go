package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "demote",
		Description: "Cabut status admin member",
		Usage:       "demote @mention",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			targets := serialize.ExtractTargets(ptz.Message, ptz.Args)
			if len(targets) == 0 {
				return ptz.ReplyText("Tag admin yang ingin di-demote.")
			}
			if _, err := serialize.DemoteParticipant(ptz.Bot.Client, ptz.Chat, targets); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Berhasil demote ke member.")
		},
	})
}
