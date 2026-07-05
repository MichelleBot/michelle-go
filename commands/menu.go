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

	// Jika user mengetik .menu [kategori]
	if len(ptz.Args) > 0 {
		categoryName := strings.ToLower(ptz.Args[0])
		cmds, ok := byCat[categoryName]
		if !ok {
			return ptz.ReplyText("❌ Kategori tidak ditemukan.")
		}

		filteredCmds := make([]*core.Command, 0)
		for _, cmd := range cmds {
			isHidden := false
			for _, h := range cmd.Hidden {
				if h == "run" {
					isHidden = true
					break
				}
			}
			if !isHidden {
				filteredCmds = append(filteredCmds, cmd)
			}
		}

		lines := make([]string, 0, len(filteredCmds))
		count := len(filteredCmds)
		for i, cmd := range filteredCmds {
			prefix := "│  ◦  "
			if i == 0 {
				prefix = "┌  ◦  "
			} else if i == count-1 {
				prefix = "└  ◦  "
			}
			lines = append(lines, fmt.Sprintf("%s%s%s", prefix, prefixUsed, cmd.Usage[0]))
		}
		return ptz.ReplyText(strings.Join(lines, "\n"))
	}

	// Tampilan menu utama
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Hai @%s 🪸\n\n", ptz.Sender.User))
	sb.WriteString("◦ *Module* : 2.0.0\n")
	sb.WriteString("◦ *Database* : SQLite\n")
	sb.WriteString("◦ *Libray* : Whatsmeow-v1.1\n\n")

	categories := make([]string, 0)
	for cat := range byCat {
		if cat != "" && cat != "menu" {
			categories = append(categories, cat)
		}
	}
	sort.Strings(categories)

	lines := make([]string, 0, len(categories))
	count := len(categories)
	for i, cat := range categories {
		prefix := "│  ◦  "
		if i == 0 {
			prefix = "┌  ◦  "
		} else if i == count-1 {
			prefix = "└  ◦  "
		}
		lines = append(lines, fmt.Sprintf("%s%smenu %s", prefix, prefixUsed, cat))
	}
	sb.WriteString(strings.Join(lines, "\n"))

	sb.WriteString("\n\n> *Simple Whatsapp bot michelle*")

	return ptz.ReplyTextMention(sb.String(), []types.JID{ptz.Sender})
}
