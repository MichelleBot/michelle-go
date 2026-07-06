package admin

import (
	"fmt"
	"strings"

	"michelle/system/core"
)

func init() {
	usage := []string{
		"adzan", "antiphishing", "antiporn", "autoreply", "antidelete",
		"antiforward", "autosticker", "adminonly", "antibot", "antilink",
		"antivirtex", "antisticker", "restrict", "left", "filter",
		"localonly", "welcome", "game", "mysterybox", "antitagsw", "captcha",
	}

	for _, cmd := range usage {
		core.Use(&core.Command{
			Usage:     []string{cmd},
			UsageHint: "on / off",
			Category:  "admin",
			GroupOnly: true,
			AdminOnly: true,
			Handler:   runSettings,
		})
	}
}

func runSettings(ptz *core.Ptz) error {
	typeField := ptz.Command
	
	// Map command name to DB field
	fieldMap := map[string]string{
		"adzan":        "adzan",
		"antiphishing": "antiphishing",
		"antiporn":     "antiporn",
		"autoreply":    "filter", // Assuming autoreply uses filter
		"antidelete":   "antidelete",
		"antiforward":  "antiforward",
		"autosticker":  "autosticker",
		"adminonly":    "adminonly",
		"antibot":      "antibot",
		"antilink":     "antilink",
		"antivirtex":   "antivirtex",
		"antisticker":  "antisticker",
		"restrict":     "restrict",
		"left":         "left",
		"filter":       "filter",
		"localonly":    "localonly",
		"welcome":      "welcome",
		"game":         "game",
		"mysterybox":   "mysterybox",
		"antitagsw":    "antitagsw",
		"captcha":      "captcha",
	}

	dbField, ok := fieldMap[typeField]
	if !ok {
		return nil
	}

	var status bool
	query := fmt.Sprintf("SELECT %s FROM groups WHERE jid = ?", dbField)
	err := ptz.Bot.DB.Conn.QueryRow(query, ptz.Chat.String()).Scan(&status)
	if err != nil {
		_, err = ptz.Bot.DB.Conn.Exec("INSERT OR IGNORE INTO groups (jid, name) VALUES (?, ?)", ptz.Chat.String(), ptz.Chat.User)
		if err != nil {
			return ptz.ReplyText("❌ Gagal inisialisasi grup: " + err.Error())
		}
		status = false
	}

	if len(ptz.Args) == 0 {
		statusStr := "OFF"
		if status {
			statusStr = "ON"
		}
		return ptz.ReplyText(fmt.Sprintf("🚩 *Status terkini* : [ %s ] (Kirim *On* atau *Off*)", statusStr))
	}

	option := strings.ToLower(ptz.Args[0])
	if option != "on" && option != "off" {
		statusStr := "OFF"
		if status {
			statusStr = "ON"
		}
		return ptz.ReplyText(fmt.Sprintf("🚩 *Status terkini* : [ %s ] (Kirim *On* atau *Off*)", statusStr))
	}

	newStatus := option == "on"
	if status == newStatus {
		return ptz.ReplyText(fmt.Sprintf("🚩 %s sudah %s sebelumnya.", strings.Title(typeField), "diaktifkan"))
	}

	updateQuery := fmt.Sprintf("UPDATE groups SET %s = ? WHERE jid = ?", dbField)
	_, err = ptz.Bot.DB.Conn.Exec(updateQuery, newStatus, ptz.Chat.String())
	if err != nil {
		return ptz.ReplyText("❌ Gagal update database: " + err.Error())
	}

	return ptz.ReplyText(fmt.Sprintf("🚩 %s berhasil %s.", strings.Title(typeField), "diaktifkan"))
}
