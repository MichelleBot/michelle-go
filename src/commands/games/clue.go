package games

import (
	"fmt"
	"strings"
	"unicode"
)

var gameTypeNames = map[string]string{
	"caklontong":      "Cak Lontong",
	"tebakbendera":    "Tebak Bendera",
	"tebakkata":       "Tebak Kata",
	"tebaklagu":       "Tebak Lagu",
	"tekateki":        "Teka-Teki",
	"asahotak":        "Asah Otak",
	"siapakahaku":     "Siapakah Aku",
	"susunkata":       "Susun Kata",
	"tebakgambar":     "Tebak Gambar",
	"tebakkimia":      "Tebak Kimia",
	"tebakkalimat":    "Tebak Kalimat",
	"lengkapikalimat": "Lengkapi Kalimat",
	"tebaktebakan":    "Tebak Tebakan",
	"tebaklogo":       "Tebak Logo",
}

func buildClue(answer string, clueCount int) string {
	runes := []rune(answer)
	n := len(runes)

	revealed := make([]bool, n)
	alphaIdx := make([]int, 0, n)

	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			revealed[i] = true
		} else {
			alphaIdx = append(alphaIdx, i)
		}
	}

	totalLetters := len(alphaIdx)
	if totalLetters == 0 {
		return answer
	}

	revealCount := 0
	switch {
	case clueCount == 1:
		revealCount = max(1, totalLetters/4)
	case clueCount == 2:
		revealCount = max(1, totalLetters/2)
	default:
		revealCount = max(1, (totalLetters*3)/4)
	}

	for _, idx := range revealPlan(alphaIdx, revealCount) {
		revealed[idx] = true
	}

	var sb strings.Builder
	for i, r := range runes {
		if revealed[i] {
			if unicode.IsSpace(r) {
				sb.WriteRune(' ')
			} else {
				sb.WriteString(strings.ToUpper(string(r)))
			}
		} else {
			sb.WriteString("_")
		}
		if i < n-1 && unicode.IsLetter(r) && unicode.IsLetter(runes[i+1]) {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}

func revealPlan(alphaIdx []int, revealCount int) []int {
	if revealCount >= len(alphaIdx) {
		return alphaIdx
	}

	plan := make([]int, 0, len(alphaIdx))
	used := make(map[int]struct{}, len(alphaIdx))

	left, right := 0, len(alphaIdx)-1
	for left <= right {
		mid := (left + right) / 2
		for _, pos := range []int{mid, left, right} {
			if pos < 0 || pos >= len(alphaIdx) {
				continue
			}
			if _, ok := used[pos]; ok {
				continue
			}
			used[pos] = struct{}{}
			plan = append(plan, alphaIdx[pos])
		}
		left++
		right--
	}

	if len(plan) > revealCount {
		return plan[:revealCount]
	}
	return plan
}

func clueLabel(clueCount int) string {
	switch clueCount {
	case 1:
		return "💡 *Clue 1* (~25% terbuka)"
	case 2:
		return "💡 *Clue 2* (~50% terbuka)"
	default:
		return "💡 *Clue 3* (~75% terbuka)"
	}
}

func formatClueMessage(sess *Session) string {
	name := gameTypeNames[sess.GameType]
	if name == "" {
		name = sess.GameType
	}
	clue := buildClue(sess.Answer, sess.ClueCount)
	remaining := getRemainingTime(sess.ChatJID, sess.StarterJID, sess.GameType)
	return fmt.Sprintf(
		"%s\n\n`%s`\n\nSisa waktu: %s\nReward jika benar sekarang: +%d balance",
		clueLabel(sess.ClueCount),
		clue,
		remaining,
		RewardForClueCount(sess.ClueCount),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
