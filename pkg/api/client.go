package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/defaults"
)

const (
	DefaultUserAgentBase = "DataCrunch-Go-SDK"

	ContentTypeJSON = "application/json"
)

type Client struct {
	baseURL      string
	userAgent    string
	clientID     string
	clientSecret string
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	httpClient   *http.Client
	log          slog.Logger
}

func NewWithCredentials(clientID, clientSecret string) Client {
	return Client{
		baseURL:      defaults.DefaultBaseURL,
		userAgent:    fmt.Sprintf("%s/%s", DefaultUserAgentBase, defaults.Version),
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   http.DefaultClient,
		log:          *slog.Default(),
	}
}

func (s *Client) SetBaseURL(url string) {
	s.baseURL = url
}

func (s *Client) Request(method, path string, body interface{}, expStatus int, target interface{}) error {
	if err := s.renewAuth(); err != nil {
		return fmt.Errorf("unable to renew auth: %w", err)
	}

	res, err := s.execRequest(method, path, body)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	if res.StatusCode != expStatus {
		return fmt.Errorf("unexpected response with status code %d: %w", res.StatusCode, s.decodeErrorResponse(res))
	}

	if strings.Contains(res.Header.Get("Content-Type"), ContentTypeJSON) && target != nil {
		if err := s.decodeJSONResponse(res, target); err != nil {
			return fmt.Errorf("unable to decode response body: %w", err)
		}
	} else if target != nil { // Normal text response
		textResp, err := s.decodeTextResponse(res)
		if err != nil {
			return fmt.Errorf("unable to decode response body: %w", err)
		}
		if strTarget, ok := target.(*string); ok {
			*strTarget = textResp
		} else {
			return fmt.Errorf("target is not a *string")
		}
	}

	return nil
}

func (s *Client) execRequest(method, path string, body interface{}) (*http.Response, error) {
	var jsonBody io.Reader
	jsonBody = nil
	if body != nil {
		jsonString, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
		jsonBody = bytes.NewReader(jsonString)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", s.baseURL, path), jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	if body != nil {
		req.Header.Set("Content-Type", ContentTypeJSON)
	}

	s.log.Debug("Making request", slog.String("method", method), slog.String("url", req.URL.String()))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("request failed with status code %d: %w", resp.StatusCode, s.decodeErrorResponse(resp))
	}

	return resp, nil
}

func (s *Client) decodeJSONResponse(resp *http.Response, target interface{}) error {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.log.Warn("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}

	return nil
}

func (s *Client) decodeTextResponse(resp *http.Response) (string, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.log.Warn("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	return string(bodyBytes), nil
}

func (s *Client) decodeErrorResponse(resp *http.Response) error {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.log.Warn("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	var errMsg []byte
	_, err := resp.Body.Read(errMsg)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading error response body: %w", err)
	}

	return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(errMsg))
}
