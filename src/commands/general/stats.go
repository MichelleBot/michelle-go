package general

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"sptzx/src/core"
	"sptzx/src/utils"
)

func init() {
	core.Use(&core.Command{
		Name:        "stats",
		Aliases:     []string{"status", "system"},
		Description: "Info sistem & runtime bot secara detail",
		Usage:       "stats",
		Category:    "general",
		Handler: func(ptz *core.Ptz) error {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			uptime := time.Since(utils.StartTime)

			hostname, _ := os.Hostname()

			heapAlloc := float64(mem.HeapAlloc) / 1024 / 1024
			heapSys := float64(mem.HeapSys) / 1024 / 1024
			heapIdle := float64(mem.HeapIdle) / 1024 / 1024
			heapInuse := float64(mem.HeapInuse) / 1024 / 1024
			stackInuse := float64(mem.StackInuse) / 1024 / 1024
			sysMem := float64(mem.Sys) / 1024 / 1024
			totalAlloc := float64(mem.TotalAlloc) / 1024 / 1024

			rss := utils.RssMemMB()
			cpuUsage := utils.CpuPercent()
			diskTotal, diskFree, diskUsed := utils.DiskGB("/")

			var sb strings.Builder
			sb.WriteString("📊 *System Stats*\n")
			sb.WriteString("─────────────────\n\n")

			sb.WriteString("🖥 *Server*\n")
			sb.WriteString(fmt.Sprintf("  Hostname  : `%s`\n", hostname))
			sb.WriteString(fmt.Sprintf("  OS        : %s/%s\n", runtime.GOOS, runtime.GOARCH))
			sb.WriteString(fmt.Sprintf("  Uptime    : %s\n\n", utils.FmtUptime(uptime)))

			sb.WriteString("⚙️ *CPU*\n")
			sb.WriteString(fmt.Sprintf("  Cores     : %d vCPU\n", runtime.NumCPU()))
			sb.WriteString(fmt.Sprintf("  Usage     : %s\n\n", cpuUsage))

			sb.WriteString("🧠 *Memory (RAM)*\n")
			sb.WriteString(fmt.Sprintf("  RSS       : %.2f MB\n", rss))
			sb.WriteString(fmt.Sprintf("  Sys Total : %.2f MB\n", sysMem))
			sb.WriteString(fmt.Sprintf("  Total Alc : %.2f MB\n\n", totalAlloc))

			sb.WriteString("💾 *Go Heap*\n")
			sb.WriteString(fmt.Sprintf("  Alloc     : %.2f MB\n", heapAlloc))
			sb.WriteString(fmt.Sprintf("  In Use    : %.2f MB\n", heapInuse))
			sb.WriteString(fmt.Sprintf("  Idle      : %.2f MB\n", heapIdle))
			sb.WriteString(fmt.Sprintf("  Reserved  : %.2f MB\n\n", heapSys))

			sb.WriteString("📦 *Stack & GC*\n")
			sb.WriteString(fmt.Sprintf("  Stack     : %.2f MB\n", stackInuse))
			sb.WriteString(fmt.Sprintf("  GC Runs   : %d\n", mem.NumGC))
			sb.WriteString(fmt.Sprintf("  Next GC   : %.2f MB\n", float64(mem.NextGC)/1024/1024))
			sb.WriteString(fmt.Sprintf("  GC CPU    : %.4f%%\n\n", mem.GCCPUFraction*100))

			sb.WriteString("⚡ *Runtime*\n")
			sb.WriteString(fmt.Sprintf("  Go Ver    : %s\n", runtime.Version()))
			sb.WriteString(fmt.Sprintf("  Goroutine : %d\n", runtime.NumGoroutine()))
			sb.WriteString(fmt.Sprintf("  CGO Calls : %d\n\n", runtime.NumCgoCall()))

			sb.WriteString("💿 *Disk (/)*\n")
			sb.WriteString(fmt.Sprintf("  Total     : %.1f GB\n", diskTotal))
			sb.WriteString(fmt.Sprintf("  Used      : %.1f GB\n", diskUsed))
			sb.WriteString(fmt.Sprintf("  Free      : %.1f GB\n\n", diskFree))

			sb.WriteString(fmt.Sprintf("_%s v%s_", core.SptzxProject, core.SptzxVersion))

			return ptz.ReplyText(sb.String())
		},
	})
}
