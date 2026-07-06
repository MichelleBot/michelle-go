package admin

import (
	"fmt"
	"strings"

	"michelle/system/core"

	"go.mau.fi/whatsmeow/types"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"tagall"},
		UsageHint: "text (optional)",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		Handler:   runTagall,
	})
}

func runTagall(ptz *core.Ptz) error {
	if err := ptz.LoadGroupInfo(); err != nil {
		return ptz.ReplyText("❌ Gagal mengambil informasi grup.")
	}

	var sb strings.Builder
	message := "Halo semuanya, admin menyebut kalian di grup " + ptz.GroupInfo.GroupName.Name
	if ptz.RawArgs != "" {
		message = ptz.RawArgs
	}

	sb.WriteString("乂  *T A G A L L*\n\n")
	sb.WriteString(fmt.Sprintf("*“%s”*\n\n", message))
	
	// Add readmore
	sb.WriteString(strings.Repeat("\u200E", 4001) + "\n")

	mentionedJIDs := make([]types.JID, 0, len(ptz.GroupInfo.Participants))
	for _, p := range ptz.GroupInfo.Participants {
		mentionedJIDs = append(mentionedJIDs, p.JID)
		sb.WriteString(fmt.Sprintf("◦  @%s\n", p.JID.User))
	}

	return ptz.ReplyTextMention(sb.String(), mentionedJIDs)
}
