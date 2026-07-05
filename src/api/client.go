package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Client struct {
	BaseURL string
	HTTP    *resty.Client
	Log     waLog.Logger
}

type Response struct {
	Status  bool        `json:"status"`
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data"`
}

func (r *Response) IsSuccess() bool {
	return r.Status || r.Success
}

func NewClient(baseURL string) *Client {
	r := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(120 * time.Second)

	return &Client{
		BaseURL: baseURL,
		HTTP:    r,
		Log:     waLog.Stdout("api", "INFO", true),
	}
}

func (c *Client) SetLogger(log waLog.Logger) {
	c.Log = log
}

func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	var apiResp Response

	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&apiResp).
		Get(endpoint)

	if err != nil {
		if c.Log != nil {
			c.Log.Errorf("API Get error [%s]: %v", endpoint, err)
		}
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode())
	}

	if !apiResp.IsSuccess() {
		if apiResp.Error != "" {
			return nil, fmt.Errorf("API error: %s", apiResp.Error)
		}
		return nil, fmt.Errorf("API returned status false")
	}

	return resp.Body(), nil
}

func (c *Client) GetRaw(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetQueryParams(params).
		Get(endpoint)

	if err != nil {
		if c.Log != nil {
			c.Log.Errorf("API GetRaw error [%s]: %v", endpoint, err)
		}
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

func Request[T any](ctx context.Context, client *Client, endpoint string, params map[string]string) (T, error) {
	var result T
	raw, err := client.GetRaw(ctx, endpoint, params)
	if err != nil {
		return result, err
	}

	var apiResp struct {
		Status  bool   `json:"status"`
		Success bool   `json:"success"`
		Data    T      `json:"data"`
		Error   string `json:"error"`
	}

	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return result, err
	}

	if !apiResp.Status && !apiResp.Success {
		if apiResp.Error != "" {
			return result, fmt.Errorf("API error: %s", apiResp.Error)
		}
		return result, fmt.Errorf("API returned status false")
	}

	return apiResp.Data, nil
}
