package serialize

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func GetPrivacySettings(client *whatsmeow.Client) types.PrivacySettings {
	return client.GetPrivacySettings(context.Background())
}

func SetPrivacySetting(client *whatsmeow.Client, name types.PrivacySettingType, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), name, value)
}

func SetPrivacyLastSeen(client *whatsmeow.Client, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), types.PrivacySettingTypeLastSeen, value)
}

func SetPrivacyProfilePhoto(client *whatsmeow.Client, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), types.PrivacySettingTypeProfile, value)
}

func SetPrivacyStatus(client *whatsmeow.Client, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), types.PrivacySettingTypeStatus, value)
}

func SetPrivacyReadReceipts(client *whatsmeow.Client, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), types.PrivacySettingTypeReadReceipts, value)
}

func SetPrivacyGroupAdd(client *whatsmeow.Client, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), types.PrivacySettingTypeGroupAdd, value)
}

func SetPrivacyOnline(client *whatsmeow.Client, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), types.PrivacySettingTypeOnline, value)
}

func SetPrivacyCallAdd(client *whatsmeow.Client, value types.PrivacySetting) (types.PrivacySettings, error) {
	return client.SetPrivacySetting(context.Background(), types.PrivacySettingTypeCallAdd, value)
}

func SetDefaultDisappearingTimer(client *whatsmeow.Client, dur time.Duration) error {
	return client.SetDefaultDisappearingTimer(context.Background(), dur)
}

func SetDisappearingTimer(client *whatsmeow.Client, chat types.JID, dur time.Duration) error {
	return client.SetDisappearingTimer(context.Background(), chat, dur, time.Time{})
}

func GetBlocklist(client *whatsmeow.Client) (*types.Blocklist, error) {
	return client.GetBlocklist(context.Background())
}

func BlockUser(client *whatsmeow.Client, jid types.JID) (*types.Blocklist, error) {
	return client.UpdateBlocklist(context.Background(), jid, events.BlocklistChangeActionBlock)
}

func UnblockUser(client *whatsmeow.Client, jid types.JID) (*types.Blocklist, error) {
	return client.UpdateBlocklist(context.Background(), jid, events.BlocklistChangeActionUnblock)
}

func SubscribePresence(client *whatsmeow.Client, jid types.JID) error {
	return client.SubscribePresence(context.Background(), jid)
}
