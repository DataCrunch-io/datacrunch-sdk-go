# DataCrunch SDK for Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/datacrunch-io/datacrunch-sdk-go.svg)](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go)

The official DataCrunch SDK for the Go programming language.

## Installation

```bash
go get github.com/datacrunch-io/datacrunch-sdk-go
```

## Quick Start

### Functional Options (Recommended)

```go
package main

import (
    "context"
    "time"

    "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
)

func main() {
    // Create client with functional options
    client := datacrunch.New(
        datacrunch.WithBaseURL("https://api.datacrunch.io/v1"),
        datacrunch.WithCredentials("your-client-id", "your-client-secret"),
        datacrunch.WithTimeout(30*time.Second),
    )
    
    ctx := context.Background()
    
    // Use the client
    // instances, err := client.Instance.ListInstances(ctx, &instance.ListInstancesInput{})
    // keys, err := client.SSHKeys.ListSSHKeys(ctx, &sshkeys.ListSSHKeysInput{})
}
```

### Environment Variables

```go
// Set environment variables:
// export DATACRUNCH_CLIENT_ID="your-client-id"
// export DATACRUNCH_CLIENT_SECRET="your-client-secret"
// export DATACRUNCH_BASE_URL="https://api.datacrunch.io/v1"

client := datacrunch.NewFromEnv()
```

### Hybrid Approach

```go
// Load from env vars, override with options
client := datacrunch.NewFromEnv(
    datacrunch.WithTimeout(60*time.Second), // Override env timeout
)
```

## Services

This SDK provides access to the following DataCrunch services:

- **Instance Management** - Create, manage, and monitor compute instances
- **Instance Types** - Query available instance configurations
- **Instance Availability** - Check instance availability across locations
- **Volume Management** - Manage storage volumes
- **Volume Types** - Query available volume types
- **SSH Keys** - Manage SSH key pairs
- **Start Scripts** - Manage instance startup scripts
- **Locations** - Query available datacenter locations

## Configuration

The SDK provides multiple ways to configure authentication and settings:

### 1. Functional Options (Recommended)

```go
client := datacrunch.New(
    datacrunch.WithBaseURL("https://api.datacrunch.io/v1"),
    datacrunch.WithCredentials("client-id", "client-secret"),
    datacrunch.WithTimeout(30*time.Second),
    datacrunch.WithRetryConfig(3, time.Second, 30*time.Second),
)
```

### 2. Environment Variables

Set these environment variables:
- `DATACRUNCH_CLIENT_ID` (required)
- `DATACRUNCH_CLIENT_SECRET` (required)  
- `DATACRUNCH_BASE_URL` (default: https://api.datacrunch.io/v1)
- `DATACRUNCH_TIMEOUT` (default: 30s, format: "30s", "1m", etc.)
- `DATACRUNCH_MAX_RETRIES` (default: 3)

```bash
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"
export DATACRUNCH_TIMEOUT="60s"
```

```go
client := datacrunch.NewFromEnv()
```

### 3. Hybrid Approach

```go
// Load from environment, override specific options
client := datacrunch.NewFromEnv(
    datacrunch.WithTimeout(60*time.Second), // Override env timeout
    datacrunch.WithRetryConfig(5, 2*time.Second, 60*time.Second),
)
```

### 4. Legacy Config Struct (Still Supported)

```go
cfg := &datacrunch.Config{
    BaseURL:      "https://api.datacrunch.io/v1",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    Timeout:      30 * time.Second,
}
client := datacrunch.NewWithConfig(cfg)
```

## Error Handling

The SDK provides structured error handling through the `dcerr` package:

```go
if err != nil {
    if dcErr, ok := err.(*dcerr.Error); ok {
        log.Printf("DataCrunch API error: %s (code: %s)", dcErr.Message, dcErr.Code)
    } else {
        log.Printf("Other error: %v", err)
    }
}
```

## Examples

### Complete Working Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
)

func main() {
    // Initialize the client with functional options
    client := datacrunch.New(
        datacrunch.WithBaseURL("https://api.datacrunch.io/v1"),
        datacrunch.WithCredentials("your-client-id", "your-client-secret"),
        datacrunch.WithTimeout(30*time.Second),
        datacrunch.WithRetryConfig(3, time.Second, 30*time.Second),
    )
    
    ctx := context.Background()
    
    // Example: Working with instances
    fmt.Println("Listing instances...")
    // instances, err := client.Instance.ListInstances(ctx, &instance.ListInstancesInput{})
    // if err != nil {
    //     if dcErr, ok := err.(*dcerr.Error); ok {
    //         log.Printf("DataCrunch API error: %s (code: %s)", dcErr.Message, dcErr.Code)
    //     } else {
    //         log.Printf("Other error: %v", err)
    //     }
    //     return
    // }
    
    // Example: Working with SSH keys
    fmt.Println("Managing SSH keys...")
    // keys, err := client.SSHKeys.ListSSHKeys(ctx, &sshkeys.ListSSHKeysInput{})
    // if err != nil {
    //     log.Fatalf("Error listing SSH keys: %v", err)
    // }
    
    // Example: Working with start scripts
    fmt.Println("Managing start scripts...")
    // scripts, err := client.StartScripts.ListStartScripts(ctx, &startscripts.ListStartScriptsInput{})
    // if err != nil {
    //     log.Fatalf("Error listing start scripts: %v", err)
    // }
    
    fmt.Println("SDK operations completed successfully!")
}
```

### Multiple Configuration Patterns

```go
// 1. Functional options (recommended)
client1 := datacrunch.New(
    datacrunch.WithCredentials("client-id", "client-secret"),
    datacrunch.WithTimeout(30*time.Second),
)

// 2. Environment variables only
client2 := datacrunch.NewFromEnv()

// 3. Environment + option overrides
client3 := datacrunch.NewFromEnv(
    datacrunch.WithTimeout(60*time.Second),
)

// 4. Legacy config struct (still supported)
cfg := datacrunch.DefaultConfig()
cfg.ClientID = "your-client-id"
cfg.ClientSecret = "your-client-secret"
client4 := datacrunch.NewWithConfig(cfg)
```

## Documentation

For detailed documentation and examples, visit [pkg.go.dev](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go).

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions, please visit our [GitHub Issues](https://github.com/datacrunch-io/datacrunch-sdk-go/issues).