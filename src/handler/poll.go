package handler

import (
	"fmt"

	"sptzx/src/core"
)

func (h *EventHandler) handlePollEvent(msg *core.NormalizedMessage) {
	if msg == nil || msg.Poll == nil {
		return
	}

	if msg.Poll.IsCreation {
		if h.bot != nil && h.bot.Polls != nil {
			h.bot.Polls.SaveCreation(msg.Chat.String(), msg.Poll, msg.Info.Timestamp)
		}
		h.bot.Log.Infof("poll created by %s in %s: %s (%d options)", msg.Sender.User, msg.Chat.String(), msg.Poll.Name, msg.Poll.OptionCount)
		return
	}

	if msg.Poll.IsUpdate {
		var info string
		if h.bot != nil && h.bot.Polls != nil {
			state := h.bot.Polls.RegisterUpdate(msg.Chat.String(), msg.Poll, msg.Info.Timestamp)
			if state != nil {
				info = fmt.Sprintf(" pada poll %s total update %d", state.Name, state.UpdateCount)
			}
		}
		h.bot.Log.Infof("poll updated by %s in %s%s", msg.Sender.User, msg.Chat.String(), info)
	}
}
