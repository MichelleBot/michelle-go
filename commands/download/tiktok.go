package download

import (
	"fmt"
	"strings"
	"time"

	"michelle/system/core"
	"michelle/system/serialize"
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

	serialize.SendReaction(ptz.Bot.Client, ptz.Chat, ptz.Info.ID, ptz.Sender, "🕒")

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

		if data.Play != "" {
			videoData, err := serialize.Fetch(data.Play)
			if err != nil {
				return ptz.ReplyText("🚩 Error fetching video: " + err.Error())
			}
			return ptz.ReplyVideo(videoData, "video/mp4", caption)
		} else {
			return ptz.ReplyText("🚩 Error: No video found in API response.")
		}
	case "tikwm":
		return ptz.ReplyVideo(nil, "video/mp4", "🍟")
	case "tikmp3":
		musicURL := data.Music
		if musicURL == "" {
			musicURL = data.MusicInfo.Play
		}
		musicData, err := serialize.Fetch(musicURL)
		if err != nil {
			return ptz.ReplyText("🚩 Error fetching audio: " + err.Error())
		}
		return ptz.ReplyAudio(musicData, "audio/mpeg")
	}

	return nil
}
