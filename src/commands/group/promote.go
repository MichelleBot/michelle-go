package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "promote",
		Description: "Jadikan member sebagai admin",
		Usage:       "promote @mention",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			targets := serialize.ExtractTargets(ptz.Message, ptz.Args)
			if len(targets) == 0 {
				return ptz.ReplyText("Tag member yang ingin di-promote.")
			}
			if _, err := serialize.PromoteParticipant(ptz.Bot.Client, ptz.Chat, targets); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Berhasil promote ke admin.")
		},
	})
}
