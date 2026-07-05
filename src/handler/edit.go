package handler

import (
	"fmt"

	"sptzx/src/core"
)

func (h *EventHandler) handleEditEvent(msg *core.NormalizedMessage) {
	if msg == nil || msg.Edit == nil {
		return
	}

	oldContent := ""
	if h.bot != nil && h.bot.Messages != nil && msg.Edit.OriginalID != "" {
		if stored, ok := h.bot.Messages.Get(msg.Chat.String(), msg.Edit.OriginalID); ok {
			oldContent = stored.Content
			_, _ = h.bot.Messages.Update(msg.Chat.String(), msg.Edit.OriginalID, msg.Edit.EditedText, msg.Kind)
		}
	}

	h.bot.Log.Infof("message edited by %s in %s: old=%q new=%q", msg.Sender.User, msg.Chat.String(), oldContent, msg.Edit.EditedText)

	if msg.IsFromMe || oldContent == "" || oldContent == msg.Edit.EditedText {
		return
	}

	ptz := core.NewPtzFromNormalizedMessage(h.bot, msg)
	if ptz == nil {
		return
	}

	_ = ptz.ReplyText(fmt.Sprintf("*Pesan diedit*\n\n- Sebelum: %s\n- Sesudah: %s", oldContent, msg.Edit.EditedText))
}
