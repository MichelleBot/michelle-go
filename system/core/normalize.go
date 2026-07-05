package core

import (
	"strings"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type NormalizedEvent struct {
	Kind      EventKind
	Name      string
	Message   *NormalizedMessage
	Timestamp time.Time
	Raw       interface{}
}

type NormalizedMessage struct {
	Kind        MessageKind
	EventKind   EventKind
	Body        string
	Text        string
	Caption     string
	Content     string
	QuotedID    string
	QuotedText  string
	Info        types.MessageInfo
	Event       *events.Message
	Message     *waE2E.Message
	RawMessage  *waE2E.Message
	IsGroup     bool
	IsFromMe    bool
	IsEdit      bool
	IsViewOnce  bool
	IsEphemeral bool
	Sender      types.JID
	SenderAlt   types.JID
	Chat        types.JID
	Reaction    *NormalizedReaction
	Edit        *NormalizedEdit
	Revoke      *NormalizedRevoke
	Poll        *NormalizedPoll
}

type NormalizedReaction struct {
	Text         string
	TargetID     string
	TargetChat   string
	TargetSender string
	IsRemoval    bool
}

type NormalizedEdit struct {
	EditedText string
	OriginalID string
}

type NormalizedRevoke struct {
	ProtocolType string
	TargetID     string
}

type NormalizedPoll struct {
	ID          string
	TargetID    string
	IsCreation  bool
	IsUpdate    bool
	Name        string
	OptionCount int
	UpdateCount int
}

func NormalizeMessageEvent(evt *events.Message) *NormalizedEvent {
	if evt == nil {
		return nil
	}

	unwrapped := *evt
	unwrapped.UnwrapRaw()

	body := ExtractBody(unwrapped.Message)
	msgKind, content, text, caption := classifyNormalizedMessage(unwrapped.Message)
	quotedID, quotedText := extractQuotedMessageMeta(unwrapped.Message)

	eventKind := EventMessage
	switch {
	case unwrapped.IsEdit:
		eventKind = EventEdit
	case unwrapped.Message.GetReactionMessage() != nil:
		eventKind = EventReaction
	case unwrapped.Message.GetPollCreationMessage() != nil || unwrapped.Message.GetPollUpdateMessage() != nil:
		eventKind = EventPoll
	case unwrapped.Message.GetProtocolMessage() != nil:
		eventKind = EventRevoke
	}

	return &NormalizedEvent{
		Kind:      eventKind,
		Name:      string(eventKind),
		Timestamp: unwrapped.Info.Timestamp,
		Raw:       evt,
		Message: &NormalizedMessage{
			Kind:        msgKind,
			EventKind:   eventKind,
			Body:        body,
			Text:        text,
			Caption:     caption,
			Content:     content,
			QuotedID:    quotedID,
			QuotedText:  quotedText,
			Info:        unwrapped.Info,
			Event:       &unwrapped,
			Message:     unwrapped.Message,
			RawMessage:  evt.Message,
			IsGroup:     unwrapped.Info.IsGroup,
			IsFromMe:    unwrapped.Info.IsFromMe,
			IsEdit:      unwrapped.IsEdit,
			IsViewOnce:  unwrapped.IsViewOnce,
			IsEphemeral: unwrapped.IsEphemeral,
			Sender:      unwrapped.Info.Sender,
			SenderAlt:   unwrapped.Info.SenderAlt,
			Chat:        unwrapped.Info.Chat,
			Reaction:    extractReaction(unwrapped.Message),
			Edit:        extractEdit(unwrapped.Message, unwrapped.Info.ID, quotedID, text),
			Revoke:      extractRevoke(unwrapped.Message),
			Poll:        extractPoll(unwrapped.Message, unwrapped.Info.ID),
		},
	}
}

func NormalizeReceiptEvent(evt *events.Receipt) *NormalizedEvent {
	if evt == nil {
		return nil
	}
	return &NormalizedEvent{Kind: EventReceipt, Name: string(EventReceipt), Timestamp: evt.Timestamp, Raw: evt}
}

func NormalizePresenceEvent(evt *events.Presence) *NormalizedEvent {
	if evt == nil {
		return nil
	}
	return &NormalizedEvent{Kind: EventPresence, Name: string(EventPresence), Timestamp: time.Now(), Raw: evt}
}

func NormalizeChatPresenceEvent(evt *events.ChatPresence) *NormalizedEvent {
	if evt == nil {
		return nil
	}
	return &NormalizedEvent{Kind: EventPresence, Name: string(EventPresence), Timestamp: time.Now(), Raw: evt}
}

func NormalizeCallEvent(name string, evt interface{}) *NormalizedEvent {
	return &NormalizedEvent{Kind: EventCall, Name: name, Timestamp: time.Now(), Raw: evt}
}

func classifyNormalizedMessage(msg *waE2E.Message) (MessageKind, string, string, string) {
	if msg == nil {
		return MessageUnknown, "", "", ""
	}

	switch {
	case msg.Conversation != nil:
		text := strings.TrimSpace(msg.GetConversation())
		return MessageText, text, text, ""
	case msg.ExtendedTextMessage != nil:
		text := strings.TrimSpace(msg.GetExtendedTextMessage().GetText())
		return MessageText, text, text, ""
	case msg.ImageMessage != nil:
		caption := strings.TrimSpace(msg.GetImageMessage().GetCaption())
		return MessageImage, caption, "", caption
	case msg.VideoMessage != nil:
		caption := strings.TrimSpace(msg.GetVideoMessage().GetCaption())
		return MessageVideo, caption, "", caption
	case msg.AudioMessage != nil:
		if msg.GetAudioMessage().GetPTT() {
			return MessageVoice, "", "", ""
		}
		return MessageAudio, "", "", ""
	case msg.DocumentMessage != nil:
		caption := strings.TrimSpace(msg.GetDocumentMessage().GetCaption())
		name := strings.TrimSpace(msg.GetDocumentMessage().GetFileName())
		content := caption
		if content == "" {
			content = name
		}
		return MessageDocument, content, "", caption
	case msg.StickerMessage != nil:
		return MessageSticker, "", "", ""
	case msg.ContactMessage != nil:
		name := strings.TrimSpace(msg.GetContactMessage().GetDisplayName())
		return MessageContact, name, name, ""
	case msg.ContactsArrayMessage != nil:
		return MessageContacts, "kontak", "", ""
	case msg.LiveLocationMessage != nil:
		caption := strings.TrimSpace(msg.GetLiveLocationMessage().GetCaption())
		return MessageLiveLocation, caption, "", caption
	case msg.LocationMessage != nil:
		return MessageLocation, "", "", ""
	case msg.ButtonsResponseMessage != nil:
		selected := strings.TrimSpace(msg.GetButtonsResponseMessage().GetSelectedDisplayText())
		return MessageButtonReply, selected, selected, ""
	case msg.TemplateButtonReplyMessage != nil:
		selected := strings.TrimSpace(msg.GetTemplateButtonReplyMessage().GetSelectedDisplayText())
		return MessageTemplateReply, selected, selected, ""
	case msg.ListResponseMessage != nil:
		title := strings.TrimSpace(msg.GetListResponseMessage().GetTitle())
		return MessageListReply, title, title, ""
	case msg.ButtonsMessage != nil:
		text := strings.TrimSpace(msg.GetButtonsMessage().GetContentText())
		return MessageButtons, text, text, ""
	case msg.ListMessage != nil:
		title := strings.TrimSpace(msg.GetListMessage().GetTitle())
		return MessageList, title, title, ""
	case msg.TemplateMessage != nil:
		return MessageTemplate, "template", "", ""
	case msg.InteractiveResponseMessage != nil:
		return MessageInteractiveReply, "interactive_response", "", ""
	case msg.InteractiveMessage != nil:
		return MessageInteractive, "interactive", "", ""
	case msg.GroupInviteMessage != nil:
		groupName := strings.TrimSpace(msg.GetGroupInviteMessage().GetGroupName())
		return MessageGroupInvite, groupName, groupName, ""
	case msg.ProductMessage != nil:
		return MessageProduct, "product", "", ""
	case msg.OrderMessage != nil:
		return MessageOrder, "order", "", ""
	case msg.SendPaymentMessage != nil || msg.RequestPaymentMessage != nil || msg.DeclinePaymentRequestMessage != nil || msg.CancelPaymentRequestMessage != nil || msg.PaymentInviteMessage != nil || msg.InvoiceMessage != nil:
		return MessagePayment, "payment", "", ""
	case msg.RequestPhoneNumberMessage != nil:
		return MessageRequestPhone, "request_phone", "", ""
	case msg.KeepInChatMessage != nil:
		return MessageKeepInChat, "keep_in_chat", "", ""
	case msg.HighlyStructuredMessage != nil:
		return MessageStructured, "structured", "", ""
	case msg.ReactionMessage != nil:
		reaction := strings.TrimSpace(msg.GetReactionMessage().GetText())
		return MessageReaction, reaction, reaction, ""
	case msg.PollCreationMessage != nil || msg.PollUpdateMessage != nil:
		return MessagePoll, "", "", ""
	case msg.ProtocolMessage != nil:
		return MessageProtocol, "", "", ""
	case msg.EventInviteMessage != nil:
		return MessageEventInvite, "", "", ""
	default:
		return MessageUnknown, "", "", ""
	}
}

func extractQuotedMessageMeta(msg *waE2E.Message) (string, string) {
	if msg == nil {
		return "", ""
	}

	ctx := extractContextInfo(msg)
	if ctx == nil {
		return "", ""
	}

	return ctx.GetStanzaID(), ExtractBody(ctx.GetQuotedMessage())
}

func extractContextInfo(msg *waE2E.Message) *waE2E.ContextInfo {
	if msg == nil {
		return nil
	}

	switch {
	case msg.ExtendedTextMessage != nil:
		return msg.ExtendedTextMessage.GetContextInfo()
	case msg.ImageMessage != nil:
		return msg.ImageMessage.GetContextInfo()
	case msg.VideoMessage != nil:
		return msg.VideoMessage.GetContextInfo()
	case msg.AudioMessage != nil:
		return msg.AudioMessage.GetContextInfo()
	case msg.DocumentMessage != nil:
		return msg.DocumentMessage.GetContextInfo()
	}

	return nil
}

func extractReaction(msg *waE2E.Message) *NormalizedReaction {
	if msg == nil || msg.GetReactionMessage() == nil {
		return nil
	}

	reaction := msg.GetReactionMessage()
	key := reaction.GetKey()
	return &NormalizedReaction{
		Text:         strings.TrimSpace(reaction.GetText()),
		TargetID:     key.GetID(),
		TargetChat:   key.GetRemoteJID(),
		TargetSender: key.GetParticipant(),
		IsRemoval:    strings.TrimSpace(reaction.GetText()) == "",
	}
}

func extractEdit(msg *waE2E.Message, infoID, quotedID, text string) *NormalizedEdit {
	if msg == nil || msg.GetEditedMessage() == nil {
		return nil
	}
	originalID := quotedID
	if originalID == "" {
		originalID = infoID
	}
	return &NormalizedEdit{EditedText: text, OriginalID: originalID}
}

func extractRevoke(msg *waE2E.Message) *NormalizedRevoke {
	if msg == nil || msg.GetProtocolMessage() == nil {
		return nil
	}
	pm := msg.GetProtocolMessage()
	return &NormalizedRevoke{ProtocolType: pm.GetType().String(), TargetID: pm.GetKey().GetID()}
}

func extractPoll(msg *waE2E.Message, infoID string) *NormalizedPoll {
	if msg == nil {
		return nil
	}
	if poll := msg.GetPollCreationMessage(); poll != nil {
		return &NormalizedPoll{ID: infoID, IsCreation: true, Name: strings.TrimSpace(poll.GetName()), OptionCount: len(poll.GetOptions())}
	}
	if update := msg.GetPollUpdateMessage(); update != nil {
		targetID := ""
		if key := update.GetPollCreationMessageKey(); key != nil {
			targetID = key.GetID()
		}
		return &NormalizedPoll{ID: infoID, TargetID: targetID, IsUpdate: true, UpdateCount: 1}
	}
	return nil
}
