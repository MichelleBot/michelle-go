package owner

import (
	"fmt"
	"regexp"
	"strconv"

	"sptzx/src/core"
)

var userDigitsOnly = regexp.MustCompile(`^[0-9]{7,20}$`)

func init() {
	core.Use(&core.Command{
		Name:        "addlimit",
		Description: "Tambah limit user",
		Usage:       "addlimit <nomor> <jumlah>",
		Category:    "owner",
		OwnerOnly:   true,
		Handler:     handleAddLimit,
	})

	core.Use(&core.Command{
		Name:        "addprem",
		Description: "Tambah premium user dalam hari",
		Usage:       "addprem <nomor> <hari>",
		Category:    "owner",
		OwnerOnly:   true,
		Handler:     handleAddPremium,
	})

	core.Use(&core.Command{
		Name:        "addsaldo",
		Description: "Tambah saldo kredit user",
		Usage:       "addsaldo <nomor> <jumlah>",
		Category:    "owner",
		OwnerOnly:   true,
		Handler:     handleAddCredit,
	})
}

func handleAddLimit(ptz *core.Ptz) error {
	user, amount, err := parseUserAndAmount(ptz)
	if err != nil {
		return ptz.ReplyText(err.Error())
	}

	if err := ptz.Bot.Users.AddLimit(user, amount); err != nil {
		return ptz.ReplyText("Gagal menambah limit user")
	}

	profile := ptz.Bot.Users.GetUserProfile(user)
	return ptz.ReplyText(fmt.Sprintf("Limit tambahan user %s berhasil ditambah %d. Total limit sekarang %d dengan extra limit %d", user, amount, profile.LimitBalance, profile.ExtraLimit))
}

func handleAddPremium(ptz *core.Ptz) error {
	user, days, err := parseUserAndAmount(ptz)
	if err != nil {
		return ptz.ReplyText(err.Error())
	}

	if err := ptz.Bot.Users.AddPremiumDays(user, days); err != nil {
		return ptz.ReplyText("Gagal menambah premium user")
	}

	profile := ptz.Bot.Users.GetUserProfile(user)
	return ptz.ReplyText(fmt.Sprintf("Premium user %s berhasil ditambah %d hari aktif sampai %s", user, days, profile.PremiumUntil.Format("2006-01-02 15:04:05")))
}

func handleAddCredit(ptz *core.Ptz) error {
	user, amount, err := parseUserAndAmount(ptz)
	if err != nil {
		return ptz.ReplyText(err.Error())
	}

	if err := ptz.Bot.Users.AddCredit(user, amount); err != nil {
		return ptz.ReplyText("Gagal menambah saldo user")
	}

	profile := ptz.Bot.Users.GetUserProfile(user)
	return ptz.ReplyText(fmt.Sprintf("Saldo user %s berhasil ditambah %d total sekarang %d", user, amount, profile.Credit))
}

func parseUserAndAmount(ptz *core.Ptz) (string, int, error) {
	if len(ptz.Args) < 2 {
		return "", 0, fmt.Errorf("format salah contoh %s%s", ptz.Bot.GetPrefix(), ptz.Command+" 6281234567890 10")
	}

	user := ptz.Args[0]
	if !userDigitsOnly.MatchString(user) {
		return "", 0, fmt.Errorf("nomor user harus angka saja dengan kode negara")
	}

	amount, err := strconv.Atoi(ptz.Args[1])
	if err != nil || amount <= 0 {
		return "", 0, fmt.Errorf("jumlah harus angka lebih dari 0")
	}

	return user, amount, nil
}
