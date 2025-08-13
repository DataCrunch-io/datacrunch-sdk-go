// Example: Authentication service using global logger
package main

import (
	"context"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/internal/logger"
)

// AuthService - handles authentication
type AuthService struct {
	clientID     string
	clientSecret string
}

func NewAuthService(clientID, clientSecret string) *AuthService {
	return &AuthService{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (a *AuthService) GetAccessToken(ctx context.Context) (string, error) {
	// Just import and use logger directly!
	logger.Debug("Starting OAuth token request",
		"grant_type", "client_credentials",
		"client_id", logger.SanitizeToken(a.clientID), // Mask sensitive data
	)

	// Simulate token request
	token, err := a.requestToken(ctx)
	if err != nil {
		logger.Error("Token request failed", "error", err)
		return "", err
	}

	logger.Info("Access token obtained successfully",
		"expires_in", "3600s",
	)

	return token, nil
}

func (a *AuthService) requestToken(ctx context.Context) (string, error) {
	logger.Debug("Making HTTP request to token endpoint",
		"method", "POST",
		"endpoint", "/oauth2/token",
	)

	start := time.Now()
	
	// Simulate HTTP call
	time.Sleep(100 * time.Millisecond)
	
	logger.Debug("Token request completed",
		"status_code", 200,
		"duration", time.Since(start),
	)

	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	logger.Debug("Refreshing access token",
		"refresh_token", logger.SanitizeToken(refreshToken),
	)

	// Add contextual logger for this operation
	refreshLogger := logger.With(
		"operation", "token_refresh",
		"client_id", logger.SanitizeToken(a.clientID),
	)

	refreshLogger.Debug("Starting token refresh")

	// Simulate refresh
	time.Sleep(80 * time.Millisecond)

	refreshLogger.Info("Token refreshed successfully")
	return "new_access_token_here", nil
}