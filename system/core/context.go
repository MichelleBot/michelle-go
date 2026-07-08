package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"michelle/system/serialize"
)

type Ptz struct {
	Bot       *Bot
	Client    *whatsmeow.Client
	Event     *events.Message
	Message   *waE2E.Message
	Info      types.MessageInfo
	Args      []string
	RawArgs   string
	Command   string
	Prefix    string
	IsGroup   bool
	IsFromMe  bool
	Sender    types.JID
	SenderAlt types.JID
	Chat      types.JID
	GroupInfo *types.GroupInfo
}

func NewPtz(bot *Bot, evt *events.Message) *Ptz {
	body := ExtractBody(evt.Message)
	parts, cmd, rawArgs, args, prefix := parseCommandParts(bot, body)

	_ = parts

	return &Ptz{
		Bot:       bot,
		Client:    bot.Client,
		Event:     evt,
		Message:   evt.Message,
		Info:      evt.Info,
		Args:      args,
		RawArgs:   rawArgs,
		Command:   cmd,
		Prefix:    prefix,
		IsGroup:   evt.Info.IsGroup,
		IsFromMe:  evt.Info.IsFromMe,
		Sender:    evt.Info.Sender,
		SenderAlt: evt.Info.SenderAlt,
		Chat:      evt.Info.Chat,
	}
}

func NewPtzFromNormalizedMessage(bot *Bot, client *whatsmeow.Client, msg *NormalizedMessage) *Ptz {
	if msg == nil {
		return nil
	}

	_, cmd, rawArgs, args, prefix := parseCommandParts(bot, msg.Body)

	if client == nil {
		client = bot.Client
	}

	return &Ptz{
		Bot:       bot,
		Client:    client,
		Event:     msg.Event,
		Message:   msg.Message,
		Info:      msg.Info,
		Args:      args,
		RawArgs:   rawArgs,
		Command:   cmd,
		Prefix:    prefix,
		IsGroup:   msg.IsGroup,
		IsFromMe:  msg.IsFromMe,
		Sender:    msg.Sender,
		SenderAlt: msg.SenderAlt,
		Chat:      msg.Chat,
	}
}

func parseCommandParts(bot *Bot, body string) ([]string, string, string, []string, string) {
	body = strings.TrimSpace(body)
	var cmd, rawArgs, usedPrefix string
	args := []string{}
	var parts []string

	for _, prefix := range bot.Config.Prefixes {
		if strings.HasPrefix(body, prefix) {
			usedPrefix = prefix
			trimmedBody := strings.TrimSpace(strings.TrimPrefix(body, prefix))
			parts = strings.Fields(trimmedBody)
			if len(parts) > 0 {
				cmd = parts[0] // Removed strings.ToLower() to preserve case for comparison if needed, or handle lowercase here for consistency.
				if len(parts) > 1 {
					args = parts[1:]
					// Re-extract rawArgs based on original split structure
					rawArgs = strings.TrimSpace(strings.TrimPrefix(trimmedBody, parts[0]))
				}
			}
			break
		}
	}

	return parts, cmd, rawArgs, args, usedPrefix
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
	botID := ptz.Client.Store.ID
	if botID == nil {
		return false
	}
	botSender := *botID
	botLID := ptz.Client.Store.LID
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
	info, err := ptz.Client.GetGroupInfo(context.Background(), ptz.Chat)
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
	ptz.Bot.Log.Infof("Ptz React - Client Pointer: %p", ptz.Client)
	return serialize.SendReaction(ptz.Client, ptz.Chat, ptz.Info.ID, ptz.Sender, emoji)
}

func (ptz *Ptz) Unreact() error {
	ptz.Bot.Log.Infof("Ptz Unreact - Client Pointer: %p", ptz.Client)
	return serialize.RemoveReaction(ptz.Client, ptz.Chat, ptz.Info.ID, ptz.Sender)
}

func (ptz *Ptz) ReplyText(text string) error {
	ptz.Bot.Log.Infof("ReplyText - Ptz Bot Pointer: %p, Client Pointer: %p", ptz.Bot, ptz.Client)
	return serialize.SendTextReply(ptz.Client, ptz.Chat, text, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyTextID(text string) (types.MessageID, error) {
	ptz.Bot.Log.Infof("Ptz ReplyTextID - Client Pointer: %p", ptz.Client)
	return serialize.SendTextReplyID(ptz.Client, ptz.Chat, text, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyImage(data []byte, mime, caption string) error {
	ptz.Bot.Log.Infof("Ptz ReplyImage - Client Pointer: %p", ptz.Client)
	return serialize.SendImageReply(ptz.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyImageID(data []byte, mime, caption string) (types.MessageID, error) {
	ptz.Bot.Log.Infof("Ptz ReplyImageID - Client Pointer: %p", ptz.Client)
	return serialize.SendImageReplyID(ptz.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyVideo(data []byte, mime, caption string) error {
	ptz.Bot.Log.Infof("Ptz ReplyVideo - Client Pointer: %p", ptz.Client)
	return serialize.SendVideoReply(ptz.Client, ptz.Chat, data, mime, caption, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyAudio(data []byte, mime string) error {
	ptz.Bot.Log.Infof("Ptz ReplyAudio - Client Pointer: %p", ptz.Client)
	return serialize.SendAudioReply(ptz.Client, ptz.Chat, data, mime, false, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplySticker(data []byte, mime string, animated bool) error {
	return serialize.SendStickerReply(ptz.Client, ptz.Chat, data, mime, animated, ptz.Message, ptz.Info)
}

func (ptz *Ptz) ReplyDocument(data []byte, mime, filename, caption string) error {
	return serialize.SendDocumentReply(ptz.Client, ptz.Chat, data, mime, filename, caption, ptz.Message, ptz.Info)
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
	return serialize.SendTextReplyMention(ptz.Client, ptz.Chat, text, mentionedJIDs, ptz.Message, ptz.Info)
}

func (ptz *Ptz) SendTextMention(text string, mentionedJIDs []types.JID) error {
	return serialize.SendTextMention(ptz.Client, ptz.Chat, text, mentionedJIDs)
}

func (ptz *Ptz) ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
