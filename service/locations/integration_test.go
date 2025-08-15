package locations_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/locations"
)

func setupIntegrationTest(t *testing.T) *locations.Locations {
	t.Helper()

	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return locations.New(sess)
}

func TestListLocations_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	locations, err := svc.ListLocations()
	if err != nil {
		t.Fatalf("failed to list locations: %v", err)
	}

	t.Logf("Found %d locations", len(locations))

	// validate response structure
	for _, loc := range locations {
		if loc.Code == "" {
			t.Errorf("location %s has empty code", loc.Code)
		}
		if loc.Name == "" {
			t.Errorf("location %s has empty name", loc.Code)
		}

		// Country code should be 2 characters
		if len(loc.CountryCode) != 2 {
			t.Errorf("location %s has invalid country code %s", loc.Code, loc.CountryCode)
		}
	}
}
