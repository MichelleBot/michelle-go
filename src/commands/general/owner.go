package general

import (
	"fmt"

	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:        "owner",
		Aliases:     []string{"creator", "dev"},
		Description: "Kirim daftar kontak owner bot",
		Usage:       "owner",
		Category:    "general",
		Handler:     handleOwnerContacts,
	})
}

func handleOwnerContacts(ptz *core.Ptz) error {
	if ptz.Bot == nil || ptz.Bot.Config == nil || len(ptz.Bot.Config.Owners) == 0 {
		return ptz.ReplyText("Owner bot belum diset di konfigurasi.")
	}

	contacts := make([]struct {
		Phone string
		Name  string
	}, 0, len(ptz.Bot.Config.Owners))

	for i, owner := range ptz.Bot.Config.Owners {
		contacts = append(contacts, struct {
			Phone string
			Name  string
		}{
			Phone: owner,
			Name:  fmt.Sprintf("Owner Bot %d", i+1),
		})
	}

	return serialize.SendMultipleContacts(ptz.Bot.Client, ptz.Chat, contacts)
}
