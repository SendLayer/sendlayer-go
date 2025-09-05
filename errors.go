package sendlayer

import "fmt"

type SendLayerError struct {
	Message string
}

func (e *SendLayerError) Error() string { return e.Message }

type SendLayerAPIError struct {
	Message    string
	StatusCode int
	Response   []byte
}

func (e *SendLayerAPIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

type SendLayerAuthenticationError struct{ SendLayerError }
type SendLayerValidationError struct{ SendLayerError }
type SendLayerNotFoundError struct{ SendLayerError }
type SendLayerRateLimitError struct{ SendLayerError }
type SendLayerInternalServerError struct{ SendLayerError }
