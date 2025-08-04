package instancetypes

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
)

const (
	// EndpointsID is the service identifier
	EndpointsID = "instancetypes"
)

// InstanceTypes provides the API operation methods for making requests to
// the instance types service.
type InstanceTypes struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// New creates a new instance of the instance types client with a session.
func New(cfg *client.Config) *InstanceTypes {
	svc := &InstanceTypes{
		Client: client.New(cfg, metadata.ClientInfo{
			ServiceName: EndpointsID,
			APIVersion:  "v1",
			Endpoint:    "https://api.datacrunch.io/v1",
		}, request.Handlers{}),
	}

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// NewRequest creates a new request for the instance types service.
func (c *InstanceTypes) NewRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.Client.NewRequest(op, params, data)
	req.HTTPRequest.Header.Set("Accept", "application/json")
	return req
}
