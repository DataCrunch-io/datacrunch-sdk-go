# Advanced DataCrunch SDK Example

This example demonstrates sophisticated usage patterns and production-ready implementations with the DataCrunch SDK. Perfect for understanding how to build robust, enterprise-grade applications!

## What This Example Demonstrates

### 🎯 **Smart Resource Optimization**
- **Automatic Instance Selection** - Finds optimal GPU instances based on requirements and budget
- **Cost-Aware Location Selection** - Chooses datacenter locations based on availability and pricing
- **Hardware Requirement Matching** - Matches GPU count, memory, and performance needs

### 🔧 **Production-Ready Patterns**
- **Comprehensive Error Handling** - Specific error codes with user-friendly messages
- **Retry Logic with Exponential Backoff** - Handles rate limits and temporary failures
- **Resource Lifecycle Management** - Automatic cleanup with proper error handling
- **Structured Logging** - Progress monitoring and debugging information

### 🚀 **Real-World Workflows**
- **Complete ML Training Environment** - GPU cluster with PyTorch, Jupyter, and monitoring
- **Production Inference Server** - Docker-based deployment with FastAPI and health checks
- **Multi-Service Orchestration** - Coordinated setup across instances, storage, networking

### 📊 **Advanced Features**
- **Batch Operations** - Multiple resources created and managed together  
- **Environment-Specific Configurations** - Dev, staging, production patterns
- **Health Monitoring** - Instance startup validation and environment checks
- **Resource Tracking** - Complete audit trail of created resources

## Quick Start

### Prerequisites
```bash
# Set your credentials
export DATACRUNCH_CLIENT_ID="your-client-id"
export DATACRUNCH_CLIENT_SECRET="your-client-secret"

# Ensure you have an SSH key (for secure instance access)
ls ~/.ssh/id_*.pub
```

### Run the Advanced Example
```bash
cd examples/advanced
go run main.go
```

## What Happens When You Run It

The example executes two sophisticated scenarios:

### 📊 **Scenario 1: ML Training Cluster**
```
🎯 Scenario 1: Dev Environment
==================================================
📋 Starting workflow for: ml-training-cluster

🔍 Step 1: Analyzing requirements and selecting optimal resources...
   🔍 Finding optimal instance type (need 4 GPUs, max $5.00/hour)...
   ✅ Selected: H100.4x (4 x H100, $4.50/hour)
   🌍 Finding optimal location based on availability and cost...
   ✅ Selected: Finland (FIN-01)

🔑 Step 2: Setting up SSH access...
   🔑 Creating SSH key 'ml-training-cluster-key-1703123456' (attempt 1/3)...
   ✅ SSH key created: ssh-key-abc123

📜 Step 3: Creating startup script...
   📜 Creating startup script 'ml-training-cluster-setup-1703123456'...
   ✅ Startup script created: script-def456

💾 Step 4: Setting up persistent storage...
   💾 Creating 500GB NVMe volume 'ml-training-cluster-data-1703123456'...
   ✅ Volume created: vol-ghi789

🖥️ Step 5: Creating compute instance...
   🖥️ Creating instance 'ml-training-cluster-dev-1703123456' (attempt 1/3)...
   ✅ Instance creation initiated: inst-jkl012

📊 Step 6: Monitoring instance startup...
   ⏳ Waiting for instance inst-jkl012 to be ready...
   📊 Instance status: creating
   📊 Instance status: starting  
   📊 Instance status: running
   ✅ Instance is ready! IP: 192.168.1.100

✅ Step 7: Validating environment setup...
   🔍 Validating environment for ml-training-cluster...
   ✅ Environment validation passed
      - Instance running with IP: 192.168.1.100
      - GPU count: 4 (required: 4)
      - Cost: $4.5000/hour (budget: $5.00/hour)

🎉 Instance ready for ml-training-cluster!
   💻 SSH Access: ssh ubuntu@192.168.1.100
   🏷️ Instance ID: inst-jkl012
   🌍 Location: Finland
   💰 Cost: $4.5000/hour

✅ Scenario 1 completed successfully!
```

### 🚀 **Scenario 2: Production Inference Server**
```
⏱️ Waiting 30 seconds before next scenario...

🎯 Scenario 2: Prod Environment  
==================================================
📋 Starting workflow for: inference-server

🔍 Step 1: Analyzing requirements and selecting optimal resources...
   ✅ Using specified instance type: RTX4090.1x
   ✅ Using specified location: FIN-01

🔑 Step 2: Setting up SSH access...
   ✅ SSH key created: ssh-key-mno345

📜 Step 3: Creating startup script...
   ✅ Startup script created: script-pqr678

💾 Step 4: Setting up persistent storage...
   ✅ Volume created: vol-stu901

🖥️ Step 5: Creating compute instance...
   ✅ Instance creation initiated: inst-vwx234

📊 Step 6: Monitoring instance startup...
   ✅ Instance is ready! IP: 192.168.1.101

✅ Step 7: Validating environment setup...
   ✅ Environment validation passed

🎉 Instance ready for inference-server!
   💻 SSH Access: ssh ubuntu@192.168.1.101
   🏷️ Instance ID: inst-vwx234
   🌍 Location: Finland
   💰 Cost: $1.2000/hour

✅ Scenario 2 completed successfully!

🧹 Cleaning up resources...
   🗑️ Deleting instance: inst-jkl012
   🗑️ Deleting instance: inst-vwx234
   🗑️ Deleting volume: vol-ghi789
   🗑️ Deleting volume: vol-stu901
   🗑️ Deleting SSH key: ssh-key-abc123
   🗑️ Deleting SSH key: ssh-key-mno345
   🗑️ Deleting startup script: script-def456
   🗑️ Deleting startup script: script-pqr678
✅ All resources cleaned up successfully

🎉 All advanced scenarios completed!
```

## Key Components Deep Dive

### 🎯 **ResourceManager**
Central orchestrator that manages all DataCrunch resources:

```go
type ResourceManager struct {
    sess                *session.Session
    instanceClient      *instance.Instance
    instanceTypesClient *instancetypes.InstanceTypes
    locationsClient     *locations.Locations
    sshKeysClient       *sshkeys.SSHKey
    startScriptsClient  *startscripts.StartScripts
    volumesClient       *volumes.Volumes
    volumeTypesClient   *volumetypes.VolumeTypes
    
    // Track created resources for cleanup
    createdResources *CreatedResources
}
```

### 🔧 **Smart Configuration Optimization**
Automatically finds the best resources based on your requirements:

```go
// Find optimal instance type
bestType, err := rm.selectOptimalInstanceType(instanceTypeList, config)

// Requirements matching:
// ✅ GPU count >= requirement
// ✅ Price <= budget
// ✅ Best value (GPUs per dollar)

// Find optimal location  
bestLocation, err := rm.selectOptimalLocation(locationList, instanceType)

// Location selection factors:
// ✅ Instance type availability
// ✅ Regional pricing
// ✅ Network latency
```

### 🔄 **Retry Logic with Error Handling**
Production-ready error handling patterns:

```go
maxRetries := 3
for attempt := 1; attempt <= maxRetries; attempt++ {
    result, err := apiCall(ctx, input)
    if err == nil {
        return result, nil
    }

    // Handle specific error types
    if dcErr, ok := err.(dcerr.Error); ok {
        switch dcErr.Code() {
        case "RateLimitError":
            // Wait with exponential backoff
            delay := time.Duration(attempt*2) * time.Second
            time.Sleep(delay)
            continue
        case "InsufficientCapacity":  
            // Wait longer for resource availability
            time.Sleep(30 * time.Second)
            continue
        case "ValidationError":
            // Don't retry validation errors
            return nil, fmt.Errorf("validation failed: %v", dcErr.Message())
        }
    }
}
```

### 🏗️ **Complete ML Training Environment Setup**
Automatically configured training environment:

```bash
#!/bin/bash
# Generated ML setup script includes:

# 🐍 Python & ML Libraries
- PyTorch with CUDA support
- Transformers, Datasets, Accelerate
- NumPy, Pandas, Scikit-learn
- Matplotlib, Seaborn
- Weights & Biases, TensorBoard

# 📊 Development Environment  
- Jupyter Lab with remote access
- Git, Vim, Htop, Tmux
- NVIDIA drivers and CUDA toolkit

# 🚀 Auto-start Services
- Jupyter Lab on port 8888
- SSH access configured
- GPU monitoring tools
```

### 🐳 **Production Inference Server**
Docker-based inference server with FastAPI:

```bash
#!/bin/bash  
# Generated inference setup script includes:

# 🐳 Container Platform
- Docker with NVIDIA Container Toolkit
- Docker Compose for service orchestration
- GPU-enabled PyTorch container

# 📡 API Server
- FastAPI inference server
- Health check endpoints  
- GPU utilization monitoring
- Automatic service restart

# 📊 Monitoring
- htop, nvtop for system monitoring
- Container health checks
- Performance metrics
```

### 🧹 **Resource Lifecycle Management**
Automatic cleanup with proper error handling:

```go
func (rm *ResourceManager) Cleanup() error {
    var errors []error

    // Priority cleanup order:
    // 1. Instances (most expensive)
    // 2. Volumes (data safety)  
    // 3. SSH Keys (security)
    // 4. Scripts (least critical)

    // Collect all errors, don't fail fast
    for _, instanceID := range rm.createdResources.InstanceIDs {
        if err := rm.instanceClient.PerformInstanceAction(ctx, &instance.InstanceActionInput{
            Action: instance.InstanceActionDelete,
            ID:     instanceID,
        }); err != nil {
            errors = append(errors, fmt.Errorf("failed to delete instance %s: %v", instanceID, err))
        }
    }

    // Return all errors for debugging
    if len(errors) > 0 {
        return fmt.Errorf("cleanup errors: %v", errors)
    }
    return nil
}
```

## Configuration Options

### Environment-Specific Settings
```go
scenarios := []WorkflowConfig{
    {
        Environment:    "dev",           // Development settings
        InstanceType:   "",              // Auto-select optimal
        GPURequirement: 4,               // Need 4+ GPUs  
        MaxCostPerHour: 5.00,           // Budget constraint
        StorageSize:    500,             // 500GB storage
        SSHPublicKey:   loadSSHKey(),    // Auto-load SSH key
        SetupScript:    generateMLSetupScript(), // ML environment
    },
    {
        Environment:    "prod",          // Production settings
        InstanceType:   "RTX4090.1x",    // Specific hardware
        LocationCode:   "FIN-01",        // Specific location  
        MaxCostPerHour: 2.00,           // Lower budget
        StorageSize:    100,             // Smaller storage
        SetupScript:    generateInferenceSetupScript(), // Inference server
    },
}
```

### Hardware Optimization Logic
```go
// Smart instance selection algorithm:

candidates := filterByRequirements(instanceTypes, config)
// ✅ GPU count >= requirement
// ✅ Price <= budget
// ✅ Available in target location

best := selectBestValue(candidates)  
// 🎯 Maximize: GPUs per dollar
// 🎯 Minimize: Total cost
// 🎯 Optimize: Performance per cost
```

## Error Scenarios & Handling

### **Authentication Errors**
```
❌ Authentication failed while trying to create SSH key
💡 Check your credentials:
   - Verify DATACRUNCH_CLIENT_ID and DATACRUNCH_CLIENT_SECRET
   - Get fresh credentials from: https://datacrunch.io/account/api
```

### **Resource Capacity Issues**
```  
⚠️ Insufficient capacity, waiting before retry...
📊 Instance status: waiting for capacity
💡 Automatically retries with exponential backoff
```

### **Budget/Quota Constraints**
```
❌ No instance types found matching requirements (GPUs: 4, max cost: $2.00/hour)
💡 Increase MaxCostPerHour or reduce GPURequirement
```

### **Validation Errors**
```
📝 Invalid request while trying to create instance: SSH key not found
💡 Check that SSH keys are created before instance creation
```

## Production Usage Patterns

### **CI/CD Integration**
```bash
# In your deployment pipeline
export DATACRUNCH_CLIENT_ID="${CI_DATACRUNCH_CLIENT_ID}"
export DATACRUNCH_CLIENT_SECRET="${CI_DATACRUNCH_CLIENT_SECRET}"

# Run advanced deployment
go run examples/advanced/main.go

# Resources are automatically cleaned up on completion
```

### **Multi-Environment Management**
```go
// Different configurations per environment
environments := map[string]WorkflowConfig{
    "development": {
        GPURequirement: 1,
        MaxCostPerHour: 1.00,
        StorageSize:    50,
    },
    "staging": {
        GPURequirement: 2, 
        MaxCostPerHour: 3.00,
        StorageSize:    100,
    },
    "production": {
        GPURequirement: 8,
        MaxCostPerHour: 10.00, 
        StorageSize:    1000,
    },
}
```

### **Custom Workflow Integration**
```go
// Extend ResourceManager for your use case
type MLTrainingManager struct {
    *ResourceManager
    
    // Add your custom fields
    ExperimentID string
    ModelConfig  ModelConfig
    DataSources  []DataSource
}

// Add custom methods
func (m *MLTrainingManager) SetupDistributedTraining() error {
    // Your custom logic here
}
```

## Key Learnings

### 🎯 **Production Patterns**
- Always use retry logic with exponential backoff
- Implement comprehensive error handling with specific error codes
- Track all created resources for proper cleanup
- Use structured logging for debugging and monitoring

### 🔧 **Resource Optimization**
- Automatically select optimal hardware based on requirements
- Consider cost, performance, and availability when selecting resources
- Use environment-specific configurations for different deployment stages
- Implement health checks and validation after resource creation

### 🚀 **Enterprise Workflows**
- Orchestrate multiple services together for complete solutions  
- Use startup scripts for automated environment configuration
- Implement proper security with SSH keys and access controls
- Plan for disaster recovery with proper cleanup procedures

## Next Steps

1. **Customize for Your Use Case**
   - Modify the WorkflowConfig for your specific requirements
   - Add custom validation and health checks
   - Integrate with your existing deployment pipelines

2. **Extend the Example**  
   - Add monitoring and alerting integration
   - Implement blue-green deployments
   - Add database and networking configuration

3. **Production Deployment**
   - Set up proper credential management
   - Implement logging and monitoring
   - Add automated testing and validation

## Support

- 📖 **Documentation**: [DataCrunch API Docs](https://docs.datacrunch.io)
- 💬 **Community**: [DataCrunch Discord](https://discord.gg/datacrunch)
- 🐛 **Issues**: [GitHub Issues](https://github.com/datacrunch-io/datacrunch-sdk-go/issues)

---

**🚀 Ready for production?** This example shows you everything needed to build robust, enterprise-grade applications with the DataCrunch API!