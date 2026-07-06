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

func Example(prefix string, command string, arg string) string {
	return fmt.Sprintf("• *Example* : %s%s %s", prefix, command, arg)
}

func IsBot(id string) bool {
	return id != "" && (strings.HasPrefix(id, "3EB0") || strings.HasPrefix(id, "BAE") || strings.Contains(id, "-") || strings.Contains(strings.ToLower(id), "neoxr"))
}

func DetectBadword(input string, badwords []string) bool {
	mapSimilar := map[string]string{
		"i": "l1", "l": "i1", "1": "il", "o": "0", "0": "o", "a": "4",
		"4": "a", "e": "3", "3": "e", "b": "8", "8": "b", "s": "5",
		"5": "s", "t": "7", "7": "t", "g": "9", "9": "g",
	}

	normalize := func(str string) string {
		str = strings.ToLower(str)
		var res strings.Builder
		for _, c := range str {
			char := string(c)
			found := false
			for key, val := range mapSimilar {
				if key == char || strings.Contains(val, char) {
					res.WriteString(key)
					found = true
					break
				}
			}
			if !found {
				res.WriteString(char)
			}
		}
		return res.String()
	}

	normalizedInput := normalize(input)
	for _, word := range badwords {
		if strings.Contains(normalizedInput, normalize(word)) {
			return true
		}
	}
	return false
}

func Styles(text string) string {
	mapping := map[rune]rune{
		'a': 'ᴀ', 'b': 'ʙ', 'c': 'ᴄ', 'd': 'ᴅ', 'e': 'ᴇ', 'f': 'ꜰ', 'g': 'ɢ', 'h': 'ʜ', 'i': 'ɪ',
		'j': 'ᴊ', 'k': 'ᴋ', 'l': 'ʟ', 'm': 'ᴍ', 'n': 'ɴ', 'o': 'ᴏ', 'p': 'ᴘ', 'q': 'ǫ', 'r': 'ʀ',
		's': 'ꜱ', 't': 'ᴛ', 'u': 'ᴜ', 'v': 'ᴠ', 'w': 'ᴡ', 'x': 'x', 'y': 'ʏ', 'z': 'ᴢ',
	}
	var res strings.Builder
	for _, r := range strings.ToLower(text) {
		if val, ok := mapping[r]; ok {
			res.WriteRune(val)
		} else {
			res.WriteRune(r)
		}
	}
	return res.String()
}

func Texted(style string, text string) string {
	switch style {
	case "bold":
		return "*" + text + "*"
	case "italic":
		return "_" + text + "_"
	case "monospace":
		return "```" + text + "```"
	default:
		return text
	}
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
