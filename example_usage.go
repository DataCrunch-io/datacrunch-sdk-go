package main

import (
	"log/slog"
	"os"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/config"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

func main() {
	// Example 1: Enable debug logging
	cfg := &config.Config{}
	config.WithDebug(true)(cfg)
	
	logger := logger.GetLogger(cfg.Debug, cfg.Logger)
	logger.Info("Debug logging enabled")

	// Example 2: Custom logger
	customLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	
	cfg2 := &config.Config{}
	config.WithLogger(customLogger)(cfg2)

	// Example 3: Environment variable
	// export DATACRUNCH_DEBUG=true
	cfg3 := &config.Config{}
	logger3 := logger.GetLogger(nil, nil) // Will check environment
	
	// Example 4: Integration with request handlers
	// handlers := &request.Handlers{}
	// logger.SetupRequestLogging(handlers, logger)

	// Example 5: Manual logging
	logger.LogCredentialRetrieval(logger, "EnvProvider", true, nil)
	
	// Example 6: Authentication logging
	logger.LogTokenRequest(logger, "client_credentials", "https://api.datacrunch.io/v1/oauth2/token", "client123")
	
	// Example 7: Performance logging
	// startTime := logger.LogRequestStart(logger, "CreateInstance", "POST", "/instances", 1024)
	// // ... make request ...
	// metrics := logger.PerformanceMetrics{
	//     Operation: "CreateInstance",
	//     TotalDuration: time.Since(startTime),
	//     StatusCode: 200,
	//     // ... other metrics
	// }
	// logger.LogRequestComplete(logger, metrics)
}

// Example integration in your client struct
type DataCrunchClient struct {
	config *config.Config
	logger *slog.Logger
}

func NewClient(opts ...config.Option) *DataCrunchClient {
	cfg := &config.Config{}
	
	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}
	
	// Get logger from config
	logger := logger.GetLogger(cfg.Debug, cfg.Logger)
	
	return &DataCrunchClient{
		config: cfg,
		logger: logger,
	}
}

func (c *DataCrunchClient) CreateInstance(name string) error {
	// Log the operation start
	c.logger.Info("Creating instance", "name", name)
	
	// Your existing API call logic...
	// When making HTTP requests, the logging handlers will automatically log
	
	return nil
}