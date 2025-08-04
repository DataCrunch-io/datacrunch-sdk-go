package interfaces

import (
	"context"

	"github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
)

// InstanceTypesAPI provides the interface for the instance types service
type InstanceTypesAPI interface {
	// ListInstanceTypes lists all available instance types
	ListInstanceTypes(ctx context.Context) ([]*instancetypes.InstanceTypeResponse, error)
	// GetInstanceTypePriceHistory gets the price history for instance types
	GetInstanceTypePriceHistory(ctx context.Context) (*instancetypes.PriceHistoryResponse, error)
}

var _ InstanceTypesAPI = (*instancetypes.InstanceTypes)(nil)
