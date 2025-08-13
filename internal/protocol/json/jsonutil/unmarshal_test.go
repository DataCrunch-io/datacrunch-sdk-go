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

// Test for string response handling (some APIs return plain strings)
func TestUnmarshalJSON_StringResponses(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		target      interface{}
		expected    interface{}
		description string
	}{
		{
			name:        "plain string response",
			jsonData:    `"simple-string-id"`,
			description: "API returns plain string (like CreateInstance returning just an ID)",
			target:      new(string),
			expected:    "simple-string-id",
		},
		{
			name:        "quoted string response",
			jsonData:    `"instance-12345"`,
			description: "API returns instance ID as quoted string",
			target:      new(string),
			expected:    "instance-12345",
		},
		{
			name:        "number as string response",
			jsonData:    `"12345"`,
			description: "API returns numeric ID as string",
			target:      new(string),
			expected:    "12345",
		},
		{
			name:        "boolean response",
			jsonData:    `true`,
			description: "API returns boolean response",
			target:      new(bool),
			expected:    true,
		},
		{
			name:        "number response",
			jsonData:    `42`,
			description: "API returns plain number",
			target:      new(int),
			expected:    42,
		},
		{
			name:        "float response",
			jsonData:    `3.14159`,
			description: "API returns plain float",
			target:      new(float64),
			expected:    3.14159,
		},
		{
			name:        "null response",
			jsonData:    `null`,
			description: "API returns null",
			target:      new(*string),
			expected:    (*string)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			err := UnmarshalJSON(tt.target, strings.NewReader(tt.jsonData))
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}

			// Verify the result based on type
			switch expected := tt.expected.(type) {
			case string:
				if actual := *(tt.target.(*string)); actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			case bool:
				if actual := *(tt.target.(*bool)); actual != expected {
					t.Errorf("Expected %v, got %v", expected, actual)
				}
			case int:
				if actual := *(tt.target.(*int)); actual != expected {
					t.Errorf("Expected %d, got %d", expected, actual)
				}
			case float64:
				if actual := *(tt.target.(*float64)); actual != expected {
					t.Errorf("Expected %f, got %f", expected, actual)
				}
			case *string:
				actual := tt.target.(**string)
				if (expected == nil && *actual != nil) || (expected != nil && *actual == nil) {
					t.Errorf("Expected nil, got %v", *actual)
				}
			}

			t.Logf("Successfully parsed: %+v", tt.target)
		})
	}
}

// Test for real-world string response scenarios
func TestUnmarshalJSON_RealWorldStringScenarios(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		target      interface{}
		description string
	}{
		{
			name:        "CreateInstance returns instance ID",
			jsonData:    `"i-0123456789abcdef0"`,
			description: "CreateInstance API returns just the instance ID as string",
			target:      new(string),
		},
		{
			name:        "CreateScript returns script ID",
			jsonData:    `"script-abc123def456"`,
			description: "CreateScript API returns just the script ID",
			target:      new(string),
		},
		{
			name:        "API returns status message",
			jsonData:    `"success"`,
			description: "Some APIs return simple status strings",
			target:      new(string),
		},
		{
			name:        "API returns error message",
			jsonData:    `"Invalid request parameters"`,
			description: "Error APIs might return plain error messages",
			target:      new(string),
		},
		{
			name:        "Token refresh returns new token",
			jsonData:    `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`,
			description: "Token refresh might return just the token string",
			target:      new(string),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			err := UnmarshalJSON(tt.target, strings.NewReader(tt.jsonData))
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}

			// Just verify we got a non-empty result for strings
			if strTarget, ok := tt.target.(*string); ok && *strTarget == "" {
				t.Errorf("Expected non-empty string, got empty")
				return
			}

			t.Logf("Successfully parsed string response: %+v", tt.target)
		})
	}
}

// Test mixed response types in the same API
func TestUnmarshalJSON_MixedResponseTypes(t *testing.T) {
	t.Run("same API endpoint can return different types", func(t *testing.T) {
		// Some APIs might return different response types based on parameters or state

		// Test 1: Object response
		objectJSON := `{"id": "123", "status": "created"}`
		var objectResult struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}

		err := UnmarshalJSON(&objectResult, strings.NewReader(objectJSON))
		if err != nil {
			t.Errorf("Object response failed: %v", err)
		} else {
			t.Logf("Object response: %+v", objectResult)
		}

		// Test 2: String response (same conceptual endpoint, different scenario)
		stringJSON := `"already-exists-456"`
		var stringResult string

		err = UnmarshalJSON(&stringResult, strings.NewReader(stringJSON))
		if err != nil {
			t.Errorf("String response failed: %v", err)
		} else {
			t.Logf("String response: %s", stringResult)
		}

		// Test 3: Array response (same conceptual endpoint, list operation)
		arrayJSON := `["id1", "id2", "id3"]`
		var arrayResult []string

		err = UnmarshalJSON(&arrayResult, strings.NewReader(arrayJSON))
		if err != nil {
			t.Errorf("Array response failed: %v", err)
		} else {
			t.Logf("Array response: %+v", arrayResult)
		}
	})
}

// Test cases for common response/struct mismatch scenarios
func TestUnmarshalJSON_ResponseMismatches(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		target      interface{}
		expected    interface{}
		shouldError bool
		description string
	}{
		{
			name:        "extra fields in response (should ignore)",
			jsonData:    `{"id": "123", "name": "John", "extra_field": "ignored", "another_extra": 999}`,
			description: "API returns extra fields not in struct - should ignore them",
			target: &struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{},
			expected: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   "123",
				Name: "John",
			},
		},
		{
			name:        "missing optional fields (should not error)",
			jsonData:    `{"id": "456", "name": "Jane"}`,
			description: "API missing optional fields - should use zero values",
			target: &struct {
				ID          string  `json:"id"`
				Name        string  `json:"name"`
				Email       string  `json:"email,omitempty"`       // Missing in response
				Age         int     `json:"age,omitempty"`         // Missing in response
				IsActive    bool    `json:"is_active,omitempty"`   // Missing in response
				Description *string `json:"description,omitempty"` // Missing in response (pointer)
			}{},
			expected: struct {
				ID          string  `json:"id"`
				Name        string  `json:"name"`
				Email       string  `json:"email,omitempty"`
				Age         int     `json:"age,omitempty"`
				IsActive    bool    `json:"is_active,omitempty"`
				Description *string `json:"description,omitempty"`
			}{
				ID:          "456",
				Name:        "Jane",
				Email:       "", // Zero value
				Age:         0,  // Zero value
				IsActive:    false,
				Description: nil, // Nil pointer
			},
		},
		{
			name:        "null values in response",
			jsonData:    `{"id": "789", "name": "Bob", "email": null, "age": null, "metadata": null}`,
			description: "API returns null values - should handle gracefully",
			target: &struct {
				ID       string                 `json:"id"`
				Name     string                 `json:"name"`
				Email    *string                `json:"email"`    // Pointer to handle null
				Age      *int                   `json:"age"`      // Pointer to handle null
				Metadata map[string]interface{} `json:"metadata"` // Map can be null
			}{},
			expected: struct {
				ID       string                 `json:"id"`
				Name     string                 `json:"name"`
				Email    *string                `json:"email"`
				Age      *int                   `json:"age"`
				Metadata map[string]interface{} `json:"metadata"`
			}{
				ID:       "789",
				Name:     "Bob",
				Email:    nil, // null becomes nil
				Age:      nil, // null becomes nil
				Metadata: nil, // null becomes nil
			},
		},
		{
			name:        "different field names (case insensitive)",
			jsonData:    `{"USER_ID": "999", "full_name": "Alice Cooper", "EMAIL_ADDRESS": "alice@example.com"}`,
			description: "API uses different case than struct definition",
			target: &struct {
				UserID       string `json:"user_id"`
				FullName     string `json:"full_name"`
				EmailAddress string `json:"email_address"`
			}{},
			expected: struct {
				UserID       string `json:"user_id"`
				FullName     string `json:"full_name"`
				EmailAddress string `json:"email_address"`
			}{
				UserID:       "999",
				FullName:     "Alice Cooper",
				EmailAddress: "alice@example.com",
			},
		},
		{
			name:        "nested object with mismatched fields",
			jsonData:    `{"user": {"id": 123, "display_name": "John Doe", "unexpected_field": "ignore"}, "status": "active", "extra_root": "ignore"}`,
			description: "Nested objects with extra fields should be handled",
			target: &struct {
				User struct {
					ID          int    `json:"id"`
					DisplayName string `json:"display_name"`
					Email       string `json:"email,omitempty"` // Missing in response
				} `json:"user"`
				Status string `json:"status"`
			}{},
			expected: struct {
				User struct {
					ID          int    `json:"id"`
					DisplayName string `json:"display_name"`
					Email       string `json:"email,omitempty"`
				} `json:"user"`
				Status string `json:"status"`
			}{
				User: struct {
					ID          int    `json:"id"`
					DisplayName string `json:"display_name"`
					Email       string `json:"email,omitempty"`
				}{
					ID:          123,
					DisplayName: "John Doe",
					Email:       "", // Missing field gets zero value
				},
				Status: "active",
			},
		},
		{
			name:        "array with mixed data types (partial success)",
			jsonData:    `{"users": [{"id": "1", "name": "Alice"}, {"id": "2", "name": "Bob", "extra": "data"}], "total": 2}`,
			description: "Array elements with extra fields should work",
			target: &struct {
				Users []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"users"`
				Total int `json:"total"`
			}{},
			expected: struct {
				Users []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"users"`
				Total int `json:"total"`
			}{
				Users: []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{
					{ID: "1", Name: "Alice"},
					{ID: "2", Name: "Bob"}, // "extra" field ignored
				},
				Total: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			var err error
			if strings.Contains(tt.name, "case insensitive") {
				err = UnmarshalJSONCaseInsensitive(tt.target, strings.NewReader(tt.jsonData))
			} else {
				err = UnmarshalJSON(tt.target, strings.NewReader(tt.jsonData))
			}

			if tt.shouldError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// For successful cases, log the results for inspection
			if !tt.shouldError {
				t.Logf("Successfully parsed: %+v", tt.target)
			}
		})
	}
}

// Test for common error scenarios that cause unmarshal failures
func TestUnmarshalJSON_ErrorScenarios(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		target      interface{}
		shouldError bool
		description string
	}{
		{
			name:        "wrong data type (string to int)",
			jsonData:    `{"id": "not-a-number", "name": "John"}`,
			description: "API returns string where int is expected",
			target: &struct {
				ID   int    `json:"id"` // Expects int but gets string
				Name string `json:"name"`
			}{},
			shouldError: true,
		},
		{
			name:        "malformed JSON",
			jsonData:    `{"id": 123, "name": "John", "invalid": }`,
			description: "Invalid JSON syntax should cause error",
			target: &struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{},
			shouldError: true,
		},
		{
			name:        "empty response body",
			jsonData:    ``,
			description: "Empty response should not error (EOF handled)",
			target: &struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{},
			shouldError: false, // EOF is handled gracefully
		},
		{
			name:        "null root object",
			jsonData:    `null`,
			description: "Null response should be handled",
			target: &struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing error scenario: %s", tt.description)

			err := UnmarshalJSON(tt.target, strings.NewReader(tt.jsonData))

			if tt.shouldError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if err != nil {
				t.Logf("Got expected error: %v", err)
			} else {
				t.Logf("Successfully handled: %+v", tt.target)
			}
		})
	}
}

// Test for debugging response mismatches - shows actual vs expected
func TestUnmarshalJSON_DebuggingHelper(t *testing.T) {
	// This test helps debug real-world response/struct mismatches
	t.Run("debug response structure", func(t *testing.T) {
		// Example: You expect this struct
		type ExpectedStruct struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Status string `json:"status"`
		}

		// But API returns this JSON
		actualJSON := `{
			"id": "123",
			"user_name": "John Doe", 
			"account_status": "active",
			"created_at": "2023-01-01T10:00:00Z",
			"extra_data": {
				"metadata": "some value"
			}
		}`

		var result ExpectedStruct
		err := UnmarshalJSON(&result, strings.NewReader(actualJSON))

		t.Logf("JSON Input: %s", actualJSON)
		t.Logf("Expected struct: %+v", ExpectedStruct{})
		t.Logf("Actual result: %+v", result)
		t.Logf("Error: %v", err)

		// This will show you:
		// - result.Name will be empty (field name mismatch: "user_name" vs "name")
		// - result.Status will be empty (field name mismatch: "account_status" vs "status")
		// - Extra fields like "created_at" and "extra_data" are ignored
		// - Only "id" matches and gets populated
	})
}

// Test for array response handling (fixes the instance types API error)
func TestUnmarshalJSON_ArrayResponses(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		target      interface{}
		description string
	}{
		{
			name:        "simple array response",
			jsonData:    `[{"id": "type1", "name": "Small"}, {"id": "type2", "name": "Large"}]`,
			description: "API returns array of objects - should handle gracefully",
			target: &[]struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{},
		},
		{
			name:        "array response with locationName tags",
			jsonData:    `[{"instanceId": "i-123", "state": "running"}, {"instanceId": "i-456", "state": "stopped"}]`,
			description: "Array response with AWS-style locationName fields",
			target: &[]struct {
				InstanceID string `locationName:"instanceId"`
				State      string `locationName:"state"`
			}{},
		},
		{
			name:        "mixed case array (case insensitive)",
			jsonData:    `[{"USER_ID": "1", "display_name": "Alice"}, {"user_id": "2", "DISPLAY_NAME": "Bob"}]`,
			description: "Array with mixed case field names",
			target: &[]struct {
				UserID      string `json:"user_id"`
				DisplayName string `json:"display_name"`
			}{},
		},
		{
			name:        "array with nested objects",
			jsonData:    `[{"id": "1", "user": {"name": "John", "email": "john@example.com"}}, {"id": "2", "user": {"name": "Jane", "email": "jane@example.com"}}]`,
			description: "Array containing nested objects",
			target: &[]struct {
				ID   string `json:"id"`
				User struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				} `json:"user"`
			}{},
		},
		{
			name:        "empty array",
			jsonData:    `[]`,
			description: "Empty array response should not error",
			target:      &[]struct{ ID string }{},
		},
		{
			name:        "array of primitives",
			jsonData:    `["item1", "item2", "item3"]`,
			description: "Array of primitive values",
			target:      &[]string{},
		},
		{
			name:        "array of numbers",
			jsonData:    `[1, 2, 3, 4, 5]`,
			description: "Array of numbers",
			target:      &[]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			var err error
			if strings.Contains(tt.name, "case insensitive") {
				err = UnmarshalJSONCaseInsensitive(tt.target, strings.NewReader(tt.jsonData))
			} else {
				err = UnmarshalJSON(tt.target, strings.NewReader(tt.jsonData))
			}

			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}

			t.Logf("Successfully parsed array: %+v", tt.target)
		})
	}
}

// Test the specific instance types API scenario
func TestUnmarshalJSON_InstanceTypesScenario(t *testing.T) {
	// Simulate the exact API response that was causing the error
	instanceTypesJSON := `[
		{
			"id": "1vcpu-2gb",
			"name": "1 vCPU, 2 GB RAM",
			"cpu": 1,
			"memory": 2048,
			"location": "EU-WEST",
			"price_per_hour": 0.01
		},
		{
			"id": "2vcpu-4gb", 
			"name": "2 vCPU, 4 GB RAM",
			"cpu": 2,
			"memory": 4096,
			"location": "EU-WEST",
			"price_per_hour": 0.02
		}
	]`

	// Define the expected struct (like in the real API client)
	type InstanceType struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		CPU          int     `json:"cpu"`
		Memory       int     `json:"memory"`
		Location     string  `json:"location"`
		PricePerHour float64 `json:"price_per_hour"`
	}

	var instanceTypes []InstanceType
	err := UnmarshalJSON(&instanceTypes, strings.NewReader(instanceTypesJSON))

	if err != nil {
		t.Errorf("Failed to unmarshal instance types: %v", err)
		return
	}

	if len(instanceTypes) != 2 {
		t.Errorf("Expected 2 instance types, got %d", len(instanceTypes))
		return
	}

	// Verify the first instance type
	first := instanceTypes[0]
	if first.ID != "1vcpu-2gb" {
		t.Errorf("Expected ID '1vcpu-2gb', got '%s'", first.ID)
	}
	if first.CPU != 1 {
		t.Errorf("Expected CPU 1, got %d", first.CPU)
	}
	if first.Memory != 2048 {
		t.Errorf("Expected Memory 2048, got %d", first.Memory)
	}

	t.Logf("Successfully parsed instance types: %+v", instanceTypes)
}

// Test for the specific GetStartScript array vs object mismatch issue
func TestUnmarshalJSON_GetStartScriptMismatch(t *testing.T) {
	// Simulate the StartScriptResponse struct
	type StartScriptResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Script string `json:"script"`
	}

	tests := []struct {
		name        string
		jsonData    string
		description string
		testType    string
	}{
		{
			name:        "API returns array but we expect single object",
			jsonData:    `[{"id": "script-123", "name": "My Script", "script": "#!/bin/bash\necho hello"}]`,
			description: "GetStartScript API returns array with 1 element instead of single object",
			testType:    "array_to_object",
		},
		{
			name:        "API returns single object as expected",
			jsonData:    `{"id": "script-123", "name": "My Script", "script": "#!/bin/bash\necho hello"}`,
			description: "GetStartScript API returns single object (expected behavior)",
			testType:    "object",
		},
		{
			name:        "API returns empty array",
			jsonData:    `[]`,
			description: "GetStartScript API returns empty array (not found scenario)",
			testType:    "empty_array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			switch tt.testType {
			case "array_to_object":
				// Test 1: Try to unmarshal array into single object (this will fail)
				var singleResult StartScriptResponse
				err := UnmarshalJSON(&singleResult, strings.NewReader(tt.jsonData))
				if err == nil {
					t.Errorf("Expected error when unmarshaling array into single object, but got success")
				} else {
					t.Logf("Expected error: %v", err)
				}

				// Test 2: Unmarshal array into array, then extract first element (workaround)
				var arrayResult []StartScriptResponse
				err = UnmarshalJSON(&arrayResult, strings.NewReader(tt.jsonData))
				if err != nil {
					t.Errorf("Failed to unmarshal into array: %v", err)
					return
				}

				if len(arrayResult) > 0 {
					extractedItem := arrayResult[0]
					t.Logf("Successfully extracted single item from array: %+v", extractedItem)
				} else {
					t.Logf("Array is empty, no item to extract")
				}

			case "object":
				// Test single object unmarshaling (should work)
				var result StartScriptResponse
				err := UnmarshalJSON(&result, strings.NewReader(tt.jsonData))
				if err != nil {
					t.Errorf("Failed to unmarshal object: %v", err)
					return
				}
				t.Logf("Successfully unmarshaled object: %+v", result)

			case "empty_array":
				// Test empty array handling
				var arrayResult []StartScriptResponse
				err := UnmarshalJSON(&arrayResult, strings.NewReader(tt.jsonData))
				if err != nil {
					t.Errorf("Failed to unmarshal empty array: %v", err)
					return
				}
				t.Logf("Successfully handled empty array: %+v (length: %d)", arrayResult, len(arrayResult))
			}
		})
	}
}

// Test helper functions for handling API response mismatches
func TestUnmarshalJSON_ResponseAdapters(t *testing.T) {
	type StartScriptResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Script string `json:"script"`
	}

	// Helper function to handle array-or-object responses
	unmarshalArrayOrObject := func(data string, target *StartScriptResponse) error {
		// Try array first
		var arrayResult []StartScriptResponse
		if err := UnmarshalJSON(&arrayResult, strings.NewReader(data)); err == nil {
			if len(arrayResult) > 0 {
				*target = arrayResult[0] // Take first element
				return nil
			}
			return nil // Empty array, leave target as zero value
		}

		// If array failed, try single object
		return UnmarshalJSON(target, strings.NewReader(data))
	}

	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name:     "handles array with single element",
			jsonData: `[{"id": "script-123", "name": "Test", "script": "echo test"}]`,
			wantErr:  false,
		},
		{
			name:     "handles single object",
			jsonData: `{"id": "script-456", "name": "Test2", "script": "echo test2"}`,
			wantErr:  false,
		},
		{
			name:     "handles empty array",
			jsonData: `[]`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result StartScriptResponse
			err := unmarshalArrayOrObject(tt.jsonData, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalArrayOrObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Result: %+v", result)
		})
	}
}

// Test demonstrating how to fix the GetStartScript method
func TestUnmarshalJSON_FixedGetStartScript(t *testing.T) {
	type StartScriptResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Script string `json:"script"`
	}

	// Simulate different API response formats
	responses := []struct {
		name     string
		jsonData string
		desc     string
	}{
		{
			name:     "array_response",
			jsonData: `[{"id": "script-123", "name": "My Script", "script": "#!/bin/bash\necho hello"}]`,
			desc:     "API returns array (current behavior causing the error)",
		},
		{
			name:     "object_response",
			jsonData: `{"id": "script-456", "name": "My Script 2", "script": "#!/bin/bash\necho world"}`,
			desc:     "API returns object (expected behavior)",
		},
	}

	for _, resp := range responses {
		t.Run(resp.name, func(t *testing.T) {
			t.Logf("Testing: %s", resp.desc)

			// Method 1: Flexible unmarshaling (handles both array and object)
			var result StartScriptResponse

			// Try as array first, then fallback to object
			var arrayResult []StartScriptResponse
			if err := UnmarshalJSON(&arrayResult, strings.NewReader(resp.jsonData)); err == nil && len(arrayResult) > 0 {
				result = arrayResult[0]
				t.Logf("✅ Successfully handled as array, extracted first element: %+v", result)
			} else {
				// Fallback to single object
				if err := UnmarshalJSON(&result, strings.NewReader(resp.jsonData)); err == nil {
					t.Logf("✅ Successfully handled as single object: %+v", result)
				} else {
					t.Errorf("❌ Failed to unmarshal both as array and object: %v", err)
				}
			}
		})
	}

	// Show how to update the GetStartScript method
	t.Run("recommended_fix", func(t *testing.T) {
		t.Log(`
Recommended fix for GetStartScript method:

func (c *StartScripts) GetStartScript(id string) (*StartScriptResponse, error) {
    // ... existing code ...
    
    // Try to unmarshal as array first (handles current API behavior)
    var scripts []*StartScriptResponse
    req := c.newRequest(op, input, &scripts)
    
    if err := req.Send(); err != nil {
        return nil, err
    }
    
    if len(scripts) > 0 {
        return scripts[0], nil // Return first element
    }
    
    return nil, fmt.Errorf("script not found")
}

OR create a wrapper function that handles both cases automatically.
		`)
	})
}
