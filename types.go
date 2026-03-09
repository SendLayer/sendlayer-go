package sendlayer

import (
	"encoding/json"
	"time"
)

// UnixTime is a custom time type that can unmarshal Unix timestamps
type UnixTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	ut.Time = time.Unix(timestamp, 0)
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (ut UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ut.Time.Unix())
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Attachment struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	Filename string `json:"filename,omitempty"` // optional; if empty, derived from Path
}

// SendEmailRequest is the user-facing request for sending email.
// From and To accept either a string (email) or EmailAddress (email + optional name).
// Cc, Bcc, and ReplyTo accept the same types and can be single or slice ([]string or []EmailAddress).
// At least one of Text or Html must be set.
type SendEmailRequest struct {
	From        interface{}   // string or EmailAddress
	To          interface{}   // string, EmailAddress, []string, or []EmailAddress
	Subject     string
	Text        string        // plain text body (optional if Html is set)
	Html        string        // HTML body (optional if Text is set)
	Cc          interface{}   // optional: string, EmailAddress, []string, or []EmailAddress
	Bcc         interface{}   // optional: same as Cc
	ReplyTo     interface{}   // optional: same as Cc
	Attachments []Attachment  // optional
	Headers     map[string]string
	Tags        []string
}

type EmailRequest struct {
	From         EmailAddress   `json:"From"`
	To           []EmailAddress `json:"To"`
	Subject      string         `json:"Subject"`
	ContentType  string         `json:"ContentType"`
	HTMLContent  string         `json:"HTMLContent,omitempty"`
	PlainContent string         `json:"PlainContent,omitempty"`
	CC           []EmailAddress `json:"CC,omitempty"`
	BCC          []EmailAddress `json:"BCC,omitempty"`
	ReplyTo      []EmailAddress `json:"ReplyTo,omitempty"`
	Attachments  []struct {
		Content     string `json:"Content"`
		Type        string `json:"Type"`
		Filename    string `json:"Filename"`
		Disposition string `json:"Disposition"`
		ContentId   int    `json:"ContentId"`
	} `json:"Attachments,omitempty"`
	Headers map[string]string `json:"Headers,omitempty"`
	Tags    []string          `json:"Tags,omitempty"`
}

type EmailResponse struct {
	MessageID string `json:"MessageID"`
}

type Webhook struct {
	WebhookID  string    `json:"WebhookID"`
	WebhookURL string `json:"WebhookURL"`
	Event      string `json:"Event"`
	Status     string `json:"Status"`
}

type WebhookCreateRequest struct {
	WebhookURL string `json:"WebhookURL"`
	Event      string `json:"Event"`
}

type WebhookCreateResponse struct {
	NewWebhookID int `json:"NewWebhookID"`
}

type WebhookListResponse struct {
	Webhooks []Webhook `json:"Webhooks"`
}

type GeoLocation struct {
	City    string `json:"City"`
	Region  string `json:"Region"`
	Country string `json:"Country"`
}

type Message struct {
	Headers MessageHeaders `json:"Headers"`
}

type MessageHeaders struct {
	MessageID string        `json:"MessageId,omitempty"`
	From      []interface{} `json:"From,omitempty"`
	ReplyTo   []interface{} `json:"ReplyTo,omitempty"`
	To        []interface{} `json:"To,omitempty"`
	Cc        []interface{} `json:"Cc,omitempty"`
	Bcc       []interface{} `json:"Bcc,omitempty"`
	Subject   string        `json:"Subject,omitempty"`
	Size      int           `json:"Size,omitempty"`
	Transport string        `json:"Transport,omitempty"`
}

type Event struct {
	Event       string      `json:"Event"`
	LoggedAt    UnixTime    `json:"LoggedAt"`
	LogLevel    string      `json:"LogLevel"`
	Message     Message     `json:"Message"`
	Recipient   string      `json:"Recipient"`
	Reason      string      `json:"Reason"`
	Ip          string      `json:"Ip"`
	GeoLocation GeoLocation `json:"GeoLocation"`
}

type EventsResponse struct {
	TotalRecords int     `json:"TotalRecords"`
	Events       []Event `json:"Events"`
}

// GetEventsRequest is the user-facing request for querying events.
// All fields are optional; omitting the request (nil) returns events with default server behavior.
type GetEventsRequest struct {
	StartDate     *time.Time
	EndDate       *time.Time
	Event         string
	MessageID     string
	StartFrom     *int
	RetrieveCount *int
}
