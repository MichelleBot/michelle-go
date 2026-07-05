package search

import (
	"fmt"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "newsletter",
		Description: "Dapatkan info newsletter/channel dari link",
		Usage:       "newsletter <link>",
		Category:    "search",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Format: .newsletter <link>")
			}

			if err := ptz.React("⏳"); err != nil {
				ptz.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ptz.Unreact()

			link := strings.Join(ptz.Args, " ")
			key := strings.TrimPrefix(link, "https://whatsapp.com/channel/")

			info, err := serialize.GetNewsletterInfoWithInvite(ptz.Bot.Client, key)
			if err != nil {
				return ptz.ReplyText("Error: " + err.Error())
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Info Newsletter: %s\n\n", info.ThreadMeta.Name.Text))

			if info.ThreadMeta.Description.Text != "" {
				sb.WriteString(fmt.Sprintf("Deskripsi: %s\n\n", info.ThreadMeta.Description.Text))
			}

			sb.WriteString(fmt.Sprintf("ID: %s\n", info.ID.String()))
			sb.WriteString(fmt.Sprintf("Subscriber: %d\n", info.ThreadMeta.SubscriberCount))
			sb.WriteString(fmt.Sprintf("Status: %s\n", info.State.Type))
			sb.WriteString(fmt.Sprintf("Verified: %s\n", info.ThreadMeta.VerificationState))

			return ptz.ReplyText(sb.String())
		},
	})
}
