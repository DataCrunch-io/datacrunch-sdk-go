package sshkeys

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/restjson"
)

const (
	EndpointsID = "sshkeys"
	APIVersion  = "v1"
)

// SSHKey provides the API operation methods for making requests to
type SSHKey struct {
	*client.Client
}

// Client is an alias for SSHKey to match the expected interface
type Client = *SSHKey

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// New creates a new instance of the SSHKey client with a config provider.
// If additional configuration is needed for the client instance use the optional
// client.Config parameter to add your extra config.
//
// Example:
//
//	mySession := session.Must(session.New())
//
//	// Create a SSHKey client from just a session.
//	svc := sshkeys.New(mySession)
//
//	// Create a SSHKey client with additional configuration
//	svc := sshkeys.New(mySession, &client.Config{Timeout: 60 * time.Second})
func New(p client.ConfigProvider, cfgs ...*config.Config) *SSHKey {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(c.Config, c.Handlers)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg config.Config, handlers request.Handlers) *SSHKey {

	svc := &SSHKey{
		Client: client.New(cfg, metadata.ClientInfo{
			ServiceName: EndpointsID,
			APIVersion:  APIVersion,
			Endpoint:    *cfg.BaseURL,
		}, handlers),
	}

	// Add protocol handlers for REST JSON
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.StandardUnmarshalHandler)
	svc.Handlers.Complete.PushBackNamed(restjson.UnmarshalMetaHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

func (c *SSHKey) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
