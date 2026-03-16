package tests

import (
	"os"
	"testing"

	"github.com/sendlayer/sendlayer-go"
)

func TestSendEmailValidation(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	emails := sendlayer.NewEmailsService(client)
	_, err := emails.Send(&sendlayer.SendEmailRequest{
		From:    "invalid-email",
		To:      []string{"recipient@example.com"},
		Subject: "Subject",
		Text:    "Text",
	})
	if err == nil {
		t.Error("Expected validation error for invalid sender email")
	}
}

func TestSendEmailNoContent(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	emails := sendlayer.NewEmailsService(client)
	_, err := emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Subject",
	})
	if err == nil {
		t.Error("Expected validation error for missing content")
	}
}

func TestSendEmailNilRequest(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	emails := sendlayer.NewEmailsService(client)
	_, err := emails.Send(nil)
	if err == nil {
		t.Error("Expected validation error for nil request")
	}
}

// Integration test (requires real API key)
func TestSendEmailIntegration(t *testing.T) {
	apiKey := os.Getenv("SENDLAYER_API_KEY")
	if apiKey == "" {
		t.Skip("SENDLAYER_API_KEY not set")
	}
	sl := sendlayer.New(apiKey)
	resp, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test",
		Text:    "Hello from Go SDK",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.MessageID == "" {
		t.Error("Expected MessageID in response")
	}
}
