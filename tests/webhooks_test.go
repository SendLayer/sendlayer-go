package tests

import (
	"os"
	"testing"

	"github.com/sendlayer/sendlayer-go"
)

func TestWebhookValidation(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	webhooks := sendlayer.NewWebhooksService(client)
	_, err := webhooks.Create("invalid-url", "delivery")
	if err == nil {
		t.Error("Expected validation error for invalid URL")
	}
	_, err = webhooks.Create("https://valid.com", "invalid-event")
	if err == nil {
		t.Error("Expected validation error for invalid event")
	}
}

// Integration test (requires real API key)
func TestWebhookIntegration(t *testing.T) {
	apiKey := os.Getenv("SENDLAYER_API_KEY")
	if apiKey == "" {
		t.Skip("SENDLAYER_API_KEY not set")
	}
	sl := sendlayer.New(apiKey)
	wh, err := sl.Webhooks.Create("https://yourdomain.com/webhook", "delivery")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if wh.NewWebhookID == 0 {
		t.Error("Expected WebhookID in response")
	}
	// Clean up
	err = sl.Webhooks.Delete(wh.NewWebhookID)
	if err != nil {
		t.Errorf("Failed to delete webhook: %v", err)
	}
}
