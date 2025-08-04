package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
)

func main() {
	ctx := context.Background()

	// Example 1: Using functional options (recommended)
	fmt.Println("=== Example 1: Functional Options ===")
	client1 := datacrunch.New(
		datacrunch.WithBaseURL("https://api.datacrunch.io/v1"),
		datacrunch.WithCredentials("your-client-id", "your-client-secret"),
		datacrunch.WithTimeout(30*time.Second),
		datacrunch.WithRetryConfig(3, time.Second, 30*time.Second),
	)
	fmt.Println("Client created with functional options")

	// Example 2: Using environment variables
	fmt.Println("\n=== Example 2: Environment Variables ===")
	// Set these environment variables:
	// export DATACRUNCH_CLIENT_ID="your-client-id"
	// export DATACRUNCH_CLIENT_SECRET="your-client-secret"
	// export DATACRUNCH_BASE_URL="https://api.datacrunch.io/v1"
	// export DATACRUNCH_TIMEOUT="30s"
	client2 := datacrunch.NewFromEnv()
	fmt.Println("Client created from environment variables")

	// Example 3: Hybrid approach - env vars with option overrides
	fmt.Println("\n=== Example 3: Hybrid Approach ===")
	client3 := datacrunch.NewFromEnv(
		datacrunch.WithTimeout(60 * time.Second), // Override env timeout
	)
	fmt.Println("Client created from env vars with option overrides")

	// Example 4: Legacy config struct (still supported)
	fmt.Println("\n=== Example 4: Legacy Config Struct ===")
	cfg := &datacrunch.Config{
		BaseURL:      "https://api.datacrunch.io/v1",
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		Timeout:      30 * time.Second,
	}
	client4 := datacrunch.NewWithConfig(cfg)
	fmt.Println("Client created with legacy config struct")

	// Example API usage (same for all clients)
	fmt.Println("\n=== API Usage Examples ===")

	// Example usage with the instance service
	// instances, err := client1.Instance.ListInstances(ctx, &instance.ListInstancesInput{})
	// if err != nil {
	//     log.Fatalf("Failed to list instances: %v", err)
	// }

	// Example usage with SSH keys service
	// keys, err := client1.SSHKeys.ListSSHKeys(ctx, &sshkeys.ListSSHKeysInput{})
	// if err != nil {
	//     log.Fatalf("Failed to list SSH keys: %v", err)
	// }

	// Example usage with start scripts service
	// scripts, err := client1.StartScripts.ListStartScripts(ctx, &startscripts.ListStartScriptsInput{})
	// if err != nil {
	//     log.Fatalf("Failed to list start scripts: %v", err)
	// }

	fmt.Println("All DataCrunch SDK examples completed successfully!")
	log.Println("Note: Uncomment and modify the API calls above based on your specific implementation")

	// Use variables to avoid unused variable errors
	_ = ctx
	_ = client1
	_ = client2
	_ = client3
	_ = client4
}
