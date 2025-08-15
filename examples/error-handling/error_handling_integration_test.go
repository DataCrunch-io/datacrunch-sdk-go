//go:build integration
// +build integration

package main

import (
	"testing"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/sshkeys"
)

func TestErrorHandling_InvalidSSHKeyID(t *testing.T) {
	// Create session
	sess := session.New(session.WithDebug(true))
	sshKeysClient := sshkeys.New(sess)

	// Try to get a non-existent SSH key
	invalidKeyID := "77fe26ba-e58d-4420-aab2-75e967b181b01"
	_, err := sshKeysClient.GetSSHKey(invalidKeyID)

	// We expect this to fail
	if err == nil {
		t.Fatal("Expected error for invalid SSH key ID, but got nil")
	}

	t.Logf("Got expected error: %s", err.Error())

	// Check if it's an HTTP error
	httpErr, ok := dcerr.IsHTTPError(err)
	if !ok {
		t.Fatal("Expected HTTPError, but got different error type")
	}

	// Validate error details
	if httpErr.StatusCode != 400 {
		t.Errorf("Expected HTTP 400, got %d", httpErr.StatusCode)
	}

	if httpErr.ErrorResponse == nil {
		t.Fatal("Expected structured error response, but got nil")
	}

	if httpErr.ErrorResponse.Code != "invalid_request" {
		t.Errorf("Expected error code 'invalid_request', got '%s'", httpErr.ErrorResponse.Code)
	}

	if httpErr.ErrorResponse.Message != "Invalid ssh-key ID" {
		t.Errorf("Expected message 'Invalid ssh-key ID', got '%s'", httpErr.ErrorResponse.Message)
	}

	// Test helper functions
	statusCode := dcerr.GetStatusCode(err)
	if statusCode != 400 {
		t.Errorf("dcerr.GetStatusCode() returned %d, expected 400", statusCode)
	}

	apiErrorCode := dcerr.GetAPIErrorCode(err)
	if apiErrorCode != "invalid_request" {
		t.Errorf("dcerr.GetAPIErrorCode() returned '%s', expected 'invalid_request'", apiErrorCode)
	}

	apiErrorMessage := dcerr.GetAPIErrorMessage(err)
	if apiErrorMessage != "Invalid ssh-key ID" {
		t.Errorf("dcerr.GetAPIErrorMessage() returned '%s', expected 'Invalid ssh-key ID'", apiErrorMessage)
	}

	t.Log("✅ All error handling assertions passed")
}

func TestErrorHandling_Unauthorized(t *testing.T) {
	// Create session with invalid credentials to trigger 401
	sess := session.New(
		session.WithCredentials("invalid-client-id", "invalid-client-secret"),
		session.WithDebug(false),
	)
	sshKeysClient := sshkeys.New(sess)

	// Try to list SSH keys with invalid credentials
	_, err := sshKeysClient.ListSSHKeys()

	// We expect this to fail with authentication error
	if err == nil {
		t.Fatal("Expected authentication error, but got nil")
	}

	t.Logf("Got expected error: %s", err.Error())

	// Check if it's an HTTP error
	httpErr, ok := dcerr.IsHTTPError(err)
	if !ok {
		t.Fatal("Expected HTTPError, but got different error type")
	}

	// Should be 401 Unauthorized
	if httpErr.StatusCode != 401 {
		t.Logf("Expected HTTP 401, got %d (error: %s)", httpErr.StatusCode, httpErr.Body)
		// Note: This might also be 403 depending on API implementation
		if httpErr.StatusCode != 403 {
			t.Errorf("Expected HTTP 401 or 403, got %d", httpErr.StatusCode)
		}
	}

	t.Log("✅ Authentication error handling test passed")
}

// Helper function to test different error scenarios
func TestErrorHandling_HelperFunctions(t *testing.T) {
	tests := []struct {
		name           string
		setupError     func() error
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Invalid SSH Key ID",
			setupError: func() error {
				sess := session.New(session.WithDebug(false))
				client := sshkeys.New(sess)
				_, err := client.GetSSHKey("invalid-id")
				return err
			},
			expectedStatus: 400,
			expectedCode:   "invalid_request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupError()
			if err == nil {
				t.Fatal("Expected error, but got nil")
			}

			// Test all helper functions
			statusCode := dcerr.GetStatusCode(err)
			errorCode := dcerr.GetAPIErrorCode(err)
			errorMessage := dcerr.GetAPIErrorMessage(err)
			errorBody := dcerr.GetErrorBody(err)

			t.Logf("Status Code: %d", statusCode)
			t.Logf("Error Code: %s", errorCode)
			t.Logf("Error Message: %s", errorMessage)
			t.Logf("Error Body: %s", errorBody)

			if statusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, statusCode)
			}

			if errorCode != tt.expectedCode {
				t.Errorf("Expected error code '%s', got '%s'", tt.expectedCode, errorCode)
			}

			if errorBody == "" {
				t.Error("Expected non-empty error body")
			}

			if errorMessage == "" {
				t.Error("Expected non-empty error message")
			}
		})
	}
}
