package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sendlayer/sendlayer-go"
)

// captureServer records the request it receives and replies with the given body.
type captureServer struct {
	server  *httptest.Server
	payload map[string]interface{}
	headers http.Header
	query   string
}

func newCaptureServer(t *testing.T, status int, body string) *captureServer {
	t.Helper()
	c := &captureServer{}
	c.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.headers = r.Header.Clone()
		c.query = r.URL.RawQuery
		c.payload = map[string]interface{}{}
		_ = json.NewDecoder(r.Body).Decode(&c.payload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(c.server.Close)
	return c
}

func (c *captureServer) sdk(opts ...sendlayer.ClientOption) *sendlayer.SendLayer {
	return sendlayer.New("test-key", append([]sendlayer.ClientOption{
		sendlayer.WithBaseURL(c.server.URL),
	}, opts...)...)
}

func TestSendEmailIncludesBothContentParts(t *testing.T) {
	srv := newCaptureServer(t, 200, `{"MessageID":"abc"}`)
	sl := srv.sdk()

	_, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Subject",
		Text:    "Plain fallback",
		Html:    "<p>Rich body</p>",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := srv.payload["ContentType"]; got != "HTML" {
		t.Errorf("ContentType = %v, want HTML", got)
	}
	if got := srv.payload["HTMLContent"]; got != "<p>Rich body</p>" {
		t.Errorf("HTMLContent = %v", got)
	}
	if got := srv.payload["PlainContent"]; got != "Plain fallback" {
		t.Errorf("PlainContent = %v, want the plain-text part to survive", got)
	}
}

func TestSendEmailHTMLOnlyOmitsPlainContent(t *testing.T) {
	srv := newCaptureServer(t, 200, `{"MessageID":"abc"}`)
	_, err := srv.sdk().Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Subject",
		Html:    "<p>x</p>",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := srv.payload["ContentType"]; got != "HTML" {
		t.Errorf("ContentType = %v, want HTML", got)
	}
	if _, present := srv.payload["PlainContent"]; present {
		t.Error("PlainContent should be omitted when only Html is set")
	}
}

func TestSendEmailTextOnlyOmitsHTMLContent(t *testing.T) {
	srv := newCaptureServer(t, 200, `{"MessageID":"abc"}`)
	_, err := srv.sdk().Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Subject",
		Text:    "plain",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := srv.payload["ContentType"]; got != "Text" {
		t.Errorf("ContentType = %v, want Text", got)
	}
	if _, present := srv.payload["HTMLContent"]; present {
		t.Error("HTMLContent should be omitted when only Text is set")
	}
}

func TestDefaultTimeout(t *testing.T) {
	client := sendlayer.NewClient("test-key")
	if client.Timeout != sendlayer.DefaultTimeout {
		t.Errorf("Timeout = %v, want %v", client.Timeout, sendlayer.DefaultTimeout)
	}
	if client.HTTP.Timeout != sendlayer.DefaultTimeout {
		t.Errorf("HTTP.Timeout = %v, want %v", client.HTTP.Timeout, sendlayer.DefaultTimeout)
	}
	if sendlayer.DefaultTimeout != 30*time.Second {
		t.Errorf("DefaultTimeout = %v, want 30s", sendlayer.DefaultTimeout)
	}
}

func TestTimeoutReturnsSendLayerError(t *testing.T) {
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(300 * time.Millisecond)
		fmt.Fprint(w, `{}`)
	}))
	defer slow.Close()

	sl := sendlayer.New("test-key",
		sendlayer.WithBaseURL(slow.URL),
		sendlayer.WithTimeout(30*time.Millisecond))

	_, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Subject",
		Text:    "Body",
	})
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if _, ok := sendlayer.AsError(err); !ok {
		t.Fatalf("got %T, want a SendLayer error type", err)
	}
	if want := "Request timed out after 30ms"; err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestConnectionErrorReturnsSendLayerError(t *testing.T) {
	// Port 0 on localhost never accepts a connection.
	sl := sendlayer.New("test-key", sendlayer.WithBaseURL("http://127.0.0.1:0"))
	_, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Subject",
		Text:    "Body",
	})
	if err == nil {
		t.Fatal("expected a connection error")
	}
	if _, ok := sendlayer.AsError(err); !ok {
		t.Fatalf("got %T, want a SendLayer error type", err)
	}
}

func TestAuthorizationHeaderIsSent(t *testing.T) {
	srv := newCaptureServer(t, 200, `{"MessageID":"abc"}`)
	_, err := srv.sdk().Emails.Send(&sendlayer.SendEmailRequest{
		From: "sender@example.com", To: "r@example.com", Subject: "s", Text: "t",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := srv.headers.Get("Authorization"); got != "Bearer test-key" {
		t.Errorf("Authorization = %q", got)
	}
	if got := srv.headers.Get("User-Agent"); got != "SendLayer-Go-SDK/"+sendlayer.Version {
		t.Errorf("User-Agent = %q", got)
	}
}

func TestCustomHeadersCannotDropAuthorization(t *testing.T) {
	srv := newCaptureServer(t, 200, `{"MessageID":"abc"}`)
	sl := srv.sdk(sendlayer.WithHeaders(map[string]string{
		"X-Trace":       "abc123",
		"Authorization": "Bearer hijacked",
	}))
	_, err := sl.Emails.Send(&sendlayer.SendEmailRequest{
		From: "sender@example.com", To: "r@example.com", Subject: "s", Text: "t",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := srv.headers.Get("X-Trace"); got != "abc123" {
		t.Errorf("X-Trace = %q, want the custom header to be sent", got)
	}
	if got := srv.headers.Get("Authorization"); got != "Bearer test-key" {
		t.Errorf("Authorization = %q, want the SDK's own credential to win", got)
	}
}

func TestEmptySuccessBodyDecodesCleanly(t *testing.T) {
	srv := newCaptureServer(t, 204, ``)
	if err := srv.sdk().Webhooks.Delete(7); err != nil {
		t.Errorf("unexpected error for an empty 204 body: %v", err)
	}
}

func TestInvalidJSONSuccessBodyReturnsSendLayerError(t *testing.T) {
	srv := newCaptureServer(t, 200, `<html>not json</html>`)
	_, err := srv.sdk().Emails.Send(&sendlayer.SendEmailRequest{
		From: "sender@example.com", To: "r@example.com", Subject: "s", Text: "t",
	})
	if err == nil {
		t.Fatal("expected an error for an undecodable success body")
	}
	if _, ok := sendlayer.AsError(err); !ok {
		t.Fatalf("got %T, want a SendLayer error type", err)
	}
}

func TestQueryParametersAreEscaped(t *testing.T) {
	srv := newCaptureServer(t, 200, `{"TotalRecords":0,"Events":[]}`)
	_, err := srv.sdk().Events.Get(&sendlayer.GetEventsRequest{
		MessageID: "id with spaces&ampersand",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srv.query != "MessageID=id+with+spaces%26ampersand" {
		t.Errorf("query = %q, want the value to be escaped", srv.query)
	}
}

// TestVersionMatchesReleaseFile guards the two places the version lives: the
// Version constant reported in the User-Agent, and .github/VERSION, which the
// publish workflow reads to create the release tag.
func TestVersionMatchesReleaseFile(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", ".github", "VERSION"))
	if err != nil {
		t.Fatalf("could not read .github/VERSION: %v", err)
	}
	if want := strings.TrimSpace(string(raw)); want != sendlayer.Version {
		t.Errorf("sendlayer.Version = %q, but .github/VERSION says %q", sendlayer.Version, want)
	}
}
