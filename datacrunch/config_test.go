package datacrunch

import (
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		options  []Option
		expected *Config
	}{
		{
			name:    "default configuration",
			options: []Option{},
			expected: &Config{
				BaseURL:       "https://api.datacrunch.io/v1",
				Timeout:       30 * time.Second,
				MaxRetries:    3,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "with base URL option",
			options: []Option{
				WithBaseURL("https://custom.api.com/v2"),
			},
			expected: &Config{
				BaseURL:       "https://custom.api.com/v2",
				Timeout:       30 * time.Second,
				MaxRetries:    3,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "with credentials option",
			options: []Option{
				WithCredentials("test-client-id", "test-client-secret"),
			},
			expected: &Config{
				BaseURL:       "https://api.datacrunch.io/v1",
				ClientID:      "test-client-id",
				ClientSecret:  "test-client-secret",
				Timeout:       30 * time.Second,
				MaxRetries:    3,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "with timeout option",
			options: []Option{
				WithTimeout(60 * time.Second),
			},
			expected: &Config{
				BaseURL:       "https://api.datacrunch.io/v1",
				Timeout:       60 * time.Second,
				MaxRetries:    3,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "with retry config option",
			options: []Option{
				WithRetryConfig(5, 2*time.Second, 60*time.Second),
			},
			expected: &Config{
				BaseURL:       "https://api.datacrunch.io/v1",
				Timeout:       30 * time.Second,
				MaxRetries:    5,
				RetryDelay:    2 * time.Second,
				MaxRetryDelay: 60 * time.Second,
			},
		},
		{
			name: "multiple options",
			options: []Option{
				WithBaseURL("https://staging.api.com/v1"),
				WithCredentials("staging-id", "staging-secret"),
				WithTimeout(45 * time.Second),
				WithRetryConfig(2, 500*time.Millisecond, 10*time.Second),
			},
			expected: &Config{
				BaseURL:       "https://staging.api.com/v1",
				ClientID:      "staging-id",
				ClientSecret:  "staging-secret",
				Timeout:       45 * time.Second,
				MaxRetries:    2,
				RetryDelay:    500 * time.Millisecond,
				MaxRetryDelay: 10 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.options...)

			// Verify the client was created
			if client == nil {
				t.Fatal("Expected client to be created, got nil")
			}

			// Verify configuration
			if client.config.BaseURL != tt.expected.BaseURL {
				t.Errorf("Expected BaseURL %s, got %s", tt.expected.BaseURL, client.config.BaseURL)
			}
			if client.config.ClientID != tt.expected.ClientID {
				t.Errorf("Expected ClientID %s, got %s", tt.expected.ClientID, client.config.ClientID)
			}
			if client.config.ClientSecret != tt.expected.ClientSecret {
				t.Errorf("Expected ClientSecret %s, got %s", tt.expected.ClientSecret, client.config.ClientSecret)
			}
			if client.config.Timeout != tt.expected.Timeout {
				t.Errorf("Expected Timeout %v, got %v", tt.expected.Timeout, client.config.Timeout)
			}
			if client.config.MaxRetries != tt.expected.MaxRetries {
				t.Errorf("Expected MaxRetries %d, got %d", tt.expected.MaxRetries, client.config.MaxRetries)
			}
			if client.config.RetryDelay != tt.expected.RetryDelay {
				t.Errorf("Expected RetryDelay %v, got %v", tt.expected.RetryDelay, client.config.RetryDelay)
			}
			if client.config.MaxRetryDelay != tt.expected.MaxRetryDelay {
				t.Errorf("Expected MaxRetryDelay %v, got %v", tt.expected.MaxRetryDelay, client.config.MaxRetryDelay)
			}

			// Verify service clients are initialized
			if client.Instance == nil {
				t.Error("Expected Instance service to be initialized")
			}
			if client.SSHKeys == nil {
				t.Error("Expected SSHKeys service to be initialized")
			}
			if client.StartScripts == nil {
				t.Error("Expected StartScripts service to be initialized")
			}
		})
	}
}

func TestNewFromEnv(t *testing.T) {
	// Save original environment
	originalVars := map[string]string{
		"DATACRUNCH_BASE_URL":      os.Getenv("DATACRUNCH_BASE_URL"),
		"DATACRUNCH_CLIENT_ID":     os.Getenv("DATACRUNCH_CLIENT_ID"),
		"DATACRUNCH_CLIENT_SECRET": os.Getenv("DATACRUNCH_CLIENT_SECRET"),
		"DATACRUNCH_TIMEOUT":       os.Getenv("DATACRUNCH_TIMEOUT"),
		"DATACRUNCH_MAX_RETRIES":   os.Getenv("DATACRUNCH_MAX_RETRIES"),
	}

	// Cleanup function to restore environment
	cleanup := func() {
		for key, value := range originalVars {
			if value == "" {
				_ = os.Unsetenv(key)
			} else {
				_ = os.Setenv(key, value)
			}
		}
	}
	defer cleanup()

	tests := []struct {
		name     string
		envVars  map[string]string
		options  []Option
		expected *Config
	}{
		{
			name:    "default values when no env vars set",
			envVars: map[string]string{}, // Clear all env vars
			options: []Option{},
			expected: &Config{
				BaseURL:       "https://api.datacrunch.io/v1",
				Timeout:       30 * time.Second,
				MaxRetries:    3,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "load from environment variables",
			envVars: map[string]string{
				"DATACRUNCH_BASE_URL":      "https://env.api.com/v1",
				"DATACRUNCH_CLIENT_ID":     "env-client-id",
				"DATACRUNCH_CLIENT_SECRET": "env-client-secret",
				"DATACRUNCH_TIMEOUT":       "60s",
				"DATACRUNCH_MAX_RETRIES":   "5",
			},
			options: []Option{},
			expected: &Config{
				BaseURL:       "https://env.api.com/v1",
				ClientID:      "env-client-id",
				ClientSecret:  "env-client-secret",
				Timeout:       60 * time.Second,
				MaxRetries:    5,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "options override environment variables",
			envVars: map[string]string{
				"DATACRUNCH_BASE_URL":      "https://env.api.com/v1",
				"DATACRUNCH_CLIENT_ID":     "env-client-id",
				"DATACRUNCH_CLIENT_SECRET": "env-client-secret",
				"DATACRUNCH_TIMEOUT":       "60s",
			},
			options: []Option{
				WithBaseURL("https://override.api.com/v1"),
				WithTimeout(120 * time.Second),
			},
			expected: &Config{
				BaseURL:       "https://override.api.com/v1",
				ClientID:      "env-client-id",
				ClientSecret:  "env-client-secret",
				Timeout:       120 * time.Second,
				MaxRetries:    3,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "invalid timeout format uses default",
			envVars: map[string]string{
				"DATACRUNCH_TIMEOUT": "invalid-timeout",
			},
			options: []Option{},
			expected: &Config{
				BaseURL:       "https://api.datacrunch.io/v1",
				Timeout:       30 * time.Second, // Default value
				MaxRetries:    3,
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
		{
			name: "invalid max retries format uses default",
			envVars: map[string]string{
				"DATACRUNCH_MAX_RETRIES": "invalid-number",
			},
			options: []Option{},
			expected: &Config{
				BaseURL:       "https://api.datacrunch.io/v1",
				Timeout:       30 * time.Second,
				MaxRetries:    3, // Default value
				RetryDelay:    1 * time.Second,
				MaxRetryDelay: 30 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables first
			for key := range originalVars {
				_ = os.Unsetenv(key)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				_ = os.Setenv(key, value)
			}

			client := NewFromEnv(tt.options...)

			// Verify the client was created
			if client == nil {
				t.Fatal("Expected client to be created, got nil")
			}

			// Verify configuration
			if client.config.BaseURL != tt.expected.BaseURL {
				t.Errorf("Expected BaseURL %s, got %s", tt.expected.BaseURL, client.config.BaseURL)
			}
			if client.config.ClientID != tt.expected.ClientID {
				t.Errorf("Expected ClientID %s, got %s", tt.expected.ClientID, client.config.ClientID)
			}
			if client.config.ClientSecret != tt.expected.ClientSecret {
				t.Errorf("Expected ClientSecret %s, got %s", tt.expected.ClientSecret, client.config.ClientSecret)
			}
			if client.config.Timeout != tt.expected.Timeout {
				t.Errorf("Expected Timeout %v, got %v", tt.expected.Timeout, client.config.Timeout)
			}
			if client.config.MaxRetries != tt.expected.MaxRetries {
				t.Errorf("Expected MaxRetries %d, got %d", tt.expected.MaxRetries, client.config.MaxRetries)
			}
		})
	}
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name:   "nil config uses default",
			config: nil,
		},
		{
			name: "custom config",
			config: &Config{
				BaseURL:       "https://custom.api.com/v1",
				ClientID:      "custom-id",
				ClientSecret:  "custom-secret",
				Timeout:       90 * time.Second,
				MaxRetries:    10,
				RetryDelay:    2 * time.Second,
				MaxRetryDelay: 60 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewWithConfig(tt.config)

			// Verify the client was created
			if client == nil {
				t.Fatal("Expected client to be created, got nil")
			}

			// Verify service clients are initialized
			if client.Instance == nil {
				t.Error("Expected Instance service to be initialized")
			}
			if client.SSHKeys == nil {
				t.Error("Expected SSHKeys service to be initialized")
			}
			if client.StartScripts == nil {
				t.Error("Expected StartScripts service to be initialized")
			}

			// If custom config was provided, verify values
			if tt.config != nil {
				if client.config.BaseURL != tt.config.BaseURL {
					t.Errorf("Expected BaseURL %s, got %s", tt.config.BaseURL, client.config.BaseURL)
				}
				if client.config.ClientID != tt.config.ClientID {
					t.Errorf("Expected ClientID %s, got %s", tt.config.ClientID, client.config.ClientID)
				}
				if client.config.ClientSecret != tt.config.ClientSecret {
					t.Errorf("Expected ClientSecret %s, got %s", tt.config.ClientSecret, client.config.ClientSecret)
				}
				if client.config.Timeout != tt.config.Timeout {
					t.Errorf("Expected Timeout %v, got %v", tt.config.Timeout, client.config.Timeout)
				}
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("Expected default config to be created, got nil")
	}

	// Verify default values
	expectedDefaults := &Config{
		BaseURL:       "https://api.datacrunch.io/v1",
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryDelay:    1 * time.Second,
		MaxRetryDelay: 30 * time.Second,
	}

	if config.BaseURL != expectedDefaults.BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", expectedDefaults.BaseURL, config.BaseURL)
	}
	if config.Timeout != expectedDefaults.Timeout {
		t.Errorf("Expected Timeout %v, got %v", expectedDefaults.Timeout, config.Timeout)
	}
	if config.MaxRetries != expectedDefaults.MaxRetries {
		t.Errorf("Expected MaxRetries %d, got %d", expectedDefaults.MaxRetries, config.MaxRetries)
	}
	if config.RetryDelay != expectedDefaults.RetryDelay {
		t.Errorf("Expected RetryDelay %v, got %v", expectedDefaults.RetryDelay, config.RetryDelay)
	}
	if config.MaxRetryDelay != expectedDefaults.MaxRetryDelay {
		t.Errorf("Expected MaxRetryDelay %v, got %v", expectedDefaults.MaxRetryDelay, config.MaxRetryDelay)
	}
}

func TestOptions(t *testing.T) {
	t.Run("WithBaseURL", func(t *testing.T) {
		config := &Config{}
		option := WithBaseURL("https://test.api.com")
		option(config)

		if config.BaseURL != "https://test.api.com" {
			t.Errorf("Expected BaseURL to be set to https://test.api.com, got %s", config.BaseURL)
		}
	})

	t.Run("WithCredentials", func(t *testing.T) {
		config := &Config{}
		option := WithCredentials("test-id", "test-secret")
		option(config)

		if config.ClientID != "test-id" {
			t.Errorf("Expected ClientID to be set to test-id, got %s", config.ClientID)
		}
		if config.ClientSecret != "test-secret" {
			t.Errorf("Expected ClientSecret to be set to test-secret, got %s", config.ClientSecret)
		}
	})

	t.Run("WithTimeout", func(t *testing.T) {
		config := &Config{}
		option := WithTimeout(45 * time.Second)
		option(config)

		if config.Timeout != 45*time.Second {
			t.Errorf("Expected Timeout to be set to 45s, got %v", config.Timeout)
		}
	})

	t.Run("WithRetryConfig", func(t *testing.T) {
		config := &Config{}
		option := WithRetryConfig(5, 2*time.Second, 60*time.Second)
		option(config)

		if config.MaxRetries != 5 {
			t.Errorf("Expected MaxRetries to be set to 5, got %d", config.MaxRetries)
		}
		if config.RetryDelay != 2*time.Second {
			t.Errorf("Expected RetryDelay to be set to 2s, got %v", config.RetryDelay)
		}
		if config.MaxRetryDelay != 60*time.Second {
			t.Errorf("Expected MaxRetryDelay to be set to 60s, got %v", config.MaxRetryDelay)
		}
	})

	t.Run("WithRetryer", func(t *testing.T) {
		config := &Config{}
		mockRetryer := "mock-retryer"
		option := WithRetryer(mockRetryer)
		option(config)

		if config.Retryer != mockRetryer {
			t.Errorf("Expected Retryer to be set to mock retryer, got %v", config.Retryer)
		}
	})
}
