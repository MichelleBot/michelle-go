package owner

import (
	"regexp"

	"michelle/system/core"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"join"},
		UsageHint: "link",
		Category:  "owner",
		OwnerOnly: true,
		Handler:   handleJoin,
	})
}

func handleJoin(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText("🚩 Contoh penggunaan: " + ptz.Prefix + ptz.Command + " https://chat.whatsapp.com/kodeUndangan")
	}

	linkRegex := regexp.MustCompile(`chat.whatsapp.com/([0-9A-Za-z]{20,24})`)
	match := linkRegex.FindStringSubmatch(ptz.Args[0])
	
	if len(match) < 2 {
		return ptz.ReplyText("🚩 Link tidak valid.")
	}
	
	code := match[1]
	
	jid, err := ptz.Bot.Client.JoinGroupWithLink(code)
	if err != nil {
		return ptz.ReplyText("🚩 Maaf saya tidak bisa bergabung ke grup ini :(")
	}
	
	if jid == nil {
		return ptz.ReplyText("🚩 Maaf saya tidak bisa bergabung ke grup ini :(")
	}

	return ptz.ReplyText("🚩 Bergabung!")
}
