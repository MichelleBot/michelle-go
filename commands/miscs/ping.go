package commands

import (
	"fmt"
	"time"

	"michelle/system/core"
	"michelle/system/serialize"
)

func init() {
	core.Use(&core.Command{
		Usage:       []string{"ping", "p"},
		Hidden:      []string{"p"},
		Category:    "miscs",
		Handler:     handlePing,
	})
}

func handlePing(ptz *core.Ptz) error {
	start := time.Now()
	
	// Send initial message
	msgID, err := ptz.ReplyTextID("Memeriksa ...")
	if err != nil {
		return err
	}
	
	elapsed := time.Since(start)
	
	// Edit with result
	newText := fmt.Sprintf("✨ Kecepatan : [ %dms ]", elapsed.Milliseconds())
	return serialize.EditMessage(ptz.Client, ptz.Chat, msgID, newText)
}
