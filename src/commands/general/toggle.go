package general

import (
	"fmt"
	"strings"

	"sptzx/src/commands/games"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	handler := func(ptz *core.Ptz) error {
		if len(ptz.Args) == 0 {
			return sendToggleMenu(ptz)
		}

		action := strings.ToLower(ptz.Command)
		enable := action == "enable" || action == "on"

		opt := strings.ToLower(ptz.Args[0])

		switch opt {
		case "welcome":
			if err := checkAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			s := ptz.Bot.Settings.GetGroupSettings(ptz.Chat)
			s.WelcomeEnabled = enable
			ptz.Bot.Settings.SetGroupSettings(ptz.Chat, s)
			return replyToggle(ptz, "Welcome message", enable)

		case "goodbye":
			if err := checkAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			s := ptz.Bot.Settings.GetGroupSettings(ptz.Chat)
			s.GoodbyeEnabled = enable
			ptz.Bot.Settings.SetGroupSettings(ptz.Chat, s)
			return replyToggle(ptz, "Goodbye message", enable)

		case "antidelete":
			if err := checkAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			s := ptz.Bot.Settings.GetGroupSettings(ptz.Chat)
			s.AntiDeleteEnabled = enable
			ptz.Bot.Settings.SetGroupSettings(ptz.Chat, s)
			return replyToggle(ptz, "Anti Delete", enable)

		case "announce":
			if err := checkBotAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			if err := serialize.SetGroupAnnounce(ptz.Bot.Client, ptz.Chat, enable); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return replyToggle(ptz, "Group Announce (Hanya admin kirim pesan)", enable)

		case "locked":
			if err := checkBotAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			if err := serialize.SetGroupLocked(ptz.Bot.Client, ptz.Chat, enable); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return replyToggle(ptz, "Group Locked (Hanya admin edit info)", enable)

		case "approval":
			if err := checkBotAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			if err := serialize.SetGroupJoinApprovalMode(ptz.Bot.Client, ptz.Chat, enable); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return replyToggle(ptz, "Join Approval Mode", enable)

		case "restrict":
			if err := checkBotAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			if err := serialize.SetGroupMemberAddMode(ptz.Bot.Client, ptz.Chat, enable); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return replyToggle(ptz, "Restrict Member Add (Hanya admin tambah member)", enable)

		case "ephemeral":
			if err := checkBotAdmin(ptz); err != nil {
				return ptz.ReplyText(err.Error())
			}
			var err error
			if enable {
				err = serialize.SetDisappearing7d(ptz.Bot.Client, ptz.Chat)
			} else {
				err = serialize.SetDisappearingOff(ptz.Bot.Client, ptz.Chat)
			}
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			return replyToggle(ptz, "Pesan Menghilang (7 Hari)", enable)

		case "self":
			if !ptz.IsOwner() {
				return ptz.ReplyText("🚫 Fitur ini hanya untuk Owner.")
			}
			ptz.Bot.BotConfig.SetSelfMode(enable)
			return replyToggle(ptz, "Self Mode (Bot hanya merespons Owner)", enable)

		case "public":
			if !ptz.IsOwner() {
				return ptz.ReplyText("🚫 Fitur ini hanya untuk Owner.")
			}
			ptz.Bot.BotConfig.SetSelfMode(!enable)
			return replyToggle(ptz, "Public Mode (Bot merespons semua)", enable)

		case "privateonly":
			if !ptz.IsOwner() {
				return ptz.ReplyText("🚫 Fitur ini hanya untuk Owner.")
			}
			ptz.Bot.BotConfig.SetPrivateOnly(enable)
			return replyToggle(ptz, "Private Only (Bot hanya di DM)", enable)

		case "grouponly":
			if !ptz.IsOwner() {
				return ptz.ReplyText("🚫 Fitur ini hanya untuk Owner.")
			}
			ptz.Bot.BotConfig.SetGroupOnly(enable)
			return replyToggle(ptz, "Group Only (Bot hanya di Grup)", enable)

		case "game", "games":
			if !ptz.IsOwner() {
				return ptz.ReplyText("🚫 Fitur ini hanya untuk Owner.")
			}
			games.SetGameEnabled(ptz.Chat.String(), enable)
			return replyToggle(ptz, "Game di chat ini", enable)

		default:
			return sendToggleMenu(ptz)
		}
	}

	core.Use(&core.Command{
		Name:        "enable",
		Aliases:     []string{"on"},
		Description: "Mengaktifkan fitur bot/grup",
		Usage:       "enable <fitur>",
		Category:    "general",
		Handler:     handler,
	})

	core.Use(&core.Command{
		Name:        "disable",
		Aliases:     []string{"off"},
		Description: "Menonaktifkan fitur bot/grup",
		Usage:       "disable <fitur>",
		Category:    "general",
		Handler:     handler,
	})
}

func checkAdmin(ptz *core.Ptz) error {
	if !ptz.IsGroup {
		return fmt.Errorf("🚫 Perintah ini hanya bisa digunakan di grup.")
	}
	if err := ptz.LoadGroupInfo(); err != nil {
		return fmt.Errorf("🚫 Gagal memuat data grup. Coba lagi sebentar.")
	}
	if !ptz.IsAdmin() && !ptz.IsOwner() {
		return fmt.Errorf("🚫 Kamu bukan admin grup ini.")
	}
	return nil
}

func checkBotAdmin(ptz *core.Ptz) error {
	if err := checkAdmin(ptz); err != nil {
		return err
	}
	if !ptz.IsBotAdmin() {
		return fmt.Errorf("🚫 Bot harus menjadi admin grup terlebih dahulu.")
	}
	return nil
}

func replyToggle(ptz *core.Ptz, name string, enable bool) error {
	if enable {
		return ptz.ReplyText(fmt.Sprintf("✅ *%s* berhasil *diaktifkan*.", name))
	}
	return ptz.ReplyText(fmt.Sprintf("🚫 *%s* berhasil *dinonaktifkan*.", name))
}

func sendToggleMenu(ptz *core.Ptz) error {
	msg := "🔧 *PENGATURAN ON/OFF* 🔧\n\n"
	msg += "Usage: `.enable <opsi>` atau `.disable <opsi>`\n\n"
	msg += "*Kategori Group (Butuh Admin):*\n"
	msg += "• welcome (Sambut member baru)\n"
	msg += "• goodbye (Pesan member keluar)\n"
	msg += "• antidelete (Tampilkan pesan yang dihapus)\n"
	msg += "• announce (Hanya admin yg bisa chat)\n"
	msg += "• locked (Hanya admin yg bisa edit grup)\n"
	msg += "• approval (Persetujuan saat join)\n"
	msg += "• restrict (Hanya admin yg tambah member)\n"
	msg += "• ephemeral (Pesan hilang otomatis 7h)\n\n"
	msg += "*Kategori Owner (Hanya Owner):*\n"
	msg += "• self (Bot hanya merespons Owner)\n"
	msg += "• public (Bot merespons semua orang)\n"
	msg += "• grouponly (Bot hanya aktif di Grup)\n"
	msg += "• privateonly (Bot hanya aktif di DM)\n"
	msg += "• game (Main game di chat ini)\n\n"
	msg += "Contoh: `.enable welcome`\n"
	return ptz.ReplyText(msg)
}
