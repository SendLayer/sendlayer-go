package sendlayer

// SendLayer is the root struct for the SDK, providing access to all services.
type SendLayer struct {
	Client   *Client
	Emails   *EmailsService
	Webhooks *WebhooksService
	Events   *EventsService
}

// New creates a new SendLayer SDK instance with the given API key and options.
func New(apiKey string, opts ...ClientOption) *SendLayer {
	client := NewClient(apiKey, opts...)
	return &SendLayer{
		Client:   client,
		Emails:   NewEmailsService(client),
		Webhooks: NewWebhooksService(client),
		Events:   NewEventsService(client),
	}
}
