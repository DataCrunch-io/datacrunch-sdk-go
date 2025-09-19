package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/model"
)

const (
	PathAuthToken = "/oauth2/token"
)

func (s *Client) renewAuth() error {
	// If accessToken is available and expiry is in the future, no need to renew
	if s.accessToken != "" && time.Now().Before(s.tokenExpiry.Add(30*time.Second)) {
		return nil
	}

	var reqBody interface{}
	if s.refreshToken != "" { // If refreshToken is available, use it to get a new accessToken
		reqBody = model.RefreshRequest{
			RefreshToken: s.refreshToken,
			GrantType:    model.GrantTypeRefreshToken,
		}
	} else { // Otherwise, use client credentials to get a new accessToken
		reqBody = model.AuthRequest{
			ClientID:     s.clientID,
			ClientSecret: s.clientSecret,
			GrantType:    model.GrantTypeClientCredentials,
		}
	}

	res, err := s.execRequest(http.MethodPost, PathAuthToken, reqBody)
	if err != nil {
		return fmt.Errorf("error refreshing token: %w", err)
	}

	if res.StatusCode != http.StatusOK { // Should be 201 according to docs
		return fmt.Errorf("unexpected error with status %d while refreshing token: %w", res.StatusCode, s.decodeErrorResponse(res))
	}

	var authResp model.AuthResponse
	if err := s.decodeJSONResponse(res, &authResp); err != nil {
		return fmt.Errorf("unable to decode response body: %w", err)
	}

	s.accessToken = authResp.AccessToken
	s.refreshToken = authResp.RefreshToken
	s.tokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second)

	return nil
}
