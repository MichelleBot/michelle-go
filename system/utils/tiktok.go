package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TikTokData struct {
	ID        string `json:"id"`
	Author    struct {
		Nickname string `json:"nickname"`
		UniqueID string `json:"unique_id"`
	} `json:"author"`
	PlayCount    int    `json:"play_count"`
	DiggCount    int    `json:"digg_count"`
	CommentCount int    `json:"comment_count"`
	ShareCount   int    `json:"share_count"`
	DownloadCount int    `json:"download_count"`
	CreateTime   int64  `json:"create_time"`
	Title        string `json:"title"`
	MusicInfo    struct {
		Title    string `json:"title"`
		Author   string `json:"author"`
		Duration int    `json:"duration"`
		Original bool   `json:"original"`
		Album    string `json:"album"`
		Play     string `json:"play"`
	} `json:"music_info"`
	Play   string   `json:"play"`
	WmPlay string   `json:"wmplay"`
	Music  string   `json:"music"`
	Images []string `json:"images"`
}

type TikTokResponse struct {
	Status int        `json:"status"`
	Data   TikTokData `json:"data"`
	Msg    string     `json:"msg"`
}

func ScraperTikTok(url string) (*TikTokData, error) {
	resp, err := http.Get("https://tikwm.com/api/?url=" + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TikTokResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != 1 {
		return nil, fmt.Errorf(result.Msg)
	}

	return &result.Data, nil
}
