package search

import "sptzx/src/core"

func init() {
	core.Use(&core.Command{
		Name:        "listberita",
		Aliases:     []string{"newslist", "sumberberita"},
		Description: "Tampilkan daftar sumber berita yang tersedia",
		Usage:       "listberita",
		Category:    "search",
		Handler: func(ptz *core.Ptz) error {
			ptz.React("✅")
			defer ptz.Unreact()
			return ptz.ReplyText(formatBeritaHelp())
		},
	})
}
