package admin

import (
	"context"
	"fmt"
	"strings"

	"michelle/system/core"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func init() {
	// Kick
	core.Use(&core.Command{
		Usage:     []string{"kick"},
		UsageHint: "mention atau reply",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		BotAdmin:  true,
		Handler:   runAction,
	})

	// Promote
	core.Use(&core.Command{
		Usage:     []string{"promote"},
		UsageHint: "mention atau reply",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		BotAdmin:  true,
		Handler:   runAction,
	})

	// Demote
	core.Use(&core.Command{
		Usage:     []string{"demote"},
		UsageHint: "mention atau reply",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		BotAdmin:  true,
		Handler:   runAction,
	})
}

func resolveActionLIDToPhone(ptz *core.Ptz, jid types.JID) types.JID {
	if jid.Server == types.HiddenUserServer {
		pn, err := ptz.Client.Store.LIDs.GetPNForLID(context.Background(), jid)
		if err == nil && !pn.IsEmpty() {
			return types.NewJID(pn.User, types.DefaultUserServer)
		}
	}
	return jid
}

func runAction(ptz *core.Ptz) error {
	var target string
	if ptz.Message != nil && ptz.Message.ExtendedTextMessage != nil && ptz.Message.ExtendedTextMessage.ContextInfo != nil {
		ctx := ptz.Message.ExtendedTextMessage.ContextInfo
		if len(ctx.MentionedJID) > 0 {
			target = ctx.MentionedJID[0]
		} else if ctx.Participant != nil {
			target = *ctx.Participant
		}
	}

	if target == "" && len(ptz.Args) > 0 {
		// Basic parsing, assuming JID or phone number
		target = ptz.Args[0]
		if !strings.Contains(target, "@") {
			target = target + "@s.whatsapp.net"
		}
	}

	if target == "" {
		return ptz.ReplyText("🚩 Tag atau balas chat target.")
	}

	jid, err := types.ParseJID(target)
	if err != nil {
		return ptz.ReplyText("❌ Nomor atau JID tidak valid.")
	}

	// Resolve LID to Phone JID for display
	resolvedJID := resolveActionLIDToPhone(ptz, jid)

	var action whatsmeow.ParticipantChange
	switch ptz.Command {
	case "kick":
		action = whatsmeow.ParticipantChangeRemove
	case "promote":
		action = whatsmeow.ParticipantChangePromote
	case "demote":
		action = whatsmeow.ParticipantChangeDemote
	default:
		return ptz.ReplyText("❌ Perintah tidak dikenal.")
	}

	_, err = ptz.Client.UpdateGroupParticipants(context.Background(), ptz.Chat, []types.JID{jid}, action)
	if err != nil {
		return ptz.ReplyText(fmt.Sprintf("❌ Gagal melakukan %s: %v", ptz.Command, err))
	}

	// Ensure JID is fully qualified as a phone JID for proper linking
	mentionJID := types.NewJID(resolvedJID.User, types.DefaultUserServer)

	return ptz.ReplyTextMention(fmt.Sprintf("✅ Berhasil melakukan %s terhadap @%s", ptz.Command, mentionJID.User), []types.JID{mentionJID})
}
