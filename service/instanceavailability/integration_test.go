package instanceavailability_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instanceavailability"
)

// setupIntegrationTest sets up a real client for integration testing
func setupIntegrationTest(t *testing.T) *instanceavailability.InstanceAvailability {
	t.Helper()

	// Use session for proper credential handling
	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return instanceavailability.New(sess)
}

func TestListInstanceAvailability_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	instanceAvailability, err := svc.ListInstanceAvailability()
	if err != nil {
		t.Fatalf("failed to list instance availability: %v", err)
	}

	t.Logf("Found %d instance availability(s):\n", len(instanceAvailability))
	for _, ia := range instanceAvailability {
		t.Logf("  - %s: %s\n", ia.LocationCode, ia.Availabilities)
	}

	// validate response structure
	for _, ia := range instanceAvailability {
		if ia.LocationCode == "" {
			t.Errorf("instance availability %s has empty location code", ia.LocationCode)
		}
		if ia.Availabilities == nil {
			t.Errorf("instance availability %s has empty availabilities", ia.LocationCode)
		}
		if len(ia.Availabilities) == 0 {
			t.Errorf("instance availability %s has no availabilities", ia.LocationCode)
		}
	}
}
