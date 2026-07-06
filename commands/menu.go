package commands

import (
	"fmt"
	"sort"
	"strings"

	"go.mau.fi/whatsmeow/types"
	"michelle/system/core"
)

func init() {
	core.Use(&core.Command{
		Usage:    []string{"menu", "help", "start", "allmenu"},
		Category: "menu",
		Handler: func(ptz *core.Ptz) error {
			return sendSimpleMenu(ptz)
		},
	})
}

func sendSimpleMenu(ptz *core.Ptz) error {
	body := core.ExtractBody(ptz.Message)
	prefixUsed := ""
	for _, p := range ptz.Bot.Config.Prefixes {
		if strings.HasPrefix(body, p) {
			prefixUsed = p
			break
		}
	}
	if prefixUsed == "" {
		prefixUsed = ptz.Bot.GetPrefix()
	}

	byCat := core.GlobalRegistry().ByCategory()
	categories := make([]string, 0)
	for cat := range byCat {
		if cat != "" && cat != "menu" {
			categories = append(categories, cat)
		}
	}
	sort.Strings(categories)

	// Jika user mengetik .menu [kategori]
	if len(ptz.Args) > 0 {
		categoryName := strings.ToLower(ptz.Args[0])
		cmds, ok := byCat[categoryName]
		if !ok {
			return ptz.ReplyText("❌ Kategori tidak ditemukan.")
		}

		return ptz.ReplyText(formatCategoryCommands(cmds, prefixUsed, categoryName))
	}

	// Tampilan menu utama (semua kategori)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Hai @%s 🪸\n\n", ptz.Sender.User))

	for i, cat := range categories {
		cmds := byCat[cat]
		// Add double newline to force separation
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(formatCategoryCommands(cmds, prefixUsed, cat))
	}

	sb.WriteString("\n\n> © michelle-go v1.1")

	return ptz.ReplyTextMention(sb.String(), []types.JID{ptz.Sender})
}

func formatCategoryCommands(cmds []*core.Command, prefixUsed, categoryName string) string {
	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("– *MENU %s*\n", strings.ToUpper(categoryName)))
	
	// Flatten commands to individual usage lines first
	usageLines := make([]string, 0)
	for _, cmd := range cmds {
		hint := ""
		if cmd.UsageHint != "" {
			hint = fmt.Sprintf(" *%s*", cmd.UsageHint)
		}

		// Periksa setiap usage, masukkan hanya jika tidak ada di Hidden
		for _, usage := range cmd.Usage {
			isHidden := false
			for _, h := range cmd.Hidden {
				if usage == h {
					isHidden = true
					break
				}
			}
			if !isHidden {
				usageLines = append(usageLines, fmt.Sprintf("%s%s%s", prefixUsed, usage, hint))
			}
		}
	}

	for i, usage := range usageLines {
		prefix := "│  ◦  "
		if i == 0 {
			prefix = "┌  ◦  "
		}
		if i == len(usageLines)-1 {
			prefix = "└  ◦  "
		}
		lines = append(lines, fmt.Sprintf("%s%s", prefix, usage))
	}
	
	return strings.Join(lines, "\n")
}
