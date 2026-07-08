package admin

import (
	"fmt"
	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"setname"},
		UsageHint: "text",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		BotAdmin:  true,
		Handler:   runGroupSet,
	})
	core.Use(&core.Command{
		Usage:     []string{"setdesc"},
		UsageHint: "text",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		BotAdmin:  true,
		Handler:   runGroupSet,
	})
}

func runGroupSet(ptz *core.Ptz) error {
	value := ptz.RawArgs
	if ptz.Message != nil && ptz.Message.ExtendedTextMessage != nil && ptz.Message.ExtendedTextMessage.ContextInfo != nil {
		quoted := ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage
		if quoted != nil && quoted.Conversation != nil {
			value = *quoted.Conversation
		}
	}

	if value == "" {
		return ptz.ReplyText(fmt.Sprintf("🚩 Masukkan %s yang diinginkan.", ptz.Command))
	}

	switch ptz.Command {
	case "setname":
		if len(value) > 25 {
			return ptz.ReplyText("🚩 Teks terlalu panjang, maksimal 25 karakter.")
		}
		err := serialize.SetGroupName(ptz.Client, ptz.Chat, value)
		if err != nil {
			return ptz.ReplyText(fmt.Sprintf("❌ Gagal mengubah nama grup: %v", err))
		}
		return ptz.ReplyText("✅ Nama grup berhasil diubah.")
	case "setdesc":
		err := serialize.SetGroupDescription(ptz.Client, ptz.Chat, value)
		if err != nil {
			return ptz.ReplyText(fmt.Sprintf("❌ Gagal mengubah deskripsi grup: %v", err))
		}
		return ptz.ReplyText("✅ Deskripsi grup berhasil diubah.")
	}

	return nil
}
