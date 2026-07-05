package handler

import (
	"sptzx/src/core"
	"time"

	"go.mau.fi/whatsmeow/types/events"
)

type EventHandler struct {
	bot      *core.Bot
	registry *core.Registry
	group    *GroupHandler
}

func NewEventHandler(bot *core.Bot, registry *core.Registry) *EventHandler {
	return &EventHandler{
		bot:      bot,
		registry: registry,
		group:    NewGroupHandler(bot),
	}
}

func (h *EventHandler) Handle(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.Message:
		h.runAsync("message", func() {
			h.dispatchNormalizedEvent(core.NormalizeMessageEvent(evt))
		})
	case *events.Receipt:
		h.runAsync("receipt", func() {
			h.dispatchNormalizedEvent(core.NormalizeReceiptEvent(evt))
		})
	case *events.ChatPresence:
		h.runAsync("chat-presence", func() {
			h.dispatchNormalizedEvent(core.NormalizeChatPresenceEvent(evt))
		})
	case *events.Presence:
		h.runAsync("presence", func() {
			h.dispatchNormalizedEvent(core.NormalizePresenceEvent(evt))
		})

	case *events.GroupInfo:
		h.runAsync("group-info", func() {
			h.group.OnGroupInfo(evt)
		})
	case *events.Picture:
		h.runAsync("picture", func() {
			h.group.OnPicture(evt)
		})
	case *events.JoinedGroup:
		h.runAsync("joined-group", func() {
			h.group.OnJoinedGroup(evt)
		})

	case *events.Connected:
		h.bot.Log.Infof("✅ Connected to WhatsApp")

	case *events.Disconnected:
		h.bot.Log.Warnf("⚠️ Disconnected from WhatsApp")

	case *events.LoggedOut:
		h.bot.Log.Errorf("❌ Logged out: %s", evt.Reason.String())

	case *events.PairSuccess:
		h.bot.Log.Infof("✅ Paired: %s", evt.ID.String())

	case *events.KeepAliveTimeout:
		h.bot.Log.Warnf("⚠️ Keepalive timeout — errors: %d, last ok: %s ago",
			evt.ErrorCount, evt.LastSuccess.String())

	case *events.KeepAliveRestored:
		h.bot.Log.Infof("✅ Keepalive restored")

	case *events.OfflineSyncCompleted:
		h.bot.Log.Infof("📥 Offline sync done — %d messages", evt.Count)

	case *events.TemporaryBan:
		h.bot.Log.Errorf("🚫 Temporary ban: %s (expires in %s)", evt.Code, evt.Expire)

	case *events.ClientOutdated:
		h.bot.Log.Errorf("❌ Client outdated — update whatsmeow")

	case *events.StreamReplaced:
		h.bot.Log.Warnf("⚠️ Stream replaced — another session opened")

	case *events.UndecryptableMessage:
		h.bot.Log.Warnf("⚠️ Undecryptable message from %s", evt.Info.Sender.String())

	case *events.IdentityChange:
		h.bot.Log.Warnf("🔑 Identity changed: %s (implicit: %v)", evt.JID.String(), evt.Implicit)

	case *events.PrivacySettings:
		h.bot.Log.Debugf("🔒 Privacy settings updated")

	case *events.PushName:
		h.bot.Log.Debugf("📛 Push name update: %s → %s", evt.JID.String(), evt.NewPushName)

	case *events.CallOffer:
		h.runAsync("call-offer", func() {
			h.dispatchNormalizedEvent(core.NormalizeCallEvent("call-offer", evt))
		})

	case *events.CallAccept:
		h.runAsync("call-accept", func() {
			h.dispatchNormalizedEvent(core.NormalizeCallEvent("call-accept", evt))
		})

	case *events.CallReject:
		h.runAsync("call-reject", func() {
			h.dispatchNormalizedEvent(core.NormalizeCallEvent("call-reject", evt))
		})

	case *events.CallOfferNotice:
		h.runAsync("call-offer-notice", func() {
			h.dispatchNormalizedEvent(core.NormalizeCallEvent("call-offer-notice", evt))
		})

	case *events.CallTerminate:
		h.runAsync("call-terminate", func() {
			h.dispatchNormalizedEvent(core.NormalizeCallEvent("call-terminate", evt))
		})
	}
}

func (h *EventHandler) runAsync(name string, fn func()) {
	go func() {
		start := time.Now()
		defer h.recoverPanic(name)
		fn()
		h.bot.Log.Debugf("event %s completed in %s", name, time.Since(start))
	}()
}

func (h *EventHandler) recoverPanic(scope string) {
	if r := recover(); r != nil {
		h.bot.Log.Errorf("Recovered from panic in %s handler: %v", scope, r)
	}
}
