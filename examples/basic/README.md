# Basic Example

This example demonstrates the simplest way to use the DataCrunch SDK with the unified client approach.

## What it does

- Creates a DataCrunch client using `datacrunch.New()`
- Lists available instance types
- Lists your current instances

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

- **Unified Client**: `datacrunch.New()` creates a client with all services available
- **Automatic Credentials**: SDK automatically finds credentials from environment variables or `~/.datacrunch/credentials`
- **Simple API**: All services accessible through the client (e.g., `client.Instance`, `client.InstanceTypes`)

This is the recommended approach for most applications.

💡 Get your credentials from: https://datacrunch.io/account/api