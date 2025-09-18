package sshkeys_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/dcerr"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/sshkeys"
)

func setupIntegrationTest(t *testing.T) *sshkeys.SSHKey {
	t.Helper()

	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return sshkeys.New(sess)
}

func TestCreateSSHKey_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	sshKeyID, err := svc.CreateSSHKey(&sshkeys.CreateSSHKeyInput{
		Name: "integration-test-ssh-key",
		Key:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC81234567890",
	})
	if err != nil {
		t.Fatalf("failed to create ssh key: %v", err)
	}

	t.Logf("Created ssh key with ID: %s", sshKeyID)

	// cleanup
	defer func() {
		t.Log("Cleaning up test ssh key...")
		err := svc.DeleteSSHKey(sshKeyID)
		if err != nil {
			t.Errorf("failed to delete test ssh key %s: %v", sshKeyID, err)
		} else {
			t.Log("Successfully cleaned up test ssh key")
		}
	}()
}

func TestListSSHKeys_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	// create a ssh key
	sshKeyID, err := svc.CreateSSHKey(&sshkeys.CreateSSHKeyInput{
		Name: "integration-test-ssh-key",
		Key:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC81234567890",
	})
	if err != nil {
		t.Fatalf("failed to create ssh key: %v", err)
	}

	sshKeys, err := svc.ListSSHKeys()
	if err != nil {
		t.Fatalf("failed to list ssh keys: %v", err)
	}

	// look for sshKeyID in the list
	var found bool
	for _, sshKey := range sshKeys {
		if sshKey.ID == sshKeyID {
			t.Logf("Found ssh key with ID: %s", sshKeyID)
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("ssh key with ID %s not found", sshKeyID)
	}

	t.Logf("Found ssh key with ID: %s", sshKeyID)

	// cleanup
	defer func() {
		t.Log("Cleaning up test ssh key...")
		err := svc.DeleteSSHKey(sshKeyID)
		if err != nil {
			t.Errorf("failed to delete test ssh key %s: %v", sshKeyID, err)
		} else {
			t.Log("Successfully cleaned up test ssh key")
		}
	}()
}

func TestGetSSHKeyByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	sshKeyID, err := svc.CreateSSHKey(&sshkeys.CreateSSHKeyInput{
		Name: "integration-test-ssh-key",
		Key:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC81234567890",
	})
	if err != nil {
		t.Fatalf("failed to create ssh key: %v", err)
	}

	sshKeys, err := svc.GetSSHKey(sshKeyID)
	if err != nil {
		t.Fatalf("failed to get ssh key: %v", err)
	}

	if len(sshKeys) == 0 {
		t.Fatalf("ssh key with ID %s not found", sshKeyID)
	}

	if sshKeys[0].ID != sshKeyID {
		t.Fatalf("ssh key with ID %s not found", sshKeyID)
	}

	t.Logf("Found ssh key with ID: %s", sshKeyID)

	// cleanup
	defer func() {
		t.Log("Cleaning up test ssh key...")
		err := svc.DeleteSSHKey(sshKeyID)
		if err != nil {
			t.Errorf("failed to delete test ssh key %s: %v", sshKeyID, err)
		} else {
			t.Log("Successfully cleaned up test ssh key")
		}
	}()
}

func TestDeleteSSHKeyByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	sshKeyID, err := svc.CreateSSHKey(&sshkeys.CreateSSHKeyInput{
		Name: "integration-test-ssh-key",
		Key:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC81234567890",
	})
	if err != nil {
		t.Fatalf("failed to create ssh key: %v", err)
	}

	err = svc.DeleteSSHKey(sshKeyID)
	if err != nil {
		t.Fatalf("failed to delete ssh key: %v", err)
	}

	t.Logf("Deleted ssh key with ID: %s", sshKeyID)

}

func TestDeleteSSHKeys_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	sshKeyID1, err := svc.CreateSSHKey(&sshkeys.CreateSSHKeyInput{
		Name: "integration-test-ssh-key",
		Key:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC81234567890",
	})
	t.Logf("Created ssh key with ID: %s", sshKeyID1)
	if err != nil {
		t.Fatalf("failed to create ssh key: %v", err)
	}
	// create another ssh key with the same name
	sshKeyID2, err := svc.CreateSSHKey(&sshkeys.CreateSSHKeyInput{
		Name: "integration-test-ssh-key",
		Key:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC81234567890",
	})
	t.Logf("Created ssh key with ID: %s", sshKeyID2)
	if err != nil {
		t.Fatalf("failed to create ssh key: %v", err)
	}

	err = svc.DeleteSSHKeys(&sshkeys.DeleteSSHKeysInput{
		Keys: []string{sshKeyID1, sshKeyID2},
	})
	if err != nil {
		t.Fatalf("failed to delete ssh key: %v", err)
	}

	// find by sshKeyID1 - should fail with 404 since it was deleted
	_, err = svc.GetSSHKey(sshKeyID1)
	if err != nil {
		// This is expected - check if it's a 404 HTTP error
		if httpErr, ok := dcerr.IsHTTPError(err); ok {
			t.Logf("Expected HTTP error after deletion: %v", httpErr)
			// check status == 404 (Not Found)
			if httpErr.StatusCode == 404 {
				t.Logf("✅ Correctly got 404 Not Found for deleted SSH key")
			} else {
				t.Errorf("Expected HTTP 404, got %d", httpErr.StatusCode)
			}
		} else {
			t.Errorf("Expected HTTP error, got: %v", err)
		}
	} else {
		t.Fatalf("Expected error when getting deleted ssh key %s, but got success", sshKeyID1)
	}

	// find by sshKeyID2 - should also fail with 404 since it was deleted
	_, err = svc.GetSSHKey(sshKeyID2)
	if err != nil {
		// This is expected - check if it's a 404 HTTP error
		if httpErr, ok := dcerr.IsHTTPError(err); ok {
			t.Logf("Expected HTTP error after deletion: %v", httpErr)
			// check status == 404 (Not Found)
			if httpErr.StatusCode == 404 {
				t.Logf("✅ Correctly got 404 Not Found for deleted SSH key")
			} else {
				t.Errorf("Expected HTTP 404, got %d", httpErr.StatusCode)
			}
		} else {
			t.Errorf("Expected HTTP error, got: %v", err)
		}
	} else {
		t.Fatalf("Expected error when getting deleted ssh key %s, but got success", sshKeyID2)
	}

	t.Logf("Deleted ssh keys with IDs: %s, %s", sshKeyID1, sshKeyID2)

}
