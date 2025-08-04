package datacrunch

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client/metadata"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/sshkeys"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/startscripts"
)

type RequestRetryer interface{}

// Config holds configuration for the DataCrunch SDK
type Config struct {
	// API configuration
	BaseURL      string
	ClientID     string
	ClientSecret string
	Timeout      time.Duration

	// HTTP client configuration
	Retryer       RequestRetryer
	MaxRetries    int
	RetryDelay    time.Duration
	MaxRetryDelay time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL:       "https://api.datacrunch.io/v1",
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryDelay:    1 * time.Second,
		MaxRetryDelay: 30 * time.Second,
	}
}

// Client represents the main DataCrunch SDK client
type Client struct {
	config *Config

	// HTTP client
	httpClient *client.Client

	// Service clients
	Instance     instance.Client
	SSHKeys      sshkeys.Client
	StartScripts startscripts.Client
}

// Option is a functional option for configuring the DataCrunch client
type Option func(*Config)

// WithBaseURL sets the base URL for the API
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithCredentials sets the OAuth2 client credentials
func WithCredentials(clientID, clientSecret string) Option {
	return func(c *Config) {
		c.ClientID = clientID
		c.ClientSecret = clientSecret
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithRetryConfig sets the retry configuration
func WithRetryConfig(maxRetries int, retryDelay, maxRetryDelay time.Duration) Option {
	return func(c *Config) {
		c.MaxRetries = maxRetries
		c.RetryDelay = retryDelay
		c.MaxRetryDelay = maxRetryDelay
	}
}

// WithRetryer sets a custom retryer
func WithRetryer(retryer RequestRetryer) Option {
	return func(c *Config) {
		c.Retryer = retryer
	}
}

// New creates a new DataCrunch SDK client with functional options
func New(options ...Option) *Client {
	config := DefaultConfig()

	// Apply all options
	for _, option := range options {
		option(config)
	}

	// Create HTTP client
	httpClient := client.New(config, metadata.ClientInfo{
		ServiceName: "datacrunch",
		APIVersion:  "v1",
		Endpoint:    config.BaseURL,
	}, request.Handlers{})

	// Create wrapper for service clients
	wrapper := &httpClientWrapper{client: httpClient}

	return &Client{
		config:       config,
		httpClient:   httpClient,
		Instance:     instance.NewClient(wrapper),
		SSHKeys:      sshkeys.NewClient(wrapper),
		StartScripts: startscripts.NewClient(wrapper),
	}
}

// NewFromEnv creates a new DataCrunch SDK client using environment variables
// Supported environment variables:
// - DATACRUNCH_BASE_URL (default: https://api.datacrunch.io/v1)
// - DATACRUNCH_CLIENT_ID (required)
// - DATACRUNCH_CLIENT_SECRET (required)
// - DATACRUNCH_TIMEOUT (default: 30s, format: "30s", "1m", etc.)
// - DATACRUNCH_MAX_RETRIES (default: 3)
func NewFromEnv(options ...Option) *Client {
	config := DefaultConfig()

	// Load from environment variables
	if baseURL := os.Getenv("DATACRUNCH_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if clientID := os.Getenv("DATACRUNCH_CLIENT_ID"); clientID != "" {
		config.ClientID = clientID
	}

	if clientSecret := os.Getenv("DATACRUNCH_CLIENT_SECRET"); clientSecret != "" {
		config.ClientSecret = clientSecret
	}

	if timeoutStr := os.Getenv("DATACRUNCH_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			config.Timeout = timeout
		}
	}

	if maxRetriesStr := os.Getenv("DATACRUNCH_MAX_RETRIES"); maxRetriesStr != "" {
		if maxRetries, err := strconv.Atoi(maxRetriesStr); err == nil {
			config.MaxRetries = maxRetries
		}
	}

	// Apply additional options (these can override env vars)
	for _, option := range options {
		option(config)
	}

	// Create HTTP client
	httpClient := client.New(config, metadata.ClientInfo{
		ServiceName: "datacrunch",
		APIVersion:  "v1",
		Endpoint:    config.BaseURL,
	}, request.Handlers{})

	// Create wrapper for service clients
	wrapper := &httpClientWrapper{client: httpClient}

	return &Client{
		config:       config,
		httpClient:   httpClient,
		Instance:     instance.NewClient(wrapper),
		SSHKeys:      sshkeys.NewClient(wrapper),
		StartScripts: startscripts.NewClient(wrapper),
	}
}

// NewWithConfig creates a new DataCrunch SDK client with a config struct (legacy support)
func NewWithConfig(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	// Create HTTP client
	httpClient := client.New(config, metadata.ClientInfo{
		ServiceName: "datacrunch",
		APIVersion:  "v1",
		Endpoint:    config.BaseURL,
	}, request.Handlers{})

	// Create wrapper for service clients
	wrapper := &httpClientWrapper{client: httpClient}

	return &Client{
		config:       config,
		httpClient:   httpClient,
		Instance:     instance.NewClient(wrapper),
		SSHKeys:      sshkeys.NewClient(wrapper),
		StartScripts: startscripts.NewClient(wrapper),
	}
}

// httpClientWrapper adapts the HTTP client for service clients
type httpClientWrapper struct {
	client *client.Client
}

// httpResponse adapts the HTTP response for service clients
type httpResponse struct {
	statusCode int
	body       []byte
	response   interface{}
}

// Post implements the APIClientInterface for service clients
func (w *httpClientWrapper) Post(ctx context.Context, path string, body interface{}) (client.ResponseInterface, error) {
	resp, err := w.client.Post(ctx, path, body)
	if err != nil {
		return nil, err
	}

	return &httpResponse{
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
}

// Get implements the APIClientInterface for service clients
func (w *httpClientWrapper) Get(ctx context.Context, path string) (client.ResponseInterface, error) {
	resp, err := w.client.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	return &httpResponse{
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
}

// Delete implements the APIClientInterface for service clients
func (w *httpClientWrapper) Delete(ctx context.Context, path string) (client.ResponseInterface, error) {
	resp, err := w.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}

	return &httpResponse{
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
}

// Put implements the APIClientInterface for service clients
func (w *httpClientWrapper) Put(ctx context.Context, path string, body interface{}) (client.ResponseInterface, error) {
	resp, err := w.client.Put(ctx, path, body)
	if err != nil {
		return nil, err
	}

	return &httpResponse{
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
}

// DecodeJSON implements ResponseInterface
func (r *httpResponse) DecodeJSON(target interface{}) error {
	if httpResp, ok := r.response.(*http.Response); ok {
		// Use the client's DecodeResponse method
		client := &client.Client{} // This is a placeholder - in real implementation would use the actual client
		return client.DecodeResponse(httpResp, target)
	}
	return nil
}

// GetStatusCode implements ResponseInterface
func (r *httpResponse) GetStatusCode() int {
	return r.statusCode
}

// GetBody implements ResponseInterface
func (r *httpResponse) GetBody() []byte {
	return r.body
}
