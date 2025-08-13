package main

import (
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// Client setup - configure once globally
func NewClient(cfg *config.Config) *Client {
	// Setup global logger from config once
	logger.SetupFromConfig(cfg.Debug, cfg.Logger)
	
	return &Client{config: cfg}
}

// In your service methods - use directly
func (c *Client) CreateInstance(name string) error {
	// Use native slog methods directly
	logger.Debug("Starting instance creation", 
		"operation", "CreateInstance",
		"instance_name", name,
	)
	
	// Call other packages - they can use logger directly too
	err := c.validateRequest(name)
	if err != nil {
		logger.Error("Validation failed", "error", err)
		return err
	}

	err = c.makeAPICall(name)
	if err != nil {
		logger.Error("API call failed", "error", err)
		return err
	}

	logger.Info("Instance created successfully", "name", name)
	return nil
}

// In any package - just import and use
func (c *Client) validateRequest(name string) error {
	logger.Debug("Validating request", "name", name)
	
	if name == "" {
		logger.Warn("Empty instance name provided")
		return errors.New("name required")
	}
	
	return nil
}

// In HTTP request package
func (c *Client) makeAPICall(name string) error {
	logger.Debug("Making API request",
		"method", "POST",
		"endpoint", "/instances",
	)
	
	// Your HTTP request code here...
	
	logger.Debug("API request completed", 
		"status_code", 201,
		"response_time", "150ms",
	)
	
	return nil
}

// In retry package
func retryRequest(attempt int, delay time.Duration, err error) {
	logger.Info("Retrying request",
		"attempt", attempt,
		"delay", delay,
		"error", err.Error(),
	)
}

// In auth package
func refreshToken(clientID string) {
	logger.Debug("Refreshing OAuth token",
		"client_id", logger.SanitizeToken(clientID),
		"grant_type", "refresh_token",
	)
}

// Usage
func main() {
	// Method 1: Enable debug via config
	client := NewClient(&config.Config{
		Debug: &[]bool{true}[0],
	})
	
	// Method 2: Environment variable (DATACRUNCH_DEBUG=true)
	// The logger automatically checks environment
	
	// Now use anywhere in your codebase
	client.CreateInstance("my-instance")
}