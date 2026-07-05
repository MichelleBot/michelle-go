package commands

import (
	"fmt"
	"sort"
	"strings"

	"michelle/system/core"
)

func init() {
	core.Use(&core.Command{
		Usage:    []string{"hitstat", "hitdaily"},
		Category: "miscs",
		Handler:  runHitStat,
	})
}

func runHitStat(ptz *core.Ptz) error {
	isDaily := ptz.Command == "hitdaily"

	// Deteksi prefix yang digunakan user
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

	hitStats := core.GetStats(ptz.Bot)

	type entry struct {
		cmd  string
		stat core.CmdStat
	}

	entries := make([]entry, 0, len(hitStats))
	for cmd, stat := range hitStats {
		entries = append(entries, entry{cmd, stat})
	}

	if len(entries) == 0 {
		return ptz.ReplyText("🚩 Tidak ada perintah yang digunakan.")
	}

	sort.Slice(entries, func(i, j int) bool {
		if isDaily {
			return entries[i].stat.TodayHit > entries[j].stat.TodayHit
		}
		return entries[i].stat.TotalHit > entries[j].stat.TotalHit
	})

	var sb strings.Builder
	sb.WriteString("乂  *H I T S T A T*\n\n")

	total := 0
	for _, e := range entries {
		if isDaily {
			total += e.stat.TodayHit
		} else {
			total += e.stat.TotalHit
		}
	}

	sb.WriteString(fmt.Sprintf("*Statistik total hit perintah %s %d berhasil.*\n\n",
		map[bool]string{true: "hari ini", false: "saat ini"}[isDaily], total))

	show := 10
	if len(entries) < show {
		show = len(entries)
	}

	for i := 0; i < show; i++ {
		e := entries[i]
		hit := e.stat.TotalHit
		if isDaily {
			hit = e.stat.TodayHit
		}
		sb.WriteString(fmt.Sprintf("   ┌ *Perintah* : %s%s\n", prefixUsed, e.cmd))
		sb.WriteString(fmt.Sprintf("   │ *Hit* : %d x\n", hit))
		sb.WriteString(fmt.Sprintf("   └ *Hit Terakhir* : %s\n\n", e.stat.LastHit.Format("02/01/06 15:04:05")))
	}

	return ptz.ReplyText(sb.String())
}
