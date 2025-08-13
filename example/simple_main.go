package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// Simple client example
type SimpleClient struct {
	name string
}

func NewSimpleClient(name string, cfg *config.Config) *SimpleClient {
	// 🎯 Setup logger once - this is the key step!
	logger.SetupFromConfig(cfg.Debug, cfg.Logger)
	
	logger.Info("Client initialized", "client_name", name)
	return &SimpleClient{name: name}
}

func (c *SimpleClient) CreateServer(serverName string) error {
	// ✨ Just use logger directly - no setup needed!
	logger.Info("Creating server", "server_name", serverName)
	
	// Validation
	if serverName == "" {
		logger.Error("Server name is required")
		return errors.New("server name required")
	}
	
	// Auth step
	c.authenticateRequest()
	
	// API call
	c.makeAPICall(serverName)
	
	// Success
	logger.Info("Server created successfully", "server_name", serverName)
	return nil
}

func (c *SimpleClient) authenticateRequest() {
	logger.Debug("Authenticating request")
	
	clientID := "client_abc123def456"
	logger.Debug("Using client credentials", 
		"client_id", logger.SanitizeToken(clientID), // Masks sensitive data
	)
	
	// Simulate auth
	time.Sleep(50 * time.Millisecond)
	logger.Debug("Authentication successful")
}

func (c *SimpleClient) makeAPICall(serverName string) {
	logger.Debug("Making API call", 
		"method", "POST",
		"endpoint", "/servers",
	)
	
	start := time.Now()
	
	// Simulate API call
	time.Sleep(100 * time.Millisecond)
	
	duration := time.Since(start)
	logger.Debug("API call completed", 
		"duration", duration,
		"status", 201,
	)
}

// Example in another "package" (simulate different service)
type DatabaseService struct{}

func (d *DatabaseService) SaveServerConfig(serverName string, config map[string]string) {
	// 🚀 No logger setup needed - just use it!
	logger.Debug("Saving server config to database", 
		"server", serverName,
		"config_keys", strings.Join(d.getMapKeys(config), ","),
	)
	
	// Simulate database save
	time.Sleep(30 * time.Millisecond)
	
	logger.Info("Server config saved", "server", serverName)
}

func (d *DatabaseService) getMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Another service example
func ProcessNotification(serverName string, event string) {
	// 🎯 Works in any function - just import and use!
	logger.Info("Processing notification", 
		"server", serverName,
		"event", event,
	)
	
	if event == "error" {
		logger.Error("Error event received", "server", serverName)
	} else {
		logger.Debug("Event processed successfully", "event", event)
	}
}

func main() {
	fmt.Println("🚀 Simple Global Logger Example")
	fmt.Println(strings.Repeat("=", 40))

	// Example 1: Debug mode
	fmt.Println("\n1. Debug Mode:")
	client1 := NewSimpleClient("debug-client", &config.Config{
		Debug: &[]bool{true}[0], // Enable debug
	})
	client1.CreateServer("web-server-1")
	
	// Use other services - they all use the same logger!
	db := &DatabaseService{}
	db.SaveServerConfig("web-server-1", map[string]string{
		"cpu": "2",
		"ram": "4GB",
	})
	
	ProcessNotification("web-server-1", "created")

	fmt.Println("\n" + strings.Repeat("-", 40))

	// Example 2: Production mode (less verbose)
	fmt.Println("\n2. Production Mode:")
	client2 := NewSimpleClient("prod-client", &config.Config{
		Debug: &[]bool{false}[0], // Disable debug
	})
	client2.CreateServer("prod-server-1")

	fmt.Println("\n" + strings.Repeat("-", 40))

	// Example 3: Environment variable
	fmt.Println("\n3. Environment Variable:")
	os.Setenv("DATACRUNCH_DEBUG", "true")
	client3 := NewSimpleClient("env-client", &config.Config{})
	client3.CreateServer("env-server-1")

	fmt.Println("\n" + strings.Repeat("-", 40))

	// Example 4: Custom JSON logger
	fmt.Println("\n4. Custom JSON Logger:")
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	
	client4 := NewSimpleClient("json-client", &config.Config{
		Logger: jsonLogger,
	})
	client4.CreateServer("json-server-1")

	fmt.Println("\n✅ All examples completed!")
}