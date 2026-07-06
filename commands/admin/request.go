package admin

import (
	"context"
	"fmt"
	"strings"

	"michelle/system/core"
	"michelle/system/serialize"

	"go.mau.fi/whatsmeow/types"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"requestlist", "approveall", "rejectall"},
		UsageHint: "",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		BotAdmin:  true,
		Handler:   runRequest,
	})
}

func resolveLIDToPhone(ptz *core.Ptz, jid types.JID) types.JID {
	if jid.Server == types.HiddenUserServer {
		pn, err := ptz.Bot.Client.Store.LIDs.GetPNForLID(context.Background(), jid)
		if err == nil && !pn.IsEmpty() {
			return types.NewJID(pn.User, types.DefaultUserServer)
		}
	}
	return jid
}

func runRequest(ptz *core.Ptz) error {
	requests, err := serialize.GetGroupRequestParticipants(ptz.Bot.Client, ptz.Chat)
	if err != nil {
		return ptz.ReplyText(fmt.Sprintf("❌ Gagal mendapatkan daftar permintaan: %v", err))
	}

	if len(requests) == 0 {
		return ptz.ReplyText("❌ Tidak ada permintaan bergabung yang tertunda.")
	}

	switch ptz.Command {
	case "requestlist", "reqlist":
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Grup ini memiliki %d permintaan bergabung.\n\n", len(requests)))
		
		var mentionJIDs []types.JID
		for _, req := range requests {
			resolvedJID := resolveLIDToPhone(ptz, req.JID)
			mentionJIDs = append(mentionJIDs, resolvedJID)
			sb.WriteString(fmt.Sprintf("◦ @%s\n    Pada : %s\n", resolvedJID.User, req.RequestedAt.Format("Mon, 02/01/06 15:04")))
		}
		sb.WriteString("\n> Kirim *approveall* atau *rejectall* untuk menyetujui atau menolak semua permintaan.")
		return ptz.ReplyTextMention(sb.String(), mentionJIDs)

	case "approveall":
		var jids []types.JID
		for _, req := range requests {
			jids = append(jids, req.JID)
		}
		_, err := serialize.ApproveJoinRequests(ptz.Bot.Client, ptz.Chat, jids)
		if err != nil {
			return ptz.ReplyText(fmt.Sprintf("❌ Gagal menyetujui semua: %v", err))
		}
		return ptz.ReplyText(fmt.Sprintf("✅ %d permintaan bergabung telah disetujui.", len(jids)))

	case "rejectall":
		var jids []types.JID
		for _, req := range requests {
			jids = append(jids, req.JID)
		}
		_, err := serialize.RejectJoinRequests(ptz.Bot.Client, ptz.Chat, jids)
		if err != nil {
			return ptz.ReplyText(fmt.Sprintf("❌ Gagal menolak semua: %v", err))
		}
		return ptz.ReplyText(fmt.Sprintf("✅ %d permintaan bergabung telah ditolak.", len(jids)))
	}

	return nil
}
