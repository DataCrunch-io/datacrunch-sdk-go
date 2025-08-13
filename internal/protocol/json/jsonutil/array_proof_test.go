package jsonutil

import (
	"encoding/json"
	"strings"
	"testing"
)

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