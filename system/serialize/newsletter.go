package serialize

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func GetNewsletterInfo(client *whatsmeow.Client, jid types.JID) (*types.NewsletterMetadata, error) {
	return client.GetNewsletterInfo(context.Background(), jid)
}

func GetNewsletterInfoWithInvite(client *whatsmeow.Client, key string) (*types.NewsletterMetadata, error) {
	return client.GetNewsletterInfoWithInvite(context.Background(), key)
}

func GetSubscribedNewsletters(client *whatsmeow.Client) ([]*types.NewsletterMetadata, error) {
	return client.GetSubscribedNewsletters(context.Background())
}

func CreateNewsletter(client *whatsmeow.Client, name, description string, picture []byte) (*types.NewsletterMetadata, error) {
	return client.CreateNewsletter(context.Background(), whatsmeow.CreateNewsletterParams{
		Name:        name,
		Description: description,
		Picture:     picture,
	})
}

func FollowNewsletter(client *whatsmeow.Client, jid types.JID) error {
	return client.FollowNewsletter(context.Background(), jid)
}

func UnfollowNewsletter(client *whatsmeow.Client, jid types.JID) error {
	return client.UnfollowNewsletter(context.Background(), jid)
}

func NewsletterToggleMute(client *whatsmeow.Client, jid types.JID, mute bool) error {
	return client.NewsletterToggleMute(context.Background(), jid, mute)
}

func GetNewsletterMessages(client *whatsmeow.Client, jid types.JID, count int, before types.MessageServerID) ([]*types.NewsletterMessage, error) {
	params := &whatsmeow.GetNewsletterMessagesParams{
		Count: count,
	}
	if before != 0 {
		params.Before = before
	}
	return client.GetNewsletterMessages(context.Background(), jid, params)
}

func NewsletterSubscribeLiveUpdates(client *whatsmeow.Client, jid types.JID) (time.Duration, error) {
	return client.NewsletterSubscribeLiveUpdates(context.Background(), jid)
}

func NewsletterMarkViewed(client *whatsmeow.Client, jid types.JID, serverIDs []types.MessageServerID) error {
	return client.NewsletterMarkViewed(context.Background(), jid, serverIDs)
}

func NewsletterSendReaction(client *whatsmeow.Client, jid types.JID, serverID types.MessageServerID, reaction string, messageID types.MessageID) error {
	return client.NewsletterSendReaction(context.Background(), jid, serverID, reaction, messageID)
}

func SendNewsletterText(client *whatsmeow.Client, jid types.JID, text string) error {
	_, err := client.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: proto.String(text),
	})
	return err
}

func SendNewsletterImage(client *whatsmeow.Client, jid types.JID, imageData []byte, mime, caption string) error {
	resp, err := client.UploadNewsletter(context.Background(), imageData, whatsmeow.MediaImage)
	if err != nil {
		return err
	}
	ext := mimeToExt(mime)
	thumb, _ := GenerateJPEGThumbnail(imageData, ext)
	dim, _ := GetMediaDimensions(imageData, ext)
	h := dim.Height
	w := dim.Width
	ts := time.Now().Unix()
	_, err = client.SendMessage(context.Background(), jid, &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:           proto.String(caption),
			Mimetype:          proto.String(mime),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
			JPEGThumbnail:     thumb,
			Height:            &h,
			Width:             &w,
		},
	}, whatsmeow.SendRequestExtra{MediaHandle: resp.Handle})
	return err
}

func SendNewsletterVideo(client *whatsmeow.Client, jid types.JID, videoData []byte, mime, caption string) error {
	resp, err := client.UploadNewsletter(context.Background(), videoData, whatsmeow.MediaVideo)
	if err != nil {
		return err
	}
	ext := mimeToExt(mime)
	thumb, _ := GenerateJPEGThumbnail(videoData, ext)
	dim, _ := GetMediaDimensions(videoData, ext)
	secs, _ := GetVideoDurationSeconds(videoData, ext)
	h := dim.Height
	w := dim.Width
	ts := time.Now().Unix()
	_, err = client.SendMessage(context.Background(), jid, &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			Caption:           proto.String(caption),
			Mimetype:          proto.String(mime),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
			JPEGThumbnail:     thumb,
			Height:            &h,
			Width:             &w,
			Seconds:           &secs,
		},
	}, whatsmeow.SendRequestExtra{MediaHandle: resp.Handle})
	return err
}

func SendNewsletterAudio(client *whatsmeow.Client, jid types.JID, audioData []byte, mime string) error {
	ext := mimeToExt(mime)
	oggData, err := ToOggOpus(audioData, ext)
	if err != nil {
		oggData = audioData
	}

	secs, _ := GetVideoDurationSeconds(oggData, ".ogg")
	ts := time.Now().Unix()

	resp, err := client.UploadNewsletter(context.Background(), oggData, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}

	_, err = client.SendMessage(context.Background(), jid, &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			Mimetype:          proto.String("audio/ogg; codecs=opus"),
			PTT:               proto.Bool(false),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			Seconds:           &secs,
			MediaKeyTimestamp: &ts,
		},
	}, whatsmeow.SendRequestExtra{MediaHandle: resp.Handle})
	return err
}

func SendNewsletterDocument(client *whatsmeow.Client, jid types.JID, docData []byte, mime, filename, caption string) error {
	resp, err := client.UploadNewsletter(context.Background(), docData, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}
	ts := time.Now().Unix()
	_, err = client.SendMessage(context.Background(), jid, &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			Title:             proto.String(filename),
			FileName:          proto.String(filename),
			Caption:           proto.String(caption),
			Mimetype:          proto.String(mime),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
		},
	}, whatsmeow.SendRequestExtra{MediaHandle: resp.Handle})
	return err
}

func AcceptTOSNotice(client *whatsmeow.Client, noticeID, stage string) error {
	return client.AcceptTOSNotice(context.Background(), noticeID, stage)
}

func GetNewsletterMessageUpdates(client *whatsmeow.Client, jid types.JID, count int, since time.Time, after types.MessageServerID) ([]*types.NewsletterMessage, error) {
	return client.GetNewsletterMessageUpdates(context.Background(), jid, &whatsmeow.GetNewsletterUpdatesParams{
		Count: count,
		Since: since,
		After: after,
	})
}

func CreateNewsletterFull(client *whatsmeow.Client, name, description string, picture []byte) (*types.NewsletterMetadata, error) {
	return client.CreateNewsletter(context.Background(), whatsmeow.CreateNewsletterParams{
		Name:        name,
		Description: description,
		Picture:     picture,
	})
}
