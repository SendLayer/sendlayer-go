# SendLayer Go SDK

Official Go SDK for SendLayer API — a powerful transactional email delivery service.

## Features


- **Send transactional emails**: HTML, text, attachments, CC/BCC, reply-to, custom headers, tags
- **Manage webhooks**: Create, list, and delete webhooks
- **Retrieve email events**: With filters (date range, event type, message ID, pagination)
- **Comprehensive error handling**: API, validation, authentication, rate limit, and internal server errors
- **Idiomatic Go design**: Strong typing, proper error handling, and Go conventions
- **Full feature parity**: Matches Python and Node.js SDKs

## Installation

```bash
go get github.com/sendlayer/sendlayer-go/sendlayer
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/sendlayer/sendlayer-go/sendlayer"
)

func main() {
    // Initialize the SDK
    apiKey := os.Getenv("SENDLAYER_API_KEY")
    if apiKey == "" {
        log.Fatal("SENDLAYER_API_KEY environment variable is required")
    }
    
    sendlayer := sendlayer.New(apiKey)
    
    // Send an email
    resp, err := sendlayer.Emails.Send(
        "sender@example.com",
        []string{"recipient@example.com"},
        "Welcome to SendLayer!",
        "This is a plain text email",
        "<h1>Welcome to SendLayer!</h1><p>This is an HTML email</p>",
        nil, nil, nil, nil, nil, nil,
    )
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }
    
    fmt.Printf("Email sent successfully! Message ID: %s\n", resp.MessageID)
    
    // Create a webhook
    webhook, err := sendlayer.Webhooks.Create(
        "https://your-domain.com/webhook",
        "delivered",
    )
    if err != nil {
        log.Fatalf("Failed to create webhook: %v", err)
    }
    
    fmt.Printf("Webhook created! ID: %d\n", webhook.WebhookID)
    
    // Get events
    events, err := sendlayer.Events.Get(
        &time.Now().AddDate(0, 0, -7), // 7 days ago
        &time.Now(),                   // today
        nil, nil, nil, nil,
    )
    if err != nil {
        log.Fatalf("Failed to get events: %v", err)
    }
    
    fmt.Printf("Found %d events\n", events.TotalRecords)
}
```

## API Reference

### Root Struct

```go
type SendLayer struct {
    Client   *Client
    Emails   *EmailsService
    Webhooks *WebhooksService
    Events   *EventsService
}
```

### EmailsService

```go
// Send an email with flexible recipient formats
func (e *EmailsService) Send(
    from interface{},           // string, []string, EmailAddress, or []EmailAddress
    to interface{},             // string, []string, EmailAddress, or []EmailAddress
    subject string,
    text string,                // plain text content
    html string,                // HTML content
    cc interface{},             // optional CC recipients
    bcc interface{},            // optional BCC recipients
    replyTo interface{},        // optional reply-to addresses
    attachments []Attachment,   // optional attachments
    headers map[string]string,  // optional custom headers
    tags []string,              // optional tags
) (*EmailResponse, error)
```

**Recipient Formats:**
- `"user@example.com"` - single email string
- `[]string{"user1@example.com", "user2@example.com"}` - multiple emails
- `EmailAddress{Email: "user@example.com", Name: "User Name"}` - with display name
- `[]EmailAddress{...}` - multiple with display names

**Attachment Support:**
- Local files: `Attachment{Path: "/path/to/file.pdf", Type: "application/pdf"}`
- Remote files: `Attachment{Path: "https://example.com/file.pdf", Type: "application/pdf"}`

### WebhooksService

```go
// Create a webhook
func (w *WebhooksService) Create(url string, event string) (*Webhook, error)

// List all webhooks
func (w *WebhooksService) Get() ([]Webhook, error)

// Delete a webhook
func (w *WebhooksService) Delete(webhookID int) error
```

**Supported Event Types:**
- `"delivered"`, `"bounced"`, `"opened"`, `"clicked"`, `"unsubscribed"`, `"spam"`

### EventsService

```go
// Get email events with optional filters
func (e *EventsService) Get(
    startDate *time.Time,       // optional start date
    endDate *time.Time,         // optional end date
    event *string,              // optional event type filter
    messageID *string,          // optional message ID filter
    startFrom *int,             // optional pagination offset
    retrieveCount *int,         // optional pagination limit
) (*EventsResponse, error)
```

### Custom Types

The SDK includes custom types for better API compatibility:

```go
// UnixTime handles Unix timestamps from the API
type UnixTime struct {
    time.Time
}

// Custom JSON marshaling/unmarshaling for Unix timestamps
func (ut *UnixTime) UnmarshalJSON(data []byte) error
func (ut UnixTime) MarshalJSON() ([]byte, error)
```

### Error Types

The SDK provides comprehensive error handling with specific error types:

```go
// Base error type
type SendLayerError struct {
    Message string
}

// API-specific errors
type SendLayerAPIError struct {
    Message    string
    StatusCode int
    Response   []byte
}

// Specific error types
type SendLayerAuthenticationError struct{ SendLayerError }
type SendLayerValidationError struct{ SendLayerError }
type SendLayerNotFoundError struct{ SendLayerError }
type SendLayerRateLimitError struct{ SendLayerError }
type SendLayerInternalServerError struct{ SendLayerError }
```

## Examples

See the `examples/` directory for complete working examples:

- `send_email/main.go` - Send emails with various options
- `webhooks/main.go` - Manage webhooks
- `events/main.go` - Retrieve and filter events

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v -run TestEmailSend
```

## Configuration

The SDK supports configuration options:

```go
sendlayer := sendlayer.New(apiKey,
    sendlayer.WithTimeout(60*time.Second),
    sendlayer.WithBaseURL("https://custom-api.sendlayer.com/api/v1"),
)
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run tests: `go test ./...`
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## Support

- [SendLayer Documentation](https://developers.sendlayer.com)
- [GitHub Issues](https://github.com/sendlayer/sendlayer-go/issues)
- Email: support@sendlayer.com

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 