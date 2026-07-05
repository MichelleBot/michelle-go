package serialize

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

var httpClient = resty.New()

func Fetch(rawURL string) ([]byte, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		parsed, _ = url.Parse("https://" + rawURL)
	}
	resp, err := httpClient.R().Get(parsed.String())
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode())
	}
	return resp.Body(), nil
}

func FetchWithUA(rawURL, userAgent string) ([]byte, error) {
	resp, err := httpClient.R().
		SetHeader("User-Agent", userAgent).
		Get(rawURL)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode())
	}
	return resp.Body(), nil
}

func NumFmt(n float64) string {
	if n >= 1_000_000_000 {
		return fmt.Sprintf("%.1fB", n/1_000_000_000)
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", n/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fK", n/1_000)
	}
	return fmt.Sprintf("%.0f", n)
}

func NumFmt64(n int64) string {
	if n >= 1_000_000 {
		return strconv.FormatInt(n/1_000_000, 10) + "M"
	}
	if n >= 1_000 {
		return strconv.FormatInt(n/1_000, 10) + "K"
	}
	return strconv.FormatInt(n, 10)
}

func detectMIME(data []byte) string {
	mime := http.DetectContentType(data)
	if strings.HasPrefix(mime, "text/plain") {
		if len(data) > 0 && data[0] == 0x89 && string(data[1:4]) == "PNG" {
			return "image/png"
		}
	}
	return mime
}

func mimeToExt(mime string) string {
	switch mime {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/3gpp":
		return ".3gp"
	case "video/quicktime":
		return ".mov"
	case "video/webm":
		return ".webm"
	case "audio/ogg", "audio/ogg; codecs=opus":
		return ".ogg"
	case "audio/mpeg", "audio/mp3":
		return ".mp3"
	case "audio/mp4", "audio/aac":
		return ".m4a"
	default:
		return ".bin"
	}
}

func GetMediaExtFromMIME(mime string) string {
	return mimeToExt(mime)
}

func ClassifyMessage(msg *waE2E.Message) (msgType, content string) {
	if msg == nil {
		return "nil", ""
	}
	switch {
	case msg.Conversation != nil:
		c := *msg.Conversation
		if len(c) > 60 {
			c = c[:60] + "..."
		}
		return "text", c
	case msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil:
		c := *msg.ExtendedTextMessage.Text
		if len(c) > 60 {
			c = c[:60] + "..."
		}
		return "text", c
	case msg.ImageMessage != nil:
		caption := ""
		if msg.ImageMessage.Caption != nil {
			caption = *msg.ImageMessage.Caption
		}
		return "image", caption
	case msg.VideoMessage != nil:
		caption := ""
		if msg.VideoMessage.Caption != nil {
			caption = *msg.VideoMessage.Caption
		}
		return "video", caption
	case msg.AudioMessage != nil:
		if msg.AudioMessage.PTT != nil && *msg.AudioMessage.PTT {
			return "voice", ""
		}
		return "audio", ""
	case msg.DocumentMessage != nil:
		name := ""
		if msg.DocumentMessage.FileName != nil {
			name = *msg.DocumentMessage.FileName
		}
		return "document", name
	case msg.StickerMessage != nil:
		return "sticker", ""
	case msg.ContactMessage != nil:
		name := ""
		if msg.ContactMessage.DisplayName != nil {
			name = *msg.ContactMessage.DisplayName
		}
		return "contact", name
	case msg.LocationMessage != nil:
		return "location", ""
	case msg.ReactionMessage != nil:
		emoji := ""
		if msg.ReactionMessage.Text != nil {
			emoji = *msg.ReactionMessage.Text
		}
		return "reaction", emoji
	case msg.EditedMessage != nil:
		return "edit", ""
	case msg.ProtocolMessage != nil:
		return "protocol", ""
	case msg.EventInviteMessage != nil:
		return "event_invite", ""
	default:
		return "unknown", ""
	}
}

func RandomString(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[rand.Intn(len(s))]
}

func RemoveFiles(files ...string) {
	for _, f := range files {
		if f != "" {
			_ = os.Remove(f)
		}
	}
}
