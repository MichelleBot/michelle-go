package search

import (
	"context"
	"encoding/json"
	"fmt"
	"sptzx/src/api"
	"sptzx/src/core"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "duck",
		Aliases:     []string{"ddg", "duckduckgo"},
		Description: "Search dengan DuckDuckGo",
		Usage:       "duck <query> [region] [timeframe]",
		Category:    "search",
		Handler: func(ctx *core.Ptz) error {
			if err := ctx.React("⏳"); err != nil {
				ctx.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ctx.Unreact()

			ctx.Bot.Log.Infof("DuckDuckGo command triggered with args: %v", ctx.Args)

			if len(ctx.Args) == 0 {
				return ctx.ReplyText("Format: .duck <query> [region] [timeframe]\n\nRegion: us-en, id-id, uk-en, dll (default: us-en)\nTimeframe: d, w, m, y (day/week/month/year)\n\nContoh:\n.duck openai\n.duck openai us-en\n.duck openai us-en w")
			}

			args := ctx.Args
			ctx.Bot.Log.Infof("DuckDuckGo parsing args: %v", args)

			query := strings.Join(args, " ")
			region := "us-en"
			timeframe := ""

			lastIdx := len(args) - 1

			if lastIdx >= 0 {
				tf := strings.ToLower(args[lastIdx])
				switch tf {
				case "d", "w", "m", "y":
					timeframe = tf
					lastIdx--
				}
			}

			if lastIdx >= 0 {
				potentialRegion := args[lastIdx]
				if strings.Contains(potentialRegion, "-") {
					region = potentialRegion
					lastIdx--
				}
			}

			if lastIdx >= 0 {
				query = strings.Join(args[:lastIdx+1], " ")
			}

			ctx.Bot.Log.Infof("DuckDuckGo final params: query=%s, kl=%s, df=%s", query, region, timeframe)

			params := map[string]string{
				"query": query,
				"kl":    region,
			}
			if timeframe != "" {
				params["df"] = timeframe
			}

			client := api.NewClient(ctx.Bot.Config.SiputzX.BaseURL)
			ctx.Bot.Log.Infof("DuckDuckGo API URL: %s", ctx.Bot.Config.SiputzX.BaseURL)
			ctx.Bot.Log.Infof("DuckDuckGo params: query=%s, kl=%s, df=%s", query, region, timeframe)

			rawData, err := client.GetRaw(context.Background(), "/api/s/duckduckgo", params)
			ctx.Bot.Log.Infof("DuckDuckGo raw response length: %d", len(rawData))
			if err != nil {
				ctx.Bot.Log.Errorf("DuckDuckGo API error: %v", err)
				return ctx.ReplyText("Gagal: " + err.Error())
			}

			ctx.Bot.Log.Infof("DuckDuckGo raw response: %s", string(rawData))

			var result map[string]interface{}
			if err := json.Unmarshal(rawData, &result); err != nil {
				ctx.Bot.Log.Debugf("DuckDuckGo parse error: %v", err)
				return ctx.ReplyText("Gagal parse response: " + err.Error())
			}

			status, exists := result["status"].(bool)
			ctx.Bot.Log.Debugf("DuckDuckGo status: %v (exists: %v)", status, exists)
			if !status {
				return ctx.ReplyText("Search gagal")
			}

			data, ok := result["data"].(map[string]interface{})
			if !ok {
				return ctx.ReplyText("Invalid response format")
			}

			resultsRaw, ok := data["results"].([]interface{})
			if !ok {
				return ctx.ReplyText("Tidak ada hasil ditemukan")
			}

			if len(resultsRaw) == 0 {
				return ctx.ReplyText("Tidak ada hasil ditemukan")
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("🦆 DuckDuckGo Search\n\n🔍 Query: %s\n📍 Region: %s", query, region))
			if timeframe != "" {
				sb.WriteString(fmt.Sprintf("\n⏰ Time: %s", timeframeName(timeframe)))
			}
			sb.WriteString(fmt.Sprintf("\n📊 Results: %d\n\n", len(resultsRaw)))

			for i, r := range resultsRaw {
				if i >= 10 {
					break
				}
				result, ok := r.(map[string]interface{})
				if !ok {
					continue
				}

				title, _ := result["title"].(string)
				url, _ := result["url"].(string)
				snippet, _ := result["snippet"].(string)
				timestamp, _ := result["timestamp"].(string)

				sb.WriteString(fmt.Sprintf("*%d. %s*\n", i+1, title))
				if snippet != "" {
					if len(snippet) > 150 {
						snippet = snippet[:150] + "..."
					}
					sb.WriteString(fmt.Sprintf("%s\n", snippet))
				}
				sb.WriteString(fmt.Sprintf("🔗 %s\n", url))
				if timestamp != "" {
					sb.WriteString(fmt.Sprintf("📅 %s\n", timestamp))
				}
				sb.WriteString("\n")
			}

			if len(resultsRaw) > 10 {
				sb.WriteString(fmt.Sprintf("... dan %d hasil lainnya\n", len(resultsRaw)-10))
			}

			return ctx.ReplyText(sb.String())
		},
	})
}

func timeframeName(tf string) string {
	switch tf {
	case "d":
		return "Today"
	case "w":
		return "This Week"
	case "m":
		return "This Month"
	case "y":
		return "This Year"
	default:
		return "Any Time"
	}
}
