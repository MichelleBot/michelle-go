package primbon

import (
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "zodiak",
		Aliases:     []string{"horoscope", "horoskop"},
		Description: "Ramalan zodiak",
		Usage:       "zodiak <zodiak>",
		Category:    "primbon",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*zodiak* — Ramalan zodiak\n\nUsage: .zodiak <zodiak>\nContoh: .zodiak gemini\n\nZodiak: aries, taurus, gemini, cancer, leo, virgo, libra, scorpio, sagitarius, capricorn, aquarius, pisces")
			}

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type ZodiakData struct {
				Zodiak              string `json:"zodiak"`
				NomorKeberuntungan  string `json:"nomor_keberuntungan"`
				AromaKeberuntungan  string `json:"aroma_keberuntungan"`
				PlanetYangMengitari string `json:"planet_yang_mengitari"`
				BungaKeberuntungan  string `json:"bunga_keberuntungan"`
				WarnaKeberuntungan  string `json:"warna_keberuntungan"`
				BatuKeberuntungan   string `json:"batu_keberuntungan"`
				ElemenKeberuntungan string `json:"elemen_keberuntungan"`
				PasanganZodiak      string `json:"pasangan_zodiak"`
			}

			data, err := api.Request[ZodiakData](ctx, ptz.Bot.API, "/api/primbon/zodiak", map[string]string{
				"zodiak": strings.ToLower(strings.Join(ptz.Args, " ")),
			})
			if err != nil {
				ptz.Bot.Log.Errorf("Zodiak error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server.")
			}

			var sb strings.Builder
			sb.WriteString("⭐ *Zodiak*\n\n")
			sb.WriteString("*" + data.Zodiak + "*\n\n")
			sb.WriteString("🔢 Nomor hoki: " + data.NomorKeberuntungan + "\n")
			sb.WriteString("🌸 Bunga: " + data.BungaKeberuntungan + "\n")
			sb.WriteString("🎨 Warna: " + data.WarnaKeberuntungan + "\n")
			sb.WriteString("💎 Batu: " + data.BatuKeberuntungan + "\n")
			sb.WriteString("🌬 Elemen: " + data.ElemenKeberuntungan + "\n")
			sb.WriteString("🪐 Planet: " + data.PlanetYangMengitari + "\n")
			sb.WriteString("💑 Pasangan: " + data.PasanganZodiak[:50] + "...\n")
			sb.WriteString("🌺 Aroma: " + data.AromaKeberuntungan)
			return ptz.ReplyText(sb.String())
		},
	})
}
