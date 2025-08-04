package instance

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/restjson"
)

const (
	// EndpointsID is the service identifier
	EndpointsID = "instance"
)

// Instance provides the API operation methods for making requests to
// the instance service.
type Instance struct {
	*client.Client
}

// Client is an alias for Instance to match the expected interface
type Client = *Instance

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// New creates a new instance of the instance client with a session.
func New(cfg *client.Config) *Instance {
	handlers := request.Handlers{}

	// Add protocol handlers for REST JSON
	handlers.Build.PushBackNamed(restjson.BuildHandler)
	handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	handlers.Complete.PushBackNamed(restjson.UnmarshalMetaHandler)

	svc := &Instance{
		Client: client.New(cfg, metadata.ClientInfo{
			ServiceName: EndpointsID,
			APIVersion:  "v1",
			Endpoint:    "https://api.datacrunch.io/v1",
		}, handlers),
	}

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// NewRequest creates a new request for the instance service.
func (c *Instance) NewRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.Client.NewRequest(op, params, data)
	req.HTTPRequest.Header.Set("Accept", "application/json")
	return req
}

// NewClient creates a new instance client with the provided HTTP client wrapper
func NewClient(httpClient interface{}) Client {
	// For now, return a simple instance - this would be enhanced with proper client wrapping
	return &Instance{}
}
