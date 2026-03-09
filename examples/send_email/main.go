package main

import (
	"fmt"
	"os"

	"github.com/sendlayer/sendlayer-go"
)

func main() {
	apiKey := os.Getenv("SENDLAYER_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set SENDLAYER_API_KEY environment variable")
		os.Exit(1)
	}

	sl := sendlayer.New(apiKey)

	// Example: Send a simple text email (string for From and To)
	resp, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "paulie@example.com",
		To:      []string{"pattie@example.com"},
		Subject: "Sending a Test Email With Go SDK",
		Text:    "This is a test email sent using the SendLayer Go SDK",
	})
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}
	fmt.Printf("✅ Email sent successfully! Message ID: %s\n", resp.MessageID)

	// Example: Send an email with all options (EmailAddress for From; Cc, Bcc, ReplyTo, attachments, tags)
	attachments := []sendlayer.Attachment{
		{Path: "./path/to/attachment.pdf", Type: "application/pdf"},
		{Path: "https://example.com/image.png", Type: "image/png"},
	}
	response, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:        sendlayer.EmailAddress{Email: "paulie@example.com", Name: "Paulie Paloma"},
		To:          []string{"recipient1@example.com", "recipient2@example.com"},
		Subject:     "Subject with all options",
		Html:        "<h1>Hello</h1><p>This is a test email sent using the SendLayer Go SDK</p>",
		Cc:          []string{"cc1@example.com", "cc2@example.com"},
		Bcc:         []string{"bcc1@example.com"},
		ReplyTo:     "reply-to@example.com",
		Attachments: attachments,
		Tags:        []string{"welcome", "newsletter"},
	})
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}
	fmt.Printf("✅ Complex Email sent successfully! Message ID: %s\n", response.MessageID)
}
