package tests

import (
	"os"
	"testing"

	"github.com/sendlayer/sendlayer-go/sendlayer"
)

func TestSendEmailValidation(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	emails := sendlayer.NewEmailsService(client)
	_, err := emails.Send("invalid-email", []string{"recipient@example.com"}, "Subject", "Text", "", nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Error("Expected validation error for invalid sender email")
	}
}

func TestSendEmailNoContent(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	emails := sendlayer.NewEmailsService(client)
	_, err := emails.Send("sender@example.com", []string{"recipient@example.com"}, "Subject", "", "", nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Error("Expected validation error for missing content")
	}
}

// Integration test (requires real API key)
func TestSendEmailIntegration(t *testing.T) {
	apiKey := os.Getenv("SENDLAYER_API_KEY")
	if apiKey == "" {
		t.Skip("SENDLAYER_API_KEY not set")
	}
	sl := sendlayer.New(apiKey)
	resp, err := sl.Emails.Send("sender@example.com", []string{"recipient@example.com"}, "Test", "Hello from Go SDK", "", nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.MessageID == "" {
		t.Error("Expected MessageID in response")
	}
}
