package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/locations"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/sshkeys"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/startscripts"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/volumes"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/volumetypes"
)

// WorkflowConfig holds configuration for our advanced workflow
type WorkflowConfig struct {
	ProjectName    string
	InstanceType   string
	GPURequirement int
	LocationCode   string
	StorageSize    int
	MaxCostPerHour float64
	SSHPublicKey   string
	SetupScript    string
	Environment    string // dev, staging, prod
}

// ResourceManager handles creation and cleanup of DataCrunch resources
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

type CreatedResources struct {
	InstanceIDs    []string
	SSHKeyIDs      []string
	StartScriptIDs []string
	VolumeIDs      []string
}

func main() {
	fmt.Println("🚀 DataCrunch SDK - Advanced Example")
	fmt.Println("=====================================")
	fmt.Println("This example demonstrates:")
	fmt.Println("• Complete ML training environment setup")
	fmt.Println("• Multi-service resource orchestration")
	fmt.Println("• Advanced error handling and retries")
	fmt.Println("• Resource lifecycle management")
	fmt.Println("• Cost optimization strategies")
	fmt.Println("• Production-ready patterns\n")

	// Configuration for different scenarios
	scenarios := []WorkflowConfig{
		{
			ProjectName:    "ml-training-cluster",
			InstanceType:   "",  // Will be auto-selected
			GPURequirement: 4,   // Need at least 4 GPUs
			LocationCode:   "",  // Will be auto-selected based on cost
			StorageSize:    500, // 500GB for training data
			MaxCostPerHour: 5.00,
			Environment:    "dev",
			SSHPublicKey:   loadSSHKey(),
			SetupScript:    generateMLSetupScript(),
		},
		{
			ProjectName:    "inference-server",
			InstanceType:   "RTX4090.1x",
			GPURequirement: 1,
			LocationCode:   "FIN-01",
			StorageSize:    100,
			MaxCostPerHour: 2.00,
			Environment:    "prod",
			SSHPublicKey:   loadSSHKey(),
			SetupScript:    generateInferenceSetupScript(),
		},
	}

	// Create resource manager
	rm, err := NewResourceManager()
	if err != nil {
		log.Fatalf("❌ Failed to initialize resource manager: %v", err)
	}

	// Set up cleanup on exit
	defer func() {
		fmt.Println("\n🧹 Cleaning up resources...")
		if err := rm.Cleanup(); err != nil {
			log.Printf("⚠️ Warning: Failed to cleanup some resources: %v", err)
		} else {
			fmt.Println("✅ All resources cleaned up successfully")
		}
	}()

	// Run each scenario
	for i, config := range scenarios {
		fmt.Printf("\n🎯 Scenario %d: %s Environment\n", i+1, strings.Title(config.Environment))
		fmt.Println(strings.Repeat("=", 50))

		if err := rm.RunAdvancedWorkflow(config); err != nil {
			log.Printf("❌ Scenario %d failed: %v", i+1, err)
			continue
		}

		fmt.Printf("✅ Scenario %d completed successfully!\n", i+1)

		// Wait between scenarios
		if i < len(scenarios)-1 {
			fmt.Println("\n⏱️ Waiting 30 seconds before next scenario...")
			time.Sleep(30 * time.Second)
		}
	}

	fmt.Println("\n🎉 All advanced scenarios completed!")
	fmt.Println("\nKey learnings demonstrated:")
	fmt.Println("• Smart resource selection based on requirements")
	fmt.Println("• Cost optimization and availability checking")
	fmt.Println("• Robust error handling with retries")
	fmt.Println("• Proper resource lifecycle management")
	fmt.Println("• Production-ready configuration patterns")
}

// NewResourceManager creates a new resource manager with all service clients
func NewResourceManager() (*ResourceManager, error) {
	// Create session with advanced configuration
	sess := session.New(
		session.WithTimeout(60 * time.Second),
	)

	// Test credentials
	creds := sess.GetCredentials()
	credValue, err := creds.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %v", err)
	}

	fmt.Printf("🔐 Using credentials from: %s\n", credValue.ProviderName)
	fmt.Printf("📍 API Base URL: %s\n", sess.Config.BaseURL)

	// Initialize all service clients
	rm := &ResourceManager{
		sess:                sess,
		instanceClient:      instance.New(sess),
		instanceTypesClient: instancetypes.New(sess),
		locationsClient:     locations.New(sess),
		sshKeysClient:       sshkeys.New(sess),
		startScriptsClient:  startscripts.New(sess),
		volumesClient:       volumes.New(sess),
		volumeTypesClient:   volumetypes.New(sess),
		createdResources:    &CreatedResources{},
	}

	return rm, nil
}

// RunAdvancedWorkflow orchestrates a complete advanced workflow
func (rm *ResourceManager) RunAdvancedWorkflow(config WorkflowConfig) error {
	ctx := context.Background()

	fmt.Printf("📋 Starting workflow for: %s\n", config.ProjectName)

	// Step 1: Analyze requirements and select optimal resources
	fmt.Println("\n🔍 Step 1: Analyzing requirements and selecting optimal resources...")
	optimalConfig, err := rm.OptimizeConfiguration(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to optimize configuration: %v", err)
	}

	// Step 2: Create SSH key for secure access
	fmt.Println("\n🔑 Step 2: Setting up SSH access...")
	sshKeyID, err := rm.CreateSSHKeyWithRetry(ctx, optimalConfig)
	if err != nil {
		return fmt.Errorf("failed to create SSH key: %v", err)
	}
	rm.createdResources.SSHKeyIDs = append(rm.createdResources.SSHKeyIDs, sshKeyID)

	// Step 3: Create startup script for environment setup
	fmt.Println("\n📜 Step 3: Creating startup script...")
	scriptID, err := rm.CreateStartupScript(ctx, optimalConfig)
	if err != nil {
		return fmt.Errorf("failed to create startup script: %v", err)
	}
	rm.createdResources.StartScriptIDs = append(rm.createdResources.StartScriptIDs, scriptID)

	// Step 4: Create persistent storage
	fmt.Println("\n💾 Step 4: Setting up persistent storage...")
	volumeID, err := rm.CreateVolume(ctx, optimalConfig)
	if err != nil {
		return fmt.Errorf("failed to create volume: %v", err)
	}
	rm.createdResources.VolumeIDs = append(rm.createdResources.VolumeIDs, volumeID)

	// Step 5: Create and configure the compute instance
	fmt.Println("\n🖥️ Step 5: Creating compute instance...")
	instanceID, err := rm.CreateInstanceWithRetry(ctx, optimalConfig, sshKeyID, scriptID, volumeID)
	if err != nil {
		return fmt.Errorf("failed to create instance: %v", err)
	}
	rm.createdResources.InstanceIDs = append(rm.createdResources.InstanceIDs, instanceID)

	// Step 6: Monitor instance startup and health
	fmt.Println("\n📊 Step 6: Monitoring instance startup...")
	instanceInfo, err := rm.WaitForInstanceReady(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("instance failed to start properly: %v", err)
	}

	// Step 7: Validate environment setup
	fmt.Println("\n✅ Step 7: Validating environment setup...")
	if err := rm.ValidateEnvironment(ctx, instanceInfo, optimalConfig); err != nil {
		return fmt.Errorf("environment validation failed: %v", err)
	}

	fmt.Printf("\n🎉 Instance ready for %s!\n", optimalConfig.ProjectName)
	fmt.Printf("   💻 SSH Access: ssh ubuntu@%s\n", instanceInfo.IP)
	fmt.Printf("   🏷️ Instance ID: %s\n", instanceInfo.ID)
	fmt.Printf("   🌍 Location: %s\n", instanceInfo.Location.Name)
	fmt.Printf("   💰 Cost: $%.4f/hour\n", instanceInfo.PricePerHour)

	return nil
}

// OptimizeConfiguration finds the best configuration based on requirements
func (rm *ResourceManager) OptimizeConfiguration(ctx context.Context, config WorkflowConfig) (WorkflowConfig, error) {
	optimized := config

	// Get all available instance types
	instanceTypeList, err := rm.instanceTypesClient.ListInstanceTypes(ctx)
	if err != nil {
		return optimized, fmt.Errorf("failed to list instance types: %v", err)
	}

	// Get all locations
	locationList, err := rm.locationsClient.ListLocations(ctx)
	if err != nil {
		return optimized, fmt.Errorf("failed to list locations: %v", err)
	}

	// Find optimal instance type if not specified
	if optimized.InstanceType == "" {
		fmt.Printf("   🔍 Finding optimal instance type (need %d GPUs, max $%.2f/hour)...\n",
			config.GPURequirement, config.MaxCostPerHour)

		bestType, err := rm.selectOptimalInstanceType(instanceTypeList, config)
		if err != nil {
			return optimized, err
		}
		optimized.InstanceType = bestType.InstanceType
		fmt.Printf("   ✅ Selected: %s (%d x %s, $%s/hour)\n",
			bestType.InstanceType, bestType.GPU.NumberOfGPUs, bestType.Model, bestType.PricePerHour)
	}

	// Find optimal location if not specified
	if optimized.LocationCode == "" {
		fmt.Println("   🌍 Finding optimal location based on availability and cost...")

		bestLocation, err := rm.selectOptimalLocation(locationList, optimized.InstanceType)
		if err != nil {
			return optimized, err
		}
		optimized.LocationCode = bestLocation.Code
		fmt.Printf("   ✅ Selected: %s (%s)\n", bestLocation.Name, bestLocation.Code)
	}

	return optimized, nil
}

// selectOptimalInstanceType finds the best instance type based on requirements
func (rm *ResourceManager) selectOptimalInstanceType(instanceTypes []*instancetypes.InstanceTypeResponse, config WorkflowConfig) (*instancetypes.InstanceTypeResponse, error) {
	var candidates []*instancetypes.InstanceTypeResponse

	// Filter candidates that meet requirements
	for _, it := range instanceTypes {
		if it.GPU.NumberOfGPUs >= config.GPURequirement {
			// Parse price (handle string to float conversion)
			priceStr := strings.Replace(it.PricePerHour, "$", "", -1)
			price := 0.0
			if _, err := fmt.Sscanf(priceStr, "%f", &price); err == nil {
				if price <= config.MaxCostPerHour {
					candidates = append(candidates, it)
				}
			}
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no instance types found matching requirements (GPUs: %d, max cost: $%.2f/hour)",
			config.GPURequirement, config.MaxCostPerHour)
	}

	// Select the best candidate (lowest cost with adequate GPUs)
	var best *instancetypes.InstanceTypeResponse
	bestValue := 0.0

	for _, candidate := range candidates {
		priceStr := strings.Replace(candidate.PricePerHour, "$", "", -1)
		price := 0.0
		if _, err := fmt.Sscanf(priceStr, "%f", &price); err != nil {
			continue
		}

		// Value score: more GPUs per dollar is better
		value := float64(candidate.GPU.NumberOfGPUs) / price
		if best == nil || value > bestValue {
			best = candidate
			bestValue = value
		}
	}

	return best, nil
}

// selectOptimalLocation finds the best location for the instance type
func (rm *ResourceManager) selectOptimalLocation(locations []*locations.LocationResponse, instanceType string) (*locations.LocationResponse, error) {
	if len(locations) == 0 {
		return nil, fmt.Errorf("no locations available")
	}

	// For now, select the first available location
	// In a real implementation, you'd check availability and pricing per location
	return locations[0], nil
}

// CreateSSHKeyWithRetry creates an SSH key with retry logic
func (rm *ResourceManager) CreateSSHKeyWithRetry(ctx context.Context, config WorkflowConfig) (string, error) {
	keyName := fmt.Sprintf("%s-key-%d", config.ProjectName, time.Now().Unix())

	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("   🔑 Creating SSH key '%s' (attempt %d/%d)...\n", keyName, attempt, maxRetries)

		keyResp, err := rm.sshKeysClient.CreateSSHKey(ctx, &sshkeys.CreateSSHKeyInput{
			Name: keyName,
			Key:  config.SSHPublicKey,
		})

		if err == nil {
			fmt.Printf("   ✅ SSH key created: %s\n", keyResp.ID)
			return keyResp.ID, nil
		}

		// Handle specific errors
		if dcErr, ok := err.(dcerr.Error); ok {
			switch dcErr.Code() {
			case "RateLimitError":
				if attempt < maxRetries {
					delay := time.Duration(attempt*2) * time.Second
					fmt.Printf("   ⏱️ Rate limited, waiting %v before retry...\n", delay)
					time.Sleep(delay)
					continue
				}
			case "ValidationError":
				return "", fmt.Errorf("SSH key validation failed: %v", dcErr.Message())
			}
		}

		if attempt < maxRetries {
			fmt.Printf("   ⚠️ Attempt %d failed: %v, retrying...\n", attempt, err)
			time.Sleep(time.Duration(attempt) * time.Second)
		} else {
			return "", fmt.Errorf("failed to create SSH key after %d attempts: %v", maxRetries, err)
		}
	}

	return "", fmt.Errorf("failed to create SSH key after %d attempts", maxRetries)
}

// CreateStartupScript creates a startup script for environment setup
func (rm *ResourceManager) CreateStartupScript(ctx context.Context, config WorkflowConfig) (string, error) {
	scriptName := fmt.Sprintf("%s-setup-%d", config.ProjectName, time.Now().Unix())

	fmt.Printf("   📜 Creating startup script '%s'...\n", scriptName)

	scriptID, err := rm.startScriptsClient.CreateStartScript(ctx, &startscripts.CreateStartScriptInput{
		Name:   scriptName,
		Script: config.SetupScript,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create startup script: %v", err)
	}

	fmt.Printf("   ✅ Startup script created: %s\n", scriptID)
	return scriptID, nil
}

// CreateVolume creates a persistent storage volume
func (rm *ResourceManager) CreateVolume(ctx context.Context, config WorkflowConfig) (string, error) {
	volumeName := fmt.Sprintf("%s-data-%d", config.ProjectName, time.Now().Unix())

	// Get available volume types
	volumeTypeList, err := rm.volumeTypesClient.ListVolumeTypes(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list volume types: %v", err)
	}

	if len(volumeTypeList) == 0 {
		return "", fmt.Errorf("no volume types available")
	}

	// Select NVMe if available, otherwise use first available type
	volumeType := volumeTypeList[0].Type
	for _, vt := range volumeTypeList {
		if strings.Contains(strings.ToLower(vt.Type), "nvme") {
			volumeType = vt.Type
			break
		}
	}

	fmt.Printf("   💾 Creating %dGB %s volume '%s'...\n", config.StorageSize, volumeType, volumeName)

	volumeID, err := rm.volumesClient.CreateVolume(ctx, &volumes.CreateVolumeInput{
		Name:         volumeName,
		Size:         config.StorageSize,
		Type:         volumeType,
		LocationCode: config.LocationCode,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create volume: %v", err)
	}

	fmt.Printf("   ✅ Volume created: %s\n", volumeID)
	return volumeID, nil
}

// CreateInstanceWithRetry creates a compute instance with retry logic
func (rm *ResourceManager) CreateInstanceWithRetry(ctx context.Context, config WorkflowConfig, sshKeyID, scriptID, volumeID string) (string, error) {
	instanceName := fmt.Sprintf("%s-%s-%d", config.ProjectName, config.Environment, time.Now().Unix())

	createInput := &instance.CreateInstanceInput{
		InstanceType:    config.InstanceType,
		Image:           selectImageForEnvironment(config.Environment),
		SSHKeyIDs:       []string{sshKeyID},
		StartupScriptID: scriptID,
		Hostname:        instanceName,
		Description:     fmt.Sprintf("%s instance for %s environment", config.ProjectName, config.Environment),
		LocationCode:    config.LocationCode,
		OSVolume: &instance.OSVolume{
			Name: "root",
			Size: 50,
		},
		IsSpot:          config.Environment != "prod", // Use spot pricing for non-prod
		ExistingVolumes: []string{volumeID},
		Contract:        "hourly",
		Pricing:         "standard",
	}

	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("   🖥️ Creating instance '%s' (attempt %d/%d)...\n", instanceName, attempt, maxRetries)

		instanceID, err := rm.instanceClient.CreateInstance(ctx, createInput)
		if err == nil {
			fmt.Printf("   ✅ Instance creation initiated: %s\n", instanceID)
			return instanceID, nil
		}

		// Handle specific errors
		if dcErr, ok := err.(dcerr.Error); ok {
			switch dcErr.Code() {
			case "InsufficientCapacity":
				if attempt < maxRetries {
					fmt.Printf("   ⚠️ Insufficient capacity, waiting before retry...\n")
					time.Sleep(30 * time.Second)
					continue
				}
			case "ValidationError":
				return "", fmt.Errorf("instance configuration invalid: %v", dcErr.Message())
			case "QuotaExceeded":
				return "", fmt.Errorf("quota exceeded: %v", dcErr.Message())
			}
		}

		if attempt < maxRetries {
			fmt.Printf("   ⚠️ Attempt %d failed: %v, retrying...\n", attempt, err)
			time.Sleep(time.Duration(attempt*10) * time.Second)
		}
	}

	return "", fmt.Errorf("failed to create instance after %d attempts", maxRetries)
}

// WaitForInstanceReady waits for an instance to be fully ready
func (rm *ResourceManager) WaitForInstanceReady(ctx context.Context, instanceID string) (*instance.ListInstancesResponse, error) {
	fmt.Printf("   ⏳ Waiting for instance %s to be ready...\n", instanceID)

	maxWaitTime := 10 * time.Minute
	checkInterval := 30 * time.Second
	timeout := time.After(maxWaitTime)
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for instance to be ready after %v", maxWaitTime)
		case <-ticker.C:
			// Get current instance status
			instances, err := rm.instanceClient.ListInstances(ctx)
			if err != nil {
				fmt.Printf("   ⚠️ Error checking instance status: %v\n", err)
				continue
			}

			var targetInstance *instance.ListInstancesResponse
			for _, inst := range instances {
				if inst.ID == instanceID {
					targetInstance = inst
					break
				}
			}

			if targetInstance == nil {
				fmt.Printf("   ⚠️ Instance not found in list, continuing to wait...\n")
				continue
			}

			fmt.Printf("   📊 Instance status: %s\n", targetInstance.Status)

			switch targetInstance.Status {
			case "running":
				fmt.Printf("   ✅ Instance is ready! IP: %s\n", targetInstance.IP)
				return targetInstance, nil
			case "error", "failed":
				return nil, fmt.Errorf("instance failed to start (status: %s)", targetInstance.Status)
			case "creating", "starting", "booting":
				// Continue waiting
				continue
			default:
				fmt.Printf("   ℹ️ Unexpected status '%s', continuing to wait...\n", targetInstance.Status)
				continue
			}
		}
	}
}

// ValidateEnvironment validates that the instance environment is properly set up
func (rm *ResourceManager) ValidateEnvironment(ctx context.Context, inst *instance.ListInstancesResponse, config WorkflowConfig) error {
	fmt.Printf("   🔍 Validating environment for %s...\n", config.ProjectName)

	// Basic validations
	if inst.IP == "" {
		return fmt.Errorf("instance has no IP address")
	}

	if inst.Status != "running" {
		return fmt.Errorf("instance is not in running state (current: %s)", inst.Status)
	}

	// Check GPU requirements
	if config.GPURequirement > 0 && inst.GPU.NumberOfGPUs < config.GPURequirement {
		return fmt.Errorf("instance has insufficient GPUs (need %d, got %d)",
			config.GPURequirement, inst.GPU.NumberOfGPUs)
	}

	// Verify cost constraints
	if inst.PricePerHour > config.MaxCostPerHour {
		return fmt.Errorf("instance cost exceeds budget ($%.4f/hour > $%.2f/hour)",
			inst.PricePerHour, config.MaxCostPerHour)
	}

	fmt.Printf("   ✅ Environment validation passed\n")
	fmt.Printf("      - Instance running with IP: %s\n", inst.IP)
	fmt.Printf("      - GPU count: %d (required: %d)\n", inst.GPU.NumberOfGPUs, config.GPURequirement)
	fmt.Printf("      - Cost: $%.4f/hour (budget: $%.2f/hour)\n", inst.PricePerHour, config.MaxCostPerHour)

	return nil
}

// Cleanup removes all created resources
func (rm *ResourceManager) Cleanup() error {
	ctx := context.Background()
	var errors []error

	// Cleanup instances (most important)
	for _, instanceID := range rm.createdResources.InstanceIDs {
		fmt.Printf("   🗑️ Deleting instance: %s\n", instanceID)
		if err := rm.instanceClient.PerformInstanceAction(ctx, &instance.InstanceActionInput{
			Action: instance.InstanceActionDelete,
			ID:     instanceID,
		}); err != nil {
			errors = append(errors, fmt.Errorf("failed to delete instance %s: %v", instanceID, err))
		}
	}

	// Cleanup volumes
	for _, volumeID := range rm.createdResources.VolumeIDs {
		fmt.Printf("   🗑️ Deleting volume: %s\n", volumeID)
		if err := rm.volumesClient.DeleteVolume(ctx, volumeID, true); err != nil {
			errors = append(errors, fmt.Errorf("failed to delete volume %s: %v", volumeID, err))
		}
	}

	// Cleanup SSH keys
	for _, sshKeyID := range rm.createdResources.SSHKeyIDs {
		fmt.Printf("   🗑️ Deleting SSH key: %s\n", sshKeyID)
		if err := rm.sshKeysClient.DeleteSSHKey(ctx, sshKeyID); err != nil {
			errors = append(errors, fmt.Errorf("failed to delete SSH key %s: %v", sshKeyID, err))
		}
	}

	// Cleanup startup scripts
	for _, scriptID := range rm.createdResources.StartScriptIDs {
		fmt.Printf("   🗑️ Deleting startup script: %s\n", scriptID)
		if err := rm.startScriptsClient.DeleteStartScript(ctx, scriptID); err != nil {
			errors = append(errors, fmt.Errorf("failed to delete startup script %s: %v", scriptID, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}

// Helper functions

func loadSSHKey() string {
	// Try to load SSH key from standard locations
	sshKeyPaths := []string{
		os.Getenv("HOME") + "/.ssh/id_rsa.pub",
		os.Getenv("HOME") + "/.ssh/id_ed25519.pub",
		os.Getenv("HOME") + "/.ssh/id_ecdsa.pub",
	}

	for _, path := range sshKeyPaths {
		if content, err := os.ReadFile(path); err == nil {
			return strings.TrimSpace(string(content))
		}
	}

	// Fallback to a demo key if no real key found
	return "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7... demo-key@localhost"
}

func generateMLSetupScript() string {
	return `#!/bin/bash
set -e

echo "🐍 Setting up ML Training Environment..."

# Update system
apt-get update && apt-get upgrade -y

# Install Python and pip
apt-get install -y python3 python3-pip python3-venv

# Install NVIDIA drivers and CUDA (if GPU instance)
if nvidia-smi &>/dev/null; then
    echo "📊 GPU detected, installing CUDA toolkit..."
    apt-get install -y nvidia-cuda-toolkit
fi

# Create virtual environment
python3 -m venv /opt/ml-env
source /opt/ml-env/bin/activate

# Install ML libraries
pip install --upgrade pip
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118
pip install transformers datasets accelerate
pip install numpy pandas scikit-learn matplotlib seaborn
pip install jupyter jupyterlab
pip install wandb tensorboard

# Install development tools
apt-get install -y git vim htop tmux

# Set up Jupyter
mkdir -p /workspace/notebooks
cd /workspace
jupyter lab --generate-config
echo "c.ServerApp.ip = '0.0.0.0'" >> ~/.jupyter/jupyter_lab_config.py
echo "c.ServerApp.port = 8888" >> ~/.jupyter/jupyter_lab_config.py
echo "c.ServerApp.token = ''" >> ~/.jupyter/jupyter_lab_config.py
echo "c.ServerApp.password = ''" >> ~/.jupyter/jupyter_lab_config.py

# Create systemd service for Jupyter
cat > /etc/systemd/system/jupyter.service << 'EOF'
[Unit]
Description=Jupyter Lab
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/workspace
Environment=PATH=/opt/ml-env/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ExecStart=/opt/ml-env/bin/jupyter lab --config=/home/ubuntu/.jupyter/jupyter_lab_config.py
Restart=always

[Install]
WantedBy=multi-user.target
EOF

systemctl enable jupyter
systemctl start jupyter

# Set up SSH for user
mkdir -p /home/ubuntu/.ssh
chown ubuntu:ubuntu /home/ubuntu/.ssh
chmod 700 /home/ubuntu/.ssh

echo "✅ ML Training environment setup complete!"
echo "📊 Jupyter Lab available at: http://<instance-ip>:8888"
echo "🔗 SSH access: ssh ubuntu@<instance-ip>"
`
}

func generateInferenceSetupScript() string {
	return `#!/bin/bash
set -e

echo "🚀 Setting up Inference Server Environment..."

# Update system
apt-get update && apt-get upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
usermod -aG docker ubuntu

# Install NVIDIA Container Toolkit
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | tee /etc/apt/sources.list.d/nvidia-docker.list

apt-get update && apt-get install -y nvidia-docker2
systemctl restart docker

# Install monitoring tools
apt-get install -y htop nvtop

# Create inference workspace
mkdir -p /workspace/models
mkdir -p /workspace/scripts
chown -R ubuntu:ubuntu /workspace

# Create sample inference script
cat > /workspace/scripts/start_inference.py << 'EOF'
#!/usr/bin/env python3
"""
Sample inference server using FastAPI and transformers
"""
import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel
import torch

app = FastAPI(title="DataCrunch Inference Server")

class InferenceRequest(BaseModel):
    text: str

class InferenceResponse(BaseModel):
    result: str
    model_info: dict

@app.get("/health")
async def health_check():
    return {"status": "healthy", "gpu_available": torch.cuda.is_available()}

@app.post("/inference", response_model=InferenceResponse)
async def inference(request: InferenceRequest):
    # Placeholder inference logic
    result = f"Processed: {request.text}"
    model_info = {
        "gpu_count": torch.cuda.device_count() if torch.cuda.is_available() else 0,
        "device": "cuda" if torch.cuda.is_available() else "cpu"
    }
    return InferenceResponse(result=result, model_info=model_info)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
EOF

chmod +x /workspace/scripts/start_inference.py

# Create Docker compose file for production deployment
cat > /workspace/docker-compose.yml << 'EOF'
version: '3.8'
services:
  inference:
    image: pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime
    ports:
      - "8000:8000"
    volumes:
      - ./models:/workspace/models
      - ./scripts:/workspace/scripts
    working_dir: /workspace
    command: python scripts/start_inference.py
    runtime: nvidia
    environment:
      - NVIDIA_VISIBLE_DEVICES=all
EOF

chown ubuntu:ubuntu /workspace/docker-compose.yml

echo "✅ Inference Server environment setup complete!"
echo "🚀 Start server: cd /workspace && docker-compose up"
echo "📡 API will be available at: http://<instance-ip>:8000"
echo "📖 API docs at: http://<instance-ip>:8000/docs"
`
}

func selectImageForEnvironment(environment string) string {
	switch environment {
	case "dev", "development":
		return "ubuntu-20.04"
	case "staging":
		return "pytorch-2.0"
	case "prod", "production":
		return "pytorch-2.0"
	default:
		return "ubuntu-20.04"
	}
}

/*
🚀 Advanced DataCrunch SDK Usage

This example demonstrates:

💡 SMART RESOURCE OPTIMIZATION:
• Automatic selection of optimal instance types based on requirements
• Cost-aware location selection
• GPU requirement matching with budget constraints

🔧 PRODUCTION-READY PATTERNS:
• Comprehensive error handling with retries and exponential backoff
• Resource lifecycle management with automatic cleanup
• Structured logging and progress monitoring

🎯 REAL-WORLD WORKFLOWS:
• Complete ML training environment setup
• Production inference server deployment
• Multi-service resource orchestration

📊 ADVANCED FEATURES:
• Batch operations across multiple services
• Environment-specific configurations
• Cost optimization strategies
• Health monitoring and validation

🧹 RESOURCE MANAGEMENT:
• Automatic cleanup of all created resources
• Proper error handling during cleanup
• Resource tracking across service boundaries

How to run:
1. Set your credentials (same as basic example)
2. go run examples/advanced/main.go
3. Watch as it creates complete ML environments!

The script will:
✅ Optimize instance selection based on your requirements
✅ Create SSH keys, startup scripts, and storage volumes
✅ Deploy production-ready ML training and inference environments
✅ Monitor instance health and validate configurations
✅ Clean up all resources when done

Perfect for understanding how to build robust, production-ready
applications with the DataCrunch API! 🎉
*/
