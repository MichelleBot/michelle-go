package admin

import (
	"fmt"
	"michelle/system/core"
	"michelle/system/utils"
	"strconv"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"mute"},
		UsageHint: "0 / 1",
		Category:  "admin tools",
		Handler:   runMute,
		AdminOnly: true,
		GroupOnly: true,
	})
}

func runMute(ptz *core.Ptz) error {
	var muteStatus bool
	err := ptz.Bot.DB.Conn.QueryRow("SELECT mute FROM groups WHERE jid = ?", ptz.Chat.String()).Scan(&muteStatus)
	if err != nil {
		_, err = ptz.Bot.DB.Conn.Exec("INSERT OR IGNORE INTO groups (jid, name) VALUES (?, ?)", ptz.Chat.String(), ptz.Chat.User)
		if err != nil {
			return ptz.ReplyText("❌ Gagal inisialisasi grup: " + err.Error())
		}
		muteStatus = false
	}

	if len(ptz.Args) == 0 {
		statusText := "True"
		if !muteStatus {
			statusText = "False"
		}
		msg := fmt.Sprintf("🚩 *Status terkini* : [ %s ] (Kirim *1* atau *0*)\n\n> Gunakan *.mute 1* untuk mengaktifkan dan *.mute 0* untuk mematikan mute.", statusText)
		return ptz.ReplyText(msg)
	}

	arg, err := strconv.Atoi(ptz.Args[0])
	if err != nil || (arg != 0 && arg != 1) {
		statusText := "True"
		if !muteStatus {
			statusText = "False"
		}
		msg := fmt.Sprintf("🚩 *Status terkini* : [ %s ] (Kirim *1* atau *0*)\n\n> Gunakan *.mute 1* untuk mengaktifkan dan *.mute 0* untuk mematikan mute.", statusText)
		return ptz.ReplyText(msg)
	}

	newMuteStatus := arg == 1

	if muteStatus == newMuteStatus {
		if newMuteStatus {
			return ptz.ReplyText(utils.Texted("bold", "🚩 Sebelumnya sudah dimute."))
		} else {
			return ptz.ReplyText(utils.Texted("bold", "🚩 Sebelumnya tidak dimute."))
		}
	}

	_, err = ptz.Bot.DB.Conn.Exec("UPDATE groups SET mute = ? WHERE jid = ?", newMuteStatus, ptz.Chat.String())
	if err != nil {
		return ptz.ReplyText("❌ Gagal update database: " + err.Error())
	}

	if newMuteStatus {
		return ptz.ReplyText(utils.Texted("bold", "🚩 Berhasil dimute."))
	} else {
		return ptz.ReplyText(utils.Texted("bold", "🚩 Berhasil unmute."))
	}
}
