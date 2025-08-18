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

	attachments := []sendlayer.Attachment{
		{Path: "./path/to/attachment.pdf", Type: "application/pdf"},               // local file
		{Path: "https://example.com/image.png", Type: "image/png"},            // remote file
	}

	// Example: Send a simple text email
	// resp, err := sl.Emails.Send(
	// 	"paulie@example.com",
	// 	[]string{"recipient@example.com"},
	// 	"Sending a Test Email With Go SDK", 
	// 	"This is a test email sent using the SendLayer Go SDK",
	// 	"",
	// 	nil, nil, nil, nil, nil, nil,
	// )

	// Example: Send an email with all options
	resp, err := sl.Emails.Send(
		sendlayer.EmailAddress{Email: "paulie@example.com", Name: "Paulie Paloma"},   // sender
		[]string{"recipient1@example.com", "recipient2@example.com"},                          // to
		"Subject with all options",                                              // subject
		"",                                                                       // text (optional)
		"<h1>Hello</h1><p>This is a test email sent using the SendLayer Go SDK</p>",                                     // html (optional) - provide text or html
		[]string{"cc1@example.com", "cc2@example.com"},                          // cc (string | []string | EmailAddress | []EmailAddress)
		[]string{"bcc1@example.com"},                                            // bcc (same formats as cc)
		[]string{"reply-to@example.com"},                                           // replyTo (same formats as cc)
		attachments, 
		nil,                                                            // attachments                                                                // headers
		[]string{"welcome", "newsletter"},                                       // tags
	)
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}
	fmt.Printf("✅ Email sent successfully! Message ID: %s\n", resp.MessageID)
}
