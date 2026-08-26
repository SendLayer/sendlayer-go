package sendlayer

import (
	"errors"
	"fmt"
)

// ErrorEntry is a raw entry from SendLayer's "Errors" response array.
//
// See https://developers.sendlayer.com/api-reference/error-codes
type ErrorEntry struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

// SendLayerError is the base error for the SDK. Every error the SDK returns
// embeds it, so callers can read the same fields without first asserting the
// concrete type:
//
//   - StatusCode is the HTTP status of the response, or 0 for local errors
//     (validation, timeouts, connection failures).
//   - Response is the raw response body, or nil when unavailable.
//   - Errors holds the parsed SendLayer "Errors" entries; empty for local
//     errors and for responses that aren't in that shape.
type SendLayerError struct {
	Message    string
	StatusCode int
	Response   []byte
	Errors     []ErrorEntry
}

func (e *SendLayerError) Error() string { return e.Message }

// Base returns the embedded base error. It is promoted to every SendLayer error
// type, which is what lets AsError read the shared fields off any of them.
func (e *SendLayerError) Base() *SendLayerError { return e }

// Codes returns the numeric SendLayer error codes carried by this error, for
// branching on a specific failure:
//
//	if slices.Contains(slErr.Codes(), 14) { ... }
func (e *SendLayerError) Codes() []int {
	codes := make([]int, 0, len(e.Errors))
	for _, entry := range e.Errors {
		codes = append(codes, entry.Code)
	}
	return codes
}

// baseError is satisfied by every error type in this package, via the promoted
// Base method on the embedded SendLayerError.
type baseError interface {
	error
	Base() *SendLayerError
}

// AsError reports whether err came from this SDK and, if so, returns its shared
// fields. Use it to read StatusCode, Errors or Codes without switching on the
// concrete type:
//
//	if slErr, ok := sendlayer.AsError(err); ok {
//	    log.Printf("status %d: %v", slErr.StatusCode, slErr.Codes())
//	}
//
// Use errors.As with a concrete type when the specific failure matters.
func AsError(err error) (*SendLayerError, bool) {
	var base baseError
	if !errors.As(err, &base) {
		return nil, false
	}
	return base.Base(), true
}

// SendLayerAPIError is returned for API errors not covered by a more specific
// type: any 4xx outside 400/401/404/422/429, and any 5xx other than 500.
type SendLayerAPIError struct{ SendLayerError }

// Error prefixes the message with the status, e.g. "API Error 503: Server error".
func (e *SendLayerAPIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

// SendLayerAuthenticationError is returned for HTTP 401.
type SendLayerAuthenticationError struct{ SendLayerError }

// SendLayerValidationError is returned for local validation failures and for
// HTTP 400 and 422.
type SendLayerValidationError struct{ SendLayerError }

// SendLayerNotFoundError is returned for HTTP 404.
type SendLayerNotFoundError struct{ SendLayerError }

// SendLayerRateLimitError is returned for HTTP 429.
type SendLayerRateLimitError struct{ SendLayerError }

// SendLayerInternalServerError is returned for HTTP 500.
type SendLayerInternalServerError struct{ SendLayerError }
