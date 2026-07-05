package info

import (
	"context"
	"encoding/json"
	"fmt"
	"sptzx/src/api"
	"sptzx/src/core"
	"strings"
)

type cuacaItem struct {
	Datetime      string  `json:"datetime"`
	T             float64 `json:"t"`
	WeatherDesc   string  `json:"weather_desc"`
	Hu            float64 `json:"hu"`
	Ws            float64 `json:"ws"`
	Wd            string  `json:"wd"`
	VsText        string  `json:"vs_text"`
	LocalDatetime string  `json:"local_datetime"`
}

type cuacaLokasi struct {
	Provinsi  string `json:"provinsi"`
	Kotkab    string `json:"kotkab"`
	Kecamatan string `json:"kecamatan"`
	Desa      string `json:"desa"`
}

type cuacaWilayah struct {
	Nama string `json:"nama"`
}

type cuacaResp struct {
	Status bool `json:"status"`
	Data   struct {
		Wilayah cuacaWilayah `json:"wilayah"`
		Weather []struct {
			Lokasi cuacaLokasi   `json:"lokasi"`
			Cuaca  [][]cuacaItem `json:"cuaca"`
		} `json:"weather"`
	} `json:"data"`
}

func weatherEmoji(desc string) string {
	desc = strings.ToLower(desc)
	switch {
	case strings.Contains(desc, "petir"):
		return "⛈️"
	case strings.Contains(desc, "hujan lebat"):
		return "🌧️"
	case strings.Contains(desc, "hujan"):
		return "🌦️"
	case strings.Contains(desc, "berawan"):
		return "☁️"
	case strings.Contains(desc, "cerah berawan"):
		return "⛅"
	case strings.Contains(desc, "cerah"):
		return "☀️"
	default:
		return "🌡️"
	}
}

func init() {
	core.Use(&core.Command{
		Name:        "cuaca",
		Description: "Info prakiraan cuaca dari BMKG",
		Usage:       "cuaca <lokasi>",
		Category:    "info",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*cuaca* — Prakiraan cuaca dari BMKG\n\nUsage: .cuaca <lokasi>\nContoh: .cuaca Jakarta")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			query := strings.Join(ptz.Args, " ")
			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/info/cuaca", map[string]string{"q": query})
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result cuacaResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}
			if !result.Status {
				return ptz.ReplyText("❌ Lokasi tidak ditemukan.")
			}

			if len(result.Data.Weather) == 0 {
				return ptz.ReplyText("❌ Data cuaca tidak tersedia.")
			}

			loc := result.Data.Weather[0].Lokasi
			wilayah := result.Data.Wilayah.Nama
			cuacaData := result.Data.Weather[0].Cuaca

			var sb strings.Builder
			sb.WriteString("🌤️ *Prakiraan Cuaca BMKG*\n\n")
			sb.WriteString(fmt.Sprintf("📍 *%s*\n", wilayah))
			sb.WriteString(fmt.Sprintf("%s, %s, %s\n\n", loc.Desa, loc.Kecamatan, loc.Provinsi))

			count := 0
			for _, dayGroup := range cuacaData {
				for _, item := range dayGroup {
					if count >= 6 {
						break
					}
					emoji := weatherEmoji(item.WeatherDesc)
					localTime := item.LocalDatetime
					if len(localTime) >= 16 {
						localTime = localTime[5:16]
					}
					sb.WriteString(fmt.Sprintf("%s *%s* — %s\n", emoji, localTime, item.WeatherDesc))
					sb.WriteString(fmt.Sprintf("   🌡️ %g°C  💧 %g%%  💨 %g m/s %s\n", item.T, item.Hu, item.Ws, item.Wd))
					count++
				}
				if count >= 6 {
					break
				}
			}

			return ptz.ReplyText(sb.String())
		},
	})
}
