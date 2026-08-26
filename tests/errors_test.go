package tests

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sendlayer/sendlayer-go"
)

// newTestSDK spins up a server that returns the given status and body for every
// request, and points an SDK instance at it.
func newTestSDK(t *testing.T, status int, body string) (*sendlayer.SendLayer, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	return sendlayer.New("test-key", sendlayer.WithBaseURL(server.URL)), server
}

func sendTestEmail(sl *sendlayer.SendLayer) error {
	_, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Subject",
		Text:    "Body",
	})
	return err
}

func TestErrorTypeMapping(t *testing.T) {
	cases := []struct {
		status int
		expect interface{}
	}{
		{400, &sendlayer.SendLayerValidationError{}},
		{401, &sendlayer.SendLayerAuthenticationError{}},
		{404, &sendlayer.SendLayerNotFoundError{}},
		{422, &sendlayer.SendLayerValidationError{}},
		{429, &sendlayer.SendLayerRateLimitError{}},
		{500, &sendlayer.SendLayerInternalServerError{}},
		{503, &sendlayer.SendLayerAPIError{}},
		{418, &sendlayer.SendLayerAPIError{}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("status_%d", tc.status), func(t *testing.T) {
			sl, _ := newTestSDK(t, tc.status, `{}`)
			err := sendTestEmail(sl)
			if err == nil {
				t.Fatalf("expected an error for status %d", tc.status)
			}

			switch tc.expect.(type) {
			case *sendlayer.SendLayerValidationError:
				var target *sendlayer.SendLayerValidationError
				if !errors.As(err, &target) {
					t.Errorf("status %d: got %T, want *SendLayerValidationError", tc.status, err)
				}
			case *sendlayer.SendLayerAuthenticationError:
				var target *sendlayer.SendLayerAuthenticationError
				if !errors.As(err, &target) {
					t.Errorf("status %d: got %T, want *SendLayerAuthenticationError", tc.status, err)
				}
			case *sendlayer.SendLayerNotFoundError:
				var target *sendlayer.SendLayerNotFoundError
				if !errors.As(err, &target) {
					t.Errorf("status %d: got %T, want *SendLayerNotFoundError", tc.status, err)
				}
			case *sendlayer.SendLayerRateLimitError:
				var target *sendlayer.SendLayerRateLimitError
				if !errors.As(err, &target) {
					t.Errorf("status %d: got %T, want *SendLayerRateLimitError", tc.status, err)
				}
			case *sendlayer.SendLayerInternalServerError:
				var target *sendlayer.SendLayerInternalServerError
				if !errors.As(err, &target) {
					t.Errorf("status %d: got %T, want *SendLayerInternalServerError", tc.status, err)
				}
			case *sendlayer.SendLayerAPIError:
				var target *sendlayer.SendLayerAPIError
				if !errors.As(err, &target) {
					t.Errorf("status %d: got %T, want *SendLayerAPIError", tc.status, err)
				}
			}
		})
	}
}

func TestErrorMessageFromErrorsArray(t *testing.T) {
	sl, _ := newTestSDK(t, 401, `{"Errors":[{"Code":14,"Message":"Invalid API key supplied"}]}`)
	err := sendTestEmail(sl)
	if err == nil {
		t.Fatal("expected an error")
	}
	if err.Error() != "Invalid API key supplied" {
		t.Errorf("got %q, want the API's own message", err.Error())
	}
}

func TestErrorMessagesJoined(t *testing.T) {
	body := `{"Errors":[{"Code":1,"Message":"Bad from"},{"Code":2,"Message":"Bad to"}]}`
	sl, _ := newTestSDK(t, 422, body)
	err := sendTestEmail(sl)
	if err == nil {
		t.Fatal("expected an error")
	}
	if err.Error() != "Bad from; Bad to" {
		t.Errorf("got %q, want %q", err.Error(), "Bad from; Bad to")
	}
}

func TestErrorMessageFallsBackToSingularErrorKey(t *testing.T) {
	sl, _ := newTestSDK(t, 400, `{"Error":"Missing subject"}`)
	err := sendTestEmail(sl)
	if err == nil || err.Error() != "Missing subject" {
		t.Errorf("got %v, want %q", err, "Missing subject")
	}
}

func TestErrorMessageFallsBackToDefault(t *testing.T) {
	sl, _ := newTestSDK(t, 429, `{}`)
	err := sendTestEmail(sl)
	if err == nil || err.Error() != "Rate limit exceeded" {
		t.Errorf("got %v, want %q", err, "Rate limit exceeded")
	}
}

func TestNonJSONErrorBodyDoesNotLeakDecodeError(t *testing.T) {
	sl, _ := newTestSDK(t, 502, `<html>Bad Gateway</html>`)
	err := sendTestEmail(sl)
	if err == nil {
		t.Fatal("expected an error")
	}
	slErr, ok := sendlayer.AsError(err)
	if !ok {
		t.Fatalf("got %T, want a SendLayer error type", err)
	}
	if slErr.StatusCode != 502 {
		t.Errorf("StatusCode = %d, want 502", slErr.StatusCode)
	}
	if err.Error() != "API Error 502: Bad Gateway" {
		t.Errorf("got %q, want the reason phrase", err.Error())
	}
}

func TestSharedErrorFields(t *testing.T) {
	body := `{"Errors":[{"Code":14,"Message":"nope"}]}`
	sl, _ := newTestSDK(t, 401, body)
	err := sendTestEmail(sl)
	if err == nil {
		t.Fatal("expected an error")
	}

	slErr, ok := sendlayer.AsError(err)
	if !ok {
		t.Fatalf("AsError failed for %T", err)
	}
	if slErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", slErr.StatusCode)
	}
	if string(slErr.Response) != body {
		t.Errorf("Response = %q, want the raw body", slErr.Response)
	}
	if len(slErr.Errors) != 1 || slErr.Errors[0].Code != 14 || slErr.Errors[0].Message != "nope" {
		t.Errorf("Errors = %+v, want one entry {14, nope}", slErr.Errors)
	}
	codes := slErr.Codes()
	if len(codes) != 1 || codes[0] != 14 {
		t.Errorf("Codes() = %v, want [14]", codes)
	}
}

func TestAPIErrorMessagePrefix(t *testing.T) {
	sl, _ := newTestSDK(t, 418, `{"Error":"teapot"}`)
	err := sendTestEmail(sl)
	if err == nil || err.Error() != "API Error 418: teapot" {
		t.Errorf("got %v, want %q", err, "API Error 418: teapot")
	}
}

func TestAsErrorRejectsForeignErrors(t *testing.T) {
	if _, ok := sendlayer.AsError(errors.New("not ours")); ok {
		t.Error("AsError accepted a non-SDK error")
	}
	if _, ok := sendlayer.AsError(nil); ok {
		t.Error("AsError accepted nil")
	}
}

func TestValidationErrorsCarryNoStatus(t *testing.T) {
	// Local validation fails before any request is made, so no server is needed.
	sl := sendlayer.New("test-key")

	_, localErr := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Subject",
	})
	if localErr == nil {
		t.Fatal("expected a validation error for missing content")
	}
	slErr, ok := sendlayer.AsError(localErr)
	if !ok {
		t.Fatalf("got %T, want a SendLayer error type", localErr)
	}
	if slErr.StatusCode != 0 {
		t.Errorf("StatusCode = %d, want 0 for a local error", slErr.StatusCode)
	}
	if len(slErr.Codes()) != 0 {
		t.Errorf("Codes() = %v, want empty for a local error", slErr.Codes())
	}
}
