package sendlayer

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type WebhooksService struct {
	client *Client
}

func NewWebhooksService(client *Client) *WebhooksService {
	return &WebhooksService{client: client}
}

func (w *WebhooksService) validateURL(u string) bool {
	parsed, err := url.ParseRequestURI(u)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https")
}

func (w *WebhooksService) Create(webhookURL, event string) (*WebhookCreateResponse, error) {
	if !w.validateURL(webhookURL) {
		return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid webhook URL: %s", webhookURL)}}
	}
	eventOptions := map[string]bool{
		"bounce": true, "click": true, "open": true, "unsubscribe": true, "complaint": true, "delivery": true,
	}
	if !eventOptions[event] {
		return nil, &SendLayerValidationError{SendLayerError{fmt.Sprintf("Invalid event: %s", event)}}
	}
	payload := WebhookCreateRequest{
		WebhookURL: webhookURL,
		Event:      event,
	}
	respBody, _, err := w.client.doRequest("POST", "webhooks", payload, nil)
	if err != nil {
		return nil, err
	}
	var resp WebhookCreateResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (w *WebhooksService) Get() ([]Webhook, error) {
	respBody, _, err := w.client.doRequest("GET", "webhooks", nil, nil)
	if err != nil {
		return nil, err
	}
	var resp WebhookListResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp.Webhooks, nil
}

func (w *WebhooksService) Delete(webhookID int) error {
	if webhookID <= 0 {
		return &SendLayerValidationError{SendLayerError{"WebhookID must be greater than 0"}}
	}
	endpoint := fmt.Sprintf("webhooks/%d", webhookID)
	_, _, err := w.client.doRequest("DELETE", endpoint, nil, nil)
	return err
}
