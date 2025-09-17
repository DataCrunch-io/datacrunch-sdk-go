package instanceavailability

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/request"
)

const (
	// ServiceName is the name of the service
	ServiceName = "instance-availability"
)

// InstanceAvailabilityResponse represents the availability of instance types in a location
type InstanceAvailabilityResponse struct {
	LocationCode   string   `json:"location_code" locationName:"location_code"`
	Availabilities []string `json:"availabilities" locationName:"availabilities"`
}

// ListInstanceAvailability lists all available instance types by location
func (c *InstanceAvailability) ListInstanceAvailability() ([]*InstanceAvailabilityResponse, error) {
	op := &request.Operation{
		Name:       "ListInstanceAvailability",
		HTTPMethod: "GET",
		HTTPPath:   "/instance-availability",
	}

	var availabilities []*InstanceAvailabilityResponse
	req := c.newRequest(op, nil, &availabilities)

	return availabilities, req.Send()
}

// CheckInstanceAvailability checks if a specific instance type is available
// DISABLED: This function is disabled due to inconsistent API behavior.
// The API returns different response formats based on query parameters:
// - With insufficient params: returns array
// - With instance type: returns boolean as string ("true"/"false")
// This inconsistency makes it difficult to implement reliably.
/*
func (c *InstanceAvailability) CheckInstanceAvailability(instanceType string, locationCode *string, isSpot *bool) (bool, error) {
	op := &request.Operation{
		Name:       "CheckInstanceAvailability",
		HTTPMethod: "GET",
		HTTPPath:   fmt.Sprintf("/instance-availability/%s", instanceType),
	}

	params := map[string]interface{}{}
	if locationCode != nil {
		params["location_code"] = *locationCode
	}
	if isSpot != nil {
		params["is_spot"] = *isSpot
	}

	var available bool
	req := c.newRequest(op, params, &available)

	return available, req.Send()
}
*/
