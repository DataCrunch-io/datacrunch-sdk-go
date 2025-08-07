package interfaces

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
)

// InstanceTypesAPI provides the interface for the instance types service
type InstanceTypesAPI interface {
	// ListInstanceTypes lists all available instance types
	ListInstanceTypes() ([]*instancetypes.InstanceTypeResponse, error)
	// GetInstanceTypePriceHistory gets the price history for instance types
	GetInstanceTypePriceHistory() (*instancetypes.PriceHistoryResponse, error)
}

var _ InstanceTypesAPI = (*instancetypes.InstanceTypes)(nil)
