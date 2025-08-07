package volumetypes

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
)

// Price represents the pricing details for a volume type
type Price struct {
	PricePerMonthPerGB float64 `json:"price_per_month_per_gb"`
	CPSPerGB           float64 `json:"cps_per_gb"`
	Currency           string  `json:"currency"`
}

// VolumeTypeResponse represents a volume type
type VolumeTypeResponse struct {
	Type                 string `json:"type"`
	Price                Price  `json:"price"`
	IsSharedFS           bool   `json:"is_shared_fs"`
	BurstBandwidth       int    `json:"burst_bandwidth"`
	ContinuousBandwidth  int    `json:"continuous_bandwidth"`
	InternalNetworkSpeed int    `json:"internal_network_speed"`
	IOPS                 string `json:"iops"`
}

// ListVolumeTypes lists all available volume types
func (c *VolumeTypes) ListVolumeTypes() ([]*VolumeTypeResponse, error) {
	op := &request.Operation{
		Name:       "ListVolumeTypes",
		HTTPMethod: "GET",
		HTTPPath:   "/volume-types",
	}

	var volumeTypes []*VolumeTypeResponse
	req := c.newRequest(op, nil, &volumeTypes)

	return volumeTypes, req.Send()
}
