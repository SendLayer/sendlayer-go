package main

import (
	"fmt"
	"os"

	"github.com/sendlayer/sendlayer-go/sendlayer"
)

func main() {
	apiKey := os.Getenv("SENDLAYER_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set SENDLAYER_API_KEY environment variable")
		os.Exit(1)
	}

	sl := sendlayer.New(apiKey)

	// Create a new webhook
	webhook, err := sl.Webhooks.Create("https://example.com/webhook", "click")
	if err != nil {
		fmt.Printf("Error creating webhook: %v\n", err)
		return
	}
	fmt.Printf("✅ Webhook created!: %v\n", webhook)

	// List all webhooks
	webhookList, err := sl.Webhooks.Get()
	if err != nil {
		fmt.Printf("Error listing webhooks: %v\n", err)
		return
	}
	fmt.Printf("📋 Found %d webhooks:\n", len(webhookList))
	fmt.Printf("Webhook: %v\n", webhookList)

	webhookID := 26393

	// Delete the webhook we just created
	err = sl.Webhooks.Delete(webhookID)
	if err != nil {
		fmt.Printf("Error deleting webhook: %v\n", err)
		return
	}
	fmt.Println("✅ Webhook deleted!")
}
