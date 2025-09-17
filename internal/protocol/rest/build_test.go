package rest

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/request"
)

func TestBuild_URIParameters(t *testing.T) {
	tests := []struct {
		name     string
		params   interface{}
		basePath string
		expected string
		wantErr  bool
	}{
		{
			name: "single URI parameter",
			params: &struct {
				ID string `location:"uri" locationName:"id"`
			}{ID: "instance-123"},
			basePath: "/instances/{id}",
			expected: "/instances/instance-123",
		},
		{
			name: "multiple URI parameters",
			params: &struct {
				ProjectID  string `location:"uri" locationName:"project"`
				InstanceID string `location:"uri" locationName:"instance"`
			}{
				ProjectID:  "proj-456",
				InstanceID: "inst-789",
			},
			basePath: "/projects/{project}/instances/{instance}",
			expected: "/projects/proj-456/instances/inst-789",
		},
		{
			name: "URI with special characters (should be escaped)",
			params: &struct {
				Name string `location:"uri" locationName:"name"`
			}{Name: "test name with spaces"},
			basePath: "/items/{name}",
			expected: "/items/test%20name%20with%20spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request
			httpReq, _ := http.NewRequest("GET", "", nil)
			httpReq.URL, _ = url.Parse(tt.basePath)

			req := &request.Request{
				HTTPRequest: httpReq,
				Params:      tt.params,
			}

			// Test the build function
			Build(req)

			if tt.wantErr {
				if req.Error == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if req.Error != nil {
				t.Errorf("unexpected error: %v", req.Error)
				return
			}

			if req.HTTPRequest.URL.Path != tt.expected {
				t.Errorf("expected path %q, got %q", tt.expected, req.HTTPRequest.URL.Path)
			}
		})
	}
}

func TestBuild_QueryParameters(t *testing.T) {
	tests := []struct {
		name     string
		params   interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "single query parameter",
			params: &struct {
				Limit int64 `location:"querystring" locationName:"limit"`
			}{Limit: 10},
			expected: "limit=10",
		},
		{
			name: "multiple query parameters",
			params: &struct {
				Limit  int64  `location:"querystring" locationName:"limit"`
				Filter string `location:"querystring" locationName:"filter"`
			}{
				Limit:  20,
				Filter: "active",
			},
			expected: "filter=active&limit=20", // URL-encoded, so order might vary
		},
		{
			name: "string slice parameter",
			params: &struct {
				Tags []*string `location:"querystring" locationName:"tags"`
			}{
				Tags: []*string{
					func() *string { s := "gpu"; return &s }(),
					func() *string { s := "compute"; return &s }(),
				},
			},
			expected: "tags=gpu&tags=compute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request
			httpReq, _ := http.NewRequest("GET", "http://example.com", nil)

			req := &request.Request{
				HTTPRequest: httpReq,
				Params:      tt.params,
			}

			// Test the build function
			Build(req)

			if tt.wantErr {
				if req.Error == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if req.Error != nil {
				t.Errorf("unexpected error: %v", req.Error)
				return
			}

			// Check query string (order might vary)
			query := req.HTTPRequest.URL.RawQuery
			if !containsAllParams(query, tt.expected) {
				t.Errorf("expected query %q, got %q", tt.expected, query)
			}
		})
	}
}

func TestBuild_HeaderParameters(t *testing.T) {
	tests := []struct {
		name     string
		params   interface{}
		expected map[string]string
		wantErr  bool
	}{
		{
			name: "single header",
			params: &struct {
				Token string `location:"header" locationName:"Authorization"`
			}{Token: "Bearer xyz"},
			expected: map[string]string{
				"Authorization": "Bearer xyz",
			},
		},
		{
			name: "multiple headers",
			params: &struct {
				Auth string `location:"header" locationName:"Authorization"`
				Type string `location:"header" locationName:"Content-Type"`
			}{
				Auth: "Bearer abc",
				Type: "application/json",
			},
			expected: map[string]string{
				"Authorization": "Bearer abc",
				"Content-Type":  "application/json",
			},
		},
		{
			name: "header map",
			params: &struct {
				Metadata map[string]*string `location:"headers" locationName:"X-Meta-"`
			}{
				Metadata: map[string]*string{
					"Version": func() *string { s := "1.0"; return &s }(),
					"Region":  func() *string { s := "us-east-1"; return &s }(),
				},
			},
			expected: map[string]string{
				"X-Meta-Version": "1.0",
				"X-Meta-Region":  "us-east-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request
			httpReq, _ := http.NewRequest("GET", "http://example.com", nil)

			req := &request.Request{
				HTTPRequest: httpReq,
				Params:      tt.params,
			}

			// Test the build function
			Build(req)

			if tt.wantErr {
				if req.Error == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if req.Error != nil {
				t.Errorf("unexpected error: %v", req.Error)
				return
			}

			// Check headers
			for key, expectedValue := range tt.expected {
				if actualValue := req.HTTPRequest.Header.Get(key); actualValue != expectedValue {
					t.Errorf("expected header %s=%q, got %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestBuild_PayloadBody(t *testing.T) {
	tests := []struct {
		name         string
		params       interface{}
		expectedBody string
		wantErr      bool
	}{
		{
			name: "string payload",
			params: &struct {
				_    struct{} `payload:"Body"`
				Body string   `type:"string"`
			}{Body: "test payload"},
			expectedBody: "test payload",
		},
		{
			name: "byte slice payload",
			params: &struct {
				_    struct{} `payload:"Body"`
				Body []byte   `type:"blob"`
			}{Body: []byte("binary data")},
			expectedBody: "binary data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request
			httpReq, _ := http.NewRequest("POST", "http://example.com", nil)

			req := &request.Request{
				HTTPRequest: httpReq,
				Params:      tt.params,
			}

			// Test the build function
			Build(req)

			if tt.wantErr {
				if req.Error == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if req.Error != nil {
				t.Errorf("unexpected error: %v", req.Error)
				return
			}

			// Read body
			if req.Body != nil {
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					t.Errorf("failed to read body: %v", err)
					return
				}

				if string(bodyBytes) != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, string(bodyBytes))
				}
			} else if tt.expectedBody != "" {
				t.Error("expected body but got none")
			}
		})
	}
}

// Helper function to check if query string contains all expected parameters
func containsAllParams(actual, expected string) bool {
	// Simple check - for production you'd want more robust comparison
	actualParams := strings.Split(actual, "&")
	expectedParams := strings.Split(expected, "&")

	for _, expectedParam := range expectedParams {
		found := false
		for _, actualParam := range actualParams {
			if actualParam == expectedParam {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
