package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "sptzx/src/commands/ai"
	_ "sptzx/src/commands/downloader"
	"sptzx/src/commands/games"
	_ "sptzx/src/commands/general"
	_ "sptzx/src/commands/group"
	_ "sptzx/src/commands/maker"
	_ "sptzx/src/commands/owner"
	_ "sptzx/src/commands/primbon"
	_ "sptzx/src/commands/random"
	_ "sptzx/src/commands/search"
	_ "sptzx/src/commands/stalk"
	_ "sptzx/src/commands/sticker"
	_ "sptzx/src/commands/tools"

	"sptzx/src/config"
	"sptzx/src/core"
	"sptzx/src/handler"

	"time"

	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	waCompanionReg "go.mau.fi/whatsmeow/proto/waCompanionReg"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	whatsmeow.KeepAliveIntervalMin = 20 * time.Second
	whatsmeow.KeepAliveIntervalMax = 30 * time.Second
	whatsmeow.KeepAliveResponseDeadline = 15 * time.Second
	whatsmeow.KeepAliveMaxFailTime = 3 * time.Minute
}

func configureClient(client *whatsmeow.Client) {
	store.SetOSInfo("Chrome", [3]uint32{131, 0, 0})

	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_CHROME.Enum()
	store.DeviceProps.RequireFullSync = proto.Bool(false)
	store.DeviceProps.HistorySyncConfig.StorageQuotaMb = proto.Uint32(10240)
	store.DeviceProps.HistorySyncConfig.InlineInitialPayloadInE2EeMsg = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportBotUserAgentChatHistory = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportCagReactionsAndPolls = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportBizHostedMsg = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportMessageAssociation = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportRecentSyncChunkMessageCountTuning = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportHostedGroupMsg = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportFbidBotChatHistory = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportManusHistory = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.SupportHatchHistory = proto.Bool(true)

	client.EnableAutoReconnect = true
	client.AutoTrustIdentity = true
	client.AutomaticMessageRerequestFromPhone = true
	client.SynchronousAck = false
	client.EnableDecryptedEventBuffer = true
	client.UseRetryMessageStore = true
	client.SendReportingTokens = true
	client.EmitAppStateEventsOnFullSync = false
}

func main() {
	godotenv.Load()

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "Config error:", err)
		os.Exit(1)
	}
	log := waLog.Stdout("sptzx", cfg.LogLevel, true)
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", cfg.SessionDB, log)
	if err != nil {
		fmt.Fprintln(os.Stderr, "DB error:", err)
		os.Exit(1)
	}

	sharedDB, err := sql.Open("sqlite3", cfg.SessionDB)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Shared DB error:", err)
		os.Exit(1)
	}
	sharedDB.SetMaxOpenConns(25)
	sharedDB.SetMaxIdleConns(5)
	sharedDB.SetConnMaxLifetime(5 * time.Minute)

	if err := games.Init(cfg.SessionDB); err != nil {
		fmt.Fprintln(os.Stderr, "Game DB error:", err)
		os.Exit(1)
	}
	log.Infof("✅ DB initialized")

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Device store error:", err)
		os.Exit(1)
	}

	client := whatsmeow.NewClient(deviceStore, log)
	configureClient(client)

	registry := core.GlobalRegistry()
	bot := core.NewBot(cfg, container, client, log, sharedDB)
	bot.Registry = registry

	evtHandler := handler.NewEventHandler(bot, registry)
	client.AddEventHandler(evtHandler.Handle)

	if deviceStore.ID != nil {
		log.Infof("Session ditemukan. Auto connect tanpa login ulang.")
		if err := client.Connect(); err != nil {
			fmt.Fprintln(os.Stderr, "Connect error:", err)
			os.Exit(1)
		}
	} else {
		if err := core.ResolveLoginConfigInteractive(cfg); err != nil {
			fmt.Fprintln(os.Stderr, "CLI login setup error:", err)
			os.Exit(1)
		}

		switch cfg.LoginMethod {
		case "pair", "paircode":
			if cfg.PairingPhone == "" {
				fmt.Fprintln(os.Stderr, "PAIRING_PHONE wajib diisi untuk pairing")
				os.Exit(1)
			}
			if err := core.LoginPairCode(client, cfg.PairingPhone, log); err != nil {
				fmt.Fprintln(os.Stderr, "Login error:", err)
				os.Exit(1)
			}
		default:
			if err := core.LoginQR(client, log); err != nil {
				fmt.Fprintln(os.Stderr, "Login error:", err)
				os.Exit(1)
			}
		}
	}

	log.Infof("Bot berjalan. Tekan Ctrl+C untuk stop.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Infof("Menutup koneksi...")
	sharedDB.Close()
	client.Disconnect()
}
