package serialize

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/types"
)

func MuteChat(client *whatsmeow.Client, chat types.JID, mute bool, duration time.Duration) error {
	return client.SendAppState(context.Background(), appstate.BuildMute(chat, mute, duration))
}

func MuteChatForever(client *whatsmeow.Client, chat types.JID) error {
	return client.SendAppState(context.Background(), appstate.BuildMute(chat, true, 0))
}

func UnmuteChat(client *whatsmeow.Client, chat types.JID) error {
	return client.SendAppState(context.Background(), appstate.BuildMute(chat, false, 0))
}

func PinChat(client *whatsmeow.Client, chat types.JID, pin bool) error {
	return client.SendAppState(context.Background(), appstate.BuildPin(chat, pin))
}

func ArchiveChat(client *whatsmeow.Client, chat types.JID, archive bool) error {
	return client.SendAppState(context.Background(), appstate.BuildArchive(chat, archive, time.Time{}, nil))
}

func ArchiveChatWithMessage(client *whatsmeow.Client, chat types.JID, archive bool, lastMsgTS time.Time, lastMsgKey *waCommon.MessageKey) error {
	return client.SendAppState(context.Background(), appstate.BuildArchive(chat, archive, lastMsgTS, lastMsgKey))
}

func MarkChatRead(client *whatsmeow.Client, chat types.JID, read bool) error {
	return client.SendAppState(context.Background(), appstate.BuildMarkChatAsRead(chat, read, time.Time{}, nil))
}

func MarkChatReadWithMessage(client *whatsmeow.Client, chat types.JID, read bool, lastMsgTS time.Time, lastMsgKey *waCommon.MessageKey) error {
	return client.SendAppState(context.Background(), appstate.BuildMarkChatAsRead(chat, read, lastMsgTS, lastMsgKey))
}

func StarMessage(client *whatsmeow.Client, chat, sender types.JID, msgID types.MessageID, fromMe, starred bool) error {
	return client.SendAppState(context.Background(), appstate.BuildStar(chat, sender, msgID, fromMe, starred))
}

func DeleteChat(client *whatsmeow.Client, chat types.JID, lastMsgTS time.Time, lastMsgKey *waCommon.MessageKey, deleteMedia bool) error {
	return client.SendAppState(context.Background(), appstate.BuildDeleteChat(chat, lastMsgTS, lastMsgKey, deleteMedia))
}

func LabelChat(client *whatsmeow.Client, chat types.JID, labelID string, labeled bool) error {
	return client.SendAppState(context.Background(), appstate.BuildLabelChat(chat, labelID, labeled))
}

func LabelMessage(client *whatsmeow.Client, chat types.JID, labelID, msgID string, labeled bool) error {
	return client.SendAppState(context.Background(), appstate.BuildLabelMessage(chat, labelID, msgID, labeled))
}

func EditLabel(client *whatsmeow.Client, labelID, name string, color int32, deleted bool) error {
	return client.SendAppState(context.Background(), appstate.BuildLabelEdit(labelID, name, color, deleted))
}

func SetPushName(client *whatsmeow.Client, pushName string) error {
	return client.SendAppState(context.Background(), appstate.BuildSettingPushName(pushName))
}

func SyncAppState(client *whatsmeow.Client, fullSync bool) error {
	for _, name := range appstate.AllPatchNames {
		if err := client.FetchAppState(context.Background(), name, fullSync, false); err != nil {
			return err
		}
	}
	return nil
}
