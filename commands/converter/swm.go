package converter

import (
	"fmt"
	"strings"

	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"swm"},
		UsageHint: "reply media",
		Category:  "converter",
		Handler:   runSwm,
	})
}

func runSwm(ptz *core.Ptz) error {
	input := serialize.GetInputMedia(ptz.Message, "image", "video")
	if input == nil {
		return ptz.ReplyText(fmt.Sprintf("🚩 Untuk membuat watermark pada stiker, balas media foto atau video dan gunakan format ini *%s%s packname | author*", ptz.Bot.Config.Prefixes[0], ptz.Command))
	}

	packName := ""
	author := ""
	if ptz.RawArgs != "" {
		parts := strings.SplitN(ptz.RawArgs, "|", 2)
		packName = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			author = strings.TrimSpace(parts[1])
		}
	}

	ptz.React("🕒")

	data, err := serialize.DownloadMedia(ptz.Client, input.Message)
	if err != nil {
		return ptz.ReplyText("❌ Gagal download: " + err.Error())
	}

	mime := serialize.GetMediaMIME(input.Message)
	ext := serialize.GetMediaExtFromMIME(mime)
	meta := serialize.StickerMetadata{
		PackName:   packName,
		Author:     author,
		Categories: []string{""},
	}

	if input.MsgType == "video" {
		secs, err := serialize.GetVideoDurationSeconds(data, ext)
		if err == nil && secs > 10 {
			return ptz.ReplyText("🚩 Durasi video maksimal adalah 10 detik.")
		}

		webp, err := serialize.ToAnimatedWebpExif(data, ext, true, meta)
		if err != nil {
			return ptz.ReplyText("❌ Gagal convert video ke sticker: " + err.Error())
		}
		return ptz.ReplySticker(webp, "image/webp", true)
	}

	webp, err := serialize.ToStaticWebpExif(data, ext, meta)
	if err != nil {
		return ptz.ReplyText("❌ Gagal convert image ke sticker: " + err.Error())
	}
	return ptz.ReplySticker(webp, "image/webp", false)
}
