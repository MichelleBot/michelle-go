package owner

import (
	"io"
	"mime"
	"net/url"
	"strings"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"

	"sptzx/src/core"
	"sptzx/src/utils"
)

const (
	maxFileSize       = 100 * 1024 * 1024
	documentThreshold = 10 * 1024 * 1024
)

func init() {
	core.Use(&core.Command{
		Name:        "get",
		Description: "HTTP GET request dengan browser impersonation",
		Usage:       "get <url> [browser]",
		Category:    "owner",
		OwnerOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			if err := ptz.React("⏳"); err != nil {
				ptz.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ptz.Unreact()

			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Format: .get <url> [browser]\n\nBrowser: chrome, firefox, safari, edge\nDefault: chrome\n\nContoh: .get https://api.example.com/data")
			}

			targetURL := ptz.Args[0]
			browser := "chrome"
			if len(ptz.Args) > 1 {
				browser = strings.ToLower(ptz.Args[1])
			}

			if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
				targetURL = "https://" + targetURL
			}

			parsedURL, err := url.Parse(targetURL)
			if err != nil {
				return ptz.ReplyText("URL tidak valid: " + err.Error())
			}

			ptz.Bot.Log.Infof("GET %s with %s impersonation", parsedURL.String(), browser)

			profile := utils.GetBrowserProfile(browser)
			jar := tls_client.NewCookieJar()

			options := []tls_client.HttpClientOption{
				tls_client.WithTimeoutSeconds(60),
				tls_client.WithClientProfile(profile),
				tls_client.WithCookieJar(jar),
			}

			client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
			if err != nil {
				return ptz.ReplyText("Gagal create client: " + err.Error())
			}

			req, err := fhttp.NewRequest(fhttp.MethodGet, parsedURL.String(), nil)
			if err != nil {
				return ptz.ReplyText("Gagal create request: " + err.Error())
			}

			req.Header = utils.GetDefaultHeaders(browser)

			resp, err := client.Do(req)
			if err != nil {
				return ptz.ReplyText("Request failed: " + err.Error())
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(io.LimitReader(resp.Body, maxFileSize))
			if err != nil {
				return ptz.ReplyText("Gagal baca response: " + err.Error())
			}

			contentLength := len(body)
			ptz.Bot.Log.Infof("Response size: %d bytes", contentLength)

			contentType := resp.Header.Get("Content-Type")
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			ptz.Bot.Log.Infof("Content-Type: %s", contentType)

			if contentLength > documentThreshold {
				return utils.SendAsDocument(ptz, body, contentType, parsedURL.Path)
			}

			mediaType, _, _ := mime.ParseMediaType(contentType)

			if strings.HasPrefix(mediaType, "image/") {
				return utils.SendAsImage(ptz, body, mediaType)
			} else if strings.HasPrefix(mediaType, "video/") {
				return utils.SendAsVideo(ptz, body, mediaType)
			} else if strings.HasPrefix(mediaType, "audio/") {
				return utils.SendAsAudio(ptz, body, mediaType)
			} else if strings.Contains(mediaType, "json") {
				return utils.SendAsFormattedJSON(ptz, body)
			} else if strings.Contains(mediaType, "pdf") {
				return utils.SendAsDocument(ptz, body, contentType, parsedURL.Path)
			} else if strings.Contains(mediaType, "zip") || strings.Contains(mediaType, "rar") || strings.Contains(mediaType, "7z") {
				return utils.SendAsDocument(ptz, body, contentType, parsedURL.Path)
			} else if strings.Contains(mediaType, "msword") || strings.Contains(mediaType, "officedocument") {
				return utils.SendAsDocument(ptz, body, contentType, parsedURL.Path)
			} else if strings.Contains(mediaType, "android.package-archive") || strings.HasSuffix(parsedURL.Path, ".apk") {
				return utils.SendAsDocument(ptz, body, contentType, parsedURL.Path)
			} else if strings.Contains(mediaType, "octet-stream") {
				return utils.SendAsDocument(ptz, body, contentType, parsedURL.Path)
			} else if strings.Contains(mediaType, "html") {
				return utils.SendAsText(ptz, body, "text/html")
			} else if strings.Contains(mediaType, "xml") {
				return utils.SendAsText(ptz, body, "text/xml")
			} else if strings.Contains(mediaType, "text") {
				return utils.SendAsText(ptz, body, "text/plain")
			}

			return utils.SendAsDocument(ptz, body, contentType, parsedURL.Path)
		},
	})
}
