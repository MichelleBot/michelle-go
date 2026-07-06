package owner

func init() {
	usage := []string{
		"autobackup", "autodownload", "antispam", "debug", "groupmode",
		"multiprefix", "noprefix", "online", "self", "games",
		"verify", "levelup", "notifier",
	}

	for _, cmd := range usage {
		RegisterToggleSettings(cmd)
	}
}
