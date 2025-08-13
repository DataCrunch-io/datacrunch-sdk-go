package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// Client setup
func NewClient(cfg *config.Config) *Client {
	// Setup global logger from config
	logger.SetupFromConfig(cfg.Debug, cfg.Logger)
	
	return &Client{config: cfg}
}

// In your service methods - context-based
func (c *Client) CreateInstance(ctx context.Context, name string) error {
	// Add operation context
	ctx = logger.SetLogger(ctx, logger.GetGlobalLogger().With(
		"operation", "CreateInstance",
		"instance_name", name,
	))

	// Use native slog methods
	logger.Debug(ctx, "Starting instance creation")
	
	// Call other packages
	err := c.validateRequest(ctx, name)
	if err != nil {
		logger.Error(ctx, "Validation failed", "error", err)
		return err
	}

	err = c.makeAPICall(ctx, name)
	if err != nil {
		logger.Error(ctx, "API call failed", "error", err)
		return err
	}

	logger.Info(ctx, "Instance created successfully")
	return nil
}

func (c *Client) validateRequest(ctx context.Context, name string) error {
	logger.Debug(ctx, "Validating request", "name", name)
	
	if name == "" {
		logger.Warn(ctx, "Empty instance name provided")
		return errors.New("name required")
	}
	
	return nil
}

func (c *Client) makeAPICall(ctx context.Context, name string) error {
	logger.Debug(ctx, "Making API request",
		"method", "POST",
		"endpoint", "/instances",
	)
	
	// Your HTTP request code here...
	
	logger.Debug(ctx, "API request completed", 
		"status_code", 201,
		"response_time", "150ms",
	)
	
	return nil
}

// Usage
func main() {
	// Method 1: Enable debug via config
	client := NewClient(&config.Config{
		Debug: &[]bool{true}[0], // pointer to true
	})
	
	// Method 2: Environment variable
	os.Setenv("DATACRUNCH_DEBUG", "true")
	
	// Method 3: Custom logger
	customLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client2 := NewClient(&config.Config{
		Logger: customLogger,
	})
	
	ctx := context.Background()
	client.CreateInstance(ctx, "my-instance")
}