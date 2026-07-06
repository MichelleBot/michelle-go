package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"michelle/system/core"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:    []string{"server"},
		Category: "miscs",
		Handler: func(ptz *core.Ptz) error {
			// IP info
			ipData, _ := utils.FetchAsJSON("http://ip-api.com/json")
			if ipData != nil {
				delete(ipData, "status")
				delete(ipData, "query")
			}

			// System stats
			v, _ := mem.VirtualMemory()
			d, _ := disk.Usage("/")
			c, _ := cpu.Info()
			
			cpuName := "Unknown"
			if len(c) > 0 {
				cpuName = c[0].ModelName
			}
			
			cwd, _ := os.Getwd()

			lines := []string{
				fmt.Sprintf("┌  ◦  Directory : %s", cwd),
				fmt.Sprintf("│  ◦  OS : %s (%s / %s)", runtime.GOOS, runtime.GOARCH, runtime.Version()),
				fmt.Sprintf("│  ◦  Process : %d", os.Getpid()),
				fmt.Sprintf("│  ◦  Core : %d", runtime.NumCPU()),
				fmt.Sprintf("│  ◦  RAM : %.2fGB / %.2fGB (%.2f%%)", float64(v.Used)/1024/1024/1024, float64(v.Total)/1024/1024/1024, v.UsedPercent),
				fmt.Sprintf("│  ◦  Disk : %.2fGB / %.2fGB (%.2f%%)", float64(d.Used)/1024/1024/1024, float64(d.Total)/1024/1024/1024, d.UsedPercent),
			}

			if ipData != nil {
				for k, v := range ipData {
					lines = append(lines, fmt.Sprintf("│  ◦  %s : %v", utils.UcWord(k), v))
				}
			}
			
			lines = append(lines, fmt.Sprintf("│  ◦  Uptime : %s", utils.FmtUptime(time.Since(utils.StartTime))))
			lines = append(lines, fmt.Sprintf("└  ◦  Processor : %s", cpuName))

			return ptz.ReplyText(strings.Join(lines, "\n"))
		},
	})
}
