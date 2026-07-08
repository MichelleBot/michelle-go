package converter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"michelle/system/core"
	"michelle/system/serialize"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"iqc"},
		UsageHint: "text | time",
		Category:  "converter",
		Limit:     core.PerUserLimit(1, 24*time.Hour),
		Handler:   runIqc,
	})
}

func runIqc(ptz *core.Ptz) error {
	text := ptz.RawArgs
	if text == "" {
		return ptz.ReplyText("🚩 Contoh: " + ptz.Prefix + ptz.Command + " Hai | 20:20")
	}

	ptz.React("🕒")

	parts := strings.Split(text, "|")
	msg := strings.TrimSpace(parts[0])
	timeParam := ""
	if len(parts) > 1 {
		timeParam = strings.TrimSpace(parts[1])
	}

	var buffer []byte
	if ptz.Message.ExtendedTextMessage != nil && ptz.Message.ExtendedTextMessage.ContextInfo != nil && ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil {
		quotedMsg := ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage
		if quotedMsg.ImageMessage != nil || quotedMsg.StickerMessage != nil {
			var err error
			buffer, err = serialize.DownloadMedia(ptz.Client, quotedMsg)
			if err != nil {
				return ptz.ReplyText("❌ Gagal mendownload media: " + err.Error())
			}
		}
	}

	cdnURL := ""
	if len(buffer) > 0 {
		var err error
		cdnURL, err = uploadToCDN(buffer)
		if err != nil {
			return ptz.ReplyText("❌ Gagal upload media: " + err.Error())
		}
	}

	// Construct API URL
	apiURL := fmt.Sprintf("https://pham-michelle.vercel.app/api/iqc?text=%s&time=%s", msg, timeParam)
	if cdnURL != "" {
		apiURL += "&url=" + cdnURL
	}

	json, err := utils.FetchAsJSON(apiURL)
	if err != nil {
		return ptz.ReplyText("🚩 API Error: " + err.Error())
	}
	
	if json["status"] != true {
		return ptz.ReplyText("🚩 Gagal memproses: " + fmt.Sprintf("%v", json))
	}

	// Assuming data is in json["data"]["image"] as per original logic
	data := json["data"].(map[string]interface{})
	imgBase64 := data["image"].(string)
	
	// Handle Base64 decoding
	// Remove data URI prefix if present
	if strings.Contains(imgBase64, ",") {
		imgBase64 = strings.Split(imgBase64, ",")[1]
	}
	imgData, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return ptz.ReplyText("❌ Gagal memproses gambar hasil.")
	}

	return ptz.ReplyImage(imgData, "image/jpeg", "")
}

// uploadToCDN uploads buffer to https://cdn.crypty.workers.dev/
func uploadToCDN(data []byte) (string, error) {
	client := resty.New()
	
	filename := fmt.Sprintf("%d.jpg", time.Now().Unix())

	resp, err := client.R().
		SetFileReader("file", filename, bytes.NewReader(data)).
		SetResult(map[string]interface{}{}).
		Post("https://cdn.crypty.workers.dev/")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("server error: %s", resp.Status())
	}

	result := resp.Result().(*map[string]interface{})
	
	dataField, ok := (*result)["data"]
	if !ok {
		return "", fmt.Errorf("failed to find data field")
	}
	
	dataMap, ok := dataField.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse data as map")
	}
	
	url, ok := dataMap["url"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse URL field")
	}

	return url, nil
}
