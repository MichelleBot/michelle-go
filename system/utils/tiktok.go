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
	DownloadCount int    `json:"collect_count"` // Maps to collect_count from the new JSON
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
}

type TikTokResponse struct {
	Code int        `json:"code"`
	Data TikTokData `json:"data"`
	Msg  string     `json:"msg"`
}

func ScraperTikTok(url string) (*TikTokData, error) {
	client := &http.Client{}
	// Note: The endpoint here is placeholder based on previous code. I will not change it unless instructed.
	req, err := http.NewRequest("GET", "https://tikwm.com/api/?url="+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TikTokResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API Error (Code %d): %s", result.Code, result.Msg)
	}

	return &result.Data, nil
}
