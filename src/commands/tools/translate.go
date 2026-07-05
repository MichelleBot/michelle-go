package tools

import (
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "translate",
		Aliases:     []string{"tr", "terjemah"},
		Description: "Terjemahkan teks ke bahasa lain",
		Usage:       "translate <source>:<target> <teks>\nContoh: translate en:id Hello world",
		Category:    "tools",
		Quota:       core.PerUserQuota(1),
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()

			if len(ptz.Args) < 2 {
				return ptz.ReplyText("*translate* — Terjemahkan teks\n\nUsage: .translate <source>:<target> <teks>\nContoh:\n.translate en:id Hello world\n.translate id:en Halo dunia\n\nKode bahasa: en, id, ja, ko, zh, ar, fr, de, es, dll")
			}

			langPair := strings.SplitN(ptz.Args[0], ":", 2)
			if len(langPair) != 2 {
				return ptz.ReplyText("❌ Format bahasa salah. Contoh: en:id")
			}

			source := langPair[0]
			target := langPair[1]
			text := strings.Join(ptz.Args[1:], " ")

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type TranslateData struct {
				TranslatedText string `json:"translatedText"`
			}

			data, err := api.Request[TranslateData](ctx, ptz.Bot.API, "/api/tools/translate", map[string]string{
				"text":   text,
				"source": source,
				"target": target,
			})
			if err != nil {
				ptz.Bot.Log.Errorf("Translate error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			return ptz.ReplyText("🌐 *Translate*\n\n" +
				"🔤 *" + strings.ToUpper(source) + "* → *" + strings.ToUpper(target) + "*\n\n" +
				"📥 " + text + "\n" +
				"📤 " + data.TranslatedText)
		},
	})
}
