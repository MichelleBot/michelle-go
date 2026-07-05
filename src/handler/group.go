package handler

import (
	"context"
	"strings"

	"sptzx/src/core"
	"sptzx/src/serialize"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type GroupHandler struct {
	bot *core.Bot
}

func NewGroupHandler(bot *core.Bot) *GroupHandler {
	return &GroupHandler{bot: bot}
}

func (h *GroupHandler) OnGroupInfo(evt *events.GroupInfo) {
	if len(evt.Join) > 0 {
		h.sendGroupMessage(evt.JID, evt.Join, true)
	}
	if len(evt.Leave) > 0 {
		h.sendGroupMessage(evt.JID, evt.Leave, false)
	}
}

func (h *GroupHandler) OnPicture(evt *events.Picture) {
	h.bot.Log.Debugf("Picture changed for %s by %s", evt.JID, evt.Author)
}

func (h *GroupHandler) OnJoinedGroup(evt *events.JoinedGroup) {
	h.bot.Log.Infof("Joined group: %s", evt.JID)
}

func (h *GroupHandler) resolvePhoneJID(jid types.JID) types.JID {
	if jid.Server == types.HiddenUserServer {
		pn, err := h.bot.Client.Store.LIDs.GetPNForLID(context.Background(), jid)
		if err == nil && !pn.IsEmpty() {
			return types.NewJID(pn.User, types.DefaultUserServer)
		}
		return types.NewJID(jid.User, types.DefaultUserServer)
	}
	return types.NewJID(jid.User, types.DefaultUserServer)
}

func (h *GroupHandler) sendGroupMessage(groupJID types.JID, participants []types.JID, isJoin bool) {
	settings := h.bot.Settings.GetGroupSettings(groupJID)

	if isJoin && !settings.WelcomeEnabled {
		return
	}
	if !isJoin && !settings.GoodbyeEnabled {
		return
	}

	template := settings.WelcomeMessage
	if !isJoin {
		template = settings.GoodbyeMessage
	}

	for _, p := range participants {
		phoneJID := h.resolvePhoneJID(p)
		msg := strings.ReplaceAll(template, "@user", "@"+phoneJID.User)
		if err := serialize.SendTextMention(h.bot.Client, groupJID, msg, []types.JID{phoneJID}); err != nil {
			h.bot.Log.Errorf("Failed to send group message: %v", err)
		}
	}
}
