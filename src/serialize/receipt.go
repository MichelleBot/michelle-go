package serialize

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func MarkRead(client *whatsmeow.Client, ids []types.MessageID, chat, sender types.JID) error {
	return client.MarkRead(context.Background(), ids, time.Now(), chat, sender)
}

func MarkReadSingle(client *whatsmeow.Client, id types.MessageID, chat, sender types.JID) error {
	return client.MarkRead(context.Background(), []types.MessageID{id}, time.Now(), chat, sender)
}

func MarkPlayed(client *whatsmeow.Client, id types.MessageID, chat, sender types.JID) error {
	return client.MarkRead(context.Background(), []types.MessageID{id}, time.Now(), chat, sender, types.ReceiptTypePlayed)
}

func SetForceActiveDeliveryReceipts(client *whatsmeow.Client, active bool) {
	client.SetForceActiveDeliveryReceipts(active)
}

func RejectCall(client *whatsmeow.Client, callFrom types.JID, callID string) error {
	return client.RejectCall(context.Background(), callFrom, callID)
}
