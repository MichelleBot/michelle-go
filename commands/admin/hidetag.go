package admin

import (
	"michelle/system/core"
	"michelle/system/utils"

	"go.mau.fi/whatsmeow/types"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"hidetag"},
		UsageHint: "text",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		Handler:   runHidetag,
	})
}

func runHidetag(ptz *core.Ptz) error {
	if err := ptz.LoadGroupInfo(); err != nil {
		return ptz.ReplyText("❌ Gagal mengambil informasi grup.")
	}
	if ptz.RawArgs == "" {
		return ptz.ReplyText(utils.Texted("bold", "🚩 Masukkan teks untuk hidetag."))
	}

	mentionedJIDs := make([]types.JID, 0, len(ptz.GroupInfo.Participants))
	for _, p := range ptz.GroupInfo.Participants {
		mentionedJIDs = append(mentionedJIDs, p.JID)
	}

	return ptz.SendTextMention(ptz.RawArgs, mentionedJIDs)
}
