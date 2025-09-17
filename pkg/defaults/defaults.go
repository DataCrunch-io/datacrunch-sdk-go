package defaults

import (
	"bytes"
	"fmt"
	"io"

	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
	credentials2 "github.com/datacrunch-io/datacrunch-sdk-go/pkg/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/request"
)

func Handlers() request.Handlers {
	var handlers request.Handlers

	// Add default handlers for authentication
	handlers.Validate.PushBackNamed(request.NamedHandler{
		Name: "core.ValidateCredentialsHandler",
		Fn:   ValidateCredentialsHandler,
	})

	handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "core.OAuth2AuthHandler",
		Fn:   OAuth2AuthHandler,
	})

	// Add default error handling for ALL protocols - runs FIRST in unmarshal chain
	handlers.Unmarshal.PushFront(request.NamedHandler{
		Name: "core.DefaultErrorHandler",
		Fn:   DefaultErrorHandler,
	})

	return handlers
}

// CredChain returns the default credential chain for DataCrunch
func CredChain() *credentials2.Credentials {
	return credentials2.NewChainCredentials(CredProviders())
}

// CredProviders returns the default credential providers in order of precedence
func CredProviders() []credentials2.Provider {
	return []credentials2.Provider{
		&credentials2.EnvProvider{},
		&credentials2.SharedCredentialsProvider{Filename: "", Profile: ""},
	}
}

// ValidateCredentialsHandler validates that credentials are available
func ValidateCredentialsHandler(r *request.Request) {
	if r.Config.Credentials == nil {
		r.Error = credentials2.ErrNoValidProvidersFoundInChain
	}
}

// OAuth2AuthHandler adds OAuth2 authentication to requests using credential chain
func OAuth2AuthHandler(r *request.Request) {
	// Get credentials from the request's session
	var creds *credentials2.Credentials
	var err error

	// Try to extract credentials from different sources
	if sessionCreds := r.Config.Credentials; sessionCreds != nil {
		creds = sessionCreds
	} else {
		r.Error = dcerr.New("InvalidCredentialType", "no valid credentials found in session or config", nil)
		return
	}

	// Create OAuth2Credentials wrapper for token management
	oauth2Creds := credentials2.NewOAuth2CredentialsFromProvider(creds)

	// Get a valid access token
	token, err := oauth2Creds.GetToken(r.Context())
	if err != nil {
		logger.Error("Failed to get OAuth2 token: %v", err)
		r.Error = err
		return
	}

	// Add the Authorization header
	r.HTTPRequest.Header.Set("Authorization", "Bearer "+token)
}

// DefaultErrorHandler handles HTTP error responses for ALL protocols
// This runs FIRST in the unmarshal chain, before protocol-specific unmarshaling
// When this handler sets r.Error, the request processing stops and doesn't continue to other unmarshal handlers
func DefaultErrorHandler(r *request.Request) {
	logger.Debug("DefaultErrorHandler: checking response status code %d", r.HTTPResponse.StatusCode)

	// Only handle non-success status codes
	if r.HTTPResponse.StatusCode >= 200 && r.HTTPResponse.StatusCode < 300 {
		logger.Debug("DefaultErrorHandler: success status code, skipping error handling")
		return // Continue to next handler (protocol-specific unmarshaling)
	}

	logger.Debug("DefaultErrorHandler: handling error response with status %d", r.HTTPResponse.StatusCode)

	// Read the error response body
	var errorBody string
	if r.HTTPResponse.Body != nil {
		body, err := io.ReadAll(r.HTTPResponse.Body)
		if err != nil {
			logger.Debug("DefaultErrorHandler: failed to read error response body: %v", err)
			r.Error = fmt.Errorf("status code: %d, failed to read error response body: %s", r.HTTPResponse.StatusCode, err)
			return // Stop processing - error is set
		}
		errorBody = string(body)
		logger.Debug("DefaultErrorHandler: error response body: %s", errorBody)

		// Close the original body
		if err := r.HTTPResponse.Body.Close(); err != nil {
			logger.Debug("DefaultErrorHandler: error closing response body: %v", err)
		}

		// Replace the closed body with a new reader containing the same data
		// This allows other handlers to still read the body if needed
		r.HTTPResponse.Body = io.NopCloser(bytes.NewReader(body))
	}

	// Collect request info for debugging
	requestInfo := &dcerr.RequestInfo{
		RequestURL:     r.HTTPRequest.URL.String(),
		RequestHeaders: &r.HTTPRequest.Header,
		RequestBody:    nil, // Request body is usually consumed during Build phase
	}

	// Create structured HTTP error
	r.Error = dcerr.NewHTTPError(r.HTTPResponse.StatusCode, errorBody, requestInfo)
	logger.Debug("DefaultErrorHandler: created HTTPError: %v", r.Error)
	// When r.Error is set, the request processing stops and doesn't continue to other handlers
	return
}

// SessionWithCredentials defines an interface for session-like objects with credentials
// type SessionWithCredentials interface {
// 	GetCredentials() *credentials.Credentials
// }

// // getCredentialsFromRequest extracts credentials from request config
// func getCredentialsFromRequest(r *request.Request) *credentials.Credentials {
// 	return r.Config.Credentials
// 	// // Try to extract credentials from different config types
// 	// switch cfg := r.Config.Credentials.(type) {
// 	// case SessionWithCredentials:
// 	// 	return cfg.GetCredentials()
// 	// default:
// 	// 	// Try to extract from config.Config struct using reflection
// 	// 	if cfgValue := reflect.ValueOf(r.Config); cfgValue.IsValid() {
// 	// 		// Handle both struct and pointer to struct
// 	// 		if cfgValue.Kind() == reflect.Ptr && !cfgValue.IsNil() {
// 	// 			cfgValue = cfgValue.Elem()
// 	// 		}
// 	// 		if cfgValue.Kind() == reflect.Struct {
// 	// 			// Check for Credentials field in config
// 	// 			if field := cfgValue.FieldByName("Credentials"); field.IsValid() && !field.IsNil() {
// 	// 				if creds, ok := field.Interface().(*credentials.Credentials); ok {
// 	// 					return creds
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	}
// 	// 	return nil
// 	// }
// }

// // getOAuth2CredentialsFromConfig creates OAuth2Credentials from config for backwards compatibility
// func getOAuth2CredentialsFromConfig(cfg *datacrunch.Config) *credentials.OAuth2Credentials {

// 	return NewOAuth2Credentials(clientID, clientSecret, baseURL string)
// }
