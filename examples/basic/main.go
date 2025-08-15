package main

import (
	"fmt"
	"log"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instanceavailability"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
)

func main() {
	fmt.Println("🚀 DataCrunch SDK - Basic Example")
	fmt.Println("==================================")
	fmt.Println()

	// Create a new session with debug mode enabled
	sess := session.New(session.WithDebug(false))
	// sess := session.New()

	// Verify credentials work
	creds := sess.GetCredentials()
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("❌ No credentials found. Please set DATACRUNCH_CLIENT_ID and DATACRUNCH_CLIENT_SECRET environment variables")
	}
	fmt.Println("✅ Client created successfully!")

	// List available instance types
	instanceTypesClient := instancetypes.New(sess)
	fmt.Println("\n💻 Listing instance types...")
	instanceTypes, err := instanceTypesClient.ListInstanceTypes()
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
	instanceClient := instance.New(sess)
	instances, err := instanceClient.ListInstances(nil)
	if err != nil {
		log.Fatalf("❌ Failed to list instances: %v", err)
		if httpErr, ok := dcerr.IsHTTPError(err); ok {
			fmt.Printf("❌ HTTP error: %v\n", httpErr)
		}
	}

	if len(instances) == 0 {
		fmt.Println("No instances found.")
	} else {
		fmt.Printf("Found %d instance(s):\n", len(instances))
		for _, inst := range instances {
			fmt.Printf("  - %s (%s): %s\n", inst.Hostname, inst.ID, inst.Status)
		}
	}

	// create an instanceavailability client
	instanceAvailabilityClient := instanceavailability.New(sess)
	instanceAvailability, err := instanceAvailabilityClient.ListInstanceAvailability()
	if err != nil {
		log.Fatalf("❌ Failed to list instance availability: %v", err)
	}

	fmt.Printf("Found %d instance availability(s):\n", len(instanceAvailability))
	for _, ia := range instanceAvailability {
		fmt.Printf("  - %s: %s\n", ia.LocationCode, ia.Availabilities)
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

💡 Get your credentials from: https://datacrunch.io/account/api
*/
