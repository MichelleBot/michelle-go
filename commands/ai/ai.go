package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"ai", "deepseek"},
		UsageHint: "query",
		Category:  "ai",
		Quota:     core.PerUserQuota(1),
		Handler:   handleAI,
	})
}

type APIResponse struct {
	Status  bool   `json:"status"`
	Result  string `json:"result"`
	Message string `json:"message"`
	Reply   string `json:"reply"`
	Data    struct {
		Reply string `json:"reply"`
	} `json:"data"`
}

func handleAI(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText(fmt.Sprintf("🚩 Contoh penggunaan: %s%s apa itu nodejs", ptz.Prefix, ptz.Command))
	}

	serialize.SendReaction(ptz.Bot.Client, ptz.Chat, ptz.Info.ID, ptz.Sender, "🕒")
	text := strings.Join(ptz.Args, " ")
	var result string
	var err error

	switch ptz.Command {
	case "ai":
		result, err = fetchAI("https://api.alwayscodex.my.id/api/ai/gpt5?teks=" + url.QueryEscape(text))
	case "deepseek":
		result, err = fetchDeepseek("https://www.neoapis.xyz/api/ai/deepseek?text=" + url.QueryEscape(text))
	default:
		return nil
	}

	if err != nil {
		return ptz.ReplyText("🚩 Terjadi kesalahan: " + err.Error())
	}

	return ptz.ReplyText(strings.ReplaceAll(result, "**", "*"))
}

func fetchAI(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if !res.Status {
		return "", fmt.Errorf("Tidak ada respon.")
	}
    
    if res.Result != "" { return res.Result, nil }
    if res.Reply != "" { return res.Reply, nil }
    if res.Message != "" { return res.Message, nil }
    if res.Data.Reply != "" { return res.Data.Reply, nil }
    
	return "", fmt.Errorf("Respon kosong.")
}

func fetchDeepseek(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if !res.Status {
		return "", fmt.Errorf("Tidak ada respon.")
	}
	return res.Data.Reply, nil
}
