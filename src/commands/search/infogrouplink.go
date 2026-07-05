package search

import (
	"fmt"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "grouplink",
		Description: "Dapatkan info group dari invite link",
		Usage:       "grouplink <link>",
		Category:    "search",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("Format: .grouplink <link>")
			}

			if err := ptz.React("⏳"); err != nil {
				ptz.Bot.Log.Debugf("Failed to react: %v", err)
			}
			defer ptz.Unreact()

			link := strings.Join(ptz.Args, " ")
			code := strings.TrimPrefix(link, "https://chat.whatsapp.com/")

			info, err := serialize.GetGroupInfoFromLink(ptz.Bot.Client, code)
			if err != nil {
				return ptz.ReplyText("Error: " + err.Error())
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Info Group: %s\n\n", info.Name))

			if info.Topic != "" {
				sb.WriteString(fmt.Sprintf("Topic: %s\n\n", info.Topic))
			}

			sb.WriteString(fmt.Sprintf("ID: %s\n", info.JID.String()))
			sb.WriteString(fmt.Sprintf("Owner: %s\n", info.OwnerJID.String()))
			sb.WriteString(fmt.Sprintf("Jumlah Member: %d\n", len(info.Participants)))

			return ptz.ReplyText(sb.String())
		},
	})
}
