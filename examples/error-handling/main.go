package main

import (
	"fmt"
	"log"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/sshkeys"
)

func main() {
	fmt.Println("🚀 DataCrunch SDK - Error Handling Examples")
	fmt.Println("==========================================")
	fmt.Println()

	// Create session
	sess := session.New(session.WithDebug(false))

	// Create SSH keys service
	sshKeysClient := sshkeys.New(sess)

	fmt.Println("📋 Example 1: Handling API Errors")
	fmt.Println("Attempting to get a non-existent SSH key...")

	// Try to get a non-existent SSH key (this will return an error)
	invalidKeyID := "77fe26ba-e58d-4420-aab2-75e967b181b01" // Invalid ID
	_, err := sshKeysClient.GetSSHKey(invalidKeyID)

	if err != nil {
		fmt.Printf("❌ Error occurred: %s\n", err.Error())

		// Check if it's an HTTP error from the API
		if httpErr, ok := dcerr.IsHTTPError(err); ok {
			fmt.Printf("📊 HTTP Status Code: %d\n", httpErr.StatusCode)
			fmt.Printf("📝 Raw Response Body: %s\n", httpErr.Body)

			// Check if we have structured API error response
			if httpErr.ErrorResponse != nil {
				fmt.Printf("🔍 API Error Code: %s\n", httpErr.ErrorResponse.Code)
				fmt.Printf("💬 API Error Message: %s\n", httpErr.ErrorResponse.Message)
			}
		}

		// Helper functions for common error info
		fmt.Printf("🔧 Using helper functions:\n")
		fmt.Printf("   Status Code: %d\n", dcerr.GetStatusCode(err))
		fmt.Printf("   API Error Code: %s\n", dcerr.GetAPIErrorCode(err))
		fmt.Printf("   API Error Message: %s\n", dcerr.GetAPIErrorMessage(err))

		// Handle different error types
		switch dcerr.GetStatusCode(err) {
		case 400:
			fmt.Println("👉 This is a client error (bad request)")
		case 401:
			fmt.Println("👉 This is an authentication error")
		case 403:
			fmt.Println("👉 This is an authorization error (forbidden)")
		case 404:
			fmt.Println("👉 This is a not found error")
		case 429:
			fmt.Println("👉 This is a rate limiting error")
		case 500:
			fmt.Println("👉 This is a server error")
		default:
			fmt.Println("👉 This is an unexpected error")
		}

		// Handle specific API error codes
		switch dcerr.GetAPIErrorCode(err) {
		case "invalid_request":
			fmt.Println("🎯 Specific handling: The request was malformed or invalid")
		case "resource_not_found":
			fmt.Println("🎯 Specific handling: The requested resource doesn't exist")
		case "authentication_failed":
			fmt.Println("🎯 Specific handling: Check your API credentials")
		case "rate_limit_exceeded":
			fmt.Println("🎯 Specific handling: Too many requests, please retry later")
		default:
			fmt.Printf("🎯 Unhandled API error code: %s\n", dcerr.GetAPIErrorCode(err))
		}
	} else {
		fmt.Println("✅ No error occurred (unexpected!)")
	}

	fmt.Println()
	fmt.Println("📋 Example 2: Best Practices for Error Handling")

	// Best practice function
	handleAPIOperation := func(operation string, fn func() error) {
		fmt.Printf("🔄 Executing %s...\n", operation)

		if err := fn(); err != nil {
			// Log structured error information
			log.Printf("Operation '%s' failed: %s", operation, err.Error())

			if httpErr, ok := dcerr.IsHTTPError(err); ok {
				// For different status codes, take different actions
				switch httpErr.StatusCode {
				case 400, 422:
					// Client error - user should fix input
					fmt.Printf("❌ Invalid input for %s: %s\n", operation, dcerr.GetAPIErrorMessage(err))
					return
				case 401:
					// Authentication error - user should check credentials
					fmt.Printf("🔐 Authentication failed for %s. Please check your API credentials\n", operation)
					return
				case 403:
					// Authorization error - user doesn't have permission
					fmt.Printf("🚫 Access denied for %s. Insufficient permissions\n", operation)
					return
				case 404:
					// Resource not found
					fmt.Printf("🔍 Resource not found for %s: %s\n", operation, dcerr.GetAPIErrorMessage(err))
					return
				case 429:
					// Rate limiting - could implement retry logic
					fmt.Printf("⏰ Rate limited for %s. Consider implementing exponential backoff\n", operation)
					return
				case 500, 502, 503, 504:
					// Server error - could implement retry logic
					fmt.Printf("🔧 Server error for %s. Consider retrying: %s\n", operation, dcerr.GetAPIErrorMessage(err))
					return
				default:
					// Unexpected error
					fmt.Printf("❓ Unexpected error for %s (HTTP %d): %s\n", operation, httpErr.StatusCode, err.Error())
					return
				}
			} else {
				// Non-HTTP error (network, timeout, etc.)
				fmt.Printf("🌐 Network or client error for %s: %s\n", operation, err.Error())
				return
			}
		}

		fmt.Printf("✅ %s completed successfully\n", operation)
	}

	// Example usage
	handleAPIOperation("Get Invalid SSH Key", func() error {
		_, err := sshKeysClient.GetSSHKey("invalid-key-id")
		return err
	})

	fmt.Println()
	fmt.Println("💡 Error Handling Summary:")
	fmt.Println("1. Always check for errors from SDK operations")
	fmt.Println("2. Use dcerr.IsHTTPError() to get structured error info")
	fmt.Println("3. Check HTTP status codes for different error categories")
	fmt.Println("4. Use API error codes for specific error handling")
	fmt.Println("5. Implement appropriate retry logic for recoverable errors")
	fmt.Println("6. Log errors with sufficient context for debugging")
	fmt.Println()
	fmt.Println("🎉 Error handling examples completed!")
}
