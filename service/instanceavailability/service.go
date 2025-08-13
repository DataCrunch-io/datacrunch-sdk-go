package instanceavailability

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/restjson"
)

const (
	EndpointsID = "instanceavailability"
	APIVersion  = "v1"
)

// InstanceAvailability provides the API operation methods for making requests to
type InstanceAvailability struct {
	*client.Client
}

// Client is an alias for InstanceAvailability to match the expected interface
type Client = *InstanceAvailability

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// New creates a new instance of the InstanceAvailability client with a config provider.
// If additional configuration is needed for the client instance use the optional
// client.Config parameter to add your extra config.
//
// Example:
//
//	mySession := session.Must(session.New())
//
//	// Create a InstanceAvailability client from just a session.
//	svc := instanceavailability.New(mySession)
//
//	// Create a InstanceAvailability client with additional configuration
//	svc := instanceavailability.New(mySession, &client.Config{Timeout: 60 * time.Second})
func New(p client.ConfigProvider, cfgs ...*datacrunch.Config) *InstanceAvailability {
	c := p.ClientConfig(EndpointsID, cfgs...)

	return newClient(c.Config, c.Handlers)
}

// newClientWithHandlers creates, initializes and returns a new service client instance with session handlers.
func newClient(cfg datacrunch.Config, handlers request.Handlers) *InstanceAvailability {
	svc := &InstanceAvailability{
		Client: client.New(cfg, metadata.ClientInfo{
			ServiceName: EndpointsID,
			APIVersion:  APIVersion,
			Endpoint:    *cfg.BaseURL,
		}, handlers),
	}

	// Add protocol handlers for REST JSON
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.Complete.PushBackNamed(restjson.UnmarshalMetaHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

func (c *InstanceAvailability) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
