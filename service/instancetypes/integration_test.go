package instancetypes_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instancetypes"
)

func setupIntegrationTest(t *testing.T) *instancetypes.InstanceTypes {
	t.Helper()

	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return instancetypes.New(sess)
}

func TestListInstanceTypes_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	instanceTypes, err := svc.ListInstanceTypes()
	if err != nil {
		t.Fatalf("failed to list instance types: %v", err)
	}

	t.Logf("Found %d instance types", len(instanceTypes))
}
