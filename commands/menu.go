package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"michelle/system/core"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Name:        "menu",
		Aliases:     []string{"help", "start"},
		Description: "Lihat daftar semua command",
		Usage:       "menu",
		Category:    "general",
		Limit:       core.PerUserLimit(30, time.Minute),
		Handler: func(ptz *core.Ptz) error {
			return sendFullMenu(ptz)
		},
	})
}

func sendFullMenu(ptz *core.Ptz) error {
	greeting := utils.Greeting(ptz.Bot.Config.Timezone)
	name := ptz.GetSenderName()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%s, %s*\n\n", greeting, name))
	sb.WriteString("*Daftar semua command:*\n\n")

	byCat := core.GlobalRegistry().ByCategory()
	categories := make([]string, 0, len(byCat))
	for cat := range byCat {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf("*[%s]*\n", strings.ToUpper(cat)))
		
		cmds := byCat[cat]
		for _, cmd := range cmds {
			sb.WriteString(fmt.Sprintf("- %s\n", cmd.Name))
		}
		sb.WriteString("\n")
	}

	return ptz.ReplyText(sb.String())
}
