package volumetypes_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/volumetypes"
)

func setupIntegrationTest(t *testing.T) *volumetypes.VolumeTypes {
	t.Helper()

	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return volumetypes.New(sess)
}

func TestListVolumeTypes_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	volumeTypes, err := svc.ListVolumeTypes()
	if err != nil {
		t.Fatalf("failed to list volume types: %v", err)
	}

	t.Logf("Found %d volume types", len(volumeTypes))
}
