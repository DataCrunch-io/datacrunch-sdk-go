package instance

import "context"

// HTTPClient interface for making HTTP requests
// This avoids import cycles by defining the contract without importing datacrunch package
type HTTPClient interface {
	GET(ctx context.Context, path string, queryParams map[string]string, result interface{}) error
	POST(ctx context.Context, path string, body interface{}, result interface{}) error
	PUT(ctx context.Context, path string, body interface{}, result interface{}) error
	DELETE(ctx context.Context, path string, result interface{}) error
}