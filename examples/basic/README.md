# Basic DataCrunch SDK Example

This example demonstrates the fundamentals of using the DataCrunch SDK for Go. Perfect for getting started!

## What This Example Does

1. **🔐 Credential Setup** - Shows how the AWS-style credential chain works
2. **📡 API Client Creation** - Creates session and service clients  
3. **💻 List Instance Types** - Shows available hardware configurations
4. **🖥️ List Instances** - Shows your current instances with status
5. **🔧 Error Handling** - Demonstrates proper error handling patterns

## Quick Start

### 1. Set Your Credentials

Choose one method:

**Option A - Environment Variables (Recommended for CI/CD):**
```bash
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"
```

**Option B - Credentials File (Recommended for local development):**
```bash
mkdir -p ~/.datacrunch
cat > ~/.datacrunch/credentials << EOF
[default]
client_id = your-client-id
client_secret = your-client-secret
EOF
```

💡 Get your credentials from: [https://datacrunch.io/account/api](https://datacrunch.io/account/api)

### 2. Run the Example

```bash
cd examples/basic
go run main.go
```

## Expected Output

```
🚀 DataCrunch SDK - Basic Example
==================================

📋 Step 1: Checking credential setup...
✅ Found environment variables:
   DATACRUNCH_CLIENT_ID: abc1***
   DATACRUNCH_CLIENT_SECRET: xyz9***

🔧 Step 2: Creating DataCrunch session...
✅ Using credentials from: EnvProvider
✅ Session created successfully!

📡 Step 3: Creating API clients...
✅ API clients ready!

💻 Step 4: Discovering available instance types...
✅ Found 25 available instance types:

🔥 Popular GPU Instances:
   H100.1x - 1x H100 SXM (1 x H100 GPU, 200 GB RAM) - $4.50/hour
   RTX4090.1x - 1x RTX 4090 (1 x RTX4090 GPU, 128 GB RAM) - $1.20/hour
   V100.1x - 1x V100 SXM2 (1 x V100 GPU, 100 GB RAM) - $0.80/hour

💻 CPU-Only Instances:
   c1.large - CPU Optimized Large (8 CPU cores, 32 GB RAM) - $0.15/hour
   c1.xlarge - CPU Optimized XLarge (16 CPU cores, 64 GB RAM) - $0.30/hour

🖥️ Step 5: Listing your current instances...
✅ Found 2 instance(s):

   🟢 my-training-node (inst-abc123)
      IP: 192.168.1.100 | Type: RTX4090.1x | Location: Finland | 1 x RTX4090 GPU
      Created: 2024-01-15T10:30:00Z | $1.2000/hour

   🔴 backup-server (inst-def456)  
      IP: 192.168.1.101 | Type: c1.large | Location: Finland
      Created: 2024-01-14T08:15:00Z | $0.1500/hour

🎉 You have 1 running instance(s)!
💡 Connect via SSH: ssh ubuntu@<instance-ip>

🎉 Basic example completed successfully!

Next steps:
- Check examples/credential-chain/ for advanced credential configuration
- Visit https://docs.datacrunch.io for API documentation
- Create your first instance with the CreateInstance API
```

## What You'll Learn

### 🔐 Credential Management
- How the credential chain automatically finds your credentials
- Environment variables vs credentials file
- Secure credential handling best practices

### 📡 API Client Usage  
- Creating sessions and service clients
- Context and timeout handling
- Proper error handling patterns

### 💻 DataCrunch Services
- Available instance types and pricing
- Instance status and management
- Location and hardware information

## Key Code Concepts

### Automatic Credential Discovery
```go
// SDK automatically finds credentials from:
// 1. Environment variables
// 2. ~/.datacrunch/credentials file  
// 3. Static credentials in code
sess := session.New()
```

### Service Client Creation
```go
instanceClient := instance.New(sess)
instanceTypesClient := instancetypes.New(sess)
```

### API Calls with Error Handling
```go
instances, err := client.ListInstances(ctx)
if err != nil {
    if dcErr, ok := err.(*dcerr.Error); ok {
        // Handle DataCrunch API errors
        switch dcErr.Code {
        case "AuthenticationError":
            // Handle auth errors
        case "RateLimitError":  
            // Handle rate limits
        }
    }
}
```

## Common Issues & Solutions

### "No credentials found"
```bash
# Check environment variables
env | grep DATACRUNCH

# Check credentials file
cat ~/.datacrunch/credentials

# Set environment variables
export DATACRUNCH_CLIENT_ID="your-client-id"  
export DATACRUNCH_CLIENT_SECRET="your-client-secret"
```

### "Authentication failed"
- Verify your client ID and secret are correct
- Check credentials haven't expired
- Get fresh credentials from DataCrunch dashboard

### "Network timeout"
- Check internet connection
- Verify DataCrunch API is accessible
- Try increasing timeout in session options

## Next Steps

1. **Create Your First Instance**
   ```go
   instanceID, err := client.CreateInstance(ctx, &instance.CreateInstanceInput{
       InstanceType: "RTX4090.1x",
       Image:        "ubuntu-20.04", 
       SSHKeyIDs:    []string{"your-ssh-key"},
       LocationCode: "FIN-01",
   })
   ```

2. **Explore Advanced Examples**
   - [`examples/advanced/**`](../advanced/) - Provide more API examples

3. **Read the Documentation**
   - [DataCrunch API Docs](https://docs.datacrunch.io)
   - [SDK Go Reference](https://pkg.go.dev/github.com/datacrunch-io/datacrunch-sdk-go)

## Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/datacrunch-io/datacrunch-sdk-go/issues)
- 💬 **Community**: [DataCrunch Discord](https://discord.gg/datacrunch)  
- 📧 **Support**: support@datacrunch.io

---

**Ready to build?** This example shows everything you need to get started with the DataCrunch API! 🚀