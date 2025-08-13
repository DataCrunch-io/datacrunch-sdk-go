package jsonutil

import (
	"strings"
	"testing"
	"time"
)

// Test to see if we need timestampFormat support
func TestUnmarshalJSON_TimeHandling(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		target      interface{}
		description string
		shouldWork  bool
	}{
		{
			name:        "standard ISO8601 time",
			jsonData:    `{"created_at": "2023-12-01T10:30:45Z"}`,
			description: "Most APIs use ISO8601 - should work with standard json.Unmarshal",
			target: &struct {
				CreatedAt time.Time `json:"created_at"`
			}{},
			shouldWork: true,
		},
		{
			name:        "RFC3339 time",
			jsonData:    `{"updated_at": "2023-12-01T10:30:45+02:00"}`,
			description: "RFC3339 with timezone - should work with standard json.Unmarshal",
			target: &struct {
				UpdatedAt time.Time `json:"updated_at"`
			}{},
			shouldWork: true,
		},
		{
			name:        "Unix timestamp as number",
			jsonData:    `{"timestamp": 1701426645}`,
			description: "Unix timestamp as number - standard json.Unmarshal won't work",
			target: &struct {
				Timestamp time.Time `json:"timestamp"`
			}{},
			shouldWork: false, // This would require custom handling
		},
		{
			name:        "Unix timestamp as string",
			jsonData:    `{"timestamp": "1701426645"}`,
			description: "Unix timestamp as string - would need timestampFormat support",
			target: &struct {
				Timestamp time.Time `json:"timestamp"` // Would need: timestampFormat:"unix"
			}{},
			shouldWork: false, // This would require custom handling
		},
		{
			name:        "Custom date format",
			jsonData:    `{"date": "2023-12-01 10:30:45"}`,
			description: "Custom date format - would need timestampFormat support",
			target: &struct {
				Date time.Time `json:"date"` // Would need: timestampFormat:"2006-01-02 15:04:05"
			}{},
			shouldWork: false, // This would require custom handling
		},
		{
			name:        "RFC822 format",
			jsonData:    `{"date": "Mon, 02 Jan 2006 15:04:05 MST"}`,
			description: "RFC822 format - would need timestampFormat support",
			target: &struct {
				Date time.Time `json:"date"` // Would need: timestampFormat:"rfc822"
			}{},
			shouldWork: false, // This would require custom handling
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			err := UnmarshalJSON(tt.target, strings.NewReader(tt.jsonData))

			if tt.shouldWork {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				} else {
					t.Logf("✅ Standard JSON handling works: %+v", tt.target)
				}
			} else {
				if err != nil {
					t.Logf("⚠️  Expected failure (needs timestampFormat support): %v", err)
				} else {
					t.Logf("🤔 Unexpected success: %+v", tt.target)
				}
			}
		})
	}
}

// Test to show what APIs actually return for time fields
func TestUnmarshalJSON_RealWorldTimeFormats(t *testing.T) {
	// Common time formats from real APIs
	realWorldExamples := []struct {
		api      string
		format   string
		jsonData string
	}{
		{
			api:      "GitHub API",
			format:   "ISO8601",
			jsonData: `{"created_at": "2023-12-01T10:30:45Z"}`,
		},
		{
			api:      "Stripe API",
			format:   "Unix timestamp (number)",
			jsonData: `{"created": 1701426645}`,
		},
		{
			api:      "Twitter API",
			format:   "Custom format",
			jsonData: `{"created_at": "Wed Dec 01 10:30:45 +0000 2023"}`,
		},
		{
			api:      "AWS API",
			format:   "ISO8601",
			jsonData: `{"CreationTime": "2023-12-01T10:30:45.123Z"}`,
		},
	}

	for _, example := range realWorldExamples {
		t.Run(example.api, func(t *testing.T) {
			t.Logf("API: %s, Format: %s", example.api, example.format)
			t.Logf("Example: %s", example.jsonData)

			// Most APIs (GitHub, AWS) use ISO8601 which works with standard json.Unmarshal
			// Only some APIs (Stripe, Twitter) need special handling
		})
	}

	// Conclusion for DataCrunch SDK
	t.Run("recommendation", func(t *testing.T) {
		t.Log(`
RECOMMENDATION FOR DATACRUNCH SDK:

✅ Start WITHOUT timestampFormat support:
- Most modern APIs use ISO8601/RFC3339 
- Standard json.Unmarshal handles these automatically
- Less code complexity

⚠️  Add timestampFormat support ONLY if you encounter:
- Unix timestamps as strings
- Custom date formats
- Legacy time formats

📝 If needed later, add minimal support for common cases:
- timestampFormat:"unix" (for Unix timestamps)
- timestampFormat:"iso8601" (explicit)
- timestampFormat:"rfc3339" (explicit)
		`)
	})
}

// Show how to handle Unix timestamps without timestampFormat (alternative approach)
func TestUnmarshalJSON_UnixTimestampWorkaround(t *testing.T) {
	t.Log("Alternative to timestampFormat: Use custom types with UnmarshalJSON methods")
	
	// Example of custom type approach instead of timestampFormat tag
	t.Log(`
Custom type approach:

type UnixTime struct {
    time.Time
}

func (ut *UnixTime) UnmarshalJSON(data []byte) error {
    // Custom unmarshaling for Unix timestamps
    return nil
}

type APIResponse struct {
    ID        string   json:"id"
    Timestamp UnixTime json:"timestamp" // Custom type instead of timestampFormat tag
}
	`)
}