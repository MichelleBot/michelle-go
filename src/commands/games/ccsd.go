package games

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow/types"
	"sptzx/src/api"
	"sptzx/src/core"
)

type CcsdQuestion struct {
	Pertanyaan   string
	Pilihan      []map[string]string
	JawabanBenar string
}

type CcsdSession struct {
	ChatJID       string
	SenderJID     string
	MataPelajaran string
	Questions     []CcsdQuestion
	CurrentIdx    int
	Answers       []string
	Score         int
	ExpiresAt     time.Time
	QuestionID    string
}

type CcsdDB struct {
	mu  sync.RWMutex
	mem map[string]*CcsdSession
}

var ccsdStore = &CcsdDB{
	mem: make(map[string]*CcsdSession),
}

func getCcsdKey(chatJID, senderJID string) string {
	return chatJID + "_" + senderJID
}

func SetCcsdSession(sess *CcsdSession) {
	ccsdStore.mu.Lock()
	defer ccsdStore.mu.Unlock()
	key := getCcsdKey(sess.ChatJID, sess.SenderJID)
	ccsdStore.mem[key] = sess
}

func GetCcsdSessionByMsgID(msgID string) (*CcsdSession, bool) {
	ccsdStore.mu.RLock()
	defer ccsdStore.mu.RUnlock()
	for _, sess := range ccsdStore.mem {
		if sess.QuestionID == msgID {
			return sess, true
		}
	}
	return nil, false
}

func DeleteCcsdSession(chatJID, senderJID string) {
	ccsdStore.mu.Lock()
	defer ccsdStore.mu.Unlock()
	key := getCcsdKey(chatJID, senderJID)
	delete(ccsdStore.mem, key)
}

type CcsdApiResp struct {
	Status bool `json:"status"`
	Data   struct {
		MataPelajaran string `json:"matapelajaran"`
		JumlahSoal    int    `json:"jumlah_soal"`
		Soal          []struct {
			Pertanyaan   string              `json:"pertanyaan"`
			SemuaJawaban []map[string]string `json:"semua_jawaban"`
			JawabanBenar string              `json:"jawaban_benar"`
		} `json:"soal"`
	} `json:"data"`
}

func init() {
	core.Use(&core.Command{
		Name:        "ccsd",
		Aliases:     []string{"cerdascermat"},
		Description: "Ujian Cerdas Cermat SD",
		Usage:       "ccsd <matapelajaran>",
		Category:    "games",
		Handler:     ccsdHandler,
	})
}

func ccsdHandler(ptz *core.Ptz) error {
	if !IsGameEnabled(ptz.Chat.String()) {
		return ptz.ReplyText("🚫 Fitur game sedang dinonaktifkan di sini.")
	}

	args := strings.TrimSpace(ptz.RawArgs)
	validSubjects := []string{"bindo", "tik", "pkn", "bing", "penjas", "pai", "matematika", "jawa", "ips", "ipa"}

	if args == "" {
		return ptz.ReplyText("📚 *Cerdas Cermat SD*\n\nSilakan pilih mata pelajaran:\n" + strings.Join(validSubjects, ", ") + "\n\nContoh: .ccsd matematika")
	}

	isValid := false
	for _, s := range validSubjects {
		if strings.EqualFold(s, args) {
			isValid = true
			args = strings.ToLower(s)
			break
		}
	}

	if !isValid {
		return ptz.ReplyText("❌ Mata pelajaran tidak valid. Pilih dari: " + strings.Join(validSubjects, ", "))
	}

	ccsdStore.mu.RLock()
	_, playing := ccsdStore.mem[getCcsdKey(ptz.Chat.String(), ptz.Sender.String())]
	ccsdStore.mu.RUnlock()
	if playing {
		return ptz.ReplyText("⏳ Kamu masih memiliki ujian Cerdas Cermat yang belum selesai!")
	}

	ptz.React("⏳")
	defer ptz.Unreact()

	client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
	endpoint := fmt.Sprintf("/api/games/cc-sd?matapelajaran=%s&jumlahsoal=5", args)
	raw, err := client.Get(context.Background(), endpoint, nil)
	if err != nil {
		return ptz.ReplyText("❌ Gagal mengambil soal: " + err.Error())
	}

	var apiResp CcsdApiResp
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return ptz.ReplyText("❌ Gagal memproses soal dari server.")
	}

	if len(apiResp.Data.Soal) == 0 {
		return ptz.ReplyText("❌ Soal tidak tersedia untuk mata pelajaran tersebut.")
	}

	sess := &CcsdSession{
		ChatJID:       ptz.Chat.String(),
		SenderJID:     ptz.Sender.String(),
		MataPelajaran: apiResp.Data.MataPelajaran,
		Questions:     make([]CcsdQuestion, 0, len(apiResp.Data.Soal)),
		CurrentIdx:    0,
		Answers:       make([]string, len(apiResp.Data.Soal)),
		ExpiresAt:     time.Now().Add(5 * time.Minute),
	}

	for _, s := range apiResp.Data.Soal {
		sess.Questions = append(sess.Questions, CcsdQuestion{
			Pertanyaan:   s.Pertanyaan,
			Pilihan:      s.SemuaJawaban,
			JawabanBenar: s.JawabanBenar,
		})
	}

	SetCcsdSession(sess)

	return sendCcsdQuestion(ptz, sess)
}

func sendCcsdQuestion(ptz *core.Ptz, sess *CcsdSession) error {
	if sess.CurrentIdx >= len(sess.Questions) {
		return sendCcsdReport(ptz, sess)
	}

	q := sess.Questions[sess.CurrentIdx]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🎓 *Cerdas Cermat (%s)*\n", strings.ToUpper(sess.MataPelajaran)))
	sb.WriteString(fmt.Sprintf("Soal %d dari %d\n\n", sess.CurrentIdx+1, len(sess.Questions)))
	sb.WriteString(q.Pertanyaan + "\n\n")

	var validKeys []string
	for _, p := range q.Pilihan {
		for k, v := range p {
			validKeys = append(validKeys, strings.ToUpper(k))
			sb.WriteString(fmt.Sprintf("%s. %s\n", strings.ToUpper(k), v))
		}
	}
	sort.Strings(validKeys)

	sb.WriteString(fmt.Sprintf("\n_Balas pesan ini dengan jawabanmu (%s)_", strings.Join(validKeys, "/")))

	msgID, err := ptz.ReplyTextID(sb.String())
	if err != nil {
		return err
	}

	sess.QuestionID = msgID
	return nil
}

func sendCcsdReport(ptz *core.Ptz, sess *CcsdSession) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 *RAPORT CERDAS CERMAT (%s)* 📊\n", strings.ToUpper(sess.MataPelajaran)))

	correctCount := 0
	for i, q := range sess.Questions {
		userAns := strings.TrimSpace(strings.ToLower(sess.Answers[i]))
		if userAns == getCorrectKey(q) {
			correctCount++
		}
	}

	score := (correctCount * 100) / len(sess.Questions)
	senderNum := strings.Split(sess.SenderJID, "@")[0]
	sb.WriteString(fmt.Sprintf("👤 Pemain: @%s\n", senderNum))
	sb.WriteString(fmt.Sprintf("🎯 Skor Akhir: %d/100\n", score))
	sb.WriteString(fmt.Sprintf("✅ Benar: %d | ❌ Salah: %d\n\n", correctCount, len(sess.Questions)-correctCount))

	sb.WriteString("📝 *Rincian Jawaban:*\n")
	for i, q := range sess.Questions {
		userAns := strings.TrimSpace(strings.ToLower(sess.Answers[i]))
		if userAns == "" {
			userAns = "-"
		}

		correctKey := getCorrectKey(q)
		correctText := ""
		for _, p := range q.Pilihan {
			if val, ok := p[correctKey]; ok {
				correctText = val
				break
			}
		}

		if userAns == correctKey {
			sb.WriteString(fmt.Sprintf("%d. ✅ Benar (Jawaban: %s. %s)\n", i+1, strings.ToUpper(correctKey), correctText))
		} else {
			sb.WriteString(fmt.Sprintf("%d. ❌ Salah (Kamu: %s, Seharusnya: %s. %s)\n", i+1, strings.ToUpper(userAns), strings.ToUpper(correctKey), correctText))
		}
	}

	DeleteCcsdSession(sess.ChatJID, sess.SenderJID)

	senderJIDObj, err := types.ParseJID(sess.SenderJID)
	if err == nil {
		return ptz.ReplyTextMention(sb.String(), []types.JID{senderJIDObj})
	}
	return ptz.ReplyText(sb.String())
}

func getCorrectKey(q CcsdQuestion) string {
	for _, p := range q.Pilihan {
		for k := range p {
			if strings.EqualFold(k, q.JawabanBenar) {
				return strings.ToLower(k)
			}
		}
	}
	for _, p := range q.Pilihan {
		for k, v := range p {
			if strings.EqualFold(v, q.JawabanBenar) {
				return strings.ToLower(k)
			}
		}
	}
	return strings.ToLower(q.JawabanBenar)
}

func ProcessCcsdAnswer(ptz *core.Ptz, msgID, userAnswer string) (bool, error) {
	sess, found := GetCcsdSessionByMsgID(msgID)
	if !found {
		return false, nil
	}

	ans := strings.TrimSpace(userAnswer)
	if len(ans) == 0 {
		return true, nil
	}

	ansChar := strings.ToLower(string(ans[0]))

	q := sess.Questions[sess.CurrentIdx]
	isValidChoice := false
	var validKeys []string

	for _, p := range q.Pilihan {
		for k := range p {
			validKeys = append(validKeys, strings.ToUpper(k))
			if k == ansChar {
				isValidChoice = true
			}
		}
	}

	sort.Strings(validKeys)

	if !isValidChoice {
		return true, ptz.ReplyText(fmt.Sprintf("❌ Jawaban tidak valid. Silakan balas dengan opsi: %s", strings.Join(validKeys, "/")))
	}

	sess.Answers[sess.CurrentIdx] = ansChar
	sess.CurrentIdx++

	err := sendCcsdQuestion(ptz, sess)
	return true, err
}
