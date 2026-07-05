package handler

import (
	"fmt"

	"sptzx/src/core"
)

func (h *EventHandler) handleReactionEvent(msg *core.NormalizedMessage) {
	if msg == nil || msg.Reaction == nil {
		return
	}

	action := "added"
	if msg.Reaction.IsRemoval {
		action = "removed"
	}

	h.bot.Log.Infof("reaction %s by %s on %s: %s", action, msg.Sender.User, msg.Reaction.TargetID, msg.Reaction.Text)

	if msg.IsFromMe {
		return
	}

	if msg.Reaction.IsRemoval {
		return
	}

	if msg.Reaction.Text == "❓" {
		ptz := core.NewPtzFromNormalizedMessage(h.bot, msg)
		if ptz != nil {
			_ = ptz.ReplyText(fmt.Sprintf("Reaction %s diterima pada pesan %s", msg.Reaction.Text, msg.Reaction.TargetID))
		}
	}
}
