package tools

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "ssweb",
		Aliases:     []string{"screenshot", "ss"},
		Description: "Screenshot website dengan custom device dan theme",
		Usage:       "ssweb <url> [device] [theme] [fullpage]",
		Category:    "tools",
		Quota:       core.PerUserQuota(1),
		Handler: func(ctx *core.Ptz) error {
			ctx.React("⏳")
			defer ctx.Unreact()

			if len(ctx.Args) == 0 {
				return ctx.ReplyText("Format: .ssweb <url> [device] [theme] [fullpage]\n\nDevice: desktop, mobile, tablet (default: desktop)\nTheme: light, dark (default: light)\nFullPage: true, false (default: false)\n\nContoh:\n.ssweb https://google.com\n.ssweb https://google.com mobile\n.ssweb https://google.com mobile dark\n.ssweb https://google.com desktop dark true")
			}

			targetURL := ctx.Args[0]

			if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
				targetURL = "https://" + targetURL
			}

			_, err := url.Parse(targetURL)
			if err != nil {
				return ctx.ReplyText("URL tidak valid: " + err.Error())
			}

			device := "desktop"
			theme := "light"
			fullPage := "false"

			if len(ctx.Args) > 1 {
				device = strings.ToLower(ctx.Args[1])
				if device != "desktop" && device != "mobile" && device != "tablet" {
					return ctx.ReplyText("Device tidak valid! Gunakan: desktop, mobile, atau tablet")
				}
			}

			if len(ctx.Args) > 2 {
				theme = strings.ToLower(ctx.Args[2])
				if theme != "light" && theme != "dark" {
					return ctx.ReplyText("Theme tidak valid! Gunakan: light atau dark")
				}
			}

			if len(ctx.Args) > 3 {
				fullPage = strings.ToLower(ctx.Args[3])
				if fullPage != "true" && fullPage != "false" {
					return ctx.ReplyText("FullPage tidak valid! Gunakan: true atau false")
				}
			}

			tCtx, cancel := ctx.ContextWithTimeout(45 * time.Second)
			defer cancel()

			rawData, err := ctx.Bot.API.GetRaw(tCtx, "/api/tools/ssweb", map[string]string{
				"url":      targetURL,
				"device":   device,
				"theme":    theme,
				"fullPage": fullPage,
			})
			if err != nil {
				ctx.Bot.Log.Errorf("SSWeb error: %v", err)
				return ctx.ReplyText("❌ Terjadi kesalahan saat mengambil screenshot.")
			}

			if len(rawData) < 100 {
				return ctx.ReplyText("❌ Response tidak valid, silakan coba lagi.")
			}

			caption := fmt.Sprintf("📸 Website Screenshot\n\n🌐 URL: %s\n📱 Device: %s\n🎨 Theme: %s\n📄 Full Page: %s", targetURL, device, theme, fullPage)

			return ctx.ReplyImage(rawData, "image/png", caption)
		},
	})
}
