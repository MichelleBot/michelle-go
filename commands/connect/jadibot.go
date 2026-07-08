package connect

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"michelle/system/core"
	"michelle/system/handler"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"jadibot"},
		UsageHint: "628xxxx",
		Category:  "Auth",
		Handler:   runJadibot,
	})
}

func runJadibot(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText("🚩 Masukkan nomor telepon bot.\nContoh: .jadibot 628xxxx")
	}

	// Add loading reaction
	ptz.React("🕒")
	ptz.Bot.Log.Infof("DEBUG: Jadibot response for command %s via client pointer %p", ptz.Command, ptz.Client)

	phone := strings.ReplaceAll(ptz.Args[0], "+", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Basic validation
	if len(phone) < 8 {
		return ptz.ReplyText("❌ Nomor telepon tidak valid.")
	}

	// Register in database
	_, err := ptz.Bot.DB.Conn.Exec("INSERT OR REPLACE INTO jadibot_sessions (phone, owner_jid, created_at) VALUES (?, ?, ?)", phone, ptz.Sender.String(), time.Now().Unix())
	if err != nil {
		return ptz.ReplyText(fmt.Sprintf("❌ Gagal mendaftarkan sesi ke database: %v", err))
	}

	sessionPath := filepath.Join("sessions", fmt.Sprintf("%s.db", phone))
	
	// Create container for the new session
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:"+sessionPath+"?_foreign_keys=on", waLog.Noop)
	if err != nil {
		return ptz.ReplyText(fmt.Sprintf("❌ Gagal membuat sesi: %v", err))
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return ptz.ReplyText(fmt.Sprintf("❌ Gagal mendapatkan device: %v", err))
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	
	// Create a new bot instance for jadibot
	jadibot := *ptz.Bot
	jadibot.Client = client

	// Register event handler for jadibot using the jadibot instance
	evtHandler := handler.NewEventHandler(&jadibot, ptz.Bot.Registry)
	client.AddEventHandler(evtHandler.Handle)

	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			ptz.ReplyText("✅ Bot berhasil terhubung!")
		}
	})

	// Trigger PairCode login
	go func() {
		// Connect the client first
		if err := client.Connect(); err != nil {
			ptz.Bot.Log.Errorf("Jadibot %s gagal connect: %v", phone, err)
			return
		}

		// We need the code to send to user, but current LoginPairCode prints it.
		// Let's manually trigger the pair process to get the code.
		code, err := client.PairPhone(
			context.Background(),
			phone,
			true,
			whatsmeow.PairClientChrome,
			"Chrome (Linux)",
		)
		if err != nil {
			ptz.Bot.Log.Errorf("Jadibot %s gagal pairing: %v", phone, err)
			return
		}

		// Send stylized message
		msg := fmt.Sprintf("乂  *L O G I N*\n\n1. Di layar utama WhatsApp, ketuk *( ⋮ )* dan pilih *Perangkat Tertaut*.\n2. Ketuk \"Hubungkan dengan nomor telepon saja\"\n3. Masukan kode ini: *%s*\n4. Kode ini akan kedaluwarsa dalam 1 menit.\n\nʟɪɢʜᴛᴡᴇɪɢʜᴛ ᴡᴀʙᴏᴛ ᴍᴀᴅᴇ ʙʏ ᴘʜᴀᴍ ッ", code)
		ptz.ReplyText(msg)

		// Wait for success
		done := make(chan struct{})
		var once sync.Once
		id := client.AddEventHandler(func(evt interface{}) {
			if _, ok := evt.(*events.PairSuccess); ok {
				once.Do(func() {
					close(done)
				})
			}
		})
		defer client.RemoveEventHandler(id)

		select {
		case <-done:
			ptz.ReplyText("✅ Bot berhasil terhubung!")
		case <-time.After(3 * time.Minute):
			ptz.ReplyText("❌ Pairing timeout.")
		}
	}()

	return ptz.ReplyText(fmt.Sprintf("✅ Sesi jadibot %s telah dimulai. Silakan cek chat untuk kode pairing.", phone))
}
