# Changelog

All notable changes to the DataCrunch Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-01-14

### Added
- Enhanced credential provider system with flexible configuration options
- Support for custom credential chains and single provider configurations  
- New `WithDebug()` and `WithLogger()` configuration options
- Structured logging throughout the SDK using `slog`
- Better credential provider examples in `examples/advanced/main.go`
- Support for profile-based credential management with `~/.datacrunch/credentials`
- Custom unmarshal handlers for advanced API response processing

### Changed
- **BREAKING**: Moved `datacrunch/config/config.go` to `datacrunch/config.go` (package consolidation)
- **BREAKING**: Changed `Config.Debug` from `*bool` to `bool` for simpler usage
- Improved OAuth2 authentication handler with better error handling and logging
- Enhanced session management with integrated debug support
- Updated credential chain to use structured logging instead of `log.Printf`
- Better request configuration handling and BaseURL extraction

### Improved
- More detailed error messages in credential chain failures
- Flexible credential provider configuration patterns:
  - Environment variables only: `credentials.NewEnvCredentials()`
  - Shared credentials only: `credentials.NewSharedCredentials("", "profile")`  
  - Static credentials: `WithCredentials(clientID, clientSecret)`
  - Custom chains: `credentials.NewChainCredentials([]Provider{...})`
- Enhanced example documentation with multiple credential configuration patterns
- Cleaner package structure and reduced import complexity

### Removed
- Deprecated example files (`context_example.go`, `example_usage.go`, `global_example.go`)
- Old example directory structure with outdated patterns
- Unused bearer token authentication code paths

### Fixed
- OAuth2 token refresh error handling and retry logic
- Credential extraction from request configurations
- Session initialization with proper credential chain setup
- Response body reading in custom unmarshal handlers

## [1.0.0] - 2025-08-11

### Added

#### Core Features
- **DataCrunch SDK Client**: Full-featured Go SDK for DataCrunch API
- **Session-based Architecture**: Robust session management with automatic credential handling
- **Retry Functionality**: Built-in retry mechanism for API calls with configurable backoff

#### Services
- **Instance Management**: Create, list, delete, and manage compute instances
- **Instance Types**: List available instance types with pricing information
- **Instance Availability**: Check instance availability across locations
- **Locations**: List available datacenter locations
- **SSH Keys**: Manage SSH keys for instance access
- **Start Scripts**: Manage custom startup scripts for instances
- **Volumes**: Create and manage storage volumes
- **Volume Types**: List available volume types and specifications

#### Authentication & Configuration
- **Multiple Credential Providers**: Support for environment variables, credentials file, and static credentials
- **OAuth2 Integration**: Secure authentication using OAuth2 flow
- **Flexible Configuration**: Customizable client configuration with sensible defaults
- **Credential Chain**: Automatic credential discovery from multiple sources

#### Developer Experience
- **Type Safety**: Comprehensive Go interfaces and type definitions
- **Error Handling**: Structured error types with detailed error information
- **Examples**: Complete examples demonstrating SDK usage patterns
- **Test Coverage**: Comprehensive test suite with mock scenarios

### Technical Details
- **Go Version**: Compatible with Go 1.24.0+
- **Module Path**: `github.com/datacrunch-io/datacrunch-sdk-go`
- **API Protocol**: REST API with JSON payloads
- **HTTP Client**: Custom HTTP client with retry and timeout handling

### Documentation
- README with quick start guide
- API documentation and examples
- Service-specific usage patterns
- Authentication setup instructions

---

## Release Information

This is the initial stable release (v1.0.0) of the DataCrunch Go SDK. The SDK provides a complete interface to the DataCrunch platform, allowing developers to programmatically manage compute instances, storage, and other resources.

### Getting Started

```go
import "github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"

// Create client
client := datacrunch.New()

// List instance types
instanceTypes, err := client.InstanceTypes.ListInstanceTypes()
if err != nil {
    log.Fatal(err)
}
```

### Credential Setup

Set environment variables:
```bash
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"
```

For more information, see the [README](README.md) and [examples](examples/) directory.