package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "kick",
		Description: "Kick member dari group",
		Usage:       "kick @mention",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if err := ptz.LoadGroupInfo(); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			targets := serialize.ExtractTargets(ptz.Message, ptz.Args)
			if len(targets) == 0 {
				return ptz.ReplyText("Tag member yang ingin di-kick.")
			}
			if _, err := serialize.KickParticipant(ptz.Bot.Client, ptz.Chat, targets); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Member berhasil di-kick.")
		},
	})
}
