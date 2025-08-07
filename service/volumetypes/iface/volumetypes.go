package interfaces

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/service/volumetypes"
)

// VolumeTypesAPI provides the interface for the volume types service
type VolumeTypesAPI interface {
	// ListVolumeTypes lists all available volume types
	ListVolumeTypes() ([]*volumetypes.VolumeTypeResponse, error)
}

var _ VolumeTypesAPI = (*volumetypes.VolumeTypes)(nil)
