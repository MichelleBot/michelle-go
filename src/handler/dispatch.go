package handler

import (
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
	"sptzx/src/core"
)

func (h *EventHandler) dispatchNormalizedEvent(evt *core.NormalizedEvent) {
	if evt == nil {
		return
	}

	switch evt.Kind {
	case core.EventMessage:
		h.handleMessageEvent(evt.Message)
	case core.EventReaction:
		h.handleReactionEvent(evt.Message)
	case core.EventEdit:
		h.handleEditEvent(evt.Message)
	case core.EventRevoke:
		h.handleRevokeEvent(evt.Message)
	case core.EventPoll:
		h.handlePollEvent(evt.Message)
	case core.EventReceipt:
		h.handleReceiptEvent(evt)
	case core.EventPresence:
		h.handlePresenceEvent(evt)
	case core.EventCall:
		h.handleCallEvent(evt)
	}
}

func (h *EventHandler) handleReceiptEvent(evt *core.NormalizedEvent) {
	receipt, ok := evt.Raw.(*events.Receipt)
	if !ok {
		return
	}
	h.bot.Log.Debugf("📨 Receipt [%s] from %s", receipt.Type, receipt.MessageSource.Sender.String())
}

func (h *EventHandler) handlePresenceEvent(evt *core.NormalizedEvent) {
	switch raw := evt.Raw.(type) {
	case *events.ChatPresence:
		h.bot.Log.Debugf("✍️ Chat presence %s in %s", raw.State, raw.Chat.String())
	case *events.Presence:
		state := "online"
		if raw.Unavailable {
			state = "offline"
		}
		h.bot.Log.Debugf("👀 Presence %s is %s", raw.From.String(), state)
	}
}

func (h *EventHandler) handleCallEvent(evt *core.NormalizedEvent) {
	switch raw := evt.Raw.(type) {
	case *events.CallOffer:
		h.bot.Log.Infof("📞 Incoming call from %s", raw.From.String())
	case *events.CallAccept:
		h.bot.Log.Infof("✅ Call accepted by %s", raw.From.String())
	case *events.CallReject:
		h.bot.Log.Infof("❌ Call rejected by %s", raw.From.String())
	case *events.CallOfferNotice:
		h.bot.Log.Infof("📣 Call notice type=%s media=%s", raw.Type, raw.Media)
	case *events.CallTerminate:
		h.bot.Log.Infof("📵 Call ended from %s", raw.From.String())
	default:
		h.bot.Log.Debugf("Unhandled normalized call event %T", evt.Raw)
	}
}

func (h *EventHandler) logNormalizedMessage(msg *core.NormalizedMessage) {
	if msg == nil {
		return
	}

	location := "DM"
	sender := msg.Info.Sender.User
	if msg.IsGroup {
		location = fmt.Sprintf("Group[%s]", msg.Chat.User)
		if msg.Info.PushName != "" && msg.Info.PushName != "-" {
			sender = msg.Info.PushName
		}
	}

	h.bot.Log.Infof("[%s] %s → %s/%s: %s", location, sender, msg.EventKind, msg.Kind, msg.Content)
}
