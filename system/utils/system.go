package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var StartTime = time.Now()

func UcWord(str string) string {
	return strings.Title(str)
}

func FetchAsJSON(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var target map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&target)
	return target, err
}

func FmtUptime(d time.Duration) string {
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
}

func RssMemMB() float64 {
	data, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			var kb uint64
			fmt.Sscanf(strings.TrimPrefix(line, "VmRSS:"), " %d", &kb)
			return float64(kb) / 1024
		}
	}
	return 0
}

func CpuPercent() string {
	data, err := os.ReadFile("/proc/self/stat")
	if err != nil {
		return "N/A"
	}
	fields := strings.Fields(string(data))
	if len(fields) < 15 {
		return "N/A"
	}
	var utime, stime uint64
	fmt.Sscanf(fields[13], "%d", &utime)
	fmt.Sscanf(fields[14], "%d", &stime)

	uptimeSecs := time.Since(StartTime).Seconds()
	if uptimeSecs <= 0 {
		return "0.00%"
	}
	const clkTck = 100.0
	usage := (float64(utime+stime) / clkTck) / uptimeSecs * 100
	max := float64(runtime.NumCPU() * 100)
	if usage > max {
		usage = max
	}
	return fmt.Sprintf("%.2f%%", usage)
}

func DiskGB(path string) (total, free, used float64) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, 0
	}
	gb := func(blocks uint64) float64 {
		return float64(blocks) * float64(stat.Bsize) / 1024 / 1024 / 1024
	}
	total = gb(stat.Blocks)
	free = gb(stat.Bfree)
	used = total - free
	return
}
