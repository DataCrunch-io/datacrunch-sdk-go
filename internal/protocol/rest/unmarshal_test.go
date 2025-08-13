package rest

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/request"
)

func TestUnmarshalResponse_HeaderMapping(t *testing.T) {
	tests := []struct {
		name     string
		target   interface{}
		response *http.Response
		expected interface{}
		wantErr  bool
	}{
		{
			name: "header to string pointer",
			target: &struct {
				RequestID *string `location:"header" locationName:"X-Request-ID"`
			}{},
			response: &http.Response{
				Header: http.Header{
					"X-Request-Id": []string{"req-123"},
				},
			},
			expected: func() interface{} {
				s := "req-123"
				return &struct {
					RequestID *string `location:"header" locationName:"X-Request-ID"`
				}{RequestID: &s}
			}(),
		},
		{
			name: "header to int64 pointer",
			target: &struct {
				ContentLength *int64 `location:"header" locationName:"Content-Length"`
			}{},
			response: &http.Response{
				Header: http.Header{
					"Content-Length": []string{"1024"},
				},
			},
			expected: func() interface{} {
				i := int64(1024)
				return &struct {
					ContentLength *int64 `location:"header" locationName:"Content-Length"`
				}{ContentLength: &i}
			}(),
		},
		{
			name: "status code mapping",
			target: &struct {
				StatusCode *int64 `location:"statusCode"`
			}{},
			response: &http.Response{
				StatusCode: 200,
			},
			expected: func() interface{} {
				i := int64(200)
				return &struct {
					StatusCode *int64 `location:"statusCode"`
				}{StatusCode: &i}
			}(),
		},
		{
			name: "header map",
			target: &struct {
				Metadata map[string]*string `location:"headers" locationName:"X-Meta-"`
			}{},
			response: &http.Response{
				Header: http.Header{
					"X-Meta-Key1": []string{"value1"},
					"X-Meta-Key2": []string{"value2"},
					"Other-Header": []string{"ignore"},
				},
			},
			expected: func() interface{} {
				v1, v2 := "value1", "value2"
				return &struct {
					Metadata map[string]*string `location:"headers" locationName:"X-Meta-"`
				}{Metadata: map[string]*string{
					"Key1": &v1,
					"Key2": &v2,
				}}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalResponse(tt.response, tt.target, false)
			
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Compare the results using reflection
			// This is a simplified comparison - you might want more robust comparison  
			targetVal := reflect.ValueOf(tt.target).Elem()
			expectedVal := reflect.ValueOf(tt.expected).Elem()
			
			// Compare field by field for structs to handle pointer fields correctly
			if targetVal.Kind() == reflect.Struct && expectedVal.Kind() == reflect.Struct {
				for i := 0; i < targetVal.NumField(); i++ {
					targetField := targetVal.Field(i)
					expectedField := expectedVal.Field(i)
					
					if targetField.Kind() == reflect.Ptr && expectedField.Kind() == reflect.Ptr {
						// Both are pointers, compare the values they point to
						if targetField.IsNil() && expectedField.IsNil() {
							continue // both nil, equal
						}
						if targetField.IsNil() != expectedField.IsNil() {
							t.Errorf("field %d: nil mismatch - expected %v, got %v", i, expectedField.IsNil(), targetField.IsNil())
							return
						}
						if !reflect.DeepEqual(targetField.Elem().Interface(), expectedField.Elem().Interface()) {
							t.Errorf("field %d: value mismatch - expected %+v, got %+v", i, expectedField.Elem().Interface(), targetField.Elem().Interface())
							return
						}
					} else if !reflect.DeepEqual(targetField.Interface(), expectedField.Interface()) {
						t.Errorf("field %d: mismatch - expected %+v, got %+v", i, expectedField.Interface(), targetField.Interface())
						return
					}
				}
			} else if !reflect.DeepEqual(tt.target, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, tt.target)
			}
		})
	}
}

func TestUnmarshalBody_PayloadTypes(t *testing.T) {
	tests := []struct {
		name     string
		target   interface{}
		body     string
		expected interface{}
		wantErr  bool
	}{
		{
			name: "byte slice payload",
			target: &struct {
				_    struct{} `payload:"Body"`
				Body []byte   `type:"blob"`
			}{},
			body: "test data",
			expected: &struct {
				_    struct{} `payload:"Body"`
				Body []byte   `type:"blob"`
			}{Body: []byte("test data")},
		},
		{
			name: "string pointer payload",
			target: &struct {
				_    struct{} `payload:"Body"`
				Body *string  `type:"string"`
			}{},
			body: "test string",
			expected: func() interface{} {
				s := "test string"
				return &struct {
					_    struct{} `payload:"Body"`
					Body *string  `type:"string"`
				}{Body: &s}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock request with body
			req := &request.Request{
				HTTPResponse: &http.Response{
					Body: io.NopCloser(bytes.NewBufferString(tt.body)),
				},
				Data: tt.target,
			}

			// Test the unmarshal function
			v := reflect.Indirect(reflect.ValueOf(tt.target))
			err := unmarshalBody(req, v)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Compare results
			if !reflect.DeepEqual(tt.target, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, tt.target)
			}
		})
	}
}

func TestUnmarshalHeader_TypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		target   interface{}
		tag      reflect.StructTag
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "string header",
			header:   "test-value",
			target:   func() interface{} { var s *string; return &s }(),
			expected: func() interface{} { s := "test-value"; return &s }(),
		},
		{
			name:     "int64 header",
			header:   "42",
			target:   func() interface{} { var i *int64; return &i }(),
			expected: func() interface{} { i := int64(42); return &i }(),
		},
		{
			name:     "float64 header",
			header:   "3.14",
			target:   func() interface{} { var f *float64; return &f }(),
			expected: func() interface{} { f := 3.14; return &f }(),
		},
		{
			name:     "bool header",
			header:   "true",
			target:   func() interface{} { var b *bool; return &b }(),
			expected: func() interface{} { b := true; return &b }(),
		},
		{
			name:    "invalid int64",
			header:  "not-a-number",
			target:  func() interface{} { var i *int64; return &i }(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.target).Elem()
			err := unmarshalHeader(v, tt.header, tt.tag)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Compare results by dereferencing pointers for comparison
			targetPtr := reflect.ValueOf(tt.target).Elem()
			expectedPtr := reflect.ValueOf(tt.expected).Elem()
			
			// Handle pointer comparison safely
			if targetPtr.Kind() == reflect.Ptr && expectedPtr.Kind() == reflect.Ptr {
				// Both should be pointers to the same type
				if targetPtr.IsNil() && expectedPtr.IsNil() {
					return // both nil, test passes
				}
				if targetPtr.IsNil() != expectedPtr.IsNil() {
					t.Errorf("nil mismatch: expected nil=%v, got nil=%v", expectedPtr.IsNil(), targetPtr.IsNil())
					return
				}
				
				// Dereference and compare values
				targetVal := targetPtr.Elem().Interface()
				expectedVal := expectedPtr.Elem().Interface()
				if !reflect.DeepEqual(targetVal, expectedVal) {
					t.Errorf("expected %+v, got %+v", expectedVal, targetVal)
				}
			} else {
				// Direct comparison for non-pointer types
				// For this test, both target and expected should actually be pointers
				// We need to dereference both to compare the actual values
				if targetPtr.Kind() == reflect.Ptr && !targetPtr.IsNil() && expectedPtr.Kind() == reflect.Ptr && !expectedPtr.IsNil() {
					targetVal := targetPtr.Elem().Interface()
					expectedVal := expectedPtr.Elem().Interface()
					if !reflect.DeepEqual(targetVal, expectedVal) {
						t.Errorf("expected %+v, got %+v", expectedVal, targetVal)
					}
				} else if !reflect.DeepEqual(targetPtr.Interface(), expectedPtr.Interface()) {
					t.Errorf("expected %+v, got %+v", expectedPtr.Interface(), targetPtr.Interface())
				}
			}
		})
	}
}