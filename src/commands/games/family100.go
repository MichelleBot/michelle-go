package games

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow/types"
	"sptzx/src/api"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

type Family100Question struct {
	Soal    string   `json:"soal"`
	Jawaban []string `json:"jawaban"`
}

type Family100Session struct {
	ChatJID    string
	Question   string
	Answers    []string
	Revealed   []bool
	QuestionID string
	ExpiresAt  time.Time
}

type Family100DB struct {
	mu  sync.RWMutex
	mem map[string]*Family100Session
}

var f100Store = &Family100DB{
	mem: make(map[string]*Family100Session),
}

func init() {
	core.Use(&core.Command{
		Name:        "family100",
		Aliases:     []string{"f100"},
		Description: "Game Family 100",
		Usage:       "family100",
		Category:    "games",
		GroupOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			if !IsGameEnabled(ptz.Chat.String()) {
				return ptz.ReplyText("🚫 Fitur game sedang dinonaktifkan di sini.")
			}

			f100Store.mu.RLock()
			_, playing := f100Store.mem[ptz.Chat.String()]
			f100Store.mu.RUnlock()
			if playing {
				return ptz.ReplyText("⏳ Grup ini sedang bermain Family 100!")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			client := api.NewClient(ptz.Bot.Config.SiputzX.BaseURL)
			raw, err := client.Get(context.Background(), "/api/games/family100", nil)
			if err != nil {
				return ptz.ReplyText("❌ Gagal mengambil soal.")
			}

			var apiResp struct {
				Status bool              `json:"status"`
				Data   Family100Question `json:"data"`
			}
			if err := json.Unmarshal(raw, &apiResp); err != nil {
				return ptz.ReplyText("❌ Gagal memproses soal.")
			}

			sess := &Family100Session{
				ChatJID:   ptz.Chat.String(),
				Question:  apiResp.Data.Soal,
				Answers:   apiResp.Data.Jawaban,
				Revealed:  make([]bool, len(apiResp.Data.Jawaban)),
				ExpiresAt: time.Now().Add(5 * time.Minute),
			}

			f100Store.mu.Lock()
			f100Store.mem[ptz.Chat.String()] = sess
			f100Store.mu.Unlock()

			return sendFamily100Question(ptz, sess)
		},
	})
}

func sendFamily100Question(ptz *core.Ptz, sess *Family100Session) error {
	var sb strings.Builder
	sb.WriteString("👪 *Family 100*\n\n")
	sb.WriteString(sess.Question + "\n\n")

	for i, ans := range sess.Answers {
		if sess.Revealed[i] {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, ans))
		} else {
			sb.WriteString(fmt.Sprintf("%d. *[TERSEMBUNYI]*\n", i+1))
		}
	}

	sb.WriteString("\n_Tebak jawabannya!_")

	msgID, err := ptz.ReplyTextID(sb.String())
	if err != nil {
		return err
	}
	sess.QuestionID = msgID
	return nil
}

func ProcessFamily100Answer(ptz *core.Ptz, msgID, userAnswer string) (bool, error) {
	f100Store.mu.Lock()
	sess, found := f100Store.mem[ptz.Chat.String()]
	f100Store.mu.Unlock()

	if !found || sess.QuestionID != msgID {
		return false, nil
	}

	bestSim := 0.0
	bestIdx := -1

	for i, correct := range sess.Answers {
		if sess.Revealed[i] {
			continue
		}
		sim := CalculateSimilarity(correct, userAnswer)
		if sim > bestSim {
			bestSim = sim
			bestIdx = i
		}
	}

	if bestSim >= 82.0 {
		sess.Revealed[bestIdx] = true

		senderJID, _ := types.ParseJID(ptz.Sender.String())

		text := fmt.Sprintf(
			"✅ *Benar!* Selamat @%s!\n\n🎮 Jawaban ditemukan: *%s*",
			senderJID.User, sess.Answers[bestIdx],
		)

		serialize.SendTextMentionReplyToID(
			ptz.Bot.Client,
			ptz.Chat,
			text,
			[]types.JID{senderJID},
			sess.QuestionID,
			senderJID.String(),
		)

		allRevealed := true
		for _, revealed := range sess.Revealed {
			if !revealed {
				allRevealed = false
				break
			}
		}

		if allRevealed {
			ptz.ReplyText("🎉 *Selamat! Semua jawaban telah ditemukan!*")
			f100Store.mu.Lock()
			delete(f100Store.mem, sess.ChatJID)
			f100Store.mu.Unlock()
			return true, nil
		}

		return true, sendFamily100Question(ptz, sess)
	} else if bestSim >= 65.0 {
		ptz.ReplyText("⚠️ *Dikit lagi bener!* Jawaban kamu udah semakin mirip.")
		return true, nil
	}

	ptz.ReplyText("❌ *Salah!* Coba lagi.")
	return true, nil
}
