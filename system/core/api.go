package core

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type MichelleApi struct {
	BaseURL string
	APIKey  string
	client  *resty.Client
	cache   sync.Map
}

var Api *MichelleApi

func NewMichelleApi(baseURL, apiKey string) *MichelleApi {
	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		SetHeader("Accept", "application/json, text/plain, */*").
		SetTimeout(60 * time.Second)

	return &MichelleApi{
		BaseURL: baseURL,
		APIKey:  apiKey,
		client:  client,
	}
}

func (api *MichelleApi) Michelle(endpoint string, params map[string]string) (map[string]interface{}, error) {
	// Build Query
	u := url.Values{}
	u.Set("apikey", api.APIKey)
	for k, v := range params {
		u.Set(k, v)
	}

	// Simple caching
	cacheKey := endpoint + "?" + u.Encode()
	if val, ok := api.cache.Load(cacheKey); ok {
		return val.(map[string]interface{}), nil
	}

	// Request
	var result map[string]interface{}
	resp, err := api.client.R().
		SetQueryParams(map[string]string{
			"apikey": api.APIKey,
		}).
		SetQueryParams(params).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, err
	}
	
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	// Post-processing
	if result["cretor"] != nil {
		result["creator"] = result["cretor"]
		delete(result, "cretor")
	}
	result["creator"] = "@michelle.js - Darrel Parker"

	// Cache (just for 1 minute)
	api.cache.Store(cacheKey, result)
	go func() {
		time.Sleep(1 * time.Minute)
		api.cache.Delete(cacheKey)
	}()

	return result, nil
}
