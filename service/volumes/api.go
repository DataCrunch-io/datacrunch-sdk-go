package volumes

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/restjson"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/request"
)

// Instance represents an instance attached to a volume
type Instance struct {
	ID                  string `json:"id"`
	AutoRentalExtension *bool  `json:"auto_rental_extension"`
	IP                  string `json:"ip"`
	InstanceType        string `json:"instance_type"`
	Status              string `json:"status"`
	OSVolumeID          string `json:"os_volume_id"`
	Hostname            string `json:"hostname"`
}

// LongTerm represents long-term contract details
type LongTerm struct {
	EndDate             string  `json:"end_date"`
	LongTermPeriod      string  `json:"long_term_period"`
	DiscountPercentage  float64 `json:"discount_percentage"`
	AutoRentalExtension bool    `json:"auto_rental_extension"`
	NextPeriodPrice     float64 `json:"next_period_price"`
	CurrentPeriodPrice  float64 `json:"current_period_price"`
}

// VolumeResponse represents a volume
type VolumeResponse struct {
	ID                       string     `json:"id"`
	InstanceID               string     `json:"instance_id"`
	Instances                []Instance `json:"instances"`
	Name                     string     `json:"name"`
	CreatedAt                string     `json:"created_at"`
	Status                   string     `json:"status"`
	Size                     int64      `json:"size"`
	IsOSVolume               bool       `json:"is_os_volume"`
	Target                   string     `json:"target"`
	Type                     string     `json:"type"`
	Location                 string     `json:"location"`
	SSHKeyIDs                []string   `json:"ssh_key_ids"`
	PseudoPath               string     `json:"pseudo_path"`
	CreateDirectoryCommand   string     `json:"create_directory_command"`
	MountCommand             string     `json:"mount_command"`
	FilesystemToFstabCommand string     `json:"filesystem_to_fstab_command"`
	Contract                 string     `json:"contract"`
	BaseHourlyCost           float64    `json:"base_hourly_cost"`
	MonthlyPrice             float64    `json:"monthly_price"`
	Currency                 string     `json:"currency"`
	LongTerm                 *LongTerm  `json:"long_term"`
	DeletedAt                string     `json:"deleted_at,omitempty"`
}

// CreateVolumeInput represents input for creating a volume
type CreateVolumeInput struct {
	Type         string   `json:"type"`
	LocationCode string   `json:"location_code"`
	Size         int64    `json:"size"`
	InstanceID   string   `json:"instance_id,omitempty"`
	InstanceIDs  []string `json:"instance_ids,omitempty"`
	Name         string   `json:"name"`
}

// VolumeActionInput represents input for performing an action on a volume
type VolumeActionInput struct {
	Action       string   `json:"action"`
	ID           string   `json:"id"`
	Size         int64    `json:"size,omitempty"`
	InstanceID   string   `json:"instance_id,omitempty"`
	InstanceIDs  []string `json:"instance_ids,omitempty"`
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	IsPermanent  bool     `json:"is_permanent,omitempty"`
	LocationCode string   `json:"location_code,omitempty"`
}

// VolumeStatus represents the possible status values for a volume.
type VolumeStatus string

const (
	VolumeStatusOrdered   VolumeStatus = "ordered"
	VolumeStatusAttached  VolumeStatus = "attached"
	VolumeStatusAttaching VolumeStatus = "attaching"
	VolumeStatusDetached  VolumeStatus = "detached"
	VolumeStatusDeleted   VolumeStatus = "deleted"
)

type ListVolumesStatus struct {
	Status VolumeStatus `json:"status"`
}

// ListVolumes lists all volumes
func (c *Volumes) ListVolumes(status *ListVolumesStatus) ([]*VolumeResponse, error) {
	op := &request.Operation{
		Name:       "ListVolumes",
		HTTPMethod: "GET",
		HTTPPath:   "/volumes",
	}

	var volumes []*VolumeResponse
	req := c.newRequest(op, status, &volumes)

	return volumes, req.Send()
}

type GetVolumeInput struct {
	ID string `location:"uri" locationName:"id"`
}

// GetVolume gets a volume by ID
func (c *Volumes) GetVolume(id string) (*VolumeResponse, error) {
	op := &request.Operation{
		Name:       "GetVolume",
		HTTPMethod: "GET",
		HTTPPath:   "/volumes/{id}",
	}

	var volume VolumeResponse
	req := c.newRequest(op, &GetVolumeInput{ID: id}, &volume)

	return &volume, req.Send()
}

// CreateVolume creates a new volume
func (c *Volumes) CreateVolume(input *CreateVolumeInput) (string, error) {
	op := &request.Operation{
		Name:       "CreateVolume",
		HTTPMethod: "POST",
		HTTPPath:   "/volumes",
	}

	var volumeID string
	req := c.newRequest(op, input, &volumeID)

	req.Handlers.Unmarshal.RemoveByName("datacrunchsdk.restjson.Unmarshal")
	req.Handlers.Unmarshal.PushBackNamed(restjson.StringUnmarshalHandler)

	return volumeID, req.Send()
}

// PerformVolumeAction performs an action on a volume
func (c *Volumes) PerformVolumeAction(input *VolumeActionInput) error {
	op := &request.Operation{
		Name:       "PerformVolumeAction",
		HTTPMethod: "PUT",
		HTTPPath:   "/volumes",
	}

	req := c.newRequest(op, input, nil)

	return req.Send()
}

// ListTrashVolumes lists all volumes in trash
func (c *Volumes) ListTrashVolumes() ([]*VolumeResponse, error) {
	op := &request.Operation{
		Name:       "ListTrashVolumes",
		HTTPMethod: "GET",
		HTTPPath:   "/volumes/trash",
	}

	var volumes []*VolumeResponse
	req := c.newRequest(op, nil, &volumes)

	return volumes, req.Send()
}

type DeleteVolumeInput struct {
	ID string `location:"uri" locationName:"id"`
}

// DeleteVolume deletes a volume by ID
func (c *Volumes) DeleteVolume(id string, isPermanent bool) error {
	op := &request.Operation{
		Name:       "DeleteVolume",
		HTTPMethod: "DELETE",
		HTTPPath:   "/volumes/{id}",
	}

	req := c.newRequest(op, &DeleteVolumeInput{ID: id}, nil)

	return req.Send()
}
