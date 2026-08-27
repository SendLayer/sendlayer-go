package sendlayer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultTimeout is how long the SDK waits for the API before giving up.
// net/http applies no timeout of its own, which lets a hung connection block
// forever.
const DefaultTimeout = 30 * time.Second

type Client struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	Headers map[string]string
	HTTP    *http.Client
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		APIKey:  apiKey,
		BaseURL: "https://console.sendlayer.com/api/v1",
		Timeout: DefaultTimeout,
		Headers: map[string]string{},
		HTTP:    &http.Client{Timeout: DefaultTimeout},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type ClientOption func(*Client)

// WithTimeout sets how long to wait for the API. Defaults to DefaultTimeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.Timeout = timeout
		c.HTTP.Timeout = timeout
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithHeaders adds extra headers to every request. Authorization is applied
// after these, so a caller cannot accidentally drop authentication.
func WithHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// WithHTTPClient replaces the underlying *http.Client. Its Timeout is set to
// the client's configured timeout unless it already carries one.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient == nil {
			return
		}
		if httpClient.Timeout == 0 {
			httpClient.Timeout = c.Timeout
		}
		c.HTTP = httpClient
	}
}

// errorResponse is the SendLayer error envelope:
//
//	{"Errors": [{"Code": 14, "Message": "..."}]}
//
// Some responses use a singular {"Error": "..."} instead. Error is kept as raw
// JSON so a non-string value there cannot break decoding of Errors.
type errorResponse struct {
	Errors []ErrorEntry    `json:"Errors"`
	Error  json.RawMessage `json:"Error"`
}

// errorText returns the singular Error field when it holds a non-empty string.
func (r *errorResponse) errorText() string {
	if len(r.Error) == 0 {
		return ""
	}
	var text string
	if err := json.Unmarshal(r.Error, &text); err != nil {
		return ""
	}
	return text
}

// errorFactory builds one of the SDK's error types from a parsed response.
type errorFactory func(base SendLayerError) error

// errorMap maps an HTTP status onto the error type to return and the fallback
// message used when the response carries no message of its own.
var errorMap = map[int]struct {
	newError errorFactory
	fallback string
}{
	400: {func(b SendLayerError) error { return &SendLayerValidationError{b} }, "Invalid request parameters"},
	401: {func(b SendLayerError) error { return &SendLayerAuthenticationError{b} }, "Invalid API key"},
	404: {func(b SendLayerError) error { return &SendLayerNotFoundError{b} }, "Resource not found"},
	422: {func(b SendLayerError) error { return &SendLayerValidationError{b} }, "Unprocessable Entity"},
	429: {func(b SendLayerError) error { return &SendLayerRateLimitError{b} }, "Rate limit exceeded"},
	500: {func(b SendLayerError) error { return &SendLayerInternalServerError{b} }, "Internal server error"},
}

// parseErrorBody decodes an error response body. A body that isn't a JSON
// object (an HTML page from a proxy, say) yields nil rather than an error, so
// the caller falls back to a status-derived message.
func parseErrorBody(body []byte) *errorResponse {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return nil
	}
	var parsed errorResponse
	if err := json.Unmarshal(trimmed, &parsed); err != nil {
		return nil
	}
	return &parsed
}

// extractMessage builds a message from the response, preferring the API's own
// text. Multiple Errors messages are joined with "; ", then the singular Error
// key is tried, then fallback.
func extractMessage(parsed *errorResponse, fallback string) string {
	if parsed != nil {
		parts := make([]string, 0, len(parsed.Errors))
		for _, entry := range parsed.Errors {
			if entry.Message != "" {
				parts = append(parts, entry.Message)
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "; ")
		}
		if text := parsed.errorText(); text != "" {
			return text
		}
	}
	return fallback
}

// buildError maps an error response onto the appropriate SendLayer error type.
func buildError(statusCode int, body []byte) error {
	parsed := parseErrorBody(body)

	var entries []ErrorEntry
	if parsed != nil {
		entries = parsed.Errors
	}

	mapped, ok := errorMap[statusCode]
	if !ok {
		fallback := "API request failed"
		if statusCode >= 500 && statusCode < 600 {
			fallback = "Server error"
		}
		// Prefer the HTTP reason phrase over a bare generic fallback.
		if text := http.StatusText(statusCode); parsed == nil && text != "" {
			fallback = text
		}
		return &SendLayerAPIError{SendLayerError{
			Message:    extractMessage(parsed, fallback),
			StatusCode: statusCode,
			Response:   body,
			Errors:     entries,
		}}
	}

	return mapped.newError(SendLayerError{
		Message:    extractMessage(parsed, mapped.fallback),
		StatusCode: statusCode,
		Response:   body,
		Errors:     entries,
	})
}

// doRequest performs a request against the SendLayer API.
//
// It always returns a SendLayer error type -- a net/http or encoding/json error
// is never surfaced to the caller. A successful response with an empty body
// decodes to "{}" so callers can unmarshal it unconditionally.
func (c *Client) doRequest(method, endpoint string, body interface{}, query map[string]string) ([]byte, int, error) {
	requestURL := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, &SendLayerError{Message: "Failed to marshal request body: " + err.Error()}
		}
		reqBody = bytes.NewReader(b)
	}

	// Values are escaped -- concatenating them raw produced malformed URLs for
	// any parameter containing a reserved character.
	if len(query) > 0 {
		values := url.Values{}
		for k, v := range query {
			values.Set(k, v)
		}
		requestURL += "?" + values.Encode()
	}

	req, err := http.NewRequest(method, requestURL, reqBody)
	if err != nil {
		return nil, 0, &SendLayerError{Message: "Failed to create request: " + err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SendLayer-Go-SDK/"+Version)
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	// Set after the caller's headers so those cannot drop authentication.
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, 0, &SendLayerError{
				Message: fmt.Sprintf("Request timed out after %s", c.Timeout),
			}
		}
		return nil, 0, &SendLayerError{Message: "Connection error: " + err.Error()}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, &SendLayerError{
			Message:    "Failed to read response body: " + err.Error(),
			StatusCode: resp.StatusCode,
		}
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, buildError(resp.StatusCode, respBody)
	}

	// Successful responses with no body (e.g. 204 No Content) decode to an empty
	// object rather than failing the caller's json.Unmarshal.
	if len(bytes.TrimSpace(respBody)) == 0 {
		return []byte("{}"), resp.StatusCode, nil
	}

	return respBody, resp.StatusCode, nil
}

// decodeResponse unmarshals a successful response body, wrapping a decode
// failure as a SendLayerError so callers only ever see SDK error types.
func decodeResponse(body []byte, target interface{}) error {
	if err := json.Unmarshal(body, target); err != nil {
		return &SendLayerError{Message: "Invalid JSON response from API: " + err.Error()}
	}
	return nil
}
