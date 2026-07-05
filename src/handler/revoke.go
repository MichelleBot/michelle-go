package handler

import (
	"fmt"

	"sptzx/src/core"
	"sptzx/src/serialize"
)

func (h *EventHandler) handleRevokeEvent(msg *core.NormalizedMessage) {
	if msg == nil || msg.Revoke == nil {
		return
	}

	h.bot.Log.Warnf("protocol/revoke event type=%s target=%s in chat %s", msg.Revoke.ProtocolType, msg.Revoke.TargetID, msg.Chat.String())
	if msg.Revoke.ProtocolType != "REVOKE" {
		return
	}

	if msg.IsFromMe || h.bot == nil || h.bot.Messages == nil || msg.Revoke.TargetID == "" {
		return
	}

	if msg.IsGroup {
		settings := h.bot.Settings.GetGroupSettings(msg.Chat)
		if !settings.AntiDeleteEnabled {
			return
		}
	}

	stored, ok := h.bot.Messages.Get(msg.Chat.String(), msg.Revoke.TargetID)
	if !ok {
		return
	}

	name := stored.PushName
	if name == "" {
		name = stored.Sender
	}

	ptz := core.NewPtzFromNormalizedMessage(h.bot, msg)
	if ptz == nil {
		return
	}

	content := stored.Content
	if content == "" {
		content = fmt.Sprintf("[%s]", stored.Kind)
	}

	header := fmt.Sprintf("*Anti delete terdeteksi*\n\n- Pengirim: %s\n- Jenis: %s", name, stored.Kind)

	switch stored.Kind {
	case core.MessageImage:
		if len(stored.MediaData) > 0 {
			caption := header
			if stored.Caption != "" {
				caption += fmt.Sprintf("\n- Caption: %s", stored.Caption)
			}
			_ = ptz.ReplyImage(stored.MediaData, stored.MIME, caption)
			return
		}
	case core.MessageVideo:
		if len(stored.MediaData) > 0 {
			caption := header
			if stored.Caption != "" {
				caption += fmt.Sprintf("\n- Caption: %s", stored.Caption)
			}
			_ = ptz.ReplyVideo(stored.MediaData, stored.MIME, caption)
			return
		}
	case core.MessageAudio:
		if len(stored.MediaData) > 0 {
			_ = ptz.ReplyText(header)
			_ = ptz.ReplyAudio(stored.MediaData, stored.MIME)
			return
		}
	case core.MessageVoice:
		if len(stored.MediaData) > 0 {
			_ = ptz.ReplyText(header)
			_ = serialize.SendVoiceNoteReply(h.bot.Client, ptz.Chat, stored.MediaData, stored.MIME, ptz.Message, ptz.Info)
			return
		}
	case core.MessageDocument:
		if len(stored.MediaData) > 0 {
			caption := header
			if stored.Caption != "" {
				caption += fmt.Sprintf("\n- Caption: %s", stored.Caption)
			}
			_ = ptz.ReplyDocument(stored.MediaData, stored.MIME, stored.Filename, caption)
			return
		}
	case core.MessageSticker:
		if len(stored.MediaData) > 0 {
			_ = ptz.ReplyText(header)
			_ = ptz.ReplySticker(stored.MediaData, stored.MIME, stored.IsAnimated)
			return
		}
	case core.MessageLocation, core.MessageLiveLocation:
		_ = ptz.ReplyText(header)
		_ = serialize.SendLocation(h.bot.Client, ptz.Chat, stored.Latitude, stored.Longitude, stored.PlaceName, stored.Address)
		return
	case core.MessageContact, core.MessageContacts:
		if stored.ContactPhone != "" {
			_ = ptz.ReplyText(header)
			_ = serialize.SendContact(h.bot.Client, ptz.Chat, stored.ContactPhone, stored.ContactName)
			return
		}
	}

	_ = ptz.ReplyText(fmt.Sprintf("%s\n- Isi: %s", header, content))
}
