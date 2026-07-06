package commands

import (
	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:    []string{"owner", "creator", "dev"},
		Hidden:   []string{"creator", "dev"},
		Category: "miscs",
		Handler:  handleOwnerContacts,
	})
}

func handleOwnerContacts(ptz *core.Ptz) error {
	// Send multiple contacts as requested
	contacts := []struct {
		Phone string
		Name  string
	}{
		{
			Phone: ptz.Bot.Config.Owners[0],
			Name:  "Owner & Creator",
		},
		{
			Phone: "6282244425559",
			Name:  "Fahri (CO-Owner)",
		},
	}

	return serialize.SendMultipleContacts(ptz.Bot.Client, ptz.Chat, contacts)
}
