package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/defaults"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/request"
)

func main() {
	fmt.Println("🚀 DataCrunch SDK - Advanced Example")
	fmt.Println("=====================================")

	// Create a new session with debug mode enabled

	/**
		  1. Only Environment Variables (Skip Chain)

	  cfg := datacrunch.NewConfig(
	      datacrunch.WithBaseURL(baseURL),
	      datacrunch.WithCredentialsProvider(credentials.NewEnvCredentials()),
	  )

	  2. Only Shared Credentials File (Skip Chain)

	  cfg := datacrunch.NewConfig(
	      datacrunch.WithBaseURL(baseURL),
	      datacrunch.WithCredentialsProvider(credentials.NewSharedCredentials("", "default")),
	  )

	  3. Static Credentials (Hardcoded)

	  cfg := datacrunch.NewConfig(
	      datacrunch.WithBaseURL(baseURL),
	      datacrunch.WithCredentials("your-client-id", "your-client-secret"),
	  )

	  4. Chain of Providers (Order Matters)

		// default chains:
		cfg := datacrunch.NewConfig(
			datacrunch.WithBaseURL(baseURL),
			datacrunch.WithCredentialsProvider(defaults.CredChain()),
		)

		// or custom chain
	  cfg := datacrunch.NewConfig(
	      datacrunch.WithBaseURL(baseURL),
	      datacrunch.WithCredentialsProvider(credentials.NewChainCredentials([]credentials.Provider{
	          credentials.NewEnvCredentials(),
	          credentials.NewSharedCredentials("", "default"),
	      })),
	  )
	*/

	baseURL := "https://api-staging.datacrunch.io/v1"
	cfg := config.NewConfig(
		config.WithBaseURL(baseURL),
		config.WithDebug(true),
		config.WithCredentialsProvider(credentials.NewSharedCredentials("", "staging")), // use default credentials file location
	)

	clientInfo := metadata.ClientInfo{
		ServiceName: "instancetypes",
		APIVersion:  "v1",
		Endpoint:    *cfg.BaseURL,
	}

	// authentication handlers included in the default handlers
	handlers := defaults.Handlers()

	client := client.New(*cfg, clientInfo, handlers)

	// create request for actual API call
	op := &request.Operation{
		Name:       "ListInstanceTypes",
		HTTPMethod: "GET",
		HTTPPath:   "/instance-types",
	}

	req := client.NewRequest(op, nil, []interface{}{})

	req.Handlers.Unmarshal.PushBackNamed(request.NamedHandler{
		Name: "instance.ListInstanceTypesUnmarshal",
		Fn: func(r *request.Request) {
			fmt.Println("🚀 Advanced Pattern: Direct Service Creation")
			fmt.Println("============================================")

			// print request headers
			fmt.Println(r.HTTPRequest.Header)

			// print response headers
			fmt.Println(r.HTTPResponse.Header)

			// print response status code
			fmt.Println(r.HTTPResponse.StatusCode)

			// print reponse body
			// Read the response body
			if r.HTTPResponse.Body != nil {
				bodyBytes, err := io.ReadAll(r.HTTPResponse.Body)
				if err != nil {
					fmt.Println("Error reading body:", err)
					return
				}

				// Print raw body
				fmt.Println("Response Body:", string(bodyBytes))

				// Reset the body so other handlers can read it
				r.HTTPResponse.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}

		},
	})

	// send request
	err := req.Send()
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n✅ Advanced example finished successfully!")
	fmt.Println("\n💡 Advanced usage patterns shown:")
	fmt.Println("  - Direct, fine-grained service instantiation")
	fmt.Println("  - Session sharing across multiple services")
	fmt.Println("  - Efficient memory usage (instantiate only required services)")
	fmt.Println("  - Profile-based credential and environment management")
	fmt.Println("  - Custom configuration for different deployment environments")
}

/*
🚀 How to run this advanced example:

1. Set your credentials:
   export DATACRUNCH_CLIENT_ID="your-client-id"
   export DATACRUNCH_CLIENT_SECRET="your-client-secret"
   # Or create a $HOME/.datacrunch/credentials file:
   [default]
   client_id = your-client-id
   client_secret = your-client-secret
   base_url = https://api.datacrunch.io

2. Run the example:
   go run main.go

This demonstrates advanced SDK usage patterns:

1. Direct Service Instantiation:
   - Instantiate services directly from a session
   - Fine-grained control over service lifecycle
   - Only create the services you need for efficiency

2. Profile and Environment Management:
   - Test different credential profiles (default, staging, production)
   - Validate profile-based credential loading
   - Test API connectivity for multiple environments

💡 Use these patterns when you:
- Need to configure for multiple environments
- Want to optimize memory usage in microservices
- Require precise control over SDK components
- Need to test or switch between multiple credential profiles
*/
