package locations

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/restjson"
)

const (
	EndpointsID = "locations"
	APIVersion  = "v1"
)

// Locations provides the API operation methods for making requests to
type Locations struct {
	*client.Client
}

// Client is an alias for Locations to match the expected interface
type Client = *Locations

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

type ConfigProvider interface {
	ClientConfig(serviceName string, cfgs ...*interface{}) client.Config
}

// New creates a new instance of the Locations client with a config provider.
// If additional configuration is needed for the client instance use the optional
// client.Config parameter to add your extra config.
//
// Example:
//
//	mySession := session.Must(session.New())
//
//	// Create a Locations client from just a session.
//	svc := locations.New(mySession)
//
//	// Create a Locations client with additional configuration
//	svc := locations.New(mySession, &client.Config{Timeout: 60 * time.Second})
func New(p ConfigProvider, cfgs ...*interface{}) *Locations {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(c)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg client.Config) *Locations {
	handlers := request.Handlers{}

	// Add protocol handlers for REST JSON
	handlers.Build.PushBackNamed(restjson.BuildHandler)
	handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	handlers.Complete.PushBackNamed(restjson.UnmarshalMetaHandler)

	svc := &Locations{
		Client: client.New(&cfg, metadata.ClientInfo{
			ServiceName: EndpointsID,
			APIVersion:  APIVersion,
			Endpoint:    cfg.BaseURL,
		}, handlers),
	}

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

func (c *Locations) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
