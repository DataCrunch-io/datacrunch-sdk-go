package interfaces

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instanceavailability"
)

// InstanceAvailabilityAPI provides an interface to enable mocking the
// instance-availability service client's API operation.
type InstanceAvailabilityAPI interface {
	ListInstanceAvailability() ([]*instanceavailability.InstanceAvailabilityResponse, error)
	// CheckInstanceAvailability(instanceType string, locationCode *string, isSpot *bool) (bool, error)
}

var _ InstanceAvailabilityAPI = (*instanceavailability.InstanceAvailability)(nil)
