package restjson

import (
	"fmt"
	"io"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/json/jsonutil"
	"github.com/datacrunch-io/datacrunch-sdk-go/internal/protocol/rest"
)

// BuildHandler is a named request handler for building restjson protocol
// requests
var BuildHandler = request.NamedHandler{
	Name: "datacrunchsdk.restjson.Build",
	Fn:   Build,
}

// UnmarshalHandler is a named request handler for unmarshaling restjson
// protocol requests
var UnmarshalHandler = request.NamedHandler{
	Name: "datacrunchsdk.restjson.Unmarshal",
	Fn:   Unmarshal,
}

// UnmarshalMetaHandler is a named request handler for unmarshaling restjson
// protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler{
	Name: "datacrunchsdk.restjson.UnmarshalMeta",
	Fn:   UnmarshalMeta,
}

// StandardUnmarshalHandler is a named request handler for standard JSON unmarshaling
// Uses Go's standard encoding/json instead of custom jsonutil unmarshaler
var StandardUnmarshalHandler = request.NamedHandler{
	Name: "datacrunchsdk.restjson.StandardUnmarshal",
	Fn:   StandardUnmarshal,
}

// StringUnmarshalHandler is a named request handler for plain string response unmarshaling
// Used for APIs that return plain text strings instead of JSON
var StringUnmarshalHandler = request.NamedHandler{
	Name: "datacrunchsdk.restjson.StringUnmarshal",
	Fn:   StringUnmarshal,
}

// Build builds a request for the REST JSON protocol.
func Build(r *request.Request) {
	rest.Build(r)

	if t := rest.PayloadType(r.Params); t == "structure" || t == "" {
		if v := r.HTTPRequest.Header.Get("Content-Type"); len(v) == 0 {
			r.HTTPRequest.Header.Set("Content-Type", "application/json")
		}
		// Build JSON body using protocol-specific JSON utilities
		if r.ParamsFilled() {
			if body, err := jsonutil.BuildJSON(r.Params); err != nil {
				r.Error = err
			} else {
				r.SetBufferBody(body)
			}
		}
	}
}

// Unmarshal unmarshals a response body for the REST JSON protocol.
func Unmarshal(r *request.Request) {
	if t := rest.PayloadType(r.Data); t == "structure" || t == "" {
		// Unmarshal JSON using protocol-specific JSON utilities
		if r.DataFilled() && r.HTTPResponse.Body != nil {
			defer func() {
				if err := r.HTTPResponse.Body.Close(); err != nil {
					// Log the error but don't fail the function
					_ = err // Suppress unused variable warning
				}
			}()
			if err := jsonutil.UnmarshalJSON(r.Data, r.HTTPResponse.Body); err != nil {
				r.Error = err
			}
		}
	} else {
		rest.Unmarshal(r)
	}
}

// UnmarshalMeta unmarshals response headers for the REST JSON protocol.
func UnmarshalMeta(r *request.Request) {
	rest.UnmarshalMeta(r)
}

// StandardUnmarshal unmarshals a response body using Go's standard encoding/json.
// This is more flexible than the custom jsonutil unmarshaler and handles
// standard JSON types automatically.
func StandardUnmarshal(r *request.Request) {
	if t := rest.PayloadType(r.Data); t == "structure" || t == "" {
		// Unmarshal JSON using standard encoding/json via jsonutil
		if r.DataFilled() && r.HTTPResponse.Body != nil {
			defer func() {
				if err := r.HTTPResponse.Body.Close(); err != nil {
					// Log the error but don't fail the function
					_ = err // Suppress unused variable warning
				}
			}()
			if err := jsonutil.UnmarshalStandardJSON(r.Data, r.HTTPResponse.Body); err != nil {
				r.Error = dcerr.New(request.ErrCodeSerialization, "failed to decode JSON response", err)
			}
		}
	} else {
		// For non-JSON payloads, use the default REST unmarshaler
		rest.Unmarshal(r)
	}
}

// StringUnmarshal unmarshals a plain string response body.
// Used for APIs that return plain text strings instead of JSON objects.
func StringUnmarshal(r *request.Request) {
	if r.DataFilled() && r.HTTPResponse.Body != nil {
		defer func() {
			if err := r.HTTPResponse.Body.Close(); err != nil {
				// Log the error but don't fail the function
				_ = err // Suppress unused variable warning
			}
		}()

		// Read the response body as plain text
		body, err := io.ReadAll(r.HTTPResponse.Body)
		if err != nil {
			r.Error = dcerr.New(request.ErrCodeSerialization, "failed to read string response", err)
			return
		}

		// Set the string value directly
		if stringPtr, ok := r.Data.(*string); ok {
			*stringPtr = string(body)
		} else {
			r.Error = dcerr.New(request.ErrCodeSerialization,
				fmt.Sprintf("expected *string data type, got %T", r.Data), nil)
		}
	}
}
