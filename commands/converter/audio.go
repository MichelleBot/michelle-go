package converter

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	commands := []string{"bass", "blown", "chipmunk", "deep", "earrape", "fast", "fat", "nightcore", "reverse", "robot", "slow"}
	core.Use(&core.Command{
		Usage:     commands,
		UsageHint: "reply audio",
		Category:  "voice changer",
		Limit:     core.PerUserLimit(1, 24*time.Hour),
		Handler:   runVoiceChanger,
	})
}

func runVoiceChanger(ptz *core.Ptz) error {
	if ptz.Message.ExtendedTextMessage == nil || ptz.Message.ExtendedTextMessage.ContextInfo == nil || ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage == nil {
		return ptz.ReplyText("🚩 Balas dengan audio untuk menggunakan perintah ini.")
	}

	quotedMsg := ptz.Message.ExtendedTextMessage.ContextInfo.QuotedMessage
	if quotedMsg.AudioMessage == nil {
		return ptz.ReplyText("🚩 Balas dengan audio untuk menggunakan perintah ini.")
	}

	ptz.React("🕒")

	data, err := serialize.DownloadMedia(ptz.Client, quotedMsg)
	if err != nil {
		return ptz.ReplyText("❌ Gagal mendownload audio: " + err.Error())
	}

	// Determine filter
	var filter string
	switch ptz.Command {
	case "bass":
		filter = "equalizer=f=94:width_type=o:width=2:g=30"
	case "blown":
		filter = "acrusher=.1:1:64:0:log"
	case "deep":
		filter = "atempo=4/4,asetrate=44500*2/3"
	case "earrape":
		filter = "volume=12"
	case "fast":
		filter = "atempo=1.63,asetrate=44100"
	case "fat":
		filter = "atempo=1.6,asetrate=22100"
	case "nightcore":
		filter = "atempo=1.06,asetrate=44100*1.25"
	case "reverse":
		filter = "areverse"
	case "robot":
		filter = "afftfilt=real='hypot(re,im)*sin(0)':imag='hypot(re,im)*cos(0)':win_size=512:overlap=0.75"
	case "slow":
		filter = "atempo=0.7,asetrate=44100"
	case "chipmunk":
		filter = "atempo=0.5,asetrate=65100"
	default:
		return ptz.ReplyText("🚩 Filter tidak dikenal.")
	}

	// Process audio
	tmpIn := tmpFile("vc_in_", ".ogg")
	tmpOut := tmpFile("vc_out_", ".mp3")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, data, 0600); err != nil {
		return ptz.ReplyText("❌ Gagal memproses audio.")
	}
	cmd := exec.Command("ffmpeg", "-y", "-i", tmpIn, "-af", filter, tmpOut)
	if out, err := cmd.CombinedOutput(); err != nil {
		ptz.Bot.Log.Errorf("FFmpeg error: %v, output: %s", err, string(out))
		return ptz.ReplyText("❌ Konversi gagal.")
	}

	buff, err := os.ReadFile(tmpOut)
	if err != nil {
		return ptz.ReplyText("❌ Gagal membaca hasil konversi.")
	}

	ptz.Unreact()
	return ptz.ReplyAudio(buff, "audio/mpeg")
}

func tmpFile(prefix, ext string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s%d_%d%s", prefix, os.Getpid(), rand.Int63(), ext))
}
