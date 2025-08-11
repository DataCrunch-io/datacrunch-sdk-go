# Advanced Example

This example demonstrates advanced DataCrunch SDK usage patterns, including direct service creation and shared credentials testing.

## What it does

**Part 1: Shared Credentials Testing** (if `~/.datacrunch/credentials` exists)
- Tests different credential profiles (default, staging, production, development)
- Shows how to load environment-specific credentials
- Demonstrates profile-based credential management

**Part 2: Direct Service Creation**
- Creates a session using `datacrunch.NewSession()`
- Creates individual services directly from the session
- Shows all available services: Instance, InstanceTypes, Locations, SSHKeys, Volumes

## How to run

1. Set your credentials:
   ```bash
   export DATACRUNCH_CLIENT_ID="your-client-id"
   export DATACRUNCH_CLIENT_SECRET="your-client-secret"
   ```

2. Run the example:
   ```bash
   go run main.go
   ```

## Key concepts

**Shared Credentials:**
- **Profile-based credentials**: Different profiles for different environments  
- **Automatic detection**: Tests profiles only if credentials file exists
- **Environment separation**: Keep staging, production, dev credentials separate
- **Secure handling**: Masks sensitive credentials in output

**Direct Service Creation:**
- **Session-based**: Create a session once, then create individual services from it
- **Memory Efficient**: Create only the services you need
- **Fine Control**: Full control over service lifecycle and configuration
- **Shared Session**: Multiple services can share the same session

## When to use these patterns

- **Multi-environment deployments** (dev, staging, production)
- **Microservices** that only need specific DataCrunch services
- **Memory-constrained applications** 
- **Fine-grained control** over SDK component lifecycle
- **Testing different credential configurations**

## Setting up shared credentials

To test the shared credentials functionality, create `~/.datacrunch/credentials`:

```ini
[default]
client_id = your-default-client-id
client_secret = your-default-client-secret

[staging]  
client_id = your-staging-client-id
client_secret = your-staging-client-secret
base_url = https://api-staging.datacrunch.io

[production]
client_id = your-production-client-id
client_secret = your-production-client-secret
base_url = https://api.datacrunch.io
```

💡 Get your credentials from: https://datacrunch.io/account/api