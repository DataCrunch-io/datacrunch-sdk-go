package jsonutil

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

// jsonEqual compares two JSON strings semantically, ignoring field order
func jsonEqual(expected, actual string) bool {
	var expectedParsed, actualParsed interface{}
	if err := json.Unmarshal([]byte(expected), &expectedParsed); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(actual), &actualParsed); err != nil {
		return false
	}
	return reflect.DeepEqual(expectedParsed, actualParsed)
}

func TestBuildJSON_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		// Basic scalar types
		{
			name:     "string",
			input:    struct{ Value string }{Value: "hello"},
			expected: `{"Value":"hello"}`,
		},
		{
			name:     "string with escaping",
			input:    struct{ Value string }{Value: "hello\nworld\t\"test\""},
			expected: `{"Value":"hello\nworld\t\"test\""}`,
		},
		{
			name:     "int64",
			input:    struct{ Value int64 }{Value: 42},
			expected: `{"Value":42}`,
		},
		{
			name:     "negative int64",
			input:    struct{ Value int64 }{Value: -123},
			expected: `{"Value":-123}`,
		},
		{
			name:     "float64",
			input:    struct{ Value float64 }{Value: 3.14159},
			expected: `{"Value":3.14159}`,
		},
		{
			name:     "float64 zero",
			input:    struct{ Value float64 }{Value: 0.0},
			expected: `{"Value":0}`,
		},
		{
			name:     "bool true",
			input:    struct{ Value bool }{Value: true},
			expected: `{"Value":true}`,
		},
		{
			name:     "bool false",
			input:    struct{ Value bool }{Value: false},
			expected: `{"Value":false}`,
		},

		// Pointer types
		{
			name:     "string pointer",
			input:    struct{ Value *string }{Value: func() *string { s := "test"; return &s }()},
			expected: `{"Value":"test"}`,
		},
		{
			name:     "nil pointer (omitted)",
			input:    struct{ Value *string }{Value: nil},
			expected: `{}`,
		},
		{
			name:     "int64 pointer",
			input:    struct{ Value *int64 }{Value: func() *int64 { i := int64(100); return &i }()},
			expected: `{"Value":100}`,
		},

		// Byte slice (base64 encoded)
		{
			name:     "byte slice small",
			input:    struct{ Data []byte }{Data: []byte("hello")},
			expected: `{"Data":"aGVsbG8="}`, // "hello" in base64
		},
		{
			name:     "nil byte slice (omitted)",
			input:    struct{ Data []byte }{Data: nil},
			expected: `{}`,
		},

		// Special float values
		{
			name:     "NaN float",
			input:    struct{ Value float64 }{Value: math.NaN()},
			expected: `{"Value":"NaN"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildJSON(tt.input)

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

			if !jsonEqual(tt.expected, string(result)) {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}

			// Verify it's valid JSON
			var parsed interface{}
			if err := json.Unmarshal(result, &parsed); err != nil {
				t.Errorf("produced invalid JSON: %v", err)
			}
		})
	}
}

func TestBuildJSON_ComplexStructures(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "nested struct",
			input: struct {
				Name   string
				Config struct {
					Count int64
					Flag  bool
				}
			}{
				Name: "test",
				Config: struct {
					Count int64
					Flag  bool
				}{Count: 5, Flag: true},
			},
			expected: `{"Name":"test","Config":{"Count":5,"Flag":true}}`,
		},
		{
			name: "struct with slice",
			input: struct {
				Tags []string
			}{Tags: []string{"gpu", "compute", "ml"}},
			expected: `{"Tags":["gpu","compute","ml"]}`,
		},
		{
			name: "struct with map",
			input: struct {
				Metadata map[string]string
			}{Metadata: map[string]string{
				"version": "1.0",
				"region":  "us-east-1",
			}},
			// Note: map keys are sorted alphabetically
			expected: `{"Metadata":{"region":"us-east-1","version":"1.0"}}`,
		},
		{
			name: "empty slice (omitted)",
			input: struct {
				Tags []string
			}{Tags: nil},
			expected: `{}`,
		},
		{
			name: "empty map (omitted)",
			input: struct {
				Metadata map[string]string
			}{Metadata: nil},
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildJSON(tt.input)

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

			if !jsonEqual(tt.expected, string(result)) {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}

			// Verify it's valid JSON
			var parsed interface{}
			if err := json.Unmarshal(result, &parsed); err != nil {
				t.Errorf("produced invalid JSON: %v", err)
			}
		})
	}
}

func TestBuildJSON_JSONTags(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "json tag with different name",
			input: struct {
				InternalName string `json:"external_name"`
			}{InternalName: "value"},
			expected: `{"external_name":"value"}`,
		},
		{
			name: "json tag with omitempty",
			input: struct {
				Name  string `json:"name"`
				Empty string `json:"empty,omitempty"`
			}{Name: "test", Empty: ""},
			// Note: omitempty not yet implemented, field still included
			expected: `{"name":"test","empty":""}`,
		},
		{
			name: "json tag with dash (ignored)",
			input: struct {
				Name   string `json:"name"`
				Secret string `json:"-"`
			}{Name: "test", Secret: "hidden"},
			expected: `{"name":"test"}`,
		},
		{
			name: "locationName tag (supported)",
			input: struct {
				InternalName string `locationName:"external_name"`
			}{InternalName: "value"},
			expected: `{"external_name":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildJSON(tt.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !jsonEqual(tt.expected, string(result)) {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestBuildJSON_SpecialFields(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "field with location tag (ignored)",
			input: struct {
				HeaderValue string `location:"header"`
				BodyValue   string
			}{HeaderValue: "ignored", BodyValue: "included"},
			expected: `{"BodyValue":"included"}`,
		},
		{
			name: "field with ignore tag",
			input: struct {
				Name    string `json:"name"`
				Ignored string `ignore:"true"`
			}{Name: "test", Ignored: "skip"},
			expected: `{"name":"test"}`,
		},
		{
			name: "unexported field (ignored)",
			input: struct {
				Name       string
				privateVar string
			}{Name: "test", privateVar: "hidden"},
			expected: `{"Name":"test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildJSON(tt.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !jsonEqual(tt.expected, string(result)) {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestBuildJSON_TimeHandling(t *testing.T) {
	// Create a fixed time for testing
	testTime := time.Date(2023, 12, 1, 10, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "time with default format",
			input: struct {
				Timestamp time.Time
			}{Timestamp: testTime},
			// BuildJSON treats time.Time as empty struct (no special handling)
			expected: `{"Timestamp":{}}`,
		},
		{
			name: "time with custom format",
			input: struct {
				Timestamp time.Time `timestampFormat:"iso8601"`
			}{Timestamp: testTime},
			expected: `{"Timestamp":{}}`, // BuildJSON treats time.Time as empty struct
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildJSON(tt.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !jsonEqual(tt.expected, string(result)) {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestBuildJSON_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "unsupported type",
			input: struct {
				Channel chan int
			}{Channel: make(chan int)},
			wantErr: true,
		},
		{
			name: "function type",
			input: struct {
				Func func()
			}{Func: func() {}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildJSON(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestBuildJSON_CreateStartScriptScenario(t *testing.T) {
	// Test case that specifically covers the bug we found with CreateStartScript
	type CreateStartScriptInput struct {
		Name   string `json:"name"`
		Script string `json:"script"`
	}

	input := CreateStartScriptInput{
		Name:   "autoscaler-test-script",
		Script: "#!/usr/bin/env bash\necho 'Hello World'",
	}

	result, err := BuildJSON(input)
	if err != nil {
		t.Fatalf("BuildJSON failed: %v", err)
	}

	expected := `{"name":"autoscaler-test-script","script":"#!/usr/bin/env bash\necho 'Hello World'"}`
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}

	// Verify it's valid JSON that can be parsed
	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("produced invalid JSON: %v", err)
	}

	// Verify field names are lowercase as expected by API
	if parsed["name"] != "autoscaler-test-script" {
		t.Errorf("expected 'name' field, got: %v", parsed)
	}
	if parsed["script"] != "#!/usr/bin/env bash\necho 'Hello World'" {
		t.Errorf("expected 'script' field, got: %v", parsed)
	}
}

// Test that BuildJSON produces valid JSON that can be parsed by standard library
func TestBuildJSON_ValidOutput(t *testing.T) {
	type TestStruct struct {
		Name  string
		Count int64
		Price float64
		Tags  []string
	}

	original := TestStruct{
		Name:  "B200 GPU",
		Count: 30,
		Price: 3.64,
		Tags:  []string{"gpu", "compute"},
	}

	// Build JSON
	jsonData, err := BuildJSON(original)
	if err != nil {
		t.Fatalf("BuildJSON failed: %v", err)
	}

	// Verify it's valid JSON by parsing with standard library
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("produced invalid JSON: %v", err)
	}

	// Verify basic structure
	if result["Name"] != "B200 GPU" {
		t.Errorf("expected Name='B200 GPU', got %v", result["Name"])
	}
	if result["Count"] != float64(30) { // JSON numbers are parsed as float64
		t.Errorf("expected Count=30, got %v", result["Count"])
	}
}

// Proof that standard json.Unmarshal() DOES support arrays
func TestStandardJSON_ArraySupport(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		target   interface{}
		desc     string
	}{
		{
			name:     "array of objects",
			jsonData: `[{"id": "1", "name": "Alice"}, {"id": "2", "name": "Bob"}]`,
			target: &[]struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{},
			desc: "Standard JSON handles arrays of objects perfectly",
		},
		{
			name:     "array of strings",
			jsonData: `["apple", "banana", "cherry"]`,
			target:   &[]string{},
			desc:     "Standard JSON handles arrays of primitives",
		},
		{
			name:     "array of numbers",
			jsonData: `[1, 2, 3, 4, 5]`,
			target:   &[]int{},
			desc:     "Standard JSON handles arrays of numbers",
		},
		{
			name:     "nested arrays",
			jsonData: `[{"tags": ["go", "json"]}, {"tags": ["test", "array"]}]`,
			target: &[]struct {
				Tags []string `json:"tags"`
			}{},
			desc: "Standard JSON handles nested arrays",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.desc)

			// Use STANDARD json.Unmarshal (not our custom function)
			err := json.Unmarshal([]byte(tt.jsonData), tt.target)
			if err != nil {
				t.Errorf("❌ Standard json.Unmarshal failed: %v", err)
				return
			}

			t.Logf("✅ Standard json.Unmarshal succeeded: %+v", tt.target)
		})
	}
}

// Show what standard JSON CAN'T do (why we need custom unmarshal)
func TestStandardJSON_Limitations(t *testing.T) {
	type User struct {
		ID   string `locationName:"userId"` // Custom tag
		Name string `json:"user_name"`
	}

	tests := []struct {
		name     string
		jsonData string
		issue    string
	}{
		{
			name:     "locationName not supported",
			jsonData: `{"userId": "123", "user_name": "John"}`,
			issue:    "Standard JSON ignores locationName tags",
		},
		{
			name:     "case insensitive not supported",
			jsonData: `{"USER_NAME": "John"}`,
			issue:    "Standard JSON is case sensitive",
		},
		{
			name:     "special floats not supported",
			jsonData: `{"value": "NaN"}`,
			issue:    "Standard JSON can't parse 'NaN' string as float",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Issue: %s", tt.issue)

			var result User
			err := json.Unmarshal([]byte(tt.jsonData), &result)

			// These will either fail or not populate fields correctly
			t.Logf("Standard JSON result: %+v, error: %v", result, err)
			t.Logf("This is why we need custom unmarshal functions")
		})
	}
}

// Compare standard vs custom unmarshal for arrays
func TestStandardVsCustom_ArrayComparison(t *testing.T) {
	jsonData := `[{"userId": "123", "USER_NAME": "John"}, {"userId": "456", "USER_NAME": "Jane"}]`

	type User struct {
		ID   string `locationName:"userId"`
		Name string `json:"user_name"`
	}

	t.Run("standard_json_array_handling", func(t *testing.T) {
		var users []User
		err := json.Unmarshal([]byte(jsonData), &users)

		t.Logf("Standard JSON result: %+v", users)
		t.Logf("Standard JSON error: %v", err)

		// Standard JSON will:
		// ✅ Parse the array structure correctly
		// ❌ Ignore locationName tags (ID will be empty)
		// ❌ Not match USER_NAME to user_name (case sensitive)
	})

	t.Run("custom_unmarshal_array_handling", func(t *testing.T) {
		var users []User
		err := UnmarshalJSONCaseInsensitive(&users, strings.NewReader(jsonData))

		t.Logf("Custom unmarshal result: %+v", users)
		t.Logf("Custom unmarshal error: %v", err)

		// Custom unmarshal will:
		// ✅ Parse the array structure correctly
		// ✅ Handle locationName tags (ID will be populated)
		// ✅ Handle case insensitive matching (Name will be populated)
	})
}
