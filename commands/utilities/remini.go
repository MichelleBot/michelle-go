package utilities

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"michelle/system/core"
	"michelle/system/serialize"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"remini"},
		Hidden:    []string{"hd"},
		UsageHint: "reply photo",
		Category:  "utilities",
		Quota:     core.PerUserQuota(1),
		Handler:   runRemini,
	})
}

func runRemini(ptz *core.Ptz) error {
	var quotedMsg *waE2E.Message

	// Handle quoted message or viewOnce
	if ptz.Message.ExtendedTextMessage != nil && ptz.Message.ExtendedTextMessage.ContextInfo != nil {
		quotedMsg = ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage
	} else if ptz.Message.ViewOnceMessage != nil {
		quotedMsg = ptz.Message.ViewOnceMessage.Message
	}

	if quotedMsg == nil {
		// Fallback check if user just replied to a message with a photo
		quotedMsg = ptz.Message
	}

	// Validate if we have an image message
	if quotedMsg.ImageMessage == nil {
		return ptz.ReplyText(utils.Texted("bold", "🚩 Balas gambarnya."))
	}

	ptz.React("🕒")
	ptz.Bot.Log.Infof("Remini process started for %s", ptz.Sender.User)

	data, err := serialize.DownloadMedia(ptz.Client, quotedMsg)
	if err != nil {
		ptz.Bot.Log.Errorf("Download failed: %v", err)
		return ptz.ReplyText("❌ Gagal mendownload media: " + err.Error())
	}

	// Upload to CDN
	cdnURL, err := uploadToCDN(data)
	if err != nil {
		ptz.Bot.Log.Errorf("CDN upload failed: %v", err)
		return ptz.ReplyText("❌ Gagal upload media ke CDN: " + err.Error())
	}

	// Call API
	apiURL := fmt.Sprintf("https://api.lexcode.biz.id/api/tools/wink-hd?url=%s", cdnURL)
	json, err := utils.FetchAsJSON(apiURL)
	if err != nil {
		ptz.Bot.Log.Errorf("Remini API fetch error: %v", err)
		return ptz.ReplyText("🚩 Gagal memproses gambar (API Error).")
	}
	if json["success"] != true {
		ptz.Bot.Log.Errorf("Remini API error: %v", json)
		return ptz.ReplyText("🚩 Gagal memproses gambar (API Failure).")
	}

	result := json["result"].(map[string]interface{})
	imgURL := result["image"].(string)

	// Download result image
	imgData, err := utils.FetchAsBuffer(imgURL)
	if err != nil {
		ptz.Bot.Log.Errorf("Result download failed: %v", err)
		return ptz.ReplyText("❌ Gagal mendownload hasil gambar.")
	}

	// Success
	ptz.Unreact()
	
	// Reply with empty caption
	return ptz.ReplyImage(imgData, "image/jpeg", "")
}

// uploadToCDN uploads buffer to https://cdn.crypty.workers.dev/
func uploadToCDN(data []byte) (string, error) {
	client := resty.New()
	
	// Prepare filename
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
	
	// Safely access nested data
	dataField, ok := (*result)["data"]
	if !ok {
		return "", fmt.Errorf("failed to find data field. Raw response: %s", resp.String())
	}
	
	dataMap, ok := dataField.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse data as map. Raw response: %s", resp.String())
	}
	
	url, ok := dataMap["url"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse URL field. Raw response: %s", resp.String())
	}

	return url, nil
}
