// Example: Retry service using global logger
package main

import (
	"context"
	"errors"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// RetryService - handles retries
type RetryService struct {
	maxRetries int
	baseDelay  time.Duration
}

func NewRetryService() *RetryService {
	return &RetryService{
		maxRetries: 3,
		baseDelay:  1 * time.Second,
	}
}

func (r *RetryService) ExecuteWithRetry(ctx context.Context, operation string, fn func() error) error {
	var lastErr error
	
	for attempt := 1; attempt <= r.maxRetries; attempt++ {
		// Use logger directly - no setup needed!
		logger.Debug("Executing operation",
			"operation", operation,
			"attempt", attempt,
			"max_attempts", r.maxRetries,
		)

		err := fn()
		if err == nil {
			if attempt > 1 {
				logger.Info("Operation succeeded after retries",
					"operation", operation,
					"successful_attempt", attempt,
					"total_attempts", attempt,
				)
			}
			return nil
		}

		lastErr = err
		
		if attempt == r.maxRetries {
			logger.Error("All retry attempts failed",
				"operation", operation,
				"attempts", attempt,
				"final_error", err,
			)
			break
		}

		// Calculate delay
		delay := r.calculateDelay(attempt)
		
		logger.Warn("Operation failed, retrying",
			"operation", operation,
			"attempt", attempt,
			"error", err,
			"retry_delay", delay,
			"next_attempt", attempt+1,
		)

		// Wait before retry
		select {
		case <-ctx.Done():
			logger.Error("Operation cancelled during retry delay", "operation", operation)
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}

func (r *RetryService) calculateDelay(attempt int) time.Duration {
	// Exponential backoff
	delay := r.baseDelay * time.Duration(1<<uint(attempt-1))
	
	logger.Debug("Calculated retry delay",
		"attempt", attempt,
		"base_delay", r.baseDelay,
		"calculated_delay", delay,
	)
	
	return delay
}

// Example usage
func (r *RetryService) ExampleAPICall(ctx context.Context) error {
	return r.ExecuteWithRetry(ctx, "CreateInstance", func() error {
		// Simulate flaky API call
		logger.Debug("Making API call to create instance")
		
		// Simulate 50% failure rate
		if time.Now().UnixNano()%2 == 0 {
			return errors.New("API temporarily unavailable")
		}
		
		logger.Debug("API call succeeded")
		return nil
	})
}