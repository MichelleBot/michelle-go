package general

import (
	"context"
	"fmt"
	"strings"

	"github.com/showwin/speedtest-go/speedtest"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "speedtest",
		Aliases:     []string{"speed", "st"},
		Description: "Test kecepatan koneksi server",
		Usage:       "speedtest",
		Category:    "general",
		Handler: func(ptz *core.Ptz) error {
			if err := ptz.React("⏳"); err != nil {
				ptz.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ptz.Unreact()

			client := speedtest.New()

			u, err := client.FetchUserInfo()
			if err != nil {
				return ptz.ReplyText("Gagal fetch info: " + err.Error())
			}

			servers, err := client.FetchServers()
			if err != nil {
				return ptz.ReplyText("Gagal fetch server: " + err.Error())
			}
			if servers.Len() == 0 {
				return ptz.ReplyText("Tidak ada server ditemukan")
			}

			server := servers[0]

			if err = server.PingTestContext(context.Background(), nil); err != nil {
				return ptz.ReplyText("Gagal ping test: " + err.Error())
			}
			if err = server.DownloadTestContext(context.Background()); err != nil {
				return ptz.ReplyText("Gagal download test: " + err.Error())
			}
			if err = server.UploadTestContext(context.Background()); err != nil {
				return ptz.ReplyText("Gagal upload test: " + err.Error())
			}

			var sb strings.Builder
			sb.WriteString("🚀 *Speedtest*\n\n")
			sb.WriteString(fmt.Sprintf("*Server:* %s, %s\n", server.Name, server.Country))
			sb.WriteString(fmt.Sprintf("*Sponsor:* %s\n\n", server.Sponsor))
			sb.WriteString(fmt.Sprintf("*Ping:* %d ms\n", server.Latency.Milliseconds()))
			sb.WriteString(fmt.Sprintf("*Download:* %.2f Mbps\n", server.DLSpeed.Mbps()))
			sb.WriteString(fmt.Sprintf("*Upload:* %.2f Mbps\n\n", server.ULSpeed.Mbps()))
			sb.WriteString(fmt.Sprintf("*ISP:* %s\n", u.Isp))
			sb.WriteString(fmt.Sprintf("*IP:* %s", u.IP))

			return ptz.ReplyText(sb.String())
		},
	})
}
