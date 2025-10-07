package tests

import (
	"os"
	"testing"
	"time"

	"github.com/sendlayer/sendlayer-go"
)

func TestEventsValidation(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	events := sendlayer.NewEventsService(client)
	invalidStart := time.Now()
	invalidEnd := time.Now().AddDate(0, 0, -1)
	_, err := events.Get(&invalidStart, &invalidEnd, "", "", nil, nil)
	if err == nil {
		t.Error("Expected validation error for invalid date range")
	}
	_, err = events.Get(nil, nil, "invalid-event", "", nil, nil)
	if err == nil {
		t.Error("Expected validation error for invalid event")
	}
}

// Integration test (requires real API key)
func TestEventsIntegration(t *testing.T) {
	apiKey := os.Getenv("SENDLAYER_API_KEY")
	if apiKey == "" {
		t.Skip("SENDLAYER_API_KEY not set")
	}
	sl := sendlayer.New(apiKey)
	resp, err := sl.Events.Get(nil, nil, "delivered", "", nil, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.TotalRecords < 0 {
		t.Error("Expected non-negative TotalRecords")
	}
}
