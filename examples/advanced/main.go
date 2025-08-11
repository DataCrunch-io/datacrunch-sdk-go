package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/locations"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/sshkeys"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/volumes"
)

func main() {
	fmt.Println("🚀 DataCrunch SDK - Advanced Example")
	fmt.Println("=====================================")
	fmt.Println("This example shows advanced SDK usage patterns:")
	fmt.Println("• Direct service creation from session")
	fmt.Println("• Shared credentials with different profiles")
	fmt.Println("• Fine-grained control over SDK components")
	fmt.Println()

	// Test 1: Shared credentials with profiles
	if shouldTestSharedCredentials() {
		testSharedCredentials()
		fmt.Println()
	}

	// Test 2: Direct service creation
	fmt.Println("🚀 Advanced Pattern: Direct Service Creation")
	fmt.Println("============================================")
	fmt.Println("🔧 Creating session...")
	session := datacrunch.NewSession()

	// Verify credentials work
	creds := session.GetCredentials()
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("❌ No credentials found. Please set DATACRUNCH_CLIENT_ID and DATACRUNCH_CLIENT_SECRET environment variables")
	}
	fmt.Println("✅ Session created successfully!")

	// Create individual services directly from session
	fmt.Println("\n📦 Creating individual services...")

	instanceTypesService := instancetypes.New(session)
	fmt.Println("  ✅ InstanceTypes service created")

	instanceService := instance.New(session)
	fmt.Println("  ✅ Instance service created")

	locationsService := locations.New(session)
	fmt.Println("  ✅ Locations service created")

	sshKeysService := sshkeys.New(session)
	fmt.Println("  ✅ SSHKeys service created")

	volumesService := volumes.New(session)
	fmt.Println("  ✅ Volumes service created")

	// Use the individual services
	fmt.Println("\n🌍 Using Locations service...")
	locationList, err := locationsService.ListLocations()
	if err != nil {
		log.Fatalf("❌ Failed to list locations: %v", err)
	}
	fmt.Printf("Found %d locations:\n", len(locationList))
	for _, loc := range locationList {
		fmt.Printf("  - %s (%s)\n", loc.Name, loc.Code)
	}

	fmt.Println("\n💻 Using InstanceTypes service...")
	instanceTypes, err := instanceTypesService.ListInstanceTypes()
	if err != nil {
		log.Fatalf("❌ Failed to list instance types: %v", err)
	}
	fmt.Printf("Found %d instance types:\n", len(instanceTypes))
	for i, it := range instanceTypes {
		if i >= 3 { // Show first 3 only
			fmt.Printf("... and %d more\n", len(instanceTypes)-3)
			break
		}
		fmt.Printf("  - %s: %s\n", it.InstanceType, it.Name)
	}

	fmt.Println("\n🖥️ Using Instance service...")
	instances, err := instanceService.ListInstances()
	if err != nil {
		log.Fatalf("❌ Failed to list instances: %v", err)
	}
	fmt.Printf("Found %d instances\n", len(instances))

	fmt.Println("\n🔑 Using SSHKeys service...")
	sshKeys, err := sshKeysService.ListSSHKeys()
	if err != nil {
		log.Fatalf("❌ Failed to list SSH keys: %v", err)
	}
	fmt.Printf("Found %d SSH keys\n", len(sshKeys))

	fmt.Println("\n💾 Using Volumes service...")
	volumes, err := volumesService.ListVolumes()
	if err != nil {
		log.Fatalf("❌ Failed to list volumes: %v", err)
	}
	fmt.Printf("Found %d volumes\n", len(volumes))

	fmt.Println("\n🎉 Advanced example completed!")
	fmt.Println("\n💡 Advanced patterns demonstrated:")
	fmt.Println("  - Fine-grained control over service creation")
	fmt.Println("  - Ability to share session across services")
	fmt.Println("  - Memory efficiency (create only services you need)")
	fmt.Println("  - Profile-based credential management")
	fmt.Println("  - Environment-specific configurations")
}

func shouldTestSharedCredentials() bool {
	// Check if shared credentials file exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	credFile := filepath.Join(homeDir, ".datacrunch", "credentials")
	_, err = os.Stat(credFile)
	return err == nil
}

func testSharedCredentials() {
	fmt.Println("🔐 Advanced Pattern: Shared Credentials Testing")
	fmt.Println("=============================================")
	fmt.Println("📋 Testing different credential profiles...")

	// Test profiles in order of preference
	profiles := []struct {
		name    string
		profile string
	}{
		{"default", ""},
		{"staging", "staging"},
		{"production", "production"},
		{"development", "development"},
	}

	for _, p := range profiles {
		testProfile(p.profile, p.name)
	}
}

func testProfile(profile, displayName string) {
	fmt.Printf("\n🔍 Testing %s profile...\n", displayName)

	// Create shared credentials provider
	var sharedCreds *credentials.Credentials
	if profile == "" {
		sharedCreds = credentials.NewSharedCredentials("", "")
	} else {
		sharedCreds = credentials.NewSharedCredentials("", profile)
	}

	// Try to get credentials
	credValue, err := sharedCreds.Get()
	if err != nil {
		fmt.Printf("   ⚠️  %s profile not found or invalid\n", displayName)
		return
	}

	fmt.Printf("   ✅ %s credentials loaded successfully\n", displayName)
	fmt.Printf("      Provider: %s\n", credValue.ProviderName)
	fmt.Printf("      Client ID: %s\n", maskCredential(credValue.ClientID))
	if credValue.BaseURL != "" {
		fmt.Printf("      Base URL: %s\n", credValue.BaseURL)
	}

	// Test creating a session with these credentials
	fmt.Printf("   🚀 Testing session creation with %s profile...\n", displayName)
	session := datacrunch.NewSession()
	if session != nil {
		// Create a client to verify it works
		client := datacrunch.New(datacrunch.WithCredentialsProvider(sharedCreds))

		// Quick test - try to list instance types
		instanceTypes, err := client.InstanceTypes.ListInstanceTypes()
		if err != nil {
			fmt.Printf("   ⚠️  API test failed (this is expected with invalid credentials)\n")
		} else {
			fmt.Printf("   🎉 API test succeeded! Found %d instance types\n", len(instanceTypes))
		}
	}
}

func maskCredential(credential string) string {
	if len(credential) <= 8 {
		return "***"
	}
	if len(credential) <= 12 {
		return credential[:4] + "***"
	}
	return credential[:4] + "..." + credential[len(credential)-4:]
}

/*
🚀 How to run this advanced example:

1. Set your credentials
   export DATACRUNCH_CLIENT_ID="your-client-id"
   export DATACRUNCH_CLIENT_SECRET="your-client-secret"
   or create a $HOME/.datacrunch/credentials file with the following content:
   [default]
   client_id = your-client-id
   client_secret = your-client-secret
   base_url = https://api.datacrunch.io

2. Run the example:
   go run main.go

This demonstrates advanced SDK usage patterns:

1. Direct Service Creation:
   - Create services directly from a session
   - Fine-grained control over service lifecycle
   - Memory efficient (create only what you need)

2. Shared Credentials Testing:
   - Test different credential profiles (default, staging, production)
   - Verify profile-based credential loading
   - Test API connectivity with different environments

💡 Use these patterns when you:
- Need environment-specific configurations
- Want memory efficiency in microservices
- Require fine-grained control over SDK components
- Need to test multiple credential profiles
*/
