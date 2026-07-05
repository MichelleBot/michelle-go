package general

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"sptzx/src/core"
)

const limitPrice = 100

func init() {
	core.Use(&core.Command{
		Name:        "mylimit",
		Aliases:     []string{"limitku", "profile"},
		Description: "Lihat limit saldo dan premium",
		Usage:       "mylimit",
		Category:    "general",
		Handler:     handleMyLimit,
	})

	core.Use(&core.Command{
		Name:        "balance",
		Aliases:     []string{"saldo"},
		Description: "Lihat ringkasan balance dan cara topup",
		Usage:       "balance",
		Category:    "general",
		Handler:     handleBalance,
	})

	core.Use(&core.Command{
		Name:        "buylimit",
		Description: "Beli limit dengan saldo kredit",
		Usage:       "buylimit <jumlah>",
		Category:    "general",
		Handler:     handleBuyLimit,
	})

	core.Use(&core.Command{
		Name:        "leaderboard",
		Aliases:     []string{"lb", "top"},
		Description: "Top 5 user XP dan balance",
		Usage:       "leaderboard",
		Category:    "general",
		Handler:     handleLeaderboard,
	})
}

func handleMyLimit(ptz *core.Ptz) error {
	userID := ptz.GetPhoneJID().User
	profile := ptz.Bot.Users.GetUserProfile(userID)

	premiumText := "tidak aktif"
	if !profile.PremiumUntil.IsZero() && time.Now().Before(profile.PremiumUntil) {
		remaining := time.Until(profile.PremiumUntil)
		days := int(remaining.Hours()) / 24
		hours := int(remaining.Hours()) % 24
		premiumText = fmt.Sprintf("aktif sisa %d hari %d jam", days, hours)
	}

	text := fmt.Sprintf(
		"*Status akun kamu*\n\n- Total limit tersedia: *%d*\n- Jatah harian: *%d*\n- Sisa jatah harian: *%d*\n- Limit tambahan: *%d*\n- Saldo: *%d*\n- Premium: *%s*\n- Harga limit: *%d saldo* per 1 limit\n\n_Reset jatah harian: 00:00 UTC_",
		profile.LimitBalance,
		profile.DailyLimit,
		profile.DailyRemain,
		profile.ExtraLimit,
		profile.Credit,
		premiumText,
		limitPrice,
	)

	return ptz.ReplyText(text)
}

func handleBuyLimit(ptz *core.Ptz) error {
	if len(ptz.Args) < 1 {
		return ptz.ReplyText("format salah contoh buylimit 5")
	}

	qty, err := strconv.Atoi(ptz.Args[0])
	if err != nil || qty <= 0 {
		return ptz.ReplyText("jumlah limit harus angka lebih dari 0")
	}

	userID := ptz.GetPhoneJID().User
	if err := ptz.Bot.Users.BuyLimit(userID, qty, limitPrice); err != nil {
		return ptz.ReplyText("gagal beli limit " + err.Error())
	}

	profile := ptz.Bot.Users.GetUserProfile(userID)
	return ptz.ReplyText(fmt.Sprintf("berhasil beli %d limit limit sekarang %d saldo sekarang %d", qty, profile.LimitBalance, profile.Credit))
}

func handleBalance(ptz *core.Ptz) error {
	userID := ptz.GetPhoneJID().User
	profile := ptz.Bot.Users.GetUserProfile(userID)

	ownerContact := "owner tidak tersedia"
	if ptz.Bot != nil && ptz.Bot.Config != nil && len(ptz.Bot.Config.Owners) > 0 {
		ownerContact = ptz.Bot.Config.Owners[0]
	}

	text := fmt.Sprintf(
		"*Info balance*\n\n- Saldo kamu: *%d*\n- Total limit: *%d*\n- Sisa limit harian: *%d*\n- Limit tambahan: *%d*\n- XP kamu: *%d*\n\n*Cara dapat balance:*\n- reward dari game\n- topup dari admin atau owner\n\n*Info limit harian:*\n- free user: *100* per hari\n- premium user: *500* per hari\n- reset setiap *00:00 UTC*\n- limit hasil beli atau topup tetap aman\n\n*Info topup:*\n- harga limit saat ini: *%d saldo* per 1 limit\n- hubungi owner: `%s`",
		profile.Credit,
		profile.LimitBalance,
		profile.DailyRemain,
		profile.ExtraLimit,
		profile.XP,
		limitPrice,
		ownerContact,
	)

	return ptz.ReplyText(text)
}

func handleLeaderboard(ptz *core.Ptz) error {
	topXP := ptz.Bot.Users.TopByXP(5)
	topCredit := ptz.Bot.Users.TopByCredit(5)

	var sb strings.Builder
	sb.WriteString("*Leaderboard Top 5*\n\n")
	sb.WriteString("*Top XP*\n")
	if len(topXP) == 0 {
		sb.WriteString("Belum ada data\n")
	} else {
		for i, row := range topXP {
			sb.WriteString(fmt.Sprintf("%d. %s - %d XP\n", i+1, maskNumber(row.JID), row.Value))
		}
	}

	sb.WriteString("\n*Top Balance*\n")
	if len(topCredit) == 0 {
		sb.WriteString("Belum ada data")
	} else {
		for i, row := range topCredit {
			sb.WriteString(fmt.Sprintf("%d. %s - %d\n", i+1, maskNumber(row.JID), row.Value))
		}
	}

	return ptz.ReplyText(sb.String())
}

func maskNumber(number string) string {
	if len(number) <= 6 {
		return "***"
	}
	return number[:3] + strings.Repeat("*", len(number)-6) + number[len(number)-3:]
}
