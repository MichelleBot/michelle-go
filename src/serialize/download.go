package serialize

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func DownloadMedia(client *whatsmeow.Client, msg *waE2E.Message) ([]byte, error) {
	switch {
	case msg.ImageMessage != nil:
		return client.Download(context.Background(), msg.ImageMessage)
	case msg.VideoMessage != nil:
		return client.Download(context.Background(), msg.VideoMessage)
	case msg.AudioMessage != nil:
		return client.Download(context.Background(), msg.AudioMessage)
	case msg.DocumentMessage != nil:
		return client.Download(context.Background(), msg.DocumentMessage)
	case msg.StickerMessage != nil:
		return client.Download(context.Background(), msg.StickerMessage)
	default:
		return nil, whatsmeow.ErrNothingDownloadableFound
	}
}

func DownloadMediaToFile(client *whatsmeow.Client, msg *waE2E.Message, path string) error {
	data, err := DownloadMedia(client, msg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func DownloadAndSaveMedia(client *whatsmeow.Client, msg *waE2E.Message, dir string) (string, error) {
	data, err := DownloadMedia(client, msg)
	if err != nil {
		return "", err
	}

	var filename string

	switch {
	case msg.ImageMessage != nil:
		filename = "image_" + getFileHash(msg.ImageMessage.FileSHA256) + ".jpg"
	case msg.VideoMessage != nil:
		filename = "video_" + getFileHash(msg.VideoMessage.FileSHA256) + ".mp4"
	case msg.AudioMessage != nil:
		if msg.AudioMessage.GetPTT() {
			filename = "voice_" + getFileHash(msg.AudioMessage.FileSHA256) + ".ogg"
		} else {
			filename = "audio_" + getFileHash(msg.AudioMessage.FileSHA256) + ".mp3"
		}
	case msg.DocumentMessage != nil:
		filename = msg.DocumentMessage.GetFileName()
	case msg.StickerMessage != nil:
		filename = "sticker_" + getFileHash(msg.StickerMessage.FileSHA256) + ".webp"
	}

	if filename == "" {
		filename = "media_unknown"
	}

	fullPath := dir + "/" + filename
	return fullPath, os.WriteFile(fullPath, data, 0644)
}

func getFileHash(hash []byte) string {
	if len(hash) == 0 {
		return "unknown"
	}
	return fmt.Sprintf("%x", hash[:8])
}

func DownloadProfilePicture(client *whatsmeow.Client, jid types.JID) ([]byte, error) {
	info, err := client.GetProfilePictureInfo(context.Background(), jid, nil)
	if err != nil {
		return nil, err
	}
	if info == nil || info.URL == "" {
		return nil, whatsmeow.ErrProfilePictureNotSet
	}

	resp, err := httpClient.R().Get(info.URL)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode())
	}
	return resp.Body(), nil
}

func DownloadSticker(client *whatsmeow.Client, msg *waE2E.Message) ([]byte, error) {
	if msg.StickerMessage == nil {
		return nil, whatsmeow.ErrNothingDownloadableFound
	}
	return client.Download(context.Background(), msg.StickerMessage)
}

func GetMediaURL(msg *waE2E.Message) string {
	switch {
	case msg.ImageMessage != nil:
		return msg.ImageMessage.GetURL()
	case msg.VideoMessage != nil:
		return msg.VideoMessage.GetURL()
	case msg.AudioMessage != nil:
		return msg.AudioMessage.GetURL()
	case msg.DocumentMessage != nil:
		return msg.DocumentMessage.GetURL()
	case msg.StickerMessage != nil:
		return msg.StickerMessage.GetURL()
	default:
		return ""
	}
}

func GetMediaCaption(msg *waE2E.Message) string {
	switch {
	case msg.ImageMessage != nil:
		return msg.ImageMessage.GetCaption()
	case msg.VideoMessage != nil:
		return msg.VideoMessage.GetCaption()
	case msg.DocumentMessage != nil:
		return msg.DocumentMessage.GetCaption()
	default:
		return ""
	}
}

func GetMediaFilename(msg *waE2E.Message) string {
	switch {
	case msg.ImageMessage != nil:
		return "image.jpg"
	case msg.VideoMessage != nil:
		return "video.mp4"
	case msg.AudioMessage != nil:
		if msg.AudioMessage.GetPTT() {
			return "voice.ogg"
		}
		return "audio.mp3"
	case msg.DocumentMessage != nil:
		return msg.DocumentMessage.GetFileName()
	case msg.StickerMessage != nil:
		return "sticker.webp"
	default:
		return ""
	}
}

func GetMediaMIME(msg *waE2E.Message) string {
	switch {
	case msg.ImageMessage != nil:
		return msg.ImageMessage.GetMimetype()
	case msg.VideoMessage != nil:
		return msg.VideoMessage.GetMimetype()
	case msg.AudioMessage != nil:
		return msg.AudioMessage.GetMimetype()
	case msg.DocumentMessage != nil:
		return msg.DocumentMessage.GetMimetype()
	case msg.StickerMessage != nil:
		return msg.StickerMessage.GetMimetype()
	default:
		return ""
	}
}
