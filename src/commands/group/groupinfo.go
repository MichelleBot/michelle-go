package group

import (
	"fmt"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"
)

func init() {
	core.Use(&core.Command{
		Name:        "groupinfo",
		Aliases:     []string{"ginfo"},
		Description: "Lihat info lengkap group",
		Usage:       "groupinfo",
		Category:    "group",
		GroupOnly:   true,
		Handler: func(ptz *core.Ptz) error {
			ptz.React("⏳")
			defer ptz.Unreact()
			if err := ptz.LoadGroupInfo(); err != nil {
				return ptz.ReplyText("❌ " + err.Error())
			}
			info := ptz.GroupInfo
			totalAdmin := serialize.GetGroupAdminCount(info)
			totalMember := serialize.GetGroupMemberCount(info)
			disappear := "off"
			if info.IsEphemeral {
				switch info.DisappearingTimer {
				case 86400:
					disappear = "24 jam"
				case 604800:
					disappear = "7 hari"
				case 7776000:
					disappear = "90 hari"
				default:
					disappear = fmt.Sprintf("%d detik", info.DisappearingTimer)
				}
			}
			var sb strings.Builder
			sb.WriteString("📋 *Group Info*\n\n")
			sb.WriteString(fmt.Sprintf("*Nama:* %s\n", info.Name))
			if info.Topic != "" {
				sb.WriteString(fmt.Sprintf("*Deskripsi:* %s\n", info.Topic))
			}
			sb.WriteString(fmt.Sprintf("*JID:* `%s`\n", info.JID.String()))
			sb.WriteString(fmt.Sprintf("*Owner:* @%s\n", info.OwnerJID.User))
			sb.WriteString(fmt.Sprintf("*Dibuat:* %s\n", info.GroupCreated.Format("02 Jan 2006")))
			sb.WriteString(fmt.Sprintf("*Member:* %d  •  *Admin:* %d\n", totalMember, totalAdmin))
			sb.WriteString(fmt.Sprintf("*Announce:* %v  •  *Locked:* %v\n", info.IsAnnounce, info.IsLocked))
			sb.WriteString(fmt.Sprintf("*Pesan Hilang:* %s", disappear))
			return ptz.ReplyText(sb.String())
		},
	})
}
