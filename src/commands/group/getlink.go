package group

import (
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "getlink",
		Description: "Ambil invite link group",
		Usage:       "getlink",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			link, err := serialize.GetInviteLink(ptz.Bot.Client, ptz.Chat, false)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("🔗 Invite link:\n" + link)
		},
	})
}
