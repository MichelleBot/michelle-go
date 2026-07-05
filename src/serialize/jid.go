package serialize

import (
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func PhoneToJID(phone string) types.JID {
	return types.NewJID(NormalizePhone(phone), types.DefaultUserServer)
}

func NormalizePhone(phone string) string {
	phone = strings.TrimPrefix(phone, "+")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	switch {
	case strings.HasPrefix(phone, "08"):
		phone = "62" + phone[1:]
	case strings.HasPrefix(phone, "8") && len(phone) >= 9:
		phone = "62" + phone
	}
	return phone
}

func IsValidPhone(phone string) bool {
	phone = NormalizePhone(phone)
	if len(phone) < 7 || len(phone) > 15 {
		return false
	}
	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func ExtractArgsJIDs(args []string) []types.JID {
	var jids []types.JID
	for _, arg := range args {
		arg = strings.TrimPrefix(arg, "@")
		if IsValidPhone(arg) {
			jids = append(jids, PhoneToJID(arg))
		}
	}
	return jids
}

func ExtractTargets(msg *waE2E.Message, args []string) []types.JID {
	seen := map[string]struct{}{}
	var jids []types.JID

	add := func(j types.JID) {
		if _, ok := seen[j.User]; !ok {
			seen[j.User] = struct{}{}
			jids = append(jids, j)
		}
	}

	if ci := GetContextInfo(msg); ci != nil {
		for _, raw := range ci.GetMentionedJID() {
			if j, err := types.ParseJID(raw); err == nil {
				add(j)
			}
		}
	}

	for _, j := range ExtractArgsJIDs(args) {
		add(j)
	}

	return jids
}

func GetQuotedMessage(msg *waE2E.Message) *waE2E.Message {
	if ci := GetContextInfo(msg); ci != nil {
		return ci.GetQuotedMessage()
	}
	return nil
}

func GetContextInfo(msg *waE2E.Message) *waE2E.ContextInfo {
	if msg == nil {
		return nil
	}
	switch {
	case msg.ExtendedTextMessage != nil:
		return msg.ExtendedTextMessage.ContextInfo
	case msg.ImageMessage != nil:
		return msg.ImageMessage.ContextInfo
	case msg.VideoMessage != nil:
		return msg.VideoMessage.ContextInfo
	case msg.AudioMessage != nil:
		return msg.AudioMessage.ContextInfo
	case msg.DocumentMessage != nil:
		return msg.DocumentMessage.ContextInfo
	case msg.StickerMessage != nil:
		return msg.StickerMessage.ContextInfo
	}
	return nil
}

func GetMessageType(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}
	switch {
	case msg.Conversation != nil || msg.ExtendedTextMessage != nil:
		return "text"
	case msg.ImageMessage != nil:
		return "image"
	case msg.VideoMessage != nil:
		return "video"
	case msg.AudioMessage != nil:
		return "audio"
	case msg.DocumentMessage != nil:
		return "document"
	case msg.StickerMessage != nil:
		return "sticker"
	case msg.ContactMessage != nil:
		return "contact"
	case msg.LocationMessage != nil:
		return "location"
	case msg.ReactionMessage != nil:
		return "reaction"
	}
	return "unknown"
}

func IsMediaType(msg *waE2E.Message) bool {
	switch GetMessageType(msg) {
	case "image", "video", "audio", "document", "sticker":
		return true
	}
	return false
}

type InputMedia struct {
	Message  *waE2E.Message
	MimeType string
	MsgType  string
	IsQuoted bool
}

func GetInputMedia(msg *waE2E.Message, allowedTypes ...string) *InputMedia {
	allowed := map[string]struct{}{}
	for _, t := range allowedTypes {
		allowed[t] = struct{}{}
	}

	check := func(m *waE2E.Message, quoted bool) *InputMedia {
		if m == nil {
			return nil
		}
		mt := GetMessageType(m)
		if len(allowed) > 0 {
			if _, ok := allowed[mt]; !ok {
				return nil
			}
		}
		switch mt {
		case "image", "video", "audio", "document", "sticker":
			return &InputMedia{Message: m, MimeType: GetMediaMIME(m), MsgType: mt, IsQuoted: quoted}
		}
		return nil
	}

	if m := check(msg, false); m != nil {
		return m
	}

	if ci := GetContextInfo(msg); ci != nil {
		if m := check(ci.GetQuotedMessage(), true); m != nil {
			return m
		}
	}

	return nil
}
