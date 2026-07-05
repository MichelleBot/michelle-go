package core

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type StoredMessage struct {
	ID           string
	Chat         string
	Sender       string
	PushName     string
	Kind         MessageKind
	Content      string
	Timestamp    time.Time
	MediaData    []byte
	MIME         string
	Filename     string
	Caption      string
	IsAnimated   bool
	Latitude     float64
	Longitude    float64
	PlaceName    string
	Address      string
	ContactName  string
	ContactPhone string
}

type MessageStore struct {
	mu   sync.RWMutex
	data map[string]*StoredMessage
}

func NewMessageStore() *MessageStore {
	ms := &MessageStore{data: make(map[string]*StoredMessage)}
	go ms.cleanup()
	return ms
}

func (ms *MessageStore) Save(msg *NormalizedMessage) {
	ms.SaveStored(NewStoredMessage(msg))
}

func (ms *MessageStore) SaveStored(stored *StoredMessage) {
	if ms == nil || stored == nil {
		return
	}
	if stored.ID == "" {
		return
	}

	ms.mu.Lock()
	ms.data[ms.key(stored.Chat, stored.ID)] = stored
	ms.mu.Unlock()
}

func (ms *MessageStore) Get(chat, id string) (*StoredMessage, bool) {
	if ms == nil || id == "" {
		return nil, false
	}
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	stored, ok := ms.data[ms.key(chat, id)]
	return stored, ok
}

func (ms *MessageStore) Update(chat, id, content string, kind MessageKind) (*StoredMessage, bool) {
	if ms == nil || id == "" {
		return nil, false
	}
	ms.mu.Lock()
	defer ms.mu.Unlock()
	stored, ok := ms.data[ms.key(chat, id)]
	if !ok {
		return nil, false
	}
	stored.Content = content
	stored.Kind = kind
	return stored, true
}

func NewStoredMessage(msg *NormalizedMessage) *StoredMessage {
	if msg == nil {
		return nil
	}
	stored := &StoredMessage{
		ID:        msg.Info.ID,
		Chat:      msg.Chat.String(),
		Sender:    msg.Sender.User,
		PushName:  msg.Info.PushName,
		Kind:      msg.Kind,
		Content:   msg.Content,
		Timestamp: msg.Info.Timestamp,
	}

	if msg.Message != nil {
		switch {
		case msg.Message.ImageMessage != nil:
			stored.Caption = msg.Message.GetImageMessage().GetCaption()
		case msg.Message.VideoMessage != nil:
			stored.Caption = msg.Message.GetVideoMessage().GetCaption()
		case msg.Message.DocumentMessage != nil:
			stored.Caption = msg.Message.GetDocumentMessage().GetCaption()
			stored.Filename = msg.Message.GetDocumentMessage().GetFileName()
		case msg.Message.StickerMessage != nil:
			stored.IsAnimated = msg.Message.GetStickerMessage().GetIsAnimated()
		case msg.Message.LocationMessage != nil:
			loc := msg.Message.GetLocationMessage()
			stored.Latitude = loc.GetDegreesLatitude()
			stored.Longitude = loc.GetDegreesLongitude()
			stored.PlaceName = loc.GetName()
			stored.Address = loc.GetAddress()
		case msg.Message.LiveLocationMessage != nil:
			loc := msg.Message.GetLiveLocationMessage()
			stored.Latitude = loc.GetDegreesLatitude()
			stored.Longitude = loc.GetDegreesLongitude()
			stored.PlaceName = "Live Location"
			stored.Address = loc.GetCaption()
		case msg.Message.ContactMessage != nil:
			contact := msg.Message.GetContactMessage()
			stored.ContactName = contact.GetDisplayName()
			stored.ContactPhone = extractPhoneFromVCard(contact.GetVcard())
		case msg.Message.ContactsArrayMessage != nil:
			contacts := msg.Message.GetContactsArrayMessage().GetContacts()
			if len(contacts) > 0 {
				stored.ContactName = contacts[0].GetDisplayName()
				stored.ContactPhone = extractPhoneFromVCard(contacts[0].GetVcard())
			}
		}
	}

	return stored
}

func extractPhoneFromVCard(vcard string) string {
	if vcard == "" {
		return ""
	}
	for _, line := range strings.Split(vcard, "\n") {
		if idx := strings.Index(line, "waid="); idx >= 0 {
			rest := line[idx+5:]
			for i, r := range rest {
				if r < '0' || r > '9' {
					return rest[:i]
				}
			}
			return rest
		}
	}
	return ""
}

func (ms *MessageStore) key(chat, id string) string {
	return fmt.Sprintf("%s|%s", chat, id)
}

func (ms *MessageStore) cleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		cutoff := time.Now().Add(-48 * time.Hour)
		ms.mu.Lock()
		for key, stored := range ms.data {
			if stored.Timestamp.Before(cutoff) {
				delete(ms.data, key)
			}
		}
		ms.mu.Unlock()
	}
}
