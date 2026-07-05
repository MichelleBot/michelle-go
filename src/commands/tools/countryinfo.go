package tools

import (
	"fmt"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "countryinfo",
		Aliases:     []string{"country", "negara"},
		Description: "Info lengkap tentang sebuah negara",
		Usage:       "countryinfo <nama negara>",
		Category:    "tools",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*countryinfo* — Info negara\n\nUsage: .countryinfo <nama negara>\nContoh: .countryinfo Indonesia")
			}

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type CountryData struct {
				Name           string `json:"name"`
				Capital        string `json:"capital"`
				Flag           string `json:"flag"`
				PhoneCode      string `json:"phoneCode"`
				GoogleMapsLink string `json:"googleMapsLink"`
				Continent      struct {
					Name  string `json:"name"`
					Emoji string `json:"emoji"`
				} `json:"continent"`
				Area struct {
					SquareKilometers float64 `json:"squareKilometers"`
				} `json:"area"`
				Landlocked         bool   `json:"landlocked"`
				FamousFor          string `json:"famousFor"`
				ConstitutionalForm string `json:"constitutionalForm"`
				Currency           string `json:"currency"`
				DrivingSide        string `json:"drivingSide"`
				InternetTLD        string `json:"internetTLD"`
				Languages          struct {
					Native []string `json:"native"`
				} `json:"languages"`
				Neighbors []struct {
					Name string `json:"name"`
				} `json:"neighbors"`
			}

			data, err := api.Request[CountryData](ctx, ptz.Bot.API, "/api/tools/countryInfo", map[string]string{
				"name": strings.Join(ptz.Args, " "),
			})
			if err != nil {
				ptz.Bot.Log.Errorf("CountryInfo error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			var neighbors []string
			for _, n := range data.Neighbors {
				neighbors = append(neighbors, n.Name)
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("🌍 *%s* %s\n\n", data.Name, data.Continent.Emoji))
			sb.WriteString(fmt.Sprintf("🏛 Ibu Kota: *%s*\n", data.Capital))
			sb.WriteString(fmt.Sprintf("🌐 Benua: %s\n", data.Continent.Name))
			sb.WriteString(fmt.Sprintf("📞 Kode Tel: %s\n", data.PhoneCode))
			sb.WriteString(fmt.Sprintf("💰 Mata Uang: %s\n", data.Currency))
			sb.WriteString(fmt.Sprintf("🗣 Bahasa: %s\n", strings.Join(data.Languages.Native, ", ")))
			sb.WriteString(fmt.Sprintf("📐 Luas: %.0f km²\n", data.Area.SquareKilometers))
			sb.WriteString(fmt.Sprintf("🌐 TLD: %s\n", data.InternetTLD))
			sb.WriteString(fmt.Sprintf("🚗 Setir: %s\n", data.DrivingSide))
			if data.FamousFor != "" {
				sb.WriteString(fmt.Sprintf("🌟 Terkenal: %s\n", data.FamousFor))
			}
			if len(neighbors) > 0 {
				sb.WriteString(fmt.Sprintf("🗺 Berbatasan: %s\n", strings.Join(neighbors, ", ")))
			}
			sb.WriteString(fmt.Sprintf("\n🔗 %s", data.GoogleMapsLink))

			caption := sb.String()

			if data.Flag != "" {
				img, err := serialize.Fetch(data.Flag)
				if err == nil {
					return ptz.ReplyImage(img, "image/png", caption)
				}
			}

			return ptz.ReplyText(caption)
		},
	})
}
