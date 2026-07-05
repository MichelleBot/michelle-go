package games

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

type AnswerResult int

const (
	AnswerWrong AnswerResult = iota
	AnswerGettingClose
	AnswerVeryClose
	AnswerCorrect
)

type GameDef struct {
	Type     string
	Endpoint string
	Parse    func(raw []byte) (soal, jawaban, imageURL, audioURL string, err error)
	Format   func(soal string) string
}

func playGame(ptz *core.Ptz, def GameDef) error {
	if !IsGameEnabled(ptz.Chat.String()) {
		return ptz.ReplyText("🚫 Fitur game sedang dinonaktifkan di sini.")
	}

	if ptz.IsGroup {
		if sess, ok := GetActiveChatSessionAny(ptz.Chat.String()); ok {
			name := gameTypeNames[sess.GameType]
			remaining := getRemainingTimeFromSess(sess)
			return ptz.ReplyText(fmt.Sprintf(
				"*Masih ada game yang sedang berjalan*\n\n- Game: *%s*\n- Sisa waktu: *%s*\n\n*Cara lanjut:*\n- reply pesan soal untuk menjawab\n- reply pesan soal lalu ketik *%sclue* jika buntu\n- reply pesan soal lalu ketik *%snyerah* jika ingin berhenti",
				name, remaining, ptz.Bot.GetPrefix(), ptz.Bot.GetPrefix(),
			))
		}
	} else {
		if sess, ok := GetActiveUserSessionAny(ptz.Chat.String(), ptz.Sender.String()); ok {
			name := gameTypeNames[sess.GameType]
			remaining := getRemainingTimeFromSess(sess)
			return ptz.ReplyText(fmt.Sprintf(
				"*Kamu masih punya soal aktif*\n\n- Game: *%s*\n- Sisa waktu: *%s*\n\n*Cara lanjut:*\n- reply pesan soal untuk menjawab\n- reply pesan soal lalu ketik *%sclue* jika buntu\n- reply pesan soal lalu ketik *%snyerah* jika ingin berhenti",
				name, remaining, ptz.Bot.GetPrefix(), ptz.Bot.GetPrefix(),
			))
		}
	}

	ptz.React("⏳")
	defer ptz.Unreact()

	client := ptz.Bot.API
	if client == nil {
		client = api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
		client.SetLogger(ptz.Bot.Log)
	}
	raw, err := client.Get(context.Background(), def.Endpoint, nil)
	if err != nil {
		return ptz.ReplyText("❌ Gagal mengambil soal: " + err.Error())
	}

	soal, jawaban, imageURL, audioURL, err := def.Parse(raw)
	if err != nil || soal == "" || jawaban == "" {
		return ptz.ReplyText("❌ Data soal tidak valid, coba lagi.")
	}

	teks := def.Format(soal)
	footer := fmt.Sprintf("\n\n*Aturan main*\n- Waktu: *%d menit*\n- Balas pesan ini dengan jawabanmu\n- Butuh petunjuk: reply lalu ketik *%sclue*\n- Mau menyerah: reply lalu ketik *%snyerah*\n- %s", int(sessionTTL.Minutes()), ptz.Bot.GetPrefix(), ptz.Bot.GetPrefix(), RewardGuide())

	var questionID string

	switch {
	case imageURL != "":
		imgData, err := serialize.Fetch(imageURL)
		if err != nil {
			return ptz.ReplyText("❌ Gagal mengambil gambar soal.")
		}
		msgID, err := ptz.ReplyImageID(imgData, "image/jpeg", teks+footer)
		if err != nil {
			return err
		}
		questionID = msgID

	case audioURL != "":
		msgID, err := ptz.ReplyTextID(teks + footer)
		if err != nil {
			return err
		}
		questionID = msgID
		audioData, err := serialize.Fetch(audioURL)
		if err != nil {
			return ptz.ReplyText("❌ Gagal mengambil audio soal.")
		}
		if err := serialize.SendAudio(ptz.Bot.Client, ptz.Chat, audioData, "audio/mpeg", false); err != nil {
			ptz.Bot.Log.Warnf("tebaklagu: gagal kirim audio: %v", err)
		}

	default:
		msgID, err := ptz.ReplyTextID(teks + footer)
		if err != nil {
			return err
		}
		questionID = msgID
	}

	SetSession(NewSession(ptz.Chat.String(), ptz.Sender.String(), def.Type, jawaban, questionID, ptz.IsGroup))
	return nil
}

func cleanString(s string) string {
	re := regexp.MustCompile(`\([^)]*\)`)
	s = re.ReplaceAllString(s, "")
	re2 := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	s = re2.ReplaceAllString(s, "")
	return strings.ToLower(strings.TrimSpace(strings.Join(strings.Fields(s), " ")))
}

func levenshtein(s, t string) int {
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}
	}
	return d[len(s)][len(t)]
}

func CalculateSimilarity(ans, user string) float64 {
	ans = cleanString(ans)
	user = cleanString(user)

	if ans == user {
		return 100.0
	}

	ansWords := strings.Fields(ans)
	if len(ansWords) > 1 && len(user) >= 4 {
		dist := float64(levenshtein(ansWords[0], user))
		maxL := math.Max(float64(len(ansWords[0])), float64(len(user)))
		wordSim := (1.0 - (dist / maxL)) * 100.0
		if wordSim >= 85.0 {
			return wordSim
		}
	}

	maxLen := math.Max(float64(len(ans)), float64(len(user)))
	if maxLen == 0 {
		return 0.0
	}
	dist := float64(levenshtein(ans, user))
	sim := (1.0 - (dist / maxLen)) * 100.0

	if len(user) >= 4 && strings.Contains(ans, user) {
		if sim < 85.0 {
			sim = 85.0
		}
	}

	return sim
}

func CheckAnswer(questionID, userAnswer string) (AnswerResult, *Session, bool) {
	sess, found := MatchByQuestionID(questionID)
	if !found {
		return AnswerWrong, nil, false
	}

	sim := CalculateSimilarity(sess.Answer, userAnswer)

	if sim >= 82.0 {
		return AnswerCorrect, sess, true
	} else if sim >= 65.0 {
		return AnswerVeryClose, sess, true
	} else if sim >= 45.0 {
		return AnswerGettingClose, sess, true
	}
	return AnswerWrong, sess, true
}

func getRemainingTimeFromSess(sess *Session) string {
	remaining := time.Until(sess.ExpiresAt)
	if remaining <= 0 {
		return "0 detik"
	}
	mins := int(remaining.Minutes())
	secs := int(remaining.Seconds()) % 60
	if mins > 0 {
		return fmt.Sprintf("%dm %ds", mins, secs)
	}
	return fmt.Sprintf("%ds", secs)
}

func getRemainingTime(chatJID, senderJID, gameType string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	key := userSessionKey(chatJID, senderJID, gameType)
	sess, ok := db.mem[key]
	if !ok {
		return "0 detik"
	}
	return getRemainingTimeFromSess(sess)
}
