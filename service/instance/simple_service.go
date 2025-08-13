package instance

import (
	"context"
)

// SimpleService is the refactored instance service using standard JSON marshaling
type SimpleService struct {
	client HTTPClient
}

// NewSimpleService creates a new instance service with any client implementing HTTPClient
func NewSimpleService(client HTTPClient) *SimpleService {
	return &SimpleService{
		client: client,
	}
}

// Fixed input structs using STANDARD JSON tags (no AWS-style locationName)

// SimpleCreateInstanceInput uses standard JSON marshaling - fixes omitempty bug!
type SimpleCreateInstanceInput struct {
	InstanceType    string    `json:"instance_type"`
	Image           string    `json:"image"`
	SSHKeyIDs       []string  `json:"ssh_key_ids"`
	StartupScriptID string    `json:"startup_script_id,omitempty"` // ✅ Works correctly now!
	Hostname        string    `json:"hostname,omitempty"`
	Description     string    `json:"description,omitempty"`
	LocationCode    string    `json:"location_code"`
	OSVolume        *OSVolume `json:"os_volume,omitempty"`
	IsSpot          bool      `json:"is_spot"`
	Volumes         []Volume  `json:"volumes,omitempty"`
	ExistingVolumes []string  `json:"existing_volumes,omitempty"`
	Contract        string    `json:"contract"`
	Pricing         string    `json:"pricing"`
}

// SimpleListInstancesInput for query parameters
type SimpleListInstancesInput struct {
	Status string // Used for query params, not JSON body
}

// SimpleListInstancesResponse using standard JSON unmarshaling
type SimpleListInstancesResponse struct {
	ID              string   `json:"id"`
	IP              string   `json:"ip"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"created_at"`
	CPU             CPU      `json:"cpu"`
	GPU             GPU      `json:"gpu"`
	GPUMemory       Memory   `json:"gpu_memory"`
	Memory          Memory   `json:"memory"`
	Storage         Storage  `json:"storage"`
	Hostname        string   `json:"hostname"`
	Description     string   `json:"description"`
	Location        string   `json:"location"` // Standard string handling
	PricePerHour    float64  `json:"price_per_hour"`
	IsSpot          bool     `json:"is_spot"`
	InstanceType    string   `json:"instance_type"`
	Image           string   `json:"image"`
	OSName          string   `json:"os_name"`
	StartupScriptID *string  `json:"startup_script_id"` // Pointer for null values
	SSHKeyIDs       []string `json:"ssh_key_ids"`
	OSVolumeID      string   `json:"os_volume_id"`
	JupyterToken    *string  `json:"jupyter_token"` // Pointer for null values
	Contract        string   `json:"contract"`
	Pricing         string   `json:"pricing"`
}

// SimpleInstanceActionInput for instance actions
type SimpleInstanceActionInput struct {
	Action    InstanceActionType `json:"action"`
	ID        string             `json:"id"`
	VolumeIDs []string           `json:"volume_ids,omitempty"`
}

// API Methods using simplified client with standard JSON

// CreateInstance creates a new instance - now with working omitempty!
func (s *SimpleService) CreateInstance(ctx context.Context, input *SimpleCreateInstanceInput) (string, error) {
	var instanceID string
	err := s.client.POST(ctx, "/instances", input, &instanceID)
	return instanceID, err
}

// ListInstances lists instances with optional status filter
func (s *SimpleService) ListInstances(ctx context.Context, input *SimpleListInstancesInput) ([]*SimpleListInstancesResponse, error) {
	// Handle query parameters manually (clean approach)
	queryParams := make(map[string]string)
	if input != nil && input.Status != "" {
		queryParams["status"] = input.Status
	}

	var instances []*SimpleListInstancesResponse
	err := s.client.GET(ctx, "/instances", queryParams, &instances)
	return instances, err
}

// PerformInstanceAction performs an action on an instance
func (s *SimpleService) PerformInstanceAction(ctx context.Context, input *SimpleInstanceActionInput) error {
	return s.client.PUT(ctx, "/instances", input, nil)
}

// GetInstance gets a specific instance by ID
func (s *SimpleService) GetInstance(ctx context.Context, instanceID string) (*SimpleListInstancesResponse, error) {
	var instance *SimpleListInstancesResponse
	err := s.client.GET(ctx, "/instances/"+instanceID, nil, &instance)
	return instance, err
}

// DeleteInstance deletes an instance by ID
func (s *SimpleService) DeleteInstance(ctx context.Context, instanceID string) error {
	return s.client.DELETE(ctx, "/instances/"+instanceID, nil)
}
