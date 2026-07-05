package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

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

			// Memory stats
			mem := getMemoryStats()
			
			cwd, _ := os.Getwd()

			lines := []string{
				"乂  *S E R V E R*",
				"",
				fmt.Sprintf("┌  ◦  Directory : %s", cwd),
				fmt.Sprintf("│  ◦  OS : %s (%s / %s)", runtime.GOOS, runtime.GOARCH, runtime.Version()),
				fmt.Sprintf("│  ◦  Process : %d", os.Getpid()),
				fmt.Sprintf("│  ◦  Core : %d", runtime.NumCPU()),
				fmt.Sprintf("│  ◦  Heap Total : %s", mem["heapTotal"]),
				fmt.Sprintf("│  ◦  Heap Used : %s", mem["heapUsed"]),
				fmt.Sprintf("│  ◦  External : %s", mem["external"]),
				fmt.Sprintf("│  ◦  Array Buffers : %s", mem["arrayBuffers"]),
			}

			if ipData != nil {
				for k, v := range ipData {
					lines = append(lines, fmt.Sprintf("│  ◦  %s : %v", utils.UcWord(k), v))
				}
			}
			
			lines = append(lines, fmt.Sprintf("│  ◦  Uptime : %s", utils.FmtUptime(time.Since(utils.StartTime))))
			lines = append(lines, "└  ◦  Processor : (System Default)")

			return ptz.ReplyText(strings.Join(lines, "\n"))
		},
	})
}

func getMemoryStats() map[string]string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	format := func(b uint64) string {
		return fmt.Sprintf("%.2f MB", float64(b)/1024/1024)
	}

	return map[string]string{
		"heapTotal":    format(mem.HeapSys),
		"heapUsed":     format(mem.HeapAlloc),
		"external":     format(mem.OtherSys), // Best approximation in Go
		"arrayBuffers": "N/A",                 // Go doesn't expose this directly
	}
}
