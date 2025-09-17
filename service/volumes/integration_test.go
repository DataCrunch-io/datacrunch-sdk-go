package volumes_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/volumes"
)

func setupIntegrationTest(t *testing.T) *volumes.Volumes {
	t.Helper()

	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return volumes.New(sess)
}

func TestListVolumes_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	volumes, err := svc.ListVolumes(&volumes.ListVolumesStatus{
		Status: volumes.VolumeStatusOrdered,
	})
	if err != nil {
		t.Fatalf("failed to list volumes: %v", err)
	}

	t.Logf("Found %d volumes", len(volumes))

	// ignore volumes status and list all volumes
	volumes, err = svc.ListVolumes(nil)
	if err != nil {
		t.Fatalf("failed to list volumes: %v", err)
	}

	t.Logf("Found %d volumes", len(volumes))

}

func TestCreateVolume_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	// need to create an instance first

	volume, err := svc.CreateVolume(&volumes.CreateVolumeInput{
		Type:         "NVMe",
		LocationCode: "FIN-01",
		Size:         50,
		InstanceID:   "9ccf35fe-5f48-4d0d-927b-ed1014205cf3",
		InstanceIDs: []string{
			"da541b4f-6781-437c-a73c-5a4e115a3850",
			"5fba98e9-b680-4ca4-9e2b-6ee5bb5691f5",
		},
		Name: "my-volume",
	})

	if err != nil {
		t.Fatalf("failed to create volume: %v", err)
	}

	t.Logf("Created volume with ID: %s", volume)

	// cleanup

	defer func() {
		t.Log("Cleaning up test volume...")
		err := svc.DeleteVolume(volume, true)
		if err != nil {
			t.Errorf("failed to delete volume %s: %v", volume, err)
		} else {
			t.Log("Successfully cleaned up test volume")
		}
	}()
}
