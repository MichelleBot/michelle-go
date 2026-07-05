package owner

import (
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

func init() {
	core.Use(&core.Command{
		Name:      "q",
		Aliases:   []string{"inspect", "debug"},
		Category:  "owner",
		OwnerOnly: true,
		Handler: func(ptz *core.Ptz) error {
			target := ptz.Event.RawMessage

			if q := serialize.GetQuotedMessage(ptz.Message); q != nil {
				target = q
			}

			opts := protojson.MarshalOptions{
				Multiline:         true,
				EmitUnpopulated:   true,
				EmitDefaultValues: true,
				UseProtoNames:     true,
			}

			raw, err := opts.Marshal(target)
			if err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}

			tmp, _ := os.CreateTemp("", "inspect_*.json")
			defer os.Remove(tmp.Name())
			tmp.Write(raw)
			tmp.Close()

			return serialize.SendDocument(
				ptz.Bot.Client,
				ptz.Chat,
				raw,
				"application/json",
				"inspect.json",
				"",
			)
		},
	})
}
