<a href="https://sendlayer.com">
<picture>
  <source media="(prefers-color-scheme: light)" srcset="https://sendlayer.com/wp-content/themes/sendlayer-theme/assets/images/svg/logo-dark.svg">
  <source media="(prefers-color-scheme: dark)" srcset="https://sendlayer.com/wp-content/themes/sendlayer-theme/assets/images/svg/logo-light.svg">
  <img alt="SendLayer Logo" width="200px" src="https://sendlayer.com/wp-content/themes/sendlayer-theme/assets/images/svg/logo-light.svg">
</picture>
</a>

### SendLayer Go SDK

The official Go SDK for interacting with the SendLayer API, providing a simple and intuitive interface for sending emails, managing webhooks, and retrieving email events.

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE) [![Publish Go SDK](https://github.com/SendLayer/sendlayer-go/actions/workflows/publish.yml/badge.svg)](https://github.com/SendLayer/sendlayer-go/actions/workflows/publish.yml)

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

	resp, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "paulie@example.com",
		To:      "recipient@example.com",
		Subject: "Test Email",
		Text:    "This is a test email",
	})
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

Send emails using the `SendLayer` client. `From` and `To` accept a string (email) or `EmailAddress` (email + optional name). `Cc`, `Bcc`, and `ReplyTo` accept the same types and can be a single value or a slice.

```go
sl := sendlayer.New("your-api-key")

resp, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
	From:    sendlayer.EmailAddress{Email: "paulie@example.com", Name: "Paulie Paloma"},
	To:      []sendlayer.EmailAddress{
		{Email: "recipient1@example.com", Name: "Recipient 1"},
		{Email: "recipient2@example.com", Name: "Recipient 2"},
	},
	Subject: "Complex Email",
	Text:    "Plain text fallback",
	Html:    "<p>This is a <strong>test email</strong>!</p>",
	Cc:      []sendlayer.EmailAddress{{Email: "cc@example.com", Name: "CC"}},
	Bcc:     []sendlayer.EmailAddress{{Email: "bcc@example.com", Name: "BCC"}},
	ReplyTo: sendlayer.EmailAddress{Email: "reply@example.com", Name: "Reply"},
	Attachments: []sendlayer.Attachment{{Path: "path/to/file.pdf", Type: "application/pdf"}},
	Headers: map[string]string{"X-Custom-Header": "value"},
	Tags:    []string{"tag1", "tag2"},
})
if err != nil {
	log.Fatal(err)
}
```

### HTML with a plain-text fallback

Supply both `Html` and `Text` and both parts are sent. `ContentType` is reported
as `HTML`, and clients that cannot render HTML fall back to the plain-text part.

```go
resp, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
	From:    "paulie@example.com",
	To:      "pattie@example.com",
	Subject: "Welcome!",
	Html:    "<h1>Welcome!</h1><p>Welcome to our platform.</p>",
	Text:    "Welcome! Welcome to our platform.",
})
```

## Events

```go
sl := sendlayer.New("your-api-key")

// Get all events
all, err := sl.Events.Get(nil)
if err != nil {
	log.Fatal(err)
}

// Get filtered events (last 24h, opened)
end := time.Now()
start := end.Add(-24 * time.Hour)
ev := "opened"
filtered, err := sl.Events.Get(&sendlayer.GetEventsRequest{StartDate: &start, EndDate: &end, Event: ev})
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
webhook, err := sl.Webhooks.Create(&sendlayer.WebhookCreateRequest{WebhookURL: "https://your-domain.com/webhook", Event: "open"})
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

Every error the SDK returns embeds `SendLayerError`, so the same fields are
available on all of them:

| Field | Description |
|---|---|
| `Message` | Human-readable message, taken from the API's own error text when available |
| `StatusCode` | HTTP status of the response, or `0` for local errors (validation, timeouts, connection failures) |
| `Response` | Raw response body, or `nil` when unavailable |
| `Errors` | Parsed SendLayer `Errors` entries, each with the API's numeric `Code` and `Message` |
| `Codes()` | Just the numeric codes from `Errors`, for branching |

Use `errors.As` when the specific failure matters:

```go
resp, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
	From: "you@example.com", To: "them@example.com", Subject: "Hi", Text: "Hello",
})
if err != nil {
	var authErr *sendlayer.SendLayerAuthenticationError
	var valErr *sendlayer.SendLayerValidationError
	var rateErr *sendlayer.SendLayerRateLimitError

	switch {
	case errors.As(err, &authErr):
		fmt.Println("Check your API key:", authErr.Message)
	case errors.As(err, &valErr):
		fmt.Println("Validation error:", valErr.Message)
	case errors.As(err, &rateErr):
		fmt.Println("Slow down:", rateErr.Message)
	default:
		fmt.Println("SendLayer error:", err)
	}
	return
}
```

Or use `AsError` to read the shared fields without switching on the concrete
type:

```go
if slErr, ok := sendlayer.AsError(err); ok {
	fmt.Printf("status %d: %s (codes %v)\n", slErr.StatusCode, slErr.Message, slErr.Codes())
	if slices.Contains(slErr.Codes(), 14) {
		fmt.Println("That sender domain is not authorised.")
	}
}
```

See the [error code reference](https://developers.sendlayer.com/api-reference/error-codes)
for the full list of numeric codes.

### Error Types

| Type | Returned for |
|---|---|
| `SendLayerError` | Base type embedded by all others, and returned directly for local errors: timeouts, connection failures and undecodable responses |
| `SendLayerValidationError` | Invalid parameters (returned locally), and HTTP 400 / 422 |
| `SendLayerAuthenticationError` | HTTP 401 -- invalid API key |
| `SendLayerNotFoundError` | HTTP 404 |
| `SendLayerRateLimitError` | HTTP 429 |
| `SendLayerInternalServerError` | HTTP 500 |
| `SendLayerAPIError` | Any other error status, including 5xx other than 500 |

`SendLayerAPIError` is the only type whose `Error()` is prefixed -- it reads
`API Error <status>: <message>`. Every other type returns the API's message
unchanged.

A `net/http` or `encoding/json` error is never surfaced to the caller: every
error returned by the SDK is one of the types above.

## Configuration

Pass `ClientOption` values to `sendlayer.New`:

```go
sl := sendlayer.New("your-api-key",
	sendlayer.WithTimeout(60*time.Second),
	sendlayer.WithHeaders(map[string]string{"X-Request-Id": "abc123"}),
)
```

| Option | Default | Description |
|---|---|---|
| `WithTimeout(d)` | `30s` (`sendlayer.DefaultTimeout`) | How long to wait for the API before returning `SendLayerError` |
| `WithHeaders(m)` | none | Extra headers sent with every request. Cannot override `Authorization` |
| `WithBaseURL(u)` | SendLayer API v1 | Override the API base URL |
| `WithHTTPClient(c)` | `&http.Client{Timeout: 30s}` | Supply your own `*http.Client` (for a custom transport, proxy or instrumentation) |

Requests time out after 30 seconds by default. `net/http` applies no timeout of
its own, so without this a hung server would block forever.

## More Details

To learn more about using the SendLayer SDK, be sure to check our [Developer Documentation](https://developers.sendlayer.com/sdks/go).

## License

MIT License - see [LICENSE](./LICENSE) file for details 