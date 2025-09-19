package instance

import (
	"fmt"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/api"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/request"
)

const (
	EndpointsID = "instance"
	APIVersion  = "v1"
)

// Instance provides the API operation methods for making requests to
// DataCrunch Instance API
type Instance struct {
	*client.Client

	// Temporary API client for internal use
	tempAPIClient api.Client
}

// Client is an alias for Instance to match the expected interface
type Client = *Instance

// Used for custom client initialization logic
var initClient func(*client.Client)

// New creates a new instance of the Instance client with a config provider.
func New(p client.ConfigProvider, cfgs ...*config.Config) *Instance {
	c := p.ClientConfig(EndpointsID, cfgs...)

	return newClient(c.Config, c.Handlers)
}

func (c *Instance) tempGetNewAPIClient() (api.Client, error) {
	clientID, clientSecret, err := c.Config.Credentials.GetClientCredentials()
	if err != nil {
		return api.Client{}, fmt.Errorf("error getting client credentials: %w", err)
	}

	newClient := api.NewWithCredentials(clientID, clientSecret)
	newClient.SetBaseURL(*c.Config.BaseURL)

	return newClient, nil
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg config.Config, handlers request.Handlers) *Instance {

	svc := &Instance{
		Client: client.New(cfg, metadata.ClientInfo{
			ServiceName: EndpointsID,
			APIVersion:  APIVersion,
			Endpoint:    *cfg.BaseURL,
		}, handlers),
	}

	otherAPIClient, err := svc.tempGetNewAPIClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create temporary API client: %v", err))
	}
	svc.tempAPIClient = otherAPIClient

	// Add protocol handlers for REST JSON
	//svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	//svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	//svc.Handlers.Complete.PushBackNamed(restjson.UnmarshalMetaHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}
