package defaults

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
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

	return handlers
}

// CredChain returns the default credential chain for DataCrunch
func CredChain() *credentials.Credentials {
	return credentials.NewChainCredentials(CredProviders())
}

// CredProviders returns the default credential providers in order of precedence
func CredProviders() []credentials.Provider {
	return []credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
	}
}

// ValidateCredentialsHandler validates that credentials are available
func ValidateCredentialsHandler(r *request.Request) {
	if r.Config.Credentials == nil {
		r.Error = credentials.ErrNoValidProvidersFoundInChain
	}
}

// OAuth2AuthHandler adds OAuth2 authentication to requests using credential chain
func OAuth2AuthHandler(r *request.Request) {
	// Get credentials from the request's session
	var creds *credentials.Credentials
	var err error

	// Try to extract credentials from different sources
	if sessionCreds := r.Config.Credentials; sessionCreds != nil {
		creds = sessionCreds
	} else {
		r.Error = dcerr.New("InvalidCredentialType", "no valid credentials found in session or config", nil)
		return
	}

	// Create OAuth2Credentials wrapper for token management
	oauth2Creds := credentials.NewOAuth2CredentialsFromProvider(creds)

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
