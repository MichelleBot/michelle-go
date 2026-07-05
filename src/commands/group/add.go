package group

import (
	"go.mau.fi/whatsmeow/types"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "add",
		Description: "Tambah member ke group",
		Usage:       "add <nomor>",
		Category:    "group",
		GroupOnly:   true,
		AdminOnly:   true,
		BotAdmin:    true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Usage: .add <nomor>")
			}
			phone := strings.Join(ptz.Args, "")
			if !serialize.IsValidPhone(phone) {
				return ptz.ReplyText("❌ Nomor tidak valid.")
			}
			jid := serialize.PhoneToJID(phone)
			if _, err := serialize.AddParticipant(ptz.Bot.Client, ptz.Chat, []types.JID{jid}); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return ptz.ReplyText("✅ Berhasil menambahkan " + phone)
		},
	})
}
