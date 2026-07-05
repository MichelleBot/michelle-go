package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"sptzx/src/serialize"
)

type Ptz struct {
	Bot       *Bot
	Event     *events.Message
	Message   *waE2E.Message
	Info      types.MessageInfo
	Args      []string
	RawArgs   string
	Command   string
	IsGroup   bool
	IsFromMe  bool
	Sender    types.JID
	SenderAlt types.JID
	Chat      types.JID
	GroupInfo *types.GroupInfo
}

func NewPtz(bot *Bot, evt *events.Message) *Ptz {
	body := ExtractBody(evt.Message)
	parts, cmd, rawArgs, args := parseCommandParts(bot, body)

	_ = parts

	return &Ptz{
		Bot:       bot,
		Event:     evt,
		Message:   evt.Message,
		Info:      evt.Info,
		Args:      args,
		RawArgs:   rawArgs,
		Command:   cmd,
		IsGroup:   evt.Info.IsGroup,
		IsFromMe:  evt.Info.IsFromMe,
		Sender:    evt.Info.Sender,
		SenderAlt: evt.Info.SenderAlt,
		Chat:      evt.Info.Chat,
	}
}

func NewPtzFromNormalizedMessage(bot *Bot, msg *NormalizedMessage) *Ptz {
	if msg == nil {
		return nil
	}

	_, cmd, rawArgs, args := parseCommandParts(bot, msg.Body)

	return &Ptz{
		Bot:       bot,
		Event:     msg.Event,
		Message:   msg.Message,
		Info:      msg.Info,
		Args:      args,
		RawArgs:   rawArgs,
		Command:   cmd,
		IsGroup:   msg.IsGroup,
		IsFromMe:  msg.IsFromMe,
		Sender:    msg.Sender,
		SenderAlt: msg.SenderAlt,
		Chat:      msg.Chat,
	}
}

func parseCommandParts(bot *Bot, body string) ([]string, string, string, []string) {
	parts := strings.Fields(body)

	var cmd, rawArgs string
	args := []string{}

	if len(parts) > 0 {
		for _, prefix := range bot.Config.Prefixes {
			if strings.HasPrefix(parts[0], prefix) {
				cmd = strings.ToLower(strings.TrimPrefix(parts[0], prefix))
				if len(parts) > 1 {
					args = parts[1:]
					rawArgs = strings.TrimSpace(strings.TrimPrefix(body, parts[0]))
				}
				break
			}
		}
	}

	return parts, cmd, rawArgs, args
}

func ExtractBody(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}
	switch {
	case msg.Conversation != nil:
		return *msg.Conversation
	case msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil:
		return *msg.ExtendedTextMessage.Text
	case msg.ImageMessage != nil && msg.ImageMessage.Caption != nil:
		return *msg.ImageMessage.Caption
	case msg.VideoMessage != nil && msg.VideoMessage.Caption != nil:
		return *msg.VideoMessage.Caption
	case msg.DocumentMessage != nil && msg.DocumentMessage.Caption != nil:
		return *msg.DocumentMessage.Caption
	case msg.ButtonsResponseMessage != nil:
		return msg.GetButtonsResponseMessage().GetSelectedDisplayText()
	case msg.TemplateButtonReplyMessage != nil:
		return msg.GetTemplateButtonReplyMessage().GetSelectedDisplayText()
	case msg.ListResponseMessage != nil:
		return msg.GetListResponseMessage().GetTitle()
	}
	return ""
}

func matchParticipant(p types.GroupParticipant, sender, senderAlt types.JID) bool {
	if sender.Server == types.HiddenUserServer {
		if p.LID.User == sender.User {
			return true
		}
		if !senderAlt.IsEmpty() && p.PhoneNumber.User == senderAlt.User {
			return true
		}
	} else {
		if p.PhoneNumber.User == sender.User || p.JID.User == sender.User {
			return true
		}
		if !senderAlt.IsEmpty() && p.LID.User == senderAlt.User {
			return true
		}
	}
	return false
}

func (ptz *Ptz) IsOwner() bool {
	for _, owner := range ptz.Bot.Config.Owners {
		if owner == ptz.Sender.User {
			return true
		}
		if !ptz.SenderAlt.IsEmpty() && owner == ptz.SenderAlt.User {
			return true
		}
	}
	return false
}

func (ptz *Ptz) IsAdmin() bool {
	if ptz.GroupInfo == nil {
		return false
	}
	for _, p := range ptz.GroupInfo.Participants {
		if matchParticipant(p, ptz.Sender, ptz.SenderAlt) {
			return p.IsAdmin || p.IsSuperAdmin
		}
	}
	return false
}

func (ptz *Ptz) IsSuperAdmin() bool {
	if ptz.GroupInfo == nil {
		return false
	}
	for _, p := range ptz.GroupInfo.Participants {
		if matchParticipant(p, ptz.Sender, ptz.SenderAlt) {
			return p.IsSuperAdmin
		}
	}
	return false
}

func (ptz *Ptz) IsBotAdmin() bool {
	if ptz.GroupInfo == nil {
		return false
	}
	botID := ptz.Bot.Client.Store.ID
	if botID == nil {
		return false
	}
	botSender := *botID
	botLID := ptz.Bot.Client.Store.LID
	for _, p := range ptz.GroupInfo.Participants {
		if matchParticipant(p, botSender, botLID) {
			return p.IsAdmin || p.IsSuperAdmin
		}
	}
	return false
}

func (ptz *Ptz) LoadGroupInfo() error {
	if !ptz.IsGroup {
		return nil
	}
	info, err := ptz.Bot.Client.GetGroupInfo(context.Background(), ptz.Chat)
	if err != nil {
		return err
	}
	ptz.GroupInfo = info
	return nil
}

func (ptz *Ptz) GetPushName() string {
	if ptz.Info.PushName != "" && ptz.Info.PushName != "-" {
		return ptz.Info.PushName
	}
	return fmt.Sprintf("@%s", ptz.Sender.User)
}

func (ptz *Ptz) GetSenderName() string {
	if ptz.IsGroup && ptz.GroupInfo != nil {
		for _, p := range ptz.GroupInfo.Participants {
			if matchParticipant(p, ptz.Sender, ptz.SenderAlt) && p.DisplayName != "" {
				return p.DisplayName
			}
		}
	}
	if ptz.Info.PushName != "" && ptz.Info.PushName != "-" {
		return ptz.Info.PushName
	}
	return ptz.Sender.User
}

func (ptz *Ptz) React(emoji string) error {
	return serialize.SendReaction(ptz.Bot.Client, ptz.Chat, ptz.Info.ID, ptz.Sender, emoji)
}

func (ptz *Ptz) Unreact() error {
	return serialize.RemoveReaction(ptz.Bot.Client, ptz.Chat, ptz.Info.ID, ptz.Sender)
}

func (ptz *Ptz) ReplyText(text string) error {
	return serialize.SendTextReply(ptz.Bot.Client, ptz.Chat, text, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyTextID(text string) (types.MessageID, error) {
	return serialize.SendTextReplyID(ptz.Bot.Client, ptz.Chat, text, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyImage(data []byte, mime, caption string) error {
	return serialize.SendImageReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyImageID(data []byte, mime, caption string) (types.MessageID, error) {
	return serialize.SendImageReplyID(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyVideo(data []byte, mime, caption string) error {
	return serialize.SendVideoReply(ptz.Bot.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyAudio(data []byte, mime string) error {
	return serialize.SendAudioReply(ptz.Bot.Client, ptz.Chat, data, mime, false, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplySticker(data []byte, mime string, animated bool) error {
	return serialize.SendStickerReply(ptz.Bot.Client, ptz.Chat, data, mime, animated, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyDocument(data []byte, mime, filename, caption string) error {
	return serialize.SendDocumentReply(ptz.Bot.Client, ptz.Chat, data, mime, filename, caption, ptz.Message, ptz.Info)
}

func (ptz *Ptz) GetReplyText() string {
	if ptz.Message == nil || ptz.Message.ExtendedTextMessage == nil {
		return ""
	}

	ext := ptz.Message.ExtendedTextMessage
	if ext.ContextInfo == nil || ext.ContextInfo.QuotedMessage == nil {
		return ""
	}

	return ExtractBody(ext.ContextInfo.QuotedMessage)
}

func (ptz *Ptz) GetPhoneJID() types.JID {
	if ptz.Sender.Server == types.HiddenUserServer {
		if !ptz.SenderAlt.IsEmpty() && ptz.SenderAlt.Server == types.DefaultUserServer {
			return types.NewJID(ptz.SenderAlt.User, types.DefaultUserServer)
		}
	}
	return types.NewJID(ptz.Sender.User, types.DefaultUserServer)
}

func (ptz *Ptz) ReplyTextMention(text string, mentionedJIDs []types.JID) error {
	return serialize.SendTextReplyMention(ptz.Bot.Client, ptz.Chat, text, mentionedJIDs, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
