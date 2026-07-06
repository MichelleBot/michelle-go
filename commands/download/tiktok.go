package download

import (
	"fmt"
	"strings"
	"time"

	"michelle/system/core"
	"michelle/system/utils"
)

func init() {
	core.Use(&core.Command{
		Usage:     []string{"tiktok", "tikmp3", "tikwm", "tt"},
		UsageHint: "link",
		Category:  "downloader",
		Quota:     core.PerUserQuota(1),
		Handler:   handleTikTok,
	})
}

func handleTikTok(ptz *core.Ptz) error {
	if len(ptz.Args) == 0 {
		return ptz.ReplyText(fmt.Sprintf("🚩 Contoh penggunaan: %s%s https://vm.tiktok.com/ZSR7c5G6y/", ptz.Prefix, ptz.Command))
	}

	url := ptz.Args[0]
	if !strings.Contains(url, "tiktok.com") {
		return ptz.ReplyText("🚩 Link tidak valid.")
	}

	ptz.Bot.Client.SendReaction(ptz.Chat, ptz.Info.ID, "🕒")

	data, err := utils.ScraperTikTok(url)
	if err != nil {
		return ptz.ReplyText("🚩 Error: " + err.Error())
	}

	postedAt := time.Unix(data.CreateTime, 0).Format("02/01/2006 15:04")

	switch ptz.Command {
	case "tiktok", "tt":
		caption := fmt.Sprintf(`乂  *T I K T O K*

	◦  *ID* : %s
	◦  *Author* : %s (@%s)
	◦  *Views* : %d
	◦  *Likes* : %d
	◦  *Comments* : %d
	◦  *Shares* : %d
	◦  *Saved* : %d
	◦  *Posted At* : %s

乂  *M U S I C*

	◦  *Title* : %s
	◦  *Author* : %s
	◦  *Duration* : %d seconds
	◦  *Original* : %t
	◦  *Copyright* : %t

乂  *C A P T I O N*

%s`,
			data.ID, data.Author.Nickname, data.Author.UniqueID, data.PlayCount, data.DiggCount, data.CommentCount, data.ShareCount, data.DownloadCount, postedAt,
			data.MusicInfo.Title, data.MusicInfo.Author, data.MusicInfo.Duration, data.MusicInfo.Original, data.MusicInfo.Album != "", data.Title)

		if len(data.Images) > 0 {
			for _, imgURL := range data.Images {
				ptz.Bot.Client.SendImage(ptz.Chat, imgURL, "", ptz.Info.ID)
			}
		} else if data.Play != "" {
			return ptz.Bot.Client.SendVideo(ptz.Chat, data.Play, caption, ptz.Info.ID)
		}
	case "tikwm":
		return ptz.Bot.Client.SendVideo(ptz.Chat, data.WmPlay, "🍟", ptz.Info.ID)
	case "tikmp3":
		musicURL := data.Music
		if musicURL == "" {
			musicURL = data.MusicInfo.Play
		}
		return ptz.Bot.Client.SendAudio(ptz.Chat, musicURL, true, ptz.Info.ID)
	}

	return nil
}
