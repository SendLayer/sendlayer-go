package sendlayer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	HTTP    *http.Client
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		APIKey:  apiKey,
		BaseURL: "https://console.sendlayer.com/api/v1",
		Timeout: 30 * time.Second,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type ClientOption func(*Client)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.Timeout = timeout
		c.HTTP.Timeout = timeout
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

func (c *Client) doRequest(method, endpoint string, body interface{}, query map[string]string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, &SendLayerError{Message: "Failed to marshal request body"}
		}
		reqBody = bytes.NewReader(b)
	}
	if len(query) > 0 {
		url += "?"
		for k, v := range query {
			url += fmt.Sprintf("%s=%s&", k, v)
		}
		url = url[:len(url)-1]
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, 0, &SendLayerError{Message: "Failed to create request"}
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SendLayer-Go-SDK/1.0.0")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, 0, &SendLayerError{Message: "HTTP request failed: " + err.Error()}
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, &SendLayerError{Message: "Failed to read response body"}
	}
	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case 401:
			return nil, resp.StatusCode, &SendLayerAuthenticationError{SendLayerError{Message: "Invalid API key"}}
		case 400, 422:
			return nil, resp.StatusCode, &SendLayerValidationError{SendLayerError{Message: string(respBody)}}
		case 404:
			return nil, resp.StatusCode, &SendLayerNotFoundError{SendLayerError{Message: string(respBody)}}
		case 429:
			return nil, resp.StatusCode, &SendLayerRateLimitError{SendLayerError{Message: string(respBody)}}
		case 500:
			return nil, resp.StatusCode, &SendLayerInternalServerError{SendLayerError{Message: string(respBody)}}
		default:
			return nil, resp.StatusCode, &SendLayerAPIError{
				Message:    string(respBody),
				StatusCode: resp.StatusCode,
				Response:   respBody,
			}
		}
	}
	return respBody, resp.StatusCode, nil
}
