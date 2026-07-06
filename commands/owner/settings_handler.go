package owner

import (
	"fmt"
	"strings"

	"michelle/system/core"
)

func RegisterToggleSettings(cmdName string) {
	core.Use(&core.Command{
		Usage:     []string{cmdName},
		UsageHint: "on / off",
		Category:  "owner",
		OwnerOnly: true,
		Handler:   HandleToggleSettings,
	})
}

func HandleToggleSettings(ptz *core.Ptz) error {
	cmd := ptz.Command
	botSettings := ptz.Bot.BotConfig
	
	status := botSettings.GetFlag(cmd)

	if len(ptz.Args) == 0 {
		statusStr := "OFF"
		if status {
			statusStr = "ON"
		}
		return ptz.ReplyText(fmt.Sprintf("🚩 *Status terkini* : [ %s ] (Ketik *On* atau *Off*)\n\n> Gunakan *. %s on* untuk mengaktifkan dan *. %s off* untuk mematikan", statusStr, cmd, cmd))
	}

	option := strings.ToLower(ptz.Args[0])
	if option != "on" && option != "off" {
		statusStr := "OFF"
		if status {
			statusStr = "ON"
		}
		return ptz.ReplyText(fmt.Sprintf("🚩 *Status terkini* : [ %s ] (Ketik *On* atau *Off*)", statusStr))
	}

	newStatus := option == "on"
	if status == newStatus {
		statusStr := "diaktifkan"
		if !newStatus {
			statusStr = "dinonaktifkan"
		}
		return ptz.ReplyText(fmt.Sprintf("🚩 %s sebelumnya telah %s.", strings.Title(cmd), statusStr))
	}
	
	botSettings.SetFlag(cmd, newStatus)

	statusStr := "diaktifkan"
	if !newStatus {
		statusStr = "dinonaktifkan"
	}

	return ptz.ReplyText(fmt.Sprintf("🚩 %s berhasil %s.", strings.Title(cmd), statusStr))
}
