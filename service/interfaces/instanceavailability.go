package interfaces

import (
	"context"

	"github.com/datacrunch-io/datacrunch-sdk-go/service/instanceavailability"
)

// InstanceAvailabilityAPI provides an interface to enable mocking the
// instance-availability service client's API operation.
type InstanceAvailabilityAPI interface {
	ListInstanceAvailability(ctx context.Context) ([]*instanceavailability.InstanceAvailabilityResponse, error)
	CheckInstanceAvailability(ctx context.Context, instanceType string) (bool, error)
}

var _ InstanceAvailabilityAPI = (*instanceavailability.InstanceAvailability)(nil)
