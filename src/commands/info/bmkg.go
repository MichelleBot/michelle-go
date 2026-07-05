package info

import (
	"context"
	"encoding/json"
	"fmt"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

type bmkgGempa struct {
	Tanggal          string `json:"Tanggal"`
	Jam              string `json:"Jam"`
	DateTime         string `json:"DateTime"`
	Coordinates      string `json:"Coordinates"`
	Lintang          string `json:"Lintang"`
	Bujur            string `json:"Bujur"`
	Magnitude        string `json:"Magnitude"`
	Kedalaman        string `json:"Kedalaman"`
	Wilayah          string `json:"Wilayah"`
	Potensi          string `json:"Potensi"`
	Dirasakan        string `json:"Dirasakan"`
	Shakemap         string `json:"Shakemap"`
	DownloadShakemap string `json:"downloadShakemap"`
}

type bmkgResp struct {
	Status bool `json:"status"`
	Data   struct {
		Auto struct {
			Infogempa struct {
				Gempa bmkgGempa `json:"gempa"`
			} `json:"Infogempa"`
		} `json:"auto"`
		Terkini struct {
			Infogempa struct {
				Gempa []bmkgGempa `json:"gempa"`
			} `json:"Infogempa"`
		} `json:"terkini"`
		Dirasakan struct {
			Infogempa struct {
				Gempa []bmkgGempa `json:"gempa"`
			} `json:"Infogempa"`
		} `json:"dirasakan"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "bmkg",
		Description: "Info gempa terbaru dari BMKG",
		Usage:       "bmkg",
		Category:    "info",
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.GetRaw(context.Background(), "/api/info/bmkg", nil)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			var result bmkgResp
			if err := json.Unmarshal(raw, &result); err != nil {
				return ptz.ReplyText("❌ Gagal parse response.")
			}
			if !result.Status {
				return ptz.ReplyText("❌ Gagal mengambil data BMKG.")
			}

			auto := result.Data.Auto.Infogempa.Gempa
			terkini := result.Data.Terkini.Infogempa.Gempa
			dirasakan := result.Data.Dirasakan.Infogempa.Gempa

			var sb strings.Builder
			sb.WriteString("🌍 *Info Gempa BMKG*\n\n")

			sb.WriteString("⚡ *Gempa Terbaru (Dirasakan)*\n")
			sb.WriteString(fmt.Sprintf("📅 %s %s\n", auto.Tanggal, auto.Jam))
			sb.WriteString(fmt.Sprintf("📍 %s\n", auto.Wilayah))
			sb.WriteString(fmt.Sprintf("💥 *M %s* — Kedalaman %s\n", auto.Magnitude, auto.Kedalaman))
			sb.WriteString(fmt.Sprintf("📡 Koordinat: %s\n", auto.Coordinates))
			if auto.Dirasakan != "" {
				sb.WriteString(fmt.Sprintf("👥 Dirasakan: %s\n", auto.Dirasakan))
			}
			if auto.Potensi != "" {
				sb.WriteString(fmt.Sprintf("⚠️ %s\n", auto.Potensi))
			}

			if len(terkini) > 0 {
				sb.WriteString("\n📋 *5 Gempa Terkini (M≥5)*\n")
				limit := 5
				if len(terkini) < limit {
					limit = len(terkini)
				}
				for i := 0; i < limit; i++ {
					g := terkini[i]
					sb.WriteString(fmt.Sprintf("• *M%s* %s — %s\n", g.Magnitude, g.Tanggal, g.Wilayah))
				}
			}

			if len(dirasakan) > 0 {
				sb.WriteString("\n📢 *3 Gempa Dirasakan Terakhir*\n")
				limit := 3
				if len(dirasakan) < limit {
					limit = len(dirasakan)
				}
				for i := 0; i < limit; i++ {
					g := dirasakan[i]
					sb.WriteString(fmt.Sprintf("• *M%s* %s — %s\n", g.Magnitude, g.Tanggal, g.Wilayah))
					if g.Dirasakan != "" {
						sb.WriteString(fmt.Sprintf("  👥 %s\n", g.Dirasakan))
					}
				}
			}

			caption := sb.String()

			if auto.DownloadShakemap != "" {
				imgData, err := serialize.Fetch(auto.DownloadShakemap)
				if err == nil {
					return ptz.ReplyImage(imgData, "image/jpeg", caption)
				}
			}

			return ptz.ReplyText(caption)
		},
	})
}
