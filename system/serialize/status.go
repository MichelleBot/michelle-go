package serialize

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func GetStatusPrivacy(client *whatsmeow.Client) ([]types.StatusPrivacy, error) {
	return client.GetStatusPrivacy(context.Background())
}

func SendStatusText(client *whatsmeow.Client, text string) error {
	_, err := client.SendMessage(context.Background(), types.StatusBroadcastJID, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
		},
	})
	return err
}

func SendStatusImage(client *whatsmeow.Client, imageData []byte, mime, caption string) error {
	resp, err := client.Upload(context.Background(), imageData, whatsmeow.MediaImage)
	if err != nil {
		return err
	}
	ext := mimeToExt(mime)
	thumb, _ := GenerateJPEGThumbnail(imageData, ext)
	dim, _ := GetMediaDimensions(imageData, ext)
	h := dim.Height
	w := dim.Width
	ts := now()
	_, err = client.SendMessage(context.Background(), types.StatusBroadcastJID, &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:           proto.String(caption),
			Mimetype:          proto.String(mime),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			MediaKey:          resp.MediaKey,
			FileEncSHA256:     resp.FileEncSHA256,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
			JPEGThumbnail:     thumb,
			Height:            &h,
			Width:             &w,
		},
	})
	return err
}

func SendStatusVideo(client *whatsmeow.Client, videoData []byte, mime, caption string) error {
	resp, err := client.Upload(context.Background(), videoData, whatsmeow.MediaVideo)
	if err != nil {
		return err
	}
	ext := mimeToExt(mime)
	thumb, _ := GenerateJPEGThumbnail(videoData, ext)
	dim, _ := GetMediaDimensions(videoData, ext)
	secs, _ := GetVideoDurationSeconds(videoData, ext)
	h := dim.Height
	w := dim.Width
	ts := now()
	_, err = client.SendMessage(context.Background(), types.StatusBroadcastJID, &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			Caption:           proto.String(caption),
			Mimetype:          proto.String(mime),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			MediaKey:          resp.MediaKey,
			FileEncSHA256:     resp.FileEncSHA256,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
			JPEGThumbnail:     thumb,
			Height:            &h,
			Width:             &w,
			Seconds:           &secs,
		},
	})
	return err
}

func SendStatusVoiceNote(client *whatsmeow.Client, audioData []byte, mime string) error {
	oggData, err := ToOggOpus(audioData, mimeToExt(mime))
	if err != nil {
		oggData = audioData
	}
	resp, err := client.Upload(context.Background(), oggData, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}
	secs, _ := GetVideoDurationSeconds(oggData, ".ogg")
	ts := now()
	_, err = client.SendMessage(context.Background(), types.StatusBroadcastJID, &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			Mimetype:          proto.String("audio/ogg; codecs=opus"),
			PTT:               proto.Bool(true),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			MediaKey:          resp.MediaKey,
			FileEncSHA256:     resp.FileEncSHA256,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			Seconds:           &secs,
			MediaKeyTimestamp: &ts,
		},
	})
	return err
}

func SendStatusDoc(client *whatsmeow.Client, docData []byte, mime, filename, caption string) error {
	resp, err := client.Upload(context.Background(), docData, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}
	ts := now()
	thumb, _ := GenerateJPEGThumbnail(docData, mimeToExt(mime))
	_, err = client.SendMessage(context.Background(), types.StatusBroadcastJID, &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			Title:             proto.String(filename),
			FileName:          proto.String(filename),
			Caption:           proto.String(caption),
			Mimetype:          proto.String(mime),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			MediaKey:          resp.MediaKey,
			FileEncSHA256:     resp.FileEncSHA256,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
			JPEGThumbnail:     thumb,
		},
	})
	return err
}

func SendStatusTextWithExpiry(client *whatsmeow.Client, text string, expirySeconds uint32) error {
	ts := now()
	_, err := client.SendMessage(context.Background(), types.StatusBroadcastJID, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:              proto.String(text),
			MediaKeyTimestamp: &ts,
			ContextInfo: &waE2E.ContextInfo{
				Expiration: &expirySeconds,
			},
		},
	})
	return err
}
