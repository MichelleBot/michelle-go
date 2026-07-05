package core

import (
	"time"
)

type CmdStat struct {
	TotalHit int
	TodayHit int
	LastHit  time.Time
}

// TrackCommand tracks usage of a command
func TrackCommand(bot *Bot, cmdName string) {
	if bot.DB == nil {
		bot.Log.Warnf("TrackCommand: bot.DB is nil")
		return
	}
	
	now := time.Now()
	
	// Get current
	var total, today int
	var lastHit int64
	err := bot.DB.Conn.QueryRow(`SELECT total_hit, today_hit, last_hit FROM command_stats WHERE cmd = ?`, cmdName).Scan(&total, &today, &lastHit)
	
	if err != nil {
		// New
		_, err := bot.DB.Conn.Exec(`INSERT INTO command_stats (cmd, total_hit, today_hit, last_hit) VALUES (?, 1, 1, ?)`, cmdName, now.Unix())
		if err != nil {
			bot.Log.Errorf("TrackCommand: INSERT failed for %s: %v", cmdName, err)
		} else {
			bot.Log.Infof("TrackCommand: INSERT successful for %s", cmdName)
		}
		return
	}
	
	// Update
	total++
	if time.Unix(lastHit, 0).Format("020106") != now.Format("020106") {
		today = 1
	} else {
		today++
	}
	
	_, err = bot.DB.Conn.Exec(`UPDATE command_stats SET total_hit = ?, today_hit = ?, last_hit = ? WHERE cmd = ?`, total, today, now.Unix(), cmdName)
	if err != nil {
		bot.Log.Errorf("TrackCommand: UPDATE failed for %s: %v", cmdName, err)
	} else {
		bot.Log.Infof("TrackCommand: UPDATE successful for %s", cmdName)
	}
}

// GetStats returns all stats
func GetStats(bot *Bot) map[string]CmdStat {
	if bot.DB == nil {
		return make(map[string]CmdStat)
	}
	
	rows, err := bot.DB.Conn.Query(`SELECT cmd, total_hit, today_hit, last_hit FROM command_stats`)
	if err != nil {
		return make(map[string]CmdStat)
	}
	defer rows.Close()
	
	stats := make(map[string]CmdStat)
	for rows.Next() {
		var cmd string
		var total, today int
		var lastHit int64
		rows.Scan(&cmd, &total, &today, &lastHit)
		stats[cmd] = CmdStat{TotalHit: total, TodayHit: today, LastHit: time.Unix(lastHit, 0)}
	}
	return stats
}
