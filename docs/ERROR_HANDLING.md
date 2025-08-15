# Error Handling in DataCrunch SDK

This guide explains how to properly handle errors when using the DataCrunch Go SDK.

## Overview

The DataCrunch SDK provides structured error handling through the `dcerr` package. All API errors are wrapped in `HTTPError` types that provide detailed information about what went wrong.

## Error Types

### HTTPError

The primary error type returned by the SDK is `dcerr.HTTPError`, which contains:

- **StatusCode**: HTTP status code (400, 401, 404, 500, etc.)
- **Body**: Raw response body from the API
- **ErrorResponse**: Structured API error response (if JSON)
- **RequestInfo**: Original request details for debugging

### APIErrorResponse  

The DataCrunch API returns structured JSON error responses:

```json
{
  "code": "invalid_request",
  "message": "Invalid ssh-key ID"
}
```

## Basic Error Handling

### Check for Errors

Always check for errors from SDK operations:

```go
sshKey, err := sshKeysClient.GetSSHKey(&sshkeys.GetSSHKeyInput{
    ID: "some-key-id",
})
if err != nil {
    // Handle the error
    fmt.Printf("Error: %s\n", err.Error())
    return
}
// Use sshKey...
```

### Check if Error is HTTPError

Use `dcerr.IsHTTPError()` to get structured error information:

```go
if httpErr, ok := dcerr.IsHTTPError(err); ok {
    fmt.Printf("HTTP Status: %d\n", httpErr.StatusCode)
    fmt.Printf("Response Body: %s\n", httpErr.Body)
    
    if httpErr.ErrorResponse != nil {
        fmt.Printf("API Error Code: %s\n", httpErr.ErrorResponse.Code)
        fmt.Printf("API Error Message: %s\n", httpErr.ErrorResponse.Message)
    }
}
```

## Helper Functions

The `dcerr` package provides convenient helper functions:

```go
// Get HTTP status code (returns 0 if not HTTP error)
statusCode := dcerr.GetStatusCode(err)

// Get API error code (returns empty string if not available)
errorCode := dcerr.GetAPIErrorCode(err)

// Get API error message (returns empty string if not available)
errorMessage := dcerr.GetAPIErrorMessage(err)

// Get raw response body (returns empty string if not HTTP error)
body := dcerr.GetErrorBody(err)
```

## Error Handling Patterns

### Handle by HTTP Status Code

Different HTTP status codes indicate different types of errors:

```go
switch dcerr.GetStatusCode(err) {
case 400, 422:
    // Client error - invalid input
    fmt.Printf("Invalid request: %s\n", dcerr.GetAPIErrorMessage(err))
case 401:
    // Authentication error
    fmt.Println("Authentication failed. Check your API credentials.")
case 403:
    // Authorization error  
    fmt.Println("Access denied. Insufficient permissions.")
case 404:
    // Resource not found
    fmt.Printf("Resource not found: %s\n", dcerr.GetAPIErrorMessage(err))
case 429:
    // Rate limiting
    fmt.Println("Rate limited. Implement exponential backoff.")
case 500, 502, 503, 504:
    // Server error - could retry
    fmt.Printf("Server error: %s\n", dcerr.GetAPIErrorMessage(err))
default:
    fmt.Printf("Unexpected error (HTTP %d): %s\n", dcerr.GetStatusCode(err), err.Error())
}
```

### Handle by API Error Code

Handle specific API error codes for precise error handling:

```go
switch dcerr.GetAPIErrorCode(err) {
case "invalid_request":
    fmt.Println("The request was malformed or invalid")
case "resource_not_found":  
    fmt.Println("The requested resource doesn't exist")
case "authentication_failed":
    fmt.Println("Check your API credentials")
case "rate_limit_exceeded":
    fmt.Println("Too many requests, please retry later")
case "insufficient_quota":
    fmt.Println("Account quota exceeded")
default:
    fmt.Printf("Unhandled error code: %s\n", dcerr.GetAPIErrorCode(err))
}
```

## Common Error Codes

The DataCrunch API uses these common error codes:

| Error Code | HTTP Status | Description | Action |
|------------|-------------|-------------|---------|
| `invalid_request` | 400 | Malformed or invalid request | Fix request parameters |
| `authentication_failed` | 401 | Invalid API credentials | Check credentials |
| `access_denied` | 403 | Insufficient permissions | Check account permissions |
| `resource_not_found` | 404 | Resource doesn't exist | Verify resource ID |
| `rate_limit_exceeded` | 429 | Too many requests | Implement backoff |
| `quota_exceeded` | 402/403 | Account quota reached | Upgrade plan or wait |
| `internal_error` | 500 | Server-side error | Retry with backoff |
| `service_unavailable` | 503 | Service temporarily down | Retry later |

## Best Practices

### 1. Always Check Errors

Never ignore errors from SDK operations:

```go
// ❌ Bad
result, _ := client.SomeOperation(input)

// ✅ Good  
result, err := client.SomeOperation(input)
if err != nil {
    // Handle error appropriately
    return fmt.Errorf("operation failed: %w", err)
}
```

### 2. Log Structured Error Information

Include useful context when logging errors:

```go
if httpErr, ok := dcerr.IsHTTPError(err); ok {
    log.Printf("API operation failed: status=%d, code=%s, message=%s",
        httpErr.StatusCode,
        dcerr.GetAPIErrorCode(err),
        dcerr.GetAPIErrorMessage(err))
} else {
    log.Printf("Non-HTTP error: %s", err.Error())
}
```

### 3. Implement Retry Logic for Recoverable Errors

For rate limits and server errors, implement exponential backoff:

```go
func retryableOperation() error {
    maxRetries := 3
    baseDelay := time.Second
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := client.SomeOperation(input)
        if err == nil {
            return nil // Success
        }
        
        statusCode := dcerr.GetStatusCode(err)
        if statusCode == 429 || statusCode >= 500 {
            // Retriable error
            delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
            time.Sleep(delay)
            continue
        }
        
        // Non-retriable error
        return err
    }
    
    return fmt.Errorf("operation failed after %d retries", maxRetries)
}
```

### 4. Handle Authentication Errors Gracefully

For authentication errors, guide users to fix their credentials:

```go
if dcerr.GetStatusCode(err) == 401 {
    return fmt.Errorf("authentication failed: please verify your DATACRUNCH_CLIENT_ID and DATACRUNCH_CLIENT_SECRET are correct")
}
```

### 5. Provide User-Friendly Error Messages

Transform technical errors into user-friendly messages:

```go
func friendlyError(err error) string {
    switch dcerr.GetAPIErrorCode(err) {
    case "invalid_request":
        return "Invalid input provided. Please check your request parameters."
    case "resource_not_found":
        return "The requested resource was not found. Please verify the ID."
    case "rate_limit_exceeded":
        return "Too many requests. Please wait a moment and try again."
    case "quota_exceeded":
        return "Account limit reached. Please upgrade your plan or contact support."
    default:
        return fmt.Sprintf("An error occurred: %s", dcerr.GetAPIErrorMessage(err))
    }
}
```

## Error Handling Example

Here's a complete example showing proper error handling:

```go
func handleSSHKeyOperation(client *sshkeys.SSHKeys, keyID string) error {
    sshKey, err := client.GetSSHKey(&sshkeys.GetSSHKeyInput{
        ID: keyID,
    })
    
    if err != nil {
        // Log the full error for debugging
        log.Printf("GetSSHKey failed for ID %s: %s", keyID, err.Error())
        
        if httpErr, ok := dcerr.IsHTTPError(err); ok {
            switch httpErr.StatusCode {
            case 400:
                return fmt.Errorf("invalid SSH key ID format: %s", keyID)
            case 401:
                return fmt.Errorf("authentication failed: please check your API credentials")
            case 403:
                return fmt.Errorf("access denied: insufficient permissions to access SSH keys")  
            case 404:
                return fmt.Errorf("SSH key not found: %s", keyID)
            case 429:
                return fmt.Errorf("rate limited: please retry after a few seconds")
            case 500:
                return fmt.Errorf("server error: please try again later")
            default:
                return fmt.Errorf("unexpected API error (HTTP %d): %s", 
                    httpErr.StatusCode, dcerr.GetAPIErrorMessage(err))
            }
        } else {
            // Network or other non-HTTP error
            return fmt.Errorf("network error: %s", err.Error())
        }
    }
    
    // Success - use the SSH key
    fmt.Printf("SSH Key: %s (%s)\n", sshKey.Name, sshKey.ID)
    return nil
}
```

## Testing Error Handling

Test your error handling with integration tests:

```go
func TestErrorHandling(t *testing.T) {
    client := sshkeys.New(session.New())
    
    // Test invalid ID
    _, err := client.GetSSHKey(&sshkeys.GetSSHKeyInput{
        ID: "invalid-id",
    })
    
    if err == nil {
        t.Fatal("Expected error for invalid ID")
    }
    
    if dcerr.GetStatusCode(err) != 400 {
        t.Errorf("Expected 400, got %d", dcerr.GetStatusCode(err))
    }
    
    if dcerr.GetAPIErrorCode(err) != "invalid_request" {
        t.Errorf("Expected 'invalid_request', got '%s'", dcerr.GetAPIErrorCode(err))
    }
}
```

This comprehensive error handling approach will help you build robust applications with the DataCrunch SDK!