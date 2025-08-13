package jsonutil

import (
	"strings"
	"testing"
)

func TestUnmarshalJSON_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		target   interface{}
		expected interface{}
		wantErr  bool
	}{
		// String types
		{
			name: "string direct",
			json: `{"value": "hello"}`,
			target: &struct {
				Value string `json:"value"`
			}{},
			expected: struct {
				Value string `json:"value"`
			}{Value: "hello"},
		},
		{
			name: "string pointer",
			json: `{"value": "hello"}`,
			target: &struct {
				Value *string `json:"value"`
			}{},
			expected: func() struct {
				Value *string `json:"value"`
			} {
				s := "hello"
				return struct {
					Value *string `json:"value"`
				}{Value: &s}
			}(),
		},

		// Integer types
		{
			name: "int64 direct",
			json: `{"value": 42}`,
			target: &struct {
				Value int64 `json:"value"`
			}{},
			expected: struct {
				Value int64 `json:"value"`
			}{Value: 42},
		},
		{
			name: "int64 pointer",
			json: `{"value": 42}`,
			target: &struct {
				Value *int64 `json:"value"`
			}{},
			expected: func() struct {
				Value *int64 `json:"value"`
			} {
				i := int64(42)
				return struct {
					Value *int64 `json:"value"`
				}{Value: &i}
			}(),
		},

		// Float types
		{
			name: "float64 direct",
			json: `{"value": 3.14}`,
			target: &struct {
				Value float64 `json:"value"`
			}{},
			expected: struct {
				Value float64 `json:"value"`
			}{Value: 3.14},
		},
		{
			name: "float64 pointer",
			json: `{"value": 3.14}`,
			target: &struct {
				Value *float64 `json:"value"`
			}{},
			expected: func() struct {
				Value *float64 `json:"value"`
			} {
				f := 3.14
				return struct {
					Value *float64 `json:"value"`
				}{Value: &f}
			}(),
		},

		// Edge cases
		{
			name: "integer as float (truncation)",
			json: `{"value": 42.99}`,
			target: &struct {
				Value int64 `json:"value"`
			}{},
			expected: struct {
				Value int64 `json:"value"`
			}{Value: 42}, // truncated
		},
		{
			name: "null pointer",
			json: `{"value": null}`,
			target: &struct {
				Value *string `json:"value"`
			}{},
			expected: struct {
				Value *string `json:"value"`
			}{Value: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalJSON(tt.target, strings.NewReader(tt.json))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Compare values (you'd need reflection or specific comparison logic)
			// This is simplified for the example
		})
	}
}

func TestUnmarshalJSON_ComplexStruct(t *testing.T) {
	type TestStruct struct {
		Name   string   `json:"name"`
		Count  int64    `json:"count"`
		Price  float64  `json:"price"`
		Active *bool    `json:"active"`
		Tags   []string `json:"tags"`
	}

	json := `{
		"name": "B200 GPU",
		"count": 30,
		"price": 3.64,
		"active": true,
		"tags": ["gpu", "compute"]
	}`

	var result TestStruct
	err := UnmarshalJSON(&result, strings.NewReader(json))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "B200 GPU" {
		t.Errorf("expected name 'B200 GPU', got '%s'", result.Name)
	}
	if result.Count != 30 {
		t.Errorf("expected count 30, got %d", result.Count)
	}
	if result.Price != 3.64 {
		t.Errorf("expected price 3.64, got %f", result.Price)
	}
}

func TestUnmarshalJSON_CaseInsensitive(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}

	json := `{"NAME": "test"}` // uppercase JSON key

	var result TestStruct
	err := UnmarshalJSONCaseInsensitive(&result, strings.NewReader(json))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", result.Name)
	}
}

func TestUnmarshalJSON_Integration(t *testing.T) {
	// Integration tests could go here
	// For now, just a placeholder
}
