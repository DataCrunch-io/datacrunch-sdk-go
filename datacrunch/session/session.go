package session

import (
	"fmt"
	"os"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/client"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/defaults"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// Session provides a shared configuration and state for service clients
type Session struct {
	Config      *datacrunch.Config
	Handlers    request.Handlers
	Credentials *credentials.Credentials
}

// Options for configuring a session
type Options struct {
	// API configuration
	BaseURL      string
	ClientID     string
	ClientSecret string
	Timeout      time.Duration

	// Credential configuration
	Credentials                   *credentials.Credentials
	CredentialsChainVerboseErrors *bool

	// Retry configuration
	MaxRetries *int
	Retryer    interface{}

	// Logging configuration
	Debug bool
}

// DefaultOptions returns default session options with sensible retry defaults
func DefaultOptions() *Options {
	defaultMaxRetries := 3 // Provide good defaults for all users
	return &Options{
		BaseURL:    "https://api.datacrunch.io/v1",
		Timeout:    30 * time.Second,
		MaxRetries: &defaultMaxRetries, // Default to 3 retries for resilience
		Debug:      false,
	}
}

// New creates a new session with the provided options
func New(options ...func(*Options)) *Session {
	opts := DefaultOptions()

	for _, option := range options {
		option(opts)
	}

	// Setup credentials
	var creds *credentials.Credentials
	if opts.Credentials != nil {
		creds = opts.Credentials
	} else if opts.ClientID != "" && opts.ClientSecret != "" {
		// Use static credentials if provided directly
		creds = credentials.NewStaticCredentials(opts.ClientID, opts.ClientSecret, opts.BaseURL)
	} else {
		// Use default credential chain
		creds = defaults.CredChain()
	}

	var finalBaseURL string

	// Check if BaseURL was explicitly set (different from default)
	if opts.BaseURL != "" && opts.BaseURL != defaults.DefaultBaseURL {
		finalBaseURL = opts.BaseURL
	} else {
		// Determine BaseURL with correct priority: env > credential file > default
		if credValue, err := creds.Get(); err == nil && credValue.BaseURL != "" {
			finalBaseURL = credValue.BaseURL
		}
	}

	cfg := &datacrunch.Config{
		BaseURL:     &finalBaseURL,
		Timeout:     &opts.Timeout,
		MaxRetries:  opts.MaxRetries,
		Retryer:     opts.Retryer,
		Credentials: creds,
		Debug:       opts.Debug,
	}

	// setup logger
	logger.SetupFromConfig(cfg.Debug, nil)

	return &Session{
		Config:      cfg,
		Handlers:    defaults.Handlers(),
		Credentials: creds,
	}
}

// NewFromEnv creates a new session using only environment variables
// This function will panic if required environment variables are missing.
// Use New() instead if you want fallback behavior to credential files.
// Supported environment variables:
// - DATACRUNCH_BASE_URL (optional, defaults to https://api.datacrunch.io/v1)
// - DATACRUNCH_CLIENT_ID (required)
// - DATACRUNCH_CLIENT_SECRET (required)
// - DATACRUNCH_TIMEOUT (optional, default: 30s, format: "30s", "1m", etc.)
func NewFromEnv(options ...func(*Options)) *Session {
	// Only use EnvProvider - no fallback to other credential sources
	envProvider := &credentials.EnvProvider{}
	envCreds := credentials.NewCredentials(envProvider)

	// Try to get credentials from environment - fail fast if missing required vars
	credValue, err := envCreds.Get()
	if err != nil {
		panic(fmt.Sprintf("NewFromEnv requires environment variables but they are missing or invalid: %v. Required: DATACRUNCH_CLIENT_ID, DATACRUNCH_CLIENT_SECRET. Optional: DATACRUNCH_BASE_URL", err))
	}

	opts := DefaultOptions()

	// Use values from environment
	opts.BaseURL = credValue.BaseURL
	opts.ClientID = credValue.ClientID
	opts.ClientSecret = credValue.ClientSecret

	// Handle optional DATACRUNCH_TIMEOUT env var
	if timeoutStr := os.Getenv("DATACRUNCH_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			opts.Timeout = timeout
		}
	}

	// Apply additional options (these can override env vars)
	for _, option := range options {
		option(opts)
	}

	// Create session using the env-only credentials
	cfg := &datacrunch.Config{
		BaseURL:     &opts.BaseURL,
		Timeout:     &opts.Timeout,
		MaxRetries:  opts.MaxRetries,
		Retryer:     opts.Retryer,
		Credentials: envCreds,
		Debug:       opts.Debug,
	}

	return &Session{
		Config:      cfg,
		Handlers:    defaults.Handlers(),
		Credentials: envCreds,
	}
}

// WithBaseURL sets the base URL for the API
func WithBaseURL(baseURL string) func(*Options) {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

// WithCredentials sets the OAuth2 client credentials
func WithCredentials(clientID, clientSecret string) func(*Options) {
	return func(o *Options) {
		o.ClientID = clientID
		o.ClientSecret = clientSecret
	}
}

// WithMaxRetries sets the maximum number of retries for requests
func WithMaxRetries(maxRetries int) func(*Options) {
	return func(o *Options) {
		o.MaxRetries = &maxRetries
	}
}

// WithRetryer sets a custom retryer implementation
func WithRetryer(retryer interface{}) func(*Options) {
	return func(o *Options) {
		o.Retryer = retryer
	}
}

// WithNoRetries disables retry functionality entirely
func WithNoRetries() func(*Options) {
	return WithMaxRetries(0)
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) func(*Options) {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithCredentials sets custom credentials
func WithCredentialsProvider(creds *credentials.Credentials) func(*Options) {
	return func(o *Options) {
		o.Credentials = creds
	}
}

// WithCredentialChainVerboseErrors sets whether to use verbose errors in credential chain
func WithCredentialChainVerboseErrors(verbose bool) func(*Options) {
	return func(o *Options) {
		o.CredentialsChainVerboseErrors = &verbose
	}
}

// WithDebug sets the debug mode
func WithDebug(debug bool) func(*Options) {
	return func(o *Options) {
		o.Debug = debug
	}
}

// ClientConfig implements the client.ConfigProvider interface
func (s *Session) ClientConfig(serviceName string, cfgs ...*datacrunch.Config) client.Config {
	s = s.Copy(cfgs...)
	return client.Config{
		Config:   *s.Config,
		BaseURL:  *s.Config.BaseURL,
		Handlers: s.Handlers,
	}
}

func (s *Session) Copy(cfgs ...*datacrunch.Config) *Session {
	newSession := &Session{
		Config:   s.Config.Copy(cfgs...),
		Handlers: s.Handlers.Copy(),
	}

	initHandlers(newSession)

	return newSession
}

func initHandlers(s *Session) {
}

// ClientConfigNoResolveEndpoint implements the client.ConfigNoResolveEndpointProvider interface
func (s *Session) ClientConfigNoResolveEndpoint(cfgs ...*interface{}) datacrunch.Config {
	return *s.Config
}

// GetCredentials returns the session's credentials (implements SessionWithCredentials interface)
func (s *Session) GetCredentials() *credentials.Credentials {
	return s.Credentials
}
