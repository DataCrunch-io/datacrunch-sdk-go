# DataCrunch SDK for Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/datacrunch-io/datacrunch-sdk-go.svg)](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go)

The official Go SDK for the DataCrunch API. Get up and running with DataCrunch compute instances, storage, and networking in minutes.

## Installation

```bash
go get github.com/datacrunch-io/datacrunch-sdk-go
```

## Quick Start

### 1. Set your credentials

```bash
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"
```

💡 Get your credentials from: https://datacrunch.io/account/api

### 2. Basic usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
)

func main() {
    // Create client - SDK automatically finds credentials
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

## How It Works

The SDK uses an **AWS-style credential chain** that automatically finds your credentials:

1. **Environment variables** (highest priority)
   - `DATACRUNCH_CLIENT_ID` + `DATACRUNCH_CLIENT_SECRET`

2. **Shared credentials file**
   - Location: `~/.datacrunch/credentials`
   - Supports multiple profiles (default, staging, production, etc.)

3. **Static credentials** (fallback)

## Examples

Comprehensive examples are available in the [`examples/`](examples/) directory:

- **[`basic/`](examples/basic/)** - Simple unified client usage
- **[`advanced/`](examples/advanced/)** - Direct service creation and shared credentials

### Run examples

```bash
# Set credentials
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"

# Run basic example
cd examples/basic && go run main.go

# Run advanced example  
cd examples/advanced && go run main.go
```

## Available Services

| Service | Description | 
|---------|-------------|
| **Instance** | Manage compute instances |
| **InstanceTypes** | Query available hardware |
| **Volumes** | Persistent storage |
| **SSHKeys** | SSH key management |
| **StartScripts** | Startup automation |
| **Locations** | Datacenter regions |

## Two Usage Approaches

### Basic Approach (Recommended)
```go
client := datacrunch.New()
// All services available: client.Instance, client.InstanceTypes, etc.
```

### Advanced Approach  
```go
session := datacrunch.NewSession()
instanceService := instance.New(session)
// Create only services you need
```

See [`examples/`](examples/) for detailed usage patterns and credential configurations.

## Getting Help

- 📖 **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go)
- 🐛 **Issues**: [GitHub Issues](https://github.com/datacrunch-io/datacrunch-sdk-go/issues)
- 💬 **Community**: [DataCrunch Discord](https://discord.gg/datacrunch)
- 📧 **Support**: support@datacrunch.io

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

**🚀 Ready to get started?** Set your environment variables and create your first instance in under 5 minutes!