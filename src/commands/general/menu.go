package general

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"sptzx/src/core"
	"sptzx/src/utils"
)

var categoryEmoji = map[string]string{
	"general":    "🔧 General",
	"group":      "👥 Group",
	"downloader": "📥 Downloader",
	"sticker":    "🎭 Sticker",
	"owner":      "👑 Owner",
	"search":     "🔍 Search",
	"stalk":      "🕵️ Stalk",
	"tools":      "🛠 Tools",
	"primbon":    "🔮 Primbon",
	"random":     "🎲 Random",
	"maker":      "🎨 Maker",
	"ai":         "🤖 AI",
	"games":      "🎮 Games",
}

func init() {
	core.Use(&core.Command{
		Name:        "menu",
		Aliases:     []string{"help", "start"},
		Description: "Lihat daftar kategori atau command",
		Usage:       "menu | menu all | menu <kategori>",
		Category:    "general",
		Limit:       core.PerUserLimit(30, time.Minute),
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return sendGreetMenu(ptz)
			}
			arg := strings.ToLower(ptz.Args[0])
			if arg == "all" || arg == "semua" {
				return sendFullMenu(ptz)
			}
			return sendCategoryMenu(ptz, arg)
		},
	})
}

func sendGreetMenu(ptz *core.Ptz) error {
	greeting := utils.Greeting(ptz.Bot.Config.Timezone)
	name := ptz.GetSenderName()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%s, %s*\n\n", greeting, name))
	sb.WriteString("*Kategori command yang tersedia:*\n")

	cats := make(map[string]bool)
	for _, cmd := range core.GlobalRegistry().All() {
		cats[cmd.Category] = true
	}

	catList := make([]string, 0, len(cats))
	for cat := range cats {
		catList = append(catList, cat)
	}
	sort.Strings(catList)

	for _, cat := range catList {
		emoji := categoryEmoji[cat]
		if emoji == "" {
			emoji = "📌 " + strings.Title(cat)
		}
		sb.WriteString(fmt.Sprintf("- %s\n", emoji))
	}

	sb.WriteString("\n*Cara pakai:*\n")
	sb.WriteString("- ketik `.menu <kategori>` untuk melihat isi kategori\n")
	sb.WriteString("- ketik `.menu all` untuk melihat semua command")
	return ptz.ReplyText(sb.String())
}

func sendFullMenu(ptz *core.Ptz) error {
	var sb strings.Builder
	sb.WriteString("*Daftar semua command*\n\n")

	byCat := core.GlobalRegistry().ByCategory()
	categories := make([]string, 0, len(byCat))
	for cat := range byCat {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		cmds := byCat[cat]
		emoji := categoryEmoji[cat]
		if emoji == "" {
			emoji = "📌 " + strings.Title(cat)
		}
		sb.WriteString(fmt.Sprintf("*%s*\n", emoji))
		for _, cmd := range cmds {
			sb.WriteString(fmt.Sprintf("- %s\n", cmd.Name))
		}
		sb.WriteString("\n")
	}

	return ptz.ReplyText(sb.String())
}

func sendCategoryMenu(ptz *core.Ptz, category string) error {
	byCat := core.GlobalRegistry().ByCategory()
	cmds, ok := byCat[strings.ToLower(category)]
	if !ok {
		return ptz.ReplyText("❌ Kategori tidak ditemukan.")
	}

	emoji := categoryEmoji[strings.ToLower(category)]
	if emoji == "" {
		emoji = "📌 " + strings.Title(category)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%s*\n\n", emoji))
	for _, cmd := range cmds {
		sb.WriteString(fmt.Sprintf("- %s\n", cmd.Name))
	}
	sb.WriteString("\n")

	return ptz.ReplyText(sb.String())
}
