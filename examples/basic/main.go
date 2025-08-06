package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
)

func main() {
	fmt.Println("🚀 DataCrunch SDK - Basic Example")
	fmt.Println("==================================\n")

	// Step 1: Check credentials setup
	fmt.Println("📋 Step 1: Checking credential setup...")
	checkCredentialSetup()

	// Step 2: Create session (credentials loaded automatically!)
	fmt.Println("🔧 Step 2: Creating DataCrunch session...")
	sess, err := createSession()
	if err != nil {
		log.Fatalf("❌ Failed to create session: %v", err)
	}
	fmt.Println("✅ Session created successfully!")

	// Step 3: Create service clients
	fmt.Println("\n📡 Step 3: Creating API clients...")
	instanceClient := instance.New(sess)
	instanceTypesClient := instancetypes.New(sess)
	fmt.Println("✅ API clients ready!")

	// Step 4: List available instance types (helps users choose hardware)
	fmt.Println("\n💻 Step 4: Discovering available instance types...")
	listInstanceTypes(instanceTypesClient)

	// Step 5: List your current instances
	fmt.Println("\n🖥️  Step 5: Listing your current instances...")
	listInstances(instanceClient)

	fmt.Println("\n🎉 Basic example completed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("- Check examples/credential-chain/ for advanced credential configuration")
	fmt.Println("- Visit https://docs.datacrunch.io for API documentation")
	fmt.Println("- Create your first instance with the CreateInstance API")
}

// checkCredentialSetup shows users how credentials are configured
func checkCredentialSetup() {
	fmt.Println("The SDK uses an AWS-style credential chain that automatically finds your credentials:")
	fmt.Println("1. Environment variables (highest priority)")
	fmt.Println("2. ~/.datacrunch/credentials file")
	fmt.Println("3. Static credentials in code (lowest priority)")
	fmt.Println()

	// Check environment variables
	clientID := os.Getenv("DATACRUNCH_CLIENT_ID")
	clientSecret := os.Getenv("DATACRUNCH_CLIENT_SECRET")

	if clientID != "" && clientSecret != "" {
		fmt.Printf("✅ Found environment variables:\n")
		fmt.Printf("   DATACRUNCH_CLIENT_ID: %s***\n", maskCredential(clientID))
		fmt.Printf("   DATACRUNCH_CLIENT_SECRET: %s***\n", maskCredential(clientSecret))

		if baseURL := os.Getenv("DATACRUNCH_BASE_URL"); baseURL != "" {
			fmt.Printf("   DATACRUNCH_BASE_URL: %s\n", baseURL)
		}
		return
	}

	// Check credentials file
	home, _ := os.UserHomeDir()
	credFile := fmt.Sprintf("%s/.datacrunch/credentials", home)
	if _, err := os.Stat(credFile); err == nil {
		fmt.Printf("✅ Found credentials file: %s\n", credFile)
		return
	}

	// No credentials found - show setup instructions
	fmt.Println("⚠️  No credentials found! Please set up your credentials:")
	fmt.Println()
	fmt.Println("Option 1 - Environment Variables (Recommended for CI/CD):")
	fmt.Println("export DATACRUNCH_CLIENT_ID=\"your-client-id\"")
	fmt.Println("export DATACRUNCH_CLIENT_SECRET=\"your-client-secret\"")
	fmt.Println()
	fmt.Println("Option 2 - Credentials File (Recommended for local development):")
	fmt.Printf("mkdir -p %s/.datacrunch\n", home)
	fmt.Printf("cat > %s/.datacrunch/credentials << EOF\n", home)
	fmt.Println("[default]")
	fmt.Println("client_id = your-client-id")
	fmt.Println("client_secret = your-client-secret")
	fmt.Println("EOF")
	fmt.Println()
	fmt.Println("💡 Get your credentials from: https://datacrunch.io/account/api")
	fmt.Println()
}

// createSession demonstrates different ways to create a session
func createSession() (*session.Session, error) {
	// The simplest way - SDK automatically finds credentials!
	// Also gets 3 retries with exponential backoff by default - no configuration needed!
	sess := session.New()

	// Test credentials by trying to get them
	creds := sess.GetCredentials()
	credValue, err := creds.Get()
	if err != nil {
		// Show user-friendly error message
		if credErr, ok := err.(dcerr.Error); ok {
			switch credErr.Code() {
			case "NoValidProvidersFoundInChain":
				return nil, fmt.Errorf("no credentials found. Please set DATACRUNCH_CLIENT_ID and DATACRUNCH_CLIENT_SECRET environment variables, or create ~/.datacrunch/credentials file")
			default:
				return nil, fmt.Errorf("credential error: %s", credErr.Message())
			}
		}
		return nil, fmt.Errorf("failed to get credentials: %v", err)
	}

	fmt.Printf("✅ Using credentials from: %s\n", credValue.ProviderName)

	// Alternative ways to create sessions (commented out but educational):

	// Method 2: Explicitly from environment
	// sess = session.NewFromEnv()

	// Method 3: With specific options
	// sess = session.New(
	//     session.WithTimeout(60*time.Second),
	//     session.WithBaseURL("https://api.datacrunch.io/v1"),
	// )

	// Method 4: With custom retry configuration
	// sess = session.New(
	//     session.WithMaxRetries(5),        // More retries for resilience
	//     // session.WithNoRetries(),       // Or disable retries entirely
	// )

	// Method 5: With custom credential provider
	// customCreds := credentials.NewSharedCredentials("", "production")
	// sess = session.New(session.WithCredentialsProvider(customCreds))

	return sess, nil
}

// listInstanceTypes shows available hardware configurations
func listInstanceTypes(client *instancetypes.InstanceTypes) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	instanceTypeList, err := client.ListInstanceTypes(ctx)
	if err != nil {
		handleAPIError("list instance types", err)
		return
	}

	fmt.Printf("✅ Found %d available instance types:\n\n", len(instanceTypeList))

	// Show a few popular instance types
	fmt.Println("🔥 Popular GPU Instances:")
	gpuCount := 0
	for _, it := range instanceTypeList {
		if it.GPU.NumberOfGPUs > 0 && gpuCount < 3 {
			fmt.Printf("   %s - %s (%d x %s GPU, %d GB RAM) - $%s/hour\n",
				it.InstanceType,
				it.Name,
				it.GPU.NumberOfGPUs,
				it.Model,
				it.Memory.SizeInGigabytes,
				it.PricePerHour,
			)
			gpuCount++
		}
	}

	fmt.Println("\n💻 CPU-Only Instances:")
	cpuCount := 0
	for _, it := range instanceTypeList {
		if it.GPU.NumberOfGPUs == 0 && cpuCount < 2 {
			fmt.Printf("   %s - %s (%d CPU cores, %d GB RAM) - $%s/hour\n",
				it.InstanceType,
				it.Name,
				it.CPU.NumberOfCores,
				it.Memory.SizeInGigabytes,
				it.PricePerHour,
			)
			cpuCount++
		}
	}

	fmt.Printf("\n💡 See all %d instance types with: client.InstanceTypes.ListInstanceTypes()\n", len(instanceTypeList))
}

// listInstances shows your current instances
func listInstances(client *instance.Instance) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	instances, err := client.ListInstances(ctx)
	if err != nil {
		handleAPIError("list instances", err)
		return
	}

	if len(instances) == 0 {
		fmt.Println("📝 No instances found. You can create your first instance with:")
		fmt.Println("   client.Instance.CreateInstance(ctx, &instance.CreateInstanceInput{...})")
		fmt.Println()
		fmt.Println("💡 Recommended first instance:")
		fmt.Println("   - Instance Type: Pick one from the list above")
		fmt.Println("   - Image: ubuntu-20.04 or pytorch-2.0")
		fmt.Println("   - Add your SSH key for access")
		return
	}

	fmt.Printf("✅ Found %d instance(s):\n\n", len(instances))

	for _, inst := range instances {
		status := getStatusEmoji(inst.Status)
		gpuInfo := ""
		if inst.GPU.NumberOfGPUs > 0 {
			gpuInfo = fmt.Sprintf(" | %d x %s GPU", inst.GPU.NumberOfGPUs, inst.InstanceType)
		}

		fmt.Printf("   %s %s (%s)\n", status, inst.Hostname, inst.ID)
		fmt.Printf("      IP: %s | Type: %s | Location: %s%s\n",
			inst.IP, inst.InstanceType, inst.Location.Name, gpuInfo)
		fmt.Printf("      Created: %s | $%.4f/hour\n\n", inst.CreatedAt, inst.PricePerHour)
	}

	// Show helpful next steps
	runningCount := 0
	for _, inst := range instances {
		if inst.Status == "running" {
			runningCount++
		}
	}

	if runningCount > 0 {
		fmt.Printf("🎉 You have %d running instance(s)!\n", runningCount)
		fmt.Println("💡 Connect via SSH: ssh ubuntu@<instance-ip>")
	}
}

// handleAPIError provides user-friendly error messages
func handleAPIError(operation string, err error) {
	if dcErr, ok := err.(dcerr.Error); ok {
		switch dcErr.Code() {
		case "AuthenticationError":
			fmt.Printf("❌ Authentication failed while trying to %s\n", operation)
			fmt.Println("💡 Check your credentials:")
			fmt.Println("   - Verify DATACRUNCH_CLIENT_ID and DATACRUNCH_CLIENT_SECRET")
			fmt.Println("   - Get fresh credentials from: https://datacrunch.io/account/api")
		case "RateLimitError":
			fmt.Printf("⏱️  Rate limit exceeded while trying to %s\n", operation)
			fmt.Println("💡 Please wait a moment and try again")
		case "ValidationError":
			fmt.Printf("📝 Invalid request while trying to %s: %s\n", operation, dcErr.Message())
		default:
			fmt.Printf("❌ API error while trying to %s: %s (Code: %s)\n", operation, dcErr.Message(), dcErr.Code())
		}
	} else {
		fmt.Printf("❌ Network error while trying to %s: %v\n", operation, err)
		fmt.Println("💡 Check your internet connection and try again")
	}
}

// getStatusEmoji returns a friendly emoji for instance status
func getStatusEmoji(status string) string {
	switch status {
	case "running":
		return "🟢"
	case "starting", "booting":
		return "🟡"
	case "stopped", "shutdown":
		return "🔴"
	case "creating":
		return "⭕"
	default:
		return "⚪"
	}
}

// maskCredential masks sensitive credentials for safe display
func maskCredential(credential string) string {
	if len(credential) <= 8 {
		return "***"
	}
	return credential[:4] + "..."
}

/*
🚀 How to run this example:

1. Set your credentials (choose one method):

   Method A - Environment Variables:
   export DATACRUNCH_CLIENT_ID="your-client-id"
   export DATACRUNCH_CLIENT_SECRET="your-client-secret"

   Method B - Credentials File:
   mkdir -p ~/.datacrunch
   cat > ~/.datacrunch/credentials << EOF
   [default]
   client_id = your-client-id
   client_secret = your-client-secret
   EOF

2. Run the example:
   go run main.go

3. The SDK will automatically:
   ✅ Find your credentials using the credential chain
   ✅ Show you available hardware options
   ✅ List your current instances
   ✅ Give you helpful next steps

💡 What you'll learn:
- How DataCrunch credential chain works (like AWS)
- Available instance types and pricing
- Your current instances and their status
- Proper error handling for API calls

🎯 Next steps after running this example:
- Create your first instance using the API
- Set up multiple credential profiles
- Explore other services (Volumes, SSH Keys, etc.)

💬 Need help?
- Documentation: https://docs.datacrunch.io
- Discord: https://discord.gg/datacrunch
- Support: support@datacrunch.io
*/
