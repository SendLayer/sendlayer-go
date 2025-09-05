<a href="https://sendlayer.com">
<picture>
  <source media="(prefers-color-scheme: light)" srcset="https://sendlayer.com/wp-content/themes/sendlayer-theme/assets/images/svg/logo-dark.svg">
  <source media="(prefers-color-scheme: dark)" srcset="https://sendlayer.com/wp-content/themes/sendlayer-theme/assets/images/svg/logo-light.svg">
  <img alt="SendLayer Logo" width="200px" src="https://sendlayer.com/wp-content/themes/sendlayer-theme/assets/images/svg/logo-light.svg">
</picture>
</a>

### SendLayer Go SDK

The official Go SDK for interacting with the SendLayer API, providing a simple and intuitive interface for sending emails, managing webhooks, and retrieving email events.

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

## Installation

```bash
go get github.com/sendlayer/sendlayer-go
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sendlayer/sendlayer-go"
)

func main() {
	sl := sendlayer.New(os.Getenv("SENDLAYER_API_KEY"))

	resp, err := sl.Emails.Send(
		"paulie@example.com",
		"recipient@example.com",
		"Test Email",
		"This is a test email",
		"",
		nil, nil, nil,
		nil, nil, nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent! Message ID:", resp.MessageID)
}
```

## Features

- **Email Module**: Send emails with HTML/text content, attachments, CC/BCC, reply-to, custom headers, and tags
- **Webhooks Module**: Create, retrieve, and delete webhooks for various email events
- **Events Module**: Retrieve email events with filtering options
- **Error Handling**: Clear, typed errors for API and validation issues

## Email

Send emails using the `SendLayer` client:

```go
sl := sendlayer.New("your-api-key")

resp, err := sl.Emails.Send(
	from: sendlayer.EmailAddress{Email: "paulie@example.com", Name: "Paulie Paloma"},
	to: []sendlayer.EmailAddress{
		{Email: "recipient1@example.com", Name: "Recipient 1"},
		{Email: "recipient2@example.com", Name: "Recipient 2"},
	},
	subject: "Complex Email",
	text: "Plain text fallback",
	html: "<p>This is a <strong>test email</strong>!</p>",
	cc: []sendlayer.EmailAddress{{Email: "cc@example.com", Name: "CC"}},
	bcc: []sendlayer.EmailAddress{{Email: "bcc@example.com", Name: "BCC"}},
	replyTo: []sendlayer.EmailAddress{{Email: "reply@example.com", Name: "Reply"}},
	attachments: []sendlayer.Attachment{{Path: "path/to/file.pdf", Type: "application/pdf"}},
	headers: map[string]string{"X-Custom-Header": "value"},
	tags: []string{"tag1", "tag2"},
)
if err != nil {
	log.Fatal(err)
}
```

## Events

```go
sl := sendlayer.New("your-api-key")

// Get all events
all, err := sl.Events.Get(nil, nil, "", "", nil, nil)
if err != nil {
	log.Fatal(err)
}

// Get filtered events (last 24h, opened)
end := time.Now()
start := end.Add(-24 * time.Hour)
ev := "opened"
filtered, err := sl.Events.Get(&start, &end, ev, "", nil, nil)
if err != nil {
	log.Fatal(err)
}

fmt.Println("All events count:", all.TotalRecords)
fmt.Println("Filtered events count:", filtered.TotalRecords)
```

## Webhooks

```go
sl := sendlayer.New("your-api-key")

// Create a webhook
webhook, err := sl.Webhooks.Create("https://your-domain.com/webhook", "open")
if err != nil {
	log.Fatal(err)
}
fmt.Println("Webhook created:", webhook.WebhookID)

// Get all webhooks
webhooks, err := sl.Webhooks.Get()
if err != nil {
	log.Fatal(err)
}
fmt.Println("Webhooks:", webhooks)

// Delete a webhook
if err := sl.Webhooks.Delete(123); err != nil {
	log.Fatal(err)
}
```

## Error Handling

The SDK returns typed errors to help you handle different scenarios:

```go
resp, err := sl.Emails.Send(/* ... */)
if err != nil {
	var apiErr *sendlayer.SendLayerAPIError
	var valErr *sendlayer.SendLayerValidationError
	if errors.As(err, &apiErr) {
		fmt.Println("API error:", apiErr.Message, apiErr.StatusCode)
		return
	}
	if errors.As(err, &valErr) {
		fmt.Println("Validation error:", valErr.Error())
		return
	}
	fmt.Println("Unexpected error:", err)
	return
}
```

## More Details

To learn more about using the SendLayer SDK, be sure to check our [Developer Documentation](https://developers.sendlayer.com/sdks/go).

## License

MIT License - see [LICENSE](./LICENSE) file for details 