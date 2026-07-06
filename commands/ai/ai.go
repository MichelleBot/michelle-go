package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"michelle/system/core"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"ai", "blackbox", "deepseek"},
		UsageHint: "query",
		Category:  "ai",
		Quota:     core.PerUserQuota(1),
		Handler:   handleAI,
	})
}

type APIResponse struct {
	Status bool   `json:"status"`
	Result string `json:"result"`
	Data   struct {
		Reply string `json:"reply"`
	} `json:"data"`
}

func handleAI(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText(fmt.Sprintf("🚩 Contoh penggunaan: %s%s apa itu nodejs", ptz.Prefix, ptz.Command))
	}

	ptz.Bot.Client.SendReaction(ptz.Chat, ptz.Info.ID, "🕒")
	text := strings.Join(ptz.Args, " ")
	var result string
	var err error

	switch ptz.Command {
	case "ai":
		result, err = fetchAI("https://api.alwayscodex.my.id/api/ai/gpt5?teks=" + text)
	case "blackbox":
		result, err = fetchAI("https://api.alwayscodex.my.id/api/ai/blackbox?teks=" + text)
	case "deepseek":
		result, err = fetchDeepseek("https://www.neoapis.xyz/api/ai/deepseek?text=" + text)
	default:
		return nil
	}

	if err != nil {
		return ptz.ReplyText("🚩 Terjadi kesalahan: " + err.Error())
	}

	return ptz.ReplyText(strings.ReplaceAll(result, "**", "*"))
}

func fetchAI(url string) (string, error) {
	resp, err := http.Get(url)
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
	return res.Result, nil
}

func fetchDeepseek(url string) (string, error) {
	resp, err := http.Get(url)
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
`, file_path: "commands/ai/ai.go")
