# DataCrunch SDK for Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/datacrunch-io/datacrunch-sdk-go.svg)](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go)

The official Go SDK for the DataCrunch API. Get up and running with DataCrunch compute instances, storage, and networking in minutes.

## Installation

```bash
go get github.com/datacrunch-io/datacrunch-sdk-go
```

## Quick Start

The SDK uses an **AWS-style credential chain** that automatically finds your credentials from multiple sources. Just set your credentials once and start coding!

### 🚀 Basic Usage - Unified Client

```bash
# Set these environment variables
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"
```

```go
package main

import (
    "fmt"
    "log"

    "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
)

func main() {
    // Basic: Create unified client - SDK automatically finds credentials!
    client := datacrunch.New()
    
    // List instance types
    instanceTypes, err := client.InstanceTypes.ListInstanceTypes()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    fmt.Printf("Found %d instance types\n", len(instanceTypes))
    
    // List your instances
    instances, err := client.Instance.ListInstances()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    fmt.Printf("Found %d instances\n", len(instances))
}
```

### ⚡ Advanced Usage - Direct Service Creation

```go
package main

import (
    "fmt"
    "log"

    "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
    "github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
    "github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
)

func main() {
    // Advanced: Create session first, then individual services
    session := datacrunch.NewSession()
    
    // Create only the services you need
    instanceTypesService := instancetypes.New(session)
    instanceService := instance.New(session)
    
    // Use the services
    instanceTypes, err := instanceTypesService.ListInstanceTypes()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    fmt.Printf("Found %d instance types\n", len(instanceTypes))
}
```

### 📁 Alternative - Credentials File

Create `~/.datacrunch/credentials`:

```ini
[default]
client_id = your-client-id
client_secret = your-client-secret

[production]
client_id = prod-client-id
client_secret = prod-secret
```

```go
// Uses default profile automatically
client := datacrunch.New()

// Or use a specific profile
prodCreds := credentials.NewSharedCredentials("", "production")
client := datacrunch.New(datacrunch.WithCredentialsProvider(prodCreds))
```

## How Authentication Works

The SDK tries these credential sources **in order** until one succeeds:

1. **Environment Variables** (highest priority)
   - `DATACRUNCH_CLIENT_ID` + `DATACRUNCH_CLIENT_SECRET`

2. **Shared Credentials File** 
   - Location: `~/.datacrunch/credentials`
   - Supports multiple profiles like AWS

3. **Static Credentials** (fallback)
   - Hardcoded in your application

This means you can:
- 🚀 **Deploy anywhere** - Works in dev, staging, production
- 🔒 **Stay secure** - Never commit credentials to code
- 🔧 **Easy switching** - Different credentials per environment

## Real-World Examples

### Creating and Managing Instances

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
    "github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
)

func main() {
    client := datacrunch.New() // Credentials loaded automatically!
    ctx := context.Background()
    
    // Create a new instance
    createResp, err := client.Instance.CreateInstance(ctx, &instance.CreateInstanceInput{
        Name:         "my-gpu-instance",
        InstanceType: "V100.1x",
        Image:        "ubuntu-20.04",
        SSHKeyIDs:    []string{"your-ssh-key-id"},
    })
    if err != nil {
        log.Fatalf("Failed to create instance: %v", err)
    }
    
    fmt.Printf("Instance created: %s\n", createResp.ID)
    
    // Wait for instance to be ready
    for {
        instance, err := client.Instance.GetInstance(ctx, &instance.GetInstanceInput{
            ID: createResp.ID,
        })
        if err != nil {
            log.Fatalf("Error checking instance: %v", err)
        }
        
        if instance.Status == "running" {
            fmt.Printf("Instance is ready! SSH: ssh ubuntu@%s\n", instance.IPAddress)
            break
        }
        
        fmt.Println("Waiting for instance to start...")
        time.Sleep(10 * time.Second)
    }
}
```

### Working with Multiple Environments

```go
// Development
devCreds := credentials.NewSharedCredentials("", "development")
devClient := datacrunch.New(datacrunch.WithCredentialsProvider(devCreds))

// Production  
prodCreds := credentials.NewSharedCredentials("", "production")
prodClient := datacrunch.New(datacrunch.WithCredentialsProvider(prodCreds))

// Use different clients for different environments
devInstances, _ := devClient.Instance.ListInstances(ctx, nil)
prodInstances, _ := prodClient.Instance.ListInstances(ctx, nil)
```

### Complete SSH Key + Instance Workflow

```go
func main() {
    client := datacrunch.New()
    ctx := context.Background()
    
    // 1. Upload SSH key
    keyResp, err := client.SSHKeys.CreateSSHKey(ctx, &sshkeys.CreateSSHKeyInput{
        Name:      "my-laptop-key",
        PublicKey: "ssh-rsa AAAAB3NzaC1yc2E...", // Your public key
    })
    if err != nil {
        log.Fatalf("Error creating SSH key: %v", err)
    }
    
    // 2. Create instance with the SSH key
    instanceResp, err := client.Instance.CreateInstance(ctx, &instance.CreateInstanceInput{
        Name:         "my-ml-workstation",
        InstanceType: "RTX4090.1x",
        Image:        "pytorch-2.0",
        SSHKeyIDs:    []string{keyResp.ID},
        StartScript:  "#!/bin/bash\necho 'Instance ready!'",
    })
    if err != nil {
        log.Fatalf("Error creating instance: %v", err)
    }
    
    fmt.Printf("🚀 Instance created: %s\n", instanceResp.ID)
    fmt.Printf("📝 Check status: datacrunch instance get %s\n", instanceResp.ID)
}
```

## Available Services

| Service | Description | Common Use Cases |
|---------|-------------|------------------|
| **Instance** | Manage compute instances | Create GPU workstations, run training jobs |
| **InstanceTypes** | Query available hardware | Find the right GPU/CPU configuration |
| **Volumes** | Persistent storage | Data persistence, shared storage |
| **SSHKeys** | SSH key management | Secure instance access |
| **StartScripts** | Startup automation | Install software, configure environment |
| **Locations** | Datacenter regions | Choose optimal location for your workload |

## Configuration Options

### Multiple Profiles

Set up different environments in `~/.datacrunch/credentials`:

```ini
[default]
client_id = personal-account-id
client_secret = personal-secret

[company-dev]
client_id = company-dev-id
client_secret = company-dev-secret
base_url = https://dev-api.datacrunch.io

[company-prod]
client_id = company-prod-id  
client_secret = company-prod-secret
```

Use in code:
```go
// Personal account (default)
client := datacrunch.New()

// Company dev environment
devClient := datacrunch.New(
    datacrunch.WithCredentialsProvider(
        credentials.NewSharedCredentials("", "company-dev")
    ),
)
```

### Advanced - Custom Credential Chain

```go
// Create custom credential resolution order
customChain := credentials.NewChainCredentials([]credentials.Provider{
    &credentials.EnvProvider{},                                    // 1. Environment variables
    &credentials.SharedCredentialsProvider{Profile: "production"}, // 2. Production profile
    &credentials.SharedCredentialsProvider{Profile: "default"},    // 3. Default profile
})

client := datacrunch.New(datacrunch.WithCredentialsProvider(customChain))
```

### Advanced - Endpoint-Based Credentials

For enterprise setups with credential management systems:

```go
import "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"

// Fetch credentials from your credential management endpoint
endpointProvider := credentials.NewEndpointProvider(
    client.Config{BaseURL: "https://api.datacrunch.io"},
    "https://your-credential-service.com/api/v1/credentials",
    credentials.WithAuthorizationToken("Bearer your-service-token"),
)

client := datacrunch.New(datacrunch.WithCredentialsProvider(endpointProvider))
```

## Error Handling

The SDK provides structured error handling:

```go
instances, err := client.Instance.ListInstances(ctx, nil)
if err != nil {
    if dcErr, ok := err.(*dcerr.Error); ok {
        // DataCrunch API error
        fmt.Printf("API Error: %s (Code: %s)\n", dcErr.Message, dcErr.Code)
        
        switch dcErr.Code {
        case "ValidationError":
            fmt.Println("Check your input parameters")
        case "AuthenticationError": 
            fmt.Println("Check your credentials")
        case "RateLimitError":
            fmt.Println("Too many requests, please wait")
        }
    } else {
        // Network or other error
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Examples Directory

Check out practical examples in the [`examples/`](examples/) directory:

- **[`basic/`](examples/basic/)** - Simple unified client usage with `datacrunch.New()`
- **[`advanced/`](examples/advanced/)** - Direct service creation from session

### Run Examples

```bash
# Set your credentials
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"

# Run basic example (unified client)
cd examples/basic
go run main.go

# Run advanced example (direct service creation)
cd examples/advanced
go run main.go
```

## Two Approaches to Use the SDK

### 1. **Basic Approach** - Unified Client (`datacrunch.New()`)

Best for: Simple applications, getting started, when you need multiple services

```go
client := datacrunch.New()

// All services available through client
instanceTypes := client.InstanceTypes.ListInstanceTypes()
instances := client.Instance.ListInstances()
sshKeys := client.SSHKeys.ListSSHKeys()
```

**Pros:**
- Simple and straightforward
- All services available immediately
- Good for most use cases

### 2. **Advanced Approach** - Direct Service Creation

Best for: Fine-grained control, memory efficiency, service-specific applications

```go
session := datacrunch.NewSession()

// Create only the services you need
instanceService := instance.New(session)
instanceTypesService := instancetypes.New(session)

// Use individual services
instances := instanceService.ListInstances()
```

**Pros:**
- Memory efficient (create only needed services)
- Fine-grained control over service lifecycle
- Share session across multiple services
- Better for microservices architecture

## Common Use Cases

### 🤖 Machine Learning Workflows

**Basic Approach:**
```go
// Create unified client
client := datacrunch.New()

// Create GPU instance for training
instanceID, err := client.Instance.CreateInstance(&instance.CreateInstanceInput{
    Hostname:     "bert-training",
    InstanceType: "V100.1x",
    Image:        "pytorch-2.0",
    SSHKeyIDs:    []string{"your-ssh-key-id"},
})
```

**Advanced Approach:**
```go
// Create session and specific services
session := datacrunch.NewSession()
instanceService := instance.New(session)
volumeService := volumes.New(session)

// Create resources with fine-grained control
volumeID := volumeService.CreateVolume(...)
instanceID := instanceService.CreateInstance(...)
```

### 🚀 CI/CD Environments

```bash
# In your CI pipeline
export DATACRUNCH_CLIENT_ID="${CI_DATACRUNCH_CLIENT_ID}"
export DATACRUNCH_CLIENT_SECRET="${CI_DATACRUNCH_CLIENT_SECRET}"

# SDK automatically uses environment variables
go run deploy.go
```

### 🏢 Enterprise Multi-Environment

**Basic Approach:**
```go
// Different unified clients for different environments
envClients := map[string]*datacrunch.Client{
    "dev": datacrunch.New(datacrunch.WithCredentialsProvider(
        credentials.NewSharedCredentials("", "development"))),
    "prod": datacrunch.New(datacrunch.WithCredentialsProvider(
        credentials.NewSharedCredentials("", "production"))),
}

for env, client := range envClients {
    instances, _ := client.Instance.ListInstances()
    fmt.Printf("%s: %d instances\n", env, len(instances))
}
```

**Advanced Approach:**
```go
// Share session across environments with different service configurations
envSessions := map[string]*session.Session{
    "dev": datacrunch.NewSession(session.WithCredentials(...)),
    "prod": datacrunch.NewSession(session.WithCredentials(...)),
}

for env, sess := range envSessions {
    instanceService := instance.New(sess)
    instances, _ := instanceService.ListInstances()
    fmt.Printf("%s: %d instances\n", env, len(instances))
}
```

## Migration Guide

### From Legacy SDK

**Before (Legacy):**
```go
client := datacrunch.New(
    datacrunch.WithCredentials("client-id", "client-secret"),
)
```

**After - Basic Approach (Recommended):**
```bash
# Set environment variables (more secure)
export DATACRUNCH_CLIENT_ID="client-id"
export DATACRUNCH_CLIENT_SECRET="client-secret"
```
```go
client := datacrunch.New() // Credentials loaded automatically!
```

**After - Advanced Approach:**
```go
session := datacrunch.NewSession() // Credentials loaded automatically!
instanceService := instance.New(session)
```

### From Other Cloud SDKs

If you're familiar with AWS SDK, you'll feel right at home:

- Environment variables work the same way
- Credential files use similar format
- Profile-based configuration
- Automatic credential resolution

## Security Best Practices

1. ✅ **Use environment variables** in production
2. ✅ **Use credential files** for local development  
3. ✅ **Set file permissions**: `chmod 600 ~/.datacrunch/credentials`
4. ✅ **Different credentials per environment**
5. ❌ **Never commit credentials to version control**

Add to your `.gitignore`:
```gitignore
# DataCrunch credentials
.datacrunch/
**/credentials
```

## Troubleshooting

### "No valid credentials found"
```bash
# Check environment variables
env | grep DATACRUNCH

# Check credentials file
cat ~/.datacrunch/credentials

# Test credential resolution
go run -c "
creds := defaults.CredChain()
if value, err := creds.Get(); err != nil {
    log.Fatalf('Error: %v', err)  
} else {
    log.Printf('✅ Using %s provider', value.ProviderName)
}
"
```

### Permission errors
```bash
# Fix file permissions
chmod 600 ~/.datacrunch/credentials
```

### Network issues
```go
// Add timeout and retry configuration
client := datacrunch.New(
    datacrunch.WithTimeout(60*time.Second),
    datacrunch.WithRetryConfig(5, 2*time.Second, 60*time.Second),
)
```

## Getting Help

- 📖 **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go)
- 🐛 **Issues**: [GitHub Issues](https://github.com/datacrunch-io/datacrunch-sdk-go/issues)
- 💬 **Community**: [DataCrunch Discord](https://discord.gg/datacrunch)
- 📧 **Support**: support@datacrunch.io

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

**🚀 Ready to get started?** Set your environment variables and create your first instance in under 5 minutes!