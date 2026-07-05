package serialize

import (
	"context"
	"os"
	"os/exec"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func now() int64 { return time.Now().Unix() }

func buildQuotedContext(quoted *waE2E.Message, info types.MessageInfo) *waE2E.ContextInfo {
	ctx := &waE2E.ContextInfo{
		QuotedMessage: quoted,
	}
	if info.ID != "" {
		ctx.StanzaID = proto.String(info.ID)
	}
	if !info.Sender.IsEmpty() {
		ctx.Participant = proto.String(info.Sender.ToNonAD().String())
	}
	if !info.Chat.IsEmpty() {
		ctx.RemoteJID = proto.String(info.Chat.String())
	}
	return ctx
}

func SendText(client *whatsmeow.Client, to types.JID, text string) error {
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		Conversation: proto.String(text),
	})
	return err
}

func SendTextReply(client *whatsmeow.Client, to types.JID, text string, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(text),
			ContextInfo: buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	return err
}

func SendTextReplyID(client *whatsmeow.Client, to types.JID, text string, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) (types.MessageID, error) {
	resp, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(text),
			ContextInfo: buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func SendTextMention(client *whatsmeow.Client, to types.JID, text string, mentionedJIDs []types.JID) error {
	jidStrings := make([]string, len(mentionedJIDs))
	for i, jid := range mentionedJIDs {
		jidStrings[i] = jid.String()
	}
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: jidStrings,
			},
		},
	})
	return err
}

func SendTextReplyMention(client *whatsmeow.Client, to types.JID, text string, mentionedJIDs []types.JID, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
	ctx := buildQuotedContext(quotedMsg, quotedInfo)
	jidStrings := make([]string, len(mentionedJIDs))
	for i, jid := range mentionedJIDs {
		jidStrings[i] = jid.String()
	}
	ctx.MentionedJID = jidStrings
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(text),
			ContextInfo: ctx,
		},
	})
	return err
}

func SendTextMentionReplyToID(client *whatsmeow.Client, to types.JID, text string, mentionedJIDs []types.JID, quotedMsgID, quotedParticipant string) error {
	jidStrings := make([]string, len(mentionedJIDs))
	for i, jid := range mentionedJIDs {
		jidStrings[i] = jid.String()
	}
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:     proto.String(quotedMsgID),
				Participant:  proto.String(quotedParticipant),
				MentionedJID: jidStrings,
			},
		},
	})
	return err
}

func SendTextWithThumbnail(client *whatsmeow.Client, to types.JID, text, title, body, thumbnailURL string) error {
	var thumb []byte
	if thumbnailURL != "" {
		if r, err := httpClient.R().Get(thumbnailURL); err == nil {
			thumb = r.Body()
		}
	}
	ts := now()
	mediaType := waE2E.ContextInfo_ExternalAdReplyInfo_IMAGE
	renderLarge := true
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:              proto.String(text),
			MediaKeyTimestamp: &ts,
			ContextInfo: &waE2E.ContextInfo{
				ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
					Title:                 proto.String(title),
					Body:                  proto.String(body),
					MediaType:             &mediaType,
					ThumbnailURL:          proto.String(thumbnailURL),
					SourceURL:             proto.String(thumbnailURL),
					Thumbnail:             thumb,
					RenderLargerThumbnail: &renderLarge,
					ShowAdAttribution:     proto.Bool(false),
				},
			},
		},
	})
	return err
}

func SendReaction(client *whatsmeow.Client, chat types.JID, messageID types.MessageID, senderJID types.JID, emoji string) error {
	fromMe := client.Store.ID != nil && senderJID.User == client.Store.ID.User
	_, err := client.SendMessage(context.Background(), chat, &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(chat.String()),
				ID:        proto.String(messageID),
				FromMe:    proto.Bool(fromMe),
			},
			Text:              proto.String(emoji),
			SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
		},
	})
	return err
}

func RemoveReaction(client *whatsmeow.Client, chat types.JID, messageID types.MessageID, senderJID types.JID) error {
	fromMe := client.Store.ID != nil && senderJID.User == client.Store.ID.User
	_, err := client.SendMessage(context.Background(), chat, &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(chat.String()),
				ID:        proto.String(messageID),
				FromMe:    proto.Bool(fromMe),
			},
			Text:              proto.String(""),
			SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
		},
	})
	return err
}

func SendLocation(client *whatsmeow.Client, to types.JID, latitude, longitude float64, name, address string) error {
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(latitude),
			DegreesLongitude: proto.Float64(longitude),
			Name:             proto.String(name),
			Address:          proto.String(address),
		},
	})
	return err
}

func SendContact(client *whatsmeow.Client, to types.JID, phone, name string) error {
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: proto.String(name),
			Vcard: proto.String("BEGIN:VCARD\n" +
				"VERSION:3.0\n" +
				"FN:" + name + "\n" +
				"TEL;type=CELL;waid=" + phone + ":+" + phone + "\n" +
				"END:VCARD"),
		},
	})
	return err
}

func SendMultipleContacts(client *whatsmeow.Client, to types.JID, contacts []struct {
	Phone string
	Name  string
}) error {
	for _, c := range contacts {
		if err := SendContact(client, to, c.Phone, c.Name); err != nil {
			return err
		}
	}
	return nil
}

func SendImage(client *whatsmeow.Client, to types.JID, imageData []byte, mime, caption string) error {
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
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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

func SendImageReply(client *whatsmeow.Client, to types.JID, imageData []byte, mime, caption string, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
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
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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
			ContextInfo:       buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	return err
}

func SendImageReplyID(client *whatsmeow.Client, to types.JID, imageData []byte, mime, caption string, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) (types.MessageID, error) {
	resp, err := client.Upload(context.Background(), imageData, whatsmeow.MediaImage)
	if err != nil {
		return "", err
	}
	ext := mimeToExt(mime)
	thumb, _ := GenerateJPEGThumbnail(imageData, ext)
	dim, _ := GetMediaDimensions(imageData, ext)
	h := dim.Height
	w := dim.Width
	ts := now()
	sendResp, err := client.SendMessage(context.Background(), to, &waE2E.Message{
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
			ContextInfo:       buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	if err != nil {
		return "", err
	}
	return sendResp.ID, nil
}

func SendImageFile(client *whatsmeow.Client, to types.JID, path, mime, caption string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return SendImage(client, to, data, mime, caption)
}

func SendImageWithURL(client *whatsmeow.Client, to types.JID, url, caption string) error {
	data, err := Fetch(url)
	if err != nil {
		return err
	}
	mime := detectMIME(data)
	return SendImage(client, to, data, mime, caption)
}

func SendVideo(client *whatsmeow.Client, to types.JID, videoData []byte, mime, caption string) error {
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
	if h == 0 {
		h = 720
	}
	if w == 0 {
		w = 1280
	}
	ts := now()
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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

func SendVideoReply(client *whatsmeow.Client, to types.JID, videoData []byte, mime, caption string, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
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
	if h == 0 {
		h = 720
	}
	if w == 0 {
		w = 1280
	}
	ts := now()
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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
			ContextInfo:       buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	return err
}

func SendVideoFile(client *whatsmeow.Client, to types.JID, path, mime, caption string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return SendVideo(client, to, data, mime, caption)
}

func SendAudio(client *whatsmeow.Client, to types.JID, audioData []byte, mime string, ptt bool) error {
	resp, err := client.Upload(context.Background(), audioData, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}
	ts := now()
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			Mimetype:          proto.String(mime),
			PTT:               proto.Bool(ptt),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			MediaKey:          resp.MediaKey,
			FileEncSHA256:     resp.FileEncSHA256,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
		},
	})
	return err
}

func SendAudioReply(client *whatsmeow.Client, to types.JID, audioData []byte, mime string, ptt bool, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
	resp, err := client.Upload(context.Background(), audioData, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}
	ts := now()
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			Mimetype:          proto.String(mime),
			PTT:               proto.Bool(ptt),
			URL:               &resp.URL,
			DirectPath:        &resp.DirectPath,
			MediaKey:          resp.MediaKey,
			FileEncSHA256:     resp.FileEncSHA256,
			FileSHA256:        resp.FileSHA256,
			FileLength:        &resp.FileLength,
			MediaKeyTimestamp: &ts,
			ContextInfo:       buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	return err
}

func SendAudioFile(client *whatsmeow.Client, to types.JID, path, mime string, ptt bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return SendAudio(client, to, data, mime, ptt)
}

func SendDocument(client *whatsmeow.Client, to types.JID, docData []byte, mime, filename, caption string) error {
	resp, err := client.Upload(context.Background(), docData, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}
	ts := now()
	thumb, _ := GenerateJPEGThumbnail(docData, mimeToExt(mime))
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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

func SendDocumentReply(client *whatsmeow.Client, to types.JID, docData []byte, mime, filename, caption string, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
	resp, err := client.Upload(context.Background(), docData, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}
	ts := now()
	thumb, _ := GenerateJPEGThumbnail(docData, mimeToExt(mime))
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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
			ContextInfo:       buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	return err
}

func SendDocumentFile(client *whatsmeow.Client, to types.JID, path, filename, caption string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	mime := detectMIME(data)
	return SendDocument(client, to, data, mime, filename, caption)
}

func SendDocumentWithURL(client *whatsmeow.Client, to types.JID, url, filename, caption string) error {
	data, err := Fetch(url)
	if err != nil {
		return err
	}
	mime := detectMIME(data)
	return SendDocument(client, to, data, mime, filename, caption)
}

func generateStickerSidecar(data []byte) (firstFrameData []byte, length uint32) {
	tmpIn := tmpFile("stk_sc_", ".webp")
	tmpOut := tmpFile("stk_sc_", ".jpg")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	_ = os.WriteFile(tmpIn, data, 0600)
	cmd := exec.Command("ffmpeg", "-y", "-i", tmpIn, "-vframes", "1", "-vf",
		"scale=80:80:force_original_aspect_ratio=decrease,pad=80:80:(ow-iw)/2:(oh-ih)/2",
		"-f", "image2", "-vcodec", "mjpeg", "-q:v", "10", tmpOut)
	_ = cmd.Run()
	frame, err := os.ReadFile(tmpOut)
	if err != nil || len(frame) == 0 {
		return nil, 0
	}
	return frame, uint32(len(frame))
}

func buildStickerMsg(resp whatsmeow.UploadResponse, stickerData []byte, mime string, animated bool) *waE2E.StickerMessage {
	thumb, _ := GeneratePNGThumbnail(stickerData, ".webp")

	var h, w uint32 = 512, 512
	if !animated {
		if dim, err := GetMediaDimensions(stickerData, ".webp"); err == nil && dim.Width > 0 {
			h = dim.Height
			w = dim.Width
		}
	}

	ts := now()

	msg := &waE2E.StickerMessage{
		Mimetype:          proto.String(mime),
		IsAnimated:        proto.Bool(animated),
		URL:               &resp.URL,
		DirectPath:        &resp.DirectPath,
		MediaKey:          resp.MediaKey,
		FileEncSHA256:     resp.FileEncSHA256,
		FileSHA256:        resp.FileSHA256,
		FileLength:        &resp.FileLength,
		MediaKeyTimestamp: &ts,
		PngThumbnail:      thumb,
		Height:            &h,
		Width:             &w,
	}

	if animated {
		frame, length := generateStickerSidecar(stickerData)
		if length > 0 {
			msg.FirstFrameSidecar = frame
			msg.FirstFrameLength = proto.Uint32(length)
		}
	}

	return msg
}

func SendSticker(client *whatsmeow.Client, to types.JID, stickerData []byte, mime string, animated bool) error {
	resp, err := client.Upload(context.Background(), stickerData, whatsmeow.MediaImage)
	if err != nil {
		return err
	}
	stickerMsg := buildStickerMsg(resp, stickerData, mime, animated)
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
		StickerMessage: stickerMsg,
	})
	return err
}

func SendStickerReply(client *whatsmeow.Client, to types.JID, stickerData []byte, mime string, animated bool, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
	resp, err := client.Upload(context.Background(), stickerData, whatsmeow.MediaImage)
	if err != nil {
		return err
	}
	stickerMsg := buildStickerMsg(resp, stickerData, mime, animated)
	stickerMsg.ContextInfo = buildQuotedContext(quotedMsg, quotedInfo)
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
		StickerMessage: stickerMsg,
	})
	return err
}

func SendStickerFile(client *whatsmeow.Client, to types.JID, path string, animated bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	mime := detectMIME(data)
	return SendSticker(client, to, data, mime, animated)
}

func SendTyping(client *whatsmeow.Client, chat types.JID, composing bool) error {
	state := types.ChatPresenceComposing
	if !composing {
		state = types.ChatPresencePaused
	}
	return client.SendChatPresence(context.Background(), chat, state, types.ChatPresenceMediaText)
}

func SendRecording(client *whatsmeow.Client, chat types.JID) error {
	return client.SendChatPresence(context.Background(), chat, types.ChatPresenceComposing, types.ChatPresenceMediaAudio)
}

func SendEdit(client *whatsmeow.Client, chat types.JID, originalMsgID types.MessageID, newMsg *waE2E.Message) error {
	_, err := client.SendMessage(context.Background(), chat, client.BuildEdit(chat, originalMsgID, newMsg))
	return err
}

func SendRevoke(client *whatsmeow.Client, chat types.JID, senderJID types.JID, msgID types.MessageID) error {
	_, err := client.SendMessage(context.Background(), chat, client.BuildRevoke(chat, senderJID, msgID))
	return err
}

func BuildMessageKey(client *whatsmeow.Client, chat, sender types.JID, id types.MessageID) *waCommon.MessageKey {
	return client.BuildMessageKey(chat, sender, id)
}

func RevokeMyMessage(client *whatsmeow.Client, chat types.JID, msgID types.MessageID) error {
	_, err := client.SendMessage(context.Background(), chat, client.BuildRevoke(chat, types.EmptyJID, msgID))
	return err
}

func RevokeOtherMessage(client *whatsmeow.Client, chat types.JID, sender types.JID, msgID types.MessageID) error {
	_, err := client.SendMessage(context.Background(), chat, client.BuildRevoke(chat, sender, msgID))
	return err
}

func EditMessage(client *whatsmeow.Client, chat types.JID, originalMsgID types.MessageID, newText string) error {
	newMsg := &waE2E.Message{
		Conversation: proto.String(newText),
	}
	_, err := client.SendMessage(context.Background(), chat, client.BuildEdit(chat, originalMsgID, newMsg))
	return err
}

func EditImageCaption(client *whatsmeow.Client, chat types.JID, originalMsgID types.MessageID, newCaption string) error {
	newMsg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption: proto.String(newCaption),
		},
	}
	_, err := client.SendMessage(context.Background(), chat, client.BuildEdit(chat, originalMsgID, newMsg))
	return err
}

func BuildReaction(client *whatsmeow.Client, chat, sender types.JID, msgID types.MessageID, emoji string) *waE2E.Message {
	return client.BuildReaction(chat, sender, msgID, emoji)
}

func SendMessageWithID(client *whatsmeow.Client, to types.JID, msg *waE2E.Message, customID types.MessageID) error {
	_, err := client.SendMessage(context.Background(), to, msg, whatsmeow.SendRequestExtra{ID: customID})
	return err
}

func SetDisappearingTimerDuration(client *whatsmeow.Client, chat types.JID, dur time.Duration) error {
	return client.SetDisappearingTimer(context.Background(), chat, dur, time.Time{})
}

func SetDisappearingOff(client *whatsmeow.Client, chat types.JID) error {
	return client.SetDisappearingTimer(context.Background(), chat, whatsmeow.DisappearingTimerOff, time.Time{})
}

func SetDisappearing24h(client *whatsmeow.Client, chat types.JID) error {
	return client.SetDisappearingTimer(context.Background(), chat, whatsmeow.DisappearingTimer24Hours, time.Time{})
}

func SetDisappearing7d(client *whatsmeow.Client, chat types.JID) error {
	return client.SetDisappearingTimer(context.Background(), chat, whatsmeow.DisappearingTimer7Days, time.Time{})
}

func SetDisappearing90d(client *whatsmeow.Client, chat types.JID) error {
	return client.SetDisappearingTimer(context.Background(), chat, whatsmeow.DisappearingTimer90Days, time.Time{})
}

func GenerateMessageID(client *whatsmeow.Client) types.MessageID {
	return client.GenerateMessageID()
}

func SendViewOnceImage(client *whatsmeow.Client, to types.JID, imageData []byte, mime, caption string) error {
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
	viewOnce := &waE2E.Message{
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
			ViewOnce:          proto.Bool(true),
		},
	}
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{Message: viewOnce},
	})
	return err
}

func SendViewOnceVideo(client *whatsmeow.Client, to types.JID, videoData []byte, mime, caption string) error {
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
	viewOnce := &waE2E.Message{
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
			ViewOnce:          proto.Bool(true),
		},
	}
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{Message: viewOnce},
	})
	return err
}

func SendGIF(client *whatsmeow.Client, to types.JID, videoData []byte, mime, caption string) error {
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
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			Caption:           proto.String(caption),
			Mimetype:          proto.String(mime),
			GifPlayback:       proto.Bool(true),
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

func SendVoiceNote(client *whatsmeow.Client, to types.JID, audioData []byte, mime string) error {
	resp, err := client.Upload(context.Background(), audioData, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}
	secs, _ := GetVideoDurationSeconds(audioData, mimeToExt(mime))
	ts := now()
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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

func SendVoiceNoteReply(client *whatsmeow.Client, to types.JID, audioData []byte, mime string, quotedMsg *waE2E.Message, quotedInfo types.MessageInfo) error {
	resp, err := client.Upload(context.Background(), audioData, whatsmeow.MediaAudio)
	if err != nil {
		return err
	}
	secs, _ := GetVideoDurationSeconds(audioData, mimeToExt(mime))
	ts := now()
	_, err = client.SendMessage(context.Background(), to, &waE2E.Message{
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
			ContextInfo:       buildQuotedContext(quotedMsg, quotedInfo),
		},
	})
	return err
}

func SendLiveLocation(client *whatsmeow.Client, to types.JID, lat, lon float64, caption string, sequence int64) error {
	_, err := client.SendMessage(context.Background(), to, &waE2E.Message{
		LiveLocationMessage: &waE2E.LiveLocationMessage{
			DegreesLatitude:  proto.Float64(lat),
			DegreesLongitude: proto.Float64(lon),
			Caption:          proto.String(caption),
			SequenceNumber:   proto.Int64(sequence),
		},
	})
	return err
}

func BuildGroupInviteMsg(groupJID, inviteCode string, expiration int64, groupName, caption string, thumbnail []byte) *waE2E.Message {
	return &waE2E.Message{
		GroupInviteMessage: &waE2E.GroupInviteMessage{
			GroupJID:         proto.String(groupJID),
			InviteCode:       proto.String(inviteCode),
			InviteExpiration: proto.Int64(expiration),
			GroupName:        proto.String(groupName),
			Caption:          proto.String(caption),
			JPEGThumbnail:    thumbnail,
		},
	}
}

func ForwardMessage(client *whatsmeow.Client, to types.JID, msg *waE2E.Message, forwardingScore uint32) error {
	isForwarded := true
	score := forwardingScore
	var ci *waE2E.ContextInfo
	switch {
	case msg.ExtendedTextMessage != nil:
		ci = msg.ExtendedTextMessage.ContextInfo
	case msg.ImageMessage != nil:
		ci = msg.ImageMessage.ContextInfo
	case msg.VideoMessage != nil:
		ci = msg.VideoMessage.ContextInfo
	case msg.AudioMessage != nil:
		ci = msg.AudioMessage.ContextInfo
	case msg.DocumentMessage != nil:
		ci = msg.DocumentMessage.ContextInfo
	}
	if ci == nil {
		ci = &waE2E.ContextInfo{}
	}
	ci.IsForwarded = &isForwarded
	ci.ForwardingScore = &score
	switch {
	case msg.ExtendedTextMessage != nil:
		msg.ExtendedTextMessage.ContextInfo = ci
	case msg.ImageMessage != nil:
		msg.ImageMessage.ContextInfo = ci
	case msg.VideoMessage != nil:
		msg.VideoMessage.ContextInfo = ci
	case msg.AudioMessage != nil:
		msg.AudioMessage.ContextInfo = ci
	case msg.DocumentMessage != nil:
		msg.DocumentMessage.ContextInfo = ci
	}
	_, err := client.SendMessage(context.Background(), to, msg)
	return err
}
