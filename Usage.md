# Usage Guide — SendLayer Go SDK

## Local Development Setup

### Prerequisites
- Go 1.20 or higher
- SendLayer API key
- Git

### Environment Setup
1. **Set your API key:**
   ```bash
   export SENDLAYER_API_KEY=your-api-key-here
   ```

2. **Clone the repository:**
   ```bash
   git clone https://github.com/sendlayer/sendlayer-go.git
   cd sendlayer-go
   ```

3. **Initialize the module:**
   ```bash
   go mod tidy
   ```

4. **Verify installation:**
   ```bash
   go version
   go mod verify
   ```

## Architecture Overview

The Go SDK follows a Node.js-style architecture with a root struct that provides access to all services:

```go
// Initialize the SDK
sendlayer := sendlayer.New(apiKey)

// Access services through the root struct
resp, err := sendlayer.Emails.Send(...)
webhooks, err := sendlayer.Webhooks.Get()
events, err := sendlayer.Events.Get(...)
```

### Service Structure
- **`sendlayer.Emails`** - Email sending functionality
- **`sendlayer.Webhooks`** - Webhook management
- **`sendlayer.Events`** - Event retrieval and filtering

## Testing the SDK Locally

### Manual Testing

1. **Run the example scripts:**
   ```bash
   # Send an email
   go run examples/send_email/main.go
   
   # Manage webhooks
   go run examples/webhooks/main.go
   
   # Retrieve events
   go run examples/events/main.go
   ```

2. **Create your own test script:**
   ```go
   package main
   
   import (
       "fmt"
       "log"
       "os"
       "github.com/sendlayer/sendlayer-go/sendlayer"
   )
   
   func main() {
       apiKey := os.Getenv("SENDLAYER_API_KEY")
       if apiKey == "" {
           log.Fatal("SENDLAYER_API_KEY not set")
       }
       
       sendlayer := sendlayer.New(apiKey)
       
       // Test email sending
       resp, err := sendlayer.Emails.Send(
           "test@example.com",
           []string{"recipient@example.com"},
           "Test Subject",
           "Test body",
           "<h1>Test HTML</h1>",
           nil, nil, nil, nil, nil, nil,
       )
       if err != nil {
           log.Fatalf("Email test failed: %v", err)
       }
       fmt.Printf("Email sent: %s\n", resp.MessageID)
   }
   ```

3. **Test with different configurations:**
   ```go
   // Custom timeout
   sendlayer := sendlayer.New(apiKey, sendlayer.WithTimeout(60*time.Second))
   
   // Custom base URL (for testing)
   sendlayer := sendlayer.New(apiKey, 
       sendlayer.WithBaseURL("https://test-api.sendlayer.com/api/v1"))
   ```

### Interactive Testing

1. **Start Go REPL (if available):**
   ```bash
   go install github.com/x-motemen/gore/cmd/gore@latest
   gore
   ```

2. **Test in REPL:**
   ```go
   import "github.com/sendlayer/sendlayer-go/sendlayer"
   sendlayer := sendlayer.New(os.Getenv("SENDLAYER_API_KEY"))
   resp, err := sendlayer.Emails.Send("test@example.com", "recipient@example.com", "Test", "Body", "", nil, nil, nil, nil, nil, nil)
   ```

## Running Automated Tests

### Unit Tests

1. **Run all tests:**
   ```bash
   go test ./...
   ```

2. **Run with verbose output:**
   ```bash
   go test -v ./...
   ```

3. **Run with coverage:**
   ```bash
   go test -cover ./...
   ```

4. **Generate coverage report:**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out -o coverage.html
   ```

5. **Run specific test files:**
   ```bash
   go test -v ./sendlayer
   go test -v ./tests
   ```

6. **Run specific test functions:**
   ```bash
   go test -v -run TestEmailSend
   go test -v -run TestWebhookCreate
   go test -v -run TestEventsGet
   ```

### Integration Tests

1. **Run integration tests (requires API key):**
   ```bash
   SENDLAYER_API_KEY=your-key go test -v -tags=integration ./tests
   ```

2. **Run with test data:**
   ```bash
   SENDLAYER_API_KEY=your-key go test -v -run TestIntegration ./tests
   ```

### Benchmark Tests

1. **Run benchmarks:**
   ```bash
   go test -bench=. ./...
   ```

2. **Run specific benchmarks:**
   ```bash
   go test -bench=BenchmarkEmailSend ./...
   ```

## Code Quality and Linting

### Static Analysis

1. **Run go vet:**
   ```bash
   go vet ./...
   ```

2. **Run golangci-lint (if installed):**
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Run linter
   golangci-lint run
   ```

3. **Format code:**
   ```bash
   go fmt ./...
   ```

4. **Check for unused dependencies:**
   ```bash
   go mod tidy
   go mod verify
   ```

### Code Coverage

1. **Generate coverage report:**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   ```

2. **View coverage in browser:**
   ```bash
   go tool cover -html=coverage.out
   ```

## Building and Deployment

### Local Building

1. **Build the SDK:**
   ```bash
   go build ./...
   ```

2. **Build examples:**
   ```bash
   go build -o bin/send_email examples/send_email/main.go
   go build -o bin/webhooks examples/webhooks/main.go
   go build -o bin/events examples/events/main.go
   ```

3. **Cross-compilation:**
   ```bash
   # Build for different platforms
   GOOS=linux GOARCH=amd64 go build -o sendlayer-linux-amd64 ./...
   GOOS=darwin GOARCH=amd64 go build -o sendlayer-darwin-amd64 ./...
   GOOS=windows GOARCH=amd64 go build -o sendlayer-windows-amd64.exe ./...
   ```

### Version Management

1. **Update version in go.mod:**
   ```bash
   # Edit go.mod to update version
   # Example: module github.com/sendlayer/sendlayer-go/v2
   ```

2. **Tag releases:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

## Publishing to pkg.go.dev

### Preparation

1. **Ensure code quality:**
   ```bash
   go vet ./...
   go test ./...
   go mod tidy
   ```

2. **Add documentation:**
   - Ensure all exported functions have godoc comments
   - Add examples in documentation
   - Update README.md and Usage.md

3. **Create release:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

### Publishing Steps

1. **Push to GitHub:**
   ```bash
   git push origin main
   git push origin v1.0.0
   ```

2. **Verify on pkg.go.dev:**
   - Visit https://pkg.go.dev/github.com/sendlayer/sendlayer-go
   - Check that documentation is properly displayed
   - Verify examples work correctly

3. **Update documentation:**
   - Update any external documentation
   - Notify users of new release

## Troubleshooting

### Common Issues

1. **Time Unmarshaling Errors:**
   - **Issue**: `Time.UnmarshalJSON: input is not a JSON string` when getting events
   - **Cause**: API returns Unix timestamps as numbers, not ISO strings
   - **Solution**: The SDK includes a custom `UnixTime` type that handles this automatically
   - **Example**: Events are now properly parsed with `event.LoggedAt.Format("2006-01-02 15:04:05")`

2. **API Key Issues:**
   ```bash
   # Verify API key is set
   echo $SENDLAYER_API_KEY
   
   # Test API key validity
   curl -H "Authorization: Bearer $SENDLAYER_API_KEY" \
        https://console.sendlayer.com/api/v1/emails
   ```

2. **Module Issues:**
   ```bash
   # Clean module cache
   go clean -modcache
   
   # Re-download dependencies
   go mod download
   go mod tidy
   ```

3. **Build Issues:**
   ```bash
   # Clean build cache
   go clean -cache
   
   # Rebuild
   go build ./...
   ```

4. **Test Issues:**
   ```bash
   # Run tests with more verbose output
   go test -v -count=1 ./...
   
   # Check for race conditions
   go test -race ./...
   ```

### Debug Mode

1. **Enable debug logging:**
   ```go
   // Add debug prints in your code
   fmt.Printf("Request URL: %s\n", url)
   fmt.Printf("Request body: %s\n", string(body))
   ```

2. **Use Go's built-in profiler:**
   ```bash
   go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
   go tool pprof cpu.prof
   ```

## Support and Resources

### Documentation
- [SendLayer API Documentation](https://developers.sendlayer.com)
- [Go Documentation](https://golang.org/doc/)
- [Go Modules Reference](https://golang.org/ref/mod)

### Community
- [GitHub Issues](https://github.com/sendlayer/sendlayer-go/issues)
- [Go Forum](https://forum.golangbridge.org/)
- [Stack Overflow](https://stackoverflow.com/questions/tagged/go)

### Contact
- Email: support@sendlayer.com
- GitHub: [sendlayer/sendlayer-go](https://github.com/sendlayer/sendlayer-go) 