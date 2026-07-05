package utils

import (
	"encoding/json"
	"fmt"
	"mime"
	"strings"

	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client/profiles"

	"sptzx/src/core"
	"sptzx/src/serialize"
)

func GetBrowserProfile(browser string) profiles.ClientProfile {
	switch browser {
	case "firefox":
		return profiles.Firefox_147
	case "safari":
		return profiles.Chrome_131
	case "edge":
		return profiles.Chrome_131
	case "opera":
		return profiles.Chrome_131
	case "ios":
		return profiles.Chrome_131
	default:
		return profiles.Chrome_131
	}
}

func GetDefaultHeaders(browser string) fhttp.Header {
	headers := fhttp.Header{}
	headers.Set("accept", "*/*")
	headers.Set("accept-language", "en-US,en;q=0.9")
	headers.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120"`)
	headers.Set("sec-ch-ua-mobile", "?0")
	headers.Set("sec-ch-ua-platform", `"Linux"`)
	headers.Set("sec-fetch-dest", "empty")
	headers.Set("sec-fetch-mode", "cors")
	headers.Set("sec-fetch-site", "none")

	switch browser {
	case "firefox":
		headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0")
	case "safari":
		headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15")
	case "edge":
		headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
	case "opera":
		headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 OPR/117.0.0.0")
	case "ios":
		headers.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Mobile/15E148 Safari/604.1")
	default:
		headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	}

	return headers
}

func SendAsImage(ptz *core.Ptz, data []byte, mimeType string) error {
	return serialize.SendImageReply(ptz.Bot.Client, ptz.Chat, data, mimeType, "", ptz.Message, ptz.Info)
}

func SendAsVideo(ptz *core.Ptz, data []byte, mimeType string) error {
	return serialize.SendVideoReply(ptz.Bot.Client, ptz.Chat, data, mimeType, "", ptz.Message, ptz.Info)
}

func SendAsAudio(ptz *core.Ptz, data []byte, mimeType string) error {
	return serialize.SendAudioReply(ptz.Bot.Client, ptz.Chat, data, mimeType, false, ptz.Message, ptz.Info)
}

func SendAsDocument(ptz *core.Ptz, data []byte, mimeType, urlPath string) error {
	filename := DetectFilename(mimeType, urlPath)
	return serialize.SendDocumentReply(ptz.Bot.Client, ptz.Chat, data, mimeType, filename, "", ptz.Message, ptz.Info)
}

func SendAsFormattedJSON(ptz *core.Ptz, data []byte) error {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return SendAsText(ptz, data, "application/json")
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return SendAsText(ptz, data, "application/json")
	}

	if len(prettyJSON) > 4000 {
		return SendAsDocument(ptz, prettyJSON, "application/json", "response.json")
	}

	message := "```json\n" + string(prettyJSON) + "\n```"
	return ptz.ReplyText(message)
}

func SendAsText(ptz *core.Ptz, data []byte, contentType string) error {
	text := string(data)

	if len(text) > 4000 {
		filename := DetectFilename(contentType, "")
		return SendAsDocument(ptz, data, contentType, filename)
	}

	if len(text) > 65000 {
		text = text[:65000] + "\n\n... (truncated)"
	}

	return ptz.ReplyText(text)
}

func GetExtensionFromType(contentType string) string {
	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extensions) == 0 {
		return "bin"
	}
	ext := extensions[0]
	if len(ext) > 0 && ext[0] == '.' {
		return ext[1:]
	}
	return ext
}

func GetFilenameFromURL(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "download"
}

func DetectFilename(contentType, urlPath string) string {
	ext := GetExtensionFromType(contentType)

	if ext == "bin" {
		baseName := GetFilenameFromURL(urlPath)
		if baseName != "" && baseName != "download" {
			return baseName
		}
		return "download.bin"
	}

	return fmt.Sprintf("download.%s", ext)
}
