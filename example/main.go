package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

func main() {
	fmt.Println("🚀 DataCrunch SDK Global Logger Example")
	fmt.Println("=" * 50)

	// Example 1: Basic debug logging
	runBasicExample()

	// Example 2: Multiple services using same logger
	runMultiServiceExample()

	// Example 3: Production vs Development logging
	runEnvironmentExample()

	// Example 4: Custom structured logging
	runCustomLoggerExample()
}

func runBasicExample() {
	fmt.Println("\n📝 1. Basic Debug Logging Example")
	fmt.Println("-" * 30)

	// Initialize client with debug enabled
	client := NewClient(config.WithDebug(true))

	// Use the client - logging happens automatically
	ctx := context.Background()
	_, err := client.CreateInstance(ctx, "web-server", "medium")
	if err != nil {
		logger.Error("Failed to create instance", "error", err)
	}
}

func runMultiServiceExample() {
	fmt.Println("\n🔧 2. Multi-Service Example")
	fmt.Println("-" * 30)

	// Setup global logger once
	logger.SetupFromConfig(&[]bool{true}[0], nil)

	// Create services - they all use the same global logger
	authService := NewAuthService("client123", "secret456")
	retryService := NewRetryService()
	jsonService := NewJSONService()

	ctx := context.Background()

	// All services log using the same configured logger
	logger.Info("Starting multi-service example")

	// Auth service logging
	token, _ := authService.GetAccessToken(ctx)
	logger.Debug("Got access token", "token_length", len(token))

	// Retry service logging
	retryService.ExampleAPICall(ctx)

	// JSON service logging
	jsonService.ExampleUsage()

	logger.Info("Multi-service example completed")
}

func runEnvironmentExample() {
	fmt.Println("\n🌍 3. Environment-Based Logging")
	fmt.Println("-" * 30)

	// Production mode (info level only)
	fmt.Println("Production mode:")
	os.Setenv("DATACRUNCH_DEBUG", "false")
	prodClient := NewClient()
	prodClient.CreateInstance(context.Background(), "prod-server", "xlarge")

	fmt.Println("\nDevelopment mode:")
	// Development mode (debug level)
	os.Setenv("DATACRUNCH_DEBUG", "true")
	devClient := NewClient()
	devClient.CreateInstance(context.Background(), "dev-server", "small")

	// Clean up
	os.Unsetenv("DATACRUNCH_DEBUG")
}

func runCustomLoggerExample() {
	fmt.Println("\n🎨 4. Custom Structured Logger")
	fmt.Println("-" * 30)

	// Create custom JSON logger with additional fields
	customLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true, // Add source file info
	})).With(
		"service", "datacrunch-sdk",
		"version", "1.0.0",
		"environment", "development",
	)

	// Use custom logger
	client := NewClient(config.WithLogger(customLogger))
	client.CreateInstance(context.Background(), "custom-server", "large")

	// Manual logging with context
	requestLogger := logger.With(
		"request_id", "req-12345",
		"user_id", "user-67890",
	)

	requestLogger.Info("Custom structured logging example")
	requestLogger.Debug("This shows source file information")
	requestLogger.Warn("This is a warning with full context")
}

// Helper function to repeat strings (for formatting)
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// Use the repeat function for formatting
var _ = repeat // Avoid unused variable warning