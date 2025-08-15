package rest

import (
	"reflect"
	"testing"
)

func TestPayloadMember(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "struct with string payload",
			input: &struct {
				_    struct{} `payload:"Body"`
				Body string   `type:"structure"`
			}{Body: "test payload"},
			expected: "test payload",
		},
		{
			name: "struct with byte slice payload",
			input: &struct {
				_    struct{} `payload:"Data"`
				Data []byte   `type:"structure"`
			}{Data: []byte("binary data")},
			expected: []byte("binary data"),
		},
		{
			name: "struct without payload",
			input: &struct {
				Name string
				ID   int
			}{Name: "test", ID: 42},
			expected: nil,
		},
		{
			name: "struct with string payload type",
			input: &struct {
				_    struct{} `payload:"Body"`
				Body string   `type:"string"` // string payload should now be supported
			}{Body: "test"},
			expected: "test",
		},
		{
			name: "struct with invalid payload field",
			input: &struct {
				_    struct{} `payload:"NonExistent"`
				Body string
			}{Body: "test"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PayloadMember(tt.input)
			
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v (%T), got %+v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

func TestPayloadType(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name: "struct with blob payload",
			input: &struct {
				_    struct{} `payload:"Body"`
				Body []byte   `type:"blob"`
			}{},
			expected: "blob",
		},
		{
			name: "struct with string payload",
			input: &struct {
				_    struct{} `payload:"Content"`
				Content string `type:"string"`
			}{},
			expected: "string",
		},
		{
			name: "struct with structure payload",
			input: &struct {
				_       struct{} `payload:"Data"`
				Data    struct{} `type:"structure"`
			}{},
			expected: "structure",
		},
		{
			name: "struct with list payload",
			input: &struct {
				_     struct{}   `payload:"Items"`
				Items []string `type:"list"`
			}{},
			expected: "list",
		},
		{
			name: "struct without payload",
			input: &struct {
				Name string
				ID   int
			}{},
			expected: "",
		},
		{
			name: "struct with nopayload marker",
			input: &struct {
				_ struct{} `nopayload:"true"`
			}{},
			expected: "nopayload",
		},
		{
			name: "struct with payload field but no type",
			input: &struct {
				_    struct{} `payload:"Body"`
				Body string   // no type tag
			}{},
			expected: "",
		},
		{
			name: "struct with invalid payload field name",
			input: &struct {
				_    struct{} `payload:"NonExistent"`
				Body string   `type:"string"`
			}{},
			expected: "",
		},
		{
			name:     "non-struct input",
			input:    "just a string",
			expected: "",
		},
		{
			name:     "slice input",
			input:    []string{"item1", "item2"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PayloadType(tt.input)
			
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPayloadType_WithPointers(t *testing.T) {
	// Test with pointer inputs
	data := struct {
		_    struct{} `payload:"Body"`
		Body string   `type:"string"`
	}{}

	// Test direct struct
	result1 := PayloadType(data)
	if result1 != "string" {
		t.Errorf("expected 'string' for direct struct, got %q", result1)
	}

	// Test pointer to struct  
	result2 := PayloadType(&data)
	if result2 != "string" {
		t.Errorf("expected 'string' for pointer to struct, got %q", result2)
	}

	// Test pointer to pointer (should still work)
	dataPtr := &data
	result3 := PayloadType(&dataPtr)
	if result3 != "string" {
		t.Errorf("expected 'string' for pointer to pointer, got %q", result3)
	}
}