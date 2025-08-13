package defaults

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/auth"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
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
		Fn:   auth.OAuth2AuthHandler,
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
	if r.Config == nil {
		r.Error = credentials.ErrNoValidProvidersFoundInChain
	}
}
