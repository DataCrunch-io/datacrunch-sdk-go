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

💡 Get your credentials from: <https://datacrunch.io/account/api>

### 2. Basic usage

See **[examples/basic/](examples/basic/)** for complete working examples.

## Authentication

The SDK uses an **AWS-style credential chain** with flexible configuration:

1. **Environment variables** (highest priority)
   - `DATACRUNCH_CLIENT_ID` + `DATACRUNCH_CLIENT_SECRET`
   - `DATACRUNCH_BASE_URL` (optional)

2. **Shared credentials file**
   - Location: `~/.datacrunch/credentials`
   - Multiple profiles (default, staging, production, etc.)
   - INI format with profile sections

3. **Static credentials**
   - Programmatically configured
   - Useful for testing and development

4. **Custom credential providers**
   - Environment-only, credentials-file-only, or custom chains
   - Full control over credential resolution order

## Available Services

| Service | Description |
|---------|-------------|
| **Instance** | Manage compute instances |
| **InstanceTypes** | Query available hardware |
| **Volumes** | Persistent storage |
| **SSHKeys** | SSH key management |
| **StartScripts** | Startup automation |
| **Locations** | Datacenter regions |

## Examples

Comprehensive examples are available in the [`examples/`](examples/) directory:

- **[`basic/`](examples/basic/)** - Session-based service usage
- **[`advanced/`](examples/advanced/)** - Custom credential providers and configuration

### Run examples

```bash
# Set credentials
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"

# Run examples
cd examples/basic && go run main.go

# Set up your credential files
cd examples/advanced && go run main.go
```

## Usage Patterns

The SDK uses a **session-based approach**:

- Create a session that manages credentials and configuration
- Instantiate only the services you need
- Memory efficient with fine-grained control
- Supports all credential provider types

### Configuration Options

- **Debug logging**: Enable with `session.WithDebug(true)`
- **Custom base URLs**: For different environments
- **Flexible credential providers**: Environment, shared files, static, or custom chains
- **Profile support**: Multiple credential profiles in shared files

See [`examples/`](examples/) for detailed implementation patterns.

## Getting Help

- 📖 **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go)
- 🐛 **Issues**: [GitHub Issues](https://github.com/datacrunch-io/datacrunch-sdk-go/issues)
<!-- - 💬 **Community**: [DataCrunch Discord](https://discord.gg/datacrunch) -->
- 📧 **Support**: support@datacrunch.io

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

**🚀 Ready to get started?** Set your environment variables and create your first instance in under 5 minutes!