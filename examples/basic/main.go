package main

import (
	"fmt"
	"log"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
)

func main() {
	fmt.Println("🚀 DataCrunch SDK - Basic Example")
	fmt.Println("==================================")
	fmt.Println()

	// Basic usage: Create client with datacrunch.New()
	// The SDK automatically finds credentials from environment variables or ~/.datacrunch/credentials
	fmt.Println("📦 Creating DataCrunch client...")
	client := datacrunch.New()

	// Verify credentials work
	creds := client.Session.GetCredentials()
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("❌ No credentials found. Please set DATACRUNCH_CLIENT_ID and DATACRUNCH_CLIENT_SECRET environment variables")
	}
	fmt.Println("✅ Client created successfully!")

	// List available instance types
	fmt.Println("\n💻 Listing instance types...")
	instanceTypes, err := client.InstanceTypes.ListInstanceTypes()
	if err != nil {
		log.Fatalf("❌ Failed to list instance types: %v", err)
	}

	fmt.Printf("Found %d instance types:\n", len(instanceTypes))
	for i, it := range instanceTypes {
		if i >= 5 { // Show first 5 only
			fmt.Printf("... and %d more\n", len(instanceTypes)-5)
			break
		}
		fmt.Printf("  - %s: %s ($%s/hour)\n", it.InstanceType, it.Name, it.PricePerHour)
	}

	// List current instances
	fmt.Println("\n🖥️ Listing your instances...")
	instances, err := client.Instance.ListInstances()
	if err != nil {
		log.Fatalf("❌ Failed to list instances: %v", err)
	}

	if len(instances) == 0 {
		fmt.Println("No instances found.")
	} else {
		fmt.Printf("Found %d instance(s):\n", len(instances))
		for _, inst := range instances {
			fmt.Printf("  - %s (%s): %s\n", inst.Hostname, inst.ID, inst.Status)
		}
	}

	fmt.Println("\n🎉 Basic example completed!")
}

/*
🚀 How to run this basic example:

1. Set your credentials:
   export DATACRUNCH_CLIENT_ID="your-client-id"
   export DATACRUNCH_CLIENT_SECRET="your-client-secret"

2. Run the example:
   go run main.go

This demonstrates the simplest way to use the DataCrunch SDK - just call datacrunch.New()!

💡 Get your credentials from: https://datacrunch.io/account/api
*/