package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// Your main client
type DataCrunchClient struct {
	config *config.Config
	// Don't need to store logger - it's global!
}

// NewClient - setup logger once here
func NewClient(opts ...config.Option) *DataCrunchClient {
	cfg := &config.Config{}
	
	// Apply config options
	for _, opt := range opts {
		opt(cfg)
	}
	
	// 🎯 Setup global logger from config (ONE TIME ONLY)
	logger.SetupFromConfig(cfg.Debug, cfg.Logger)
	
	logger.Info("DataCrunch SDK initialized", 
		"debug_enabled", logger.IsDebugEnabled(),
	)
	
	return &DataCrunchClient{
		config: cfg,
	}
}

// Example API method
func (c *DataCrunchClient) CreateInstance(ctx context.Context, name string, instanceType string) (*Instance, error) {
	// Use logger directly - no setup needed!
	logger.Info("Creating instance", 
		"name", name, 
		"type", instanceType,
	)
	
	// Validate input
	if err := c.validateInstanceRequest(name, instanceType); err != nil {
		logger.Error("Validation failed", "error", err)
		return nil, err
	}
	
	// Make API call
	instance, err := c.makeCreateInstanceCall(ctx, name, instanceType)
	if err != nil {
		logger.Error("Instance creation failed", 
			"name", name,
			"error", err,
		)
		return nil, err
	}
	
	logger.Info("Instance created successfully", 
		"name", name,
		"instance_id", instance.ID,
	)
	
	return instance, nil
}

// Validation function - uses logger directly
func (c *DataCrunchClient) validateInstanceRequest(name, instanceType string) error {
	logger.Debug("Validating instance request", 
		"name", name,
		"type", instanceType,
	)
	
	if name == "" {
		logger.Warn("Empty instance name provided")
		return errors.New("instance name is required")
	}
	
	if len(name) > 50 {
		logger.Warn("Instance name too long", "length", len(name))
		return errors.New("instance name must be ≤ 50 characters")
	}
	
	if instanceType == "" {
		logger.Warn("No instance type specified")
		return errors.New("instance type is required")
	}
	
	logger.Debug("Validation passed")
	return nil
}

// API call function - uses logger directly
func (c *DataCrunchClient) makeCreateInstanceCall(ctx context.Context, name, instanceType string) (*Instance, error) {
	logger.Debug("Making HTTP request", 
		"method", "POST",
		"endpoint", "/instances",
	)
	
	start := time.Now()
	
	// Simulate HTTP request
	time.Sleep(200 * time.Millisecond)
	
	duration := time.Since(start)
	logger.Debug("HTTP request completed", 
		"status", 201,
		"duration", duration,
	)
	
	// Return mock instance
	return &Instance{
		ID:   "inst-12345",
		Name: name,
		Type: instanceType,
	}, nil
}

// Example type
type Instance struct {
	ID   string
	Name string
	Type string
}

func main() {
	fmt.Println("=== DataCrunch SDK Logger Example ===")
	
	// Method 1: Enable debug via config
	fmt.Println("\n1. Debug via config:")
	client1 := NewClient(config.WithDebug(true))
	client1.CreateInstance(context.Background(), "my-server", "small")
	
	// Method 2: Environment variable
	fmt.Println("\n2. Debug via environment:")
	os.Setenv("DATACRUNCH_DEBUG", "true")
	client2 := NewClient() // Will automatically use env var
	client2.CreateInstance(context.Background(), "web-server", "medium")
	
	// Method 3: Custom logger
	fmt.Println("\n3. Custom JSON logger:")
	customLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client3 := NewClient(config.WithLogger(customLogger))
	client3.CreateInstance(context.Background(), "api-server", "large")
	
	// Method 4: Production mode (info level)
	fmt.Println("\n4. Production mode:")
	os.Unsetenv("DATACRUNCH_DEBUG")
	client4 := NewClient(config.WithDebug(false))
	client4.CreateInstance(context.Background(), "prod-server", "xlarge")
}