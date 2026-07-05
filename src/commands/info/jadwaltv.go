package info

import (
	"context"
	"encoding/json"
	"fmt"
	"sptzx/src/api"
	"sptzx/src/core"
	"strings"
)

type jadwalItem struct {
	Jam   string `json:"jam"`
	Acara string `json:"acara"`
}

type jadwalResp struct {
	Status bool         `json:"status"`
	Data   []jadwalItem `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "jadwaltv",
		Description: "Jadwal acara TV nasional",
		Usage:       "jadwaltv <channel>",
		Category:    "info",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*jadwaltv* — Jadwal acara TV\n\nUsage: .jadwaltv <channel>\nChannel: sctv, rcti, trans7, transtv, gtv, mnctv, net, tvone, metrotv, antv, indosiar, inewstv, kompastv\n\nContoh: .jadwaltv sctv")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			channel := strings.ToLower(ptz.Args[0])
			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/info/jadwaltv", map[string]string{"channel": channel})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result jadwalResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}
			if !result.Status || len(result.Data) == 0 {
				return ptz.ReplyText("❌ Jadwal tidak ditemukan. Pastikan nama channel benar.")
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("📺 *Jadwal %s*\n\n", strings.ToUpper(channel)))

			for i, item := range result.Data {
				if i >= 20 {
					sb.WriteString(fmt.Sprintf("_...dan %d acara lainnya_", len(result.Data)-20))
					break
				}
				sb.WriteString(fmt.Sprintf("🕐 *%s*\n%s\n\n", item.Jam, item.Acara))
			}

			return ptz.ReplyText(sb.String())
		},
	})
}
