# Retry Configuration Examples

The DataCrunch SDK provides automatic retry functionality with sensible defaults that work out-of-the-box.

## Default Behavior (Recommended)

By default, all clients get **3 retries** with exponential backoff:

```go
// Gets 3 retries automatically - no configuration needed!
client := datacrunch.New()

// Same with session-based approach
sess := session.New()
client := datacrunch.NewWithSession(sess)
```

**Default Retry Settings:**
- **Max Retries:** 3
- **Min Delay:** 30ms
- **Max Delay:** 300s
- **Throttle Delay:** 500ms - 300s
- **Strategy:** Exponential backoff with jitter

## Custom Retry Configuration

### Adjust Max Retries

```go
// More aggressive retries
client := datacrunch.New(datacrunch.WithRetryConfig(5, 0, 0))

// Via session
sess := session.New(session.WithMaxRetries(7))
client := datacrunch.NewWithSession(sess)
```

### Disable Retries Entirely

```go
// No retries - fail fast
client := datacrunch.New(datacrunch.WithNoRetries())

// Via session
sess := session.New(session.WithNoRetries())
client := datacrunch.NewWithSession(sess)
```

### Custom Retryer

```go
import "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"

// Create custom retryer with specific timing
customRetryer := client.NewDefaultRetryer(5) // 5 retries with defaults

// Or fully customized
customRetryer := client.DefaultRetryer{
    NumMaxRetries:    10,
    MinRetryDelay:    100 * time.Millisecond,
    MaxRetryDelay:    60 * time.Second,
    MinThrottleDelay: 1 * time.Second,
    MaxThrottleDelay: 120 * time.Second,
}

client := datacrunch.New(datacrunch.WithRetryer(customRetryer))
```

## What Gets Retried

The SDK automatically retries:

- **Network errors** (connection refused, timeouts, etc.)
- **5xx server errors** (500, 502, 503, 504)
- **429 Too Many Requests** (rate limiting)
- **Authentication token expiration** (refreshes tokens automatically)
- **Transient service unavailability**

## What Doesn't Get Retried

- **4xx client errors** (400, 401, 403, 404) - except 429
- **Invalid requests** or **malformed data**
- **Canceled requests** (context cancellation)
- **Non-retryable network errors**

## Environment-Based Configuration

You can also configure retries via environment variables in your session:

```bash
export DATACRUNCH_BASE_URL="https://api.datacrunch.io/v1"
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"
```

```go
// Uses env vars + defaults to 3 retries
sess := session.NewFromEnv(session.WithMaxRetries(5)) // Override to 5
client := datacrunch.NewWithSession(sess)
```

## Best Practices

1. **Use defaults for most cases** - 3 retries with exponential backoff works well
2. **Increase retries for batch operations** - Use 5-7 retries for less time-sensitive work
3. **Disable retries for real-time operations** - Use `WithNoRetries()` when speed matters more than reliability
4. **Custom retryer for special needs** - Implement your own retry logic if needed

The SDK's retry logic is production-ready and handles most failure scenarios automatically! 🚀