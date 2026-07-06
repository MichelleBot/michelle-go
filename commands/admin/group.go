package admin

import (
	"fmt"
	"strings"

	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"group"},
		UsageHint: "open / close",
		Category:  "admin",
		GroupOnly: true,
		AdminOnly: true,
		BotAdmin:  true,
		Handler:   runGroup,
	})
}

func runGroup(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText("🚩 Masukkan argumen, close atau open.")
	}

	option := strings.ToLower(ptz.Args[0])
	
	// 'false' means group is open (not announcement)
	// 'true' means group is closed (announcement)
	var announce bool
	if option == "open" {
		announce = false
	} else if option == "close" {
		announce = true
	} else {
		return ptz.ReplyText("🚩 Argumen tidak valid. Gunakan 'open' atau 'close'.")
	}

	err := serialize.SetGroupAnnounce(ptz.Bot.Client, ptz.Chat, announce)
	if err != nil {
		return ptz.ReplyText(fmt.Sprintf("❌ Gagal mengubah pengaturan grup: %v", err))
	}

	return ptz.ReplyText(fmt.Sprintf("🚩 Grup berhasil di-%s.", option))
}
