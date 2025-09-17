//go:build integration
// +build integration

package instance_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/pkg/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/instance"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/sshkeys"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/startscripts"
)

// setupIntegrationTest sets up a real client for integration testing
func setupIntegrationTest(t *testing.T) *instance.Instance {
	t.Helper()

	// Use session for proper credential handling
	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return instance.New(sess)
}

func TestCreateAndDeleteInstance_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	instanceClient := instance.New(sess)
	sshKeysClient := sshkeys.New(sess)
	startScriptsClient := startscripts.New(sess)

	// For integration tests, we'll use existing SSH keys from the account
	// Rather than creating new ones each time

	// get all ssh keys
	sshKeyID, err := sshKeysClient.CreateSSHKey(&sshkeys.CreateSSHKeyInput{
		Name: "integration-test-ssh-key",
		Key:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC81234567890",
	})
	if err != nil {
		t.Fatalf("failed to list ssh keys: %v", err)
	}

	// create a script
	scriptID, err := startScriptsClient.CreateStartScript(&startscripts.CreateStartScriptInput{
		Name:   "integration-test-script",
		Script: "#!/bin/bash\n\necho hello world",
	})
	if err != nil {
		t.Fatalf("failed to create script: %v", err)
	}

	// Create instance
	input := &instance.CreateInstanceInput{
		InstanceType:    "1A100.22V",
		Image:           "ubuntu-22.04-cuda-12.3",
		SSHKeyIDs:       []string{sshKeyID}, // Empty for test - will use account default
		LocationCode:    "FIN-01",
		StartupScriptID: scriptID,
		IsSpot:          true, // Cheaper option for testing
		Contract:        "PAY_AS_YOU_GO",
		Pricing:         "DYNAMIC_PRICE",
		Description:     "Integration test instance - safe to delete",
		Hostname:        "integration-test-vm",
		OSVolume: &instance.OSVolume{
			Name: "integration-test-os-volume",
			Size: 50, // Smaller size for testing
		},
	}

	instanceID, err := instanceClient.CreateInstance(input)
	if err != nil {
		t.Fatalf("failed to create instance: %v", err)
	}

	if instanceID == "" {
		t.Fatal("created instance has empty ID")
	}

	t.Logf("Created instance with ID: %s", instanceID)

	// Cleanup: Delete the instance
	defer func() {
		t.Log("Cleaning up test instance...")
		err := instanceClient.PerformInstanceAction(&instance.InstanceActionInput{
			Action: instance.InstanceActionDelete,
			ID:     instanceID,
		})
		if err != nil {
			t.Errorf("failed to delete test instance %s: %v", instanceID, err)
		} else {
			t.Log("Successfully cleaned up test instance")
		}

		// delete ssh key
		err = sshKeysClient.DeleteSSHKey(sshKeyID)
		if err != nil {
			t.Errorf("failed to delete test ssh key %s: %v", sshKeyID, err)
		} else {
			t.Log("Successfully cleaned up test ssh key")
		}
	}()

	// Wait a moment for instance to be created
	time.Sleep(5 * time.Second)

	// Verify instance appears in list
	instances, err := instanceClient.ListInstances(nil)
	if err != nil {
		t.Fatalf("failed to list instances: %v", err)
	}

	found := false
	for _, inst := range instances {
		if inst.ID == instanceID {
			found = true
			if inst.Status == "" {
				t.Error("instance status is empty")
			}
			break
		}
	}

	if !found {
		t.Errorf("created instance %s not found in list", instanceID)
	}
}

func TestListInstances_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	instances, err := svc.ListInstances(nil)
	if err != nil {
		t.Fatalf("failed to list instances: %v", err)
	}

	t.Logf("Found %d instances", len(instances))

	// Validate response structure
	for i, inst := range instances {
		if inst.ID == "" {
			t.Errorf("instance %d has empty ID", i)
		}
		if inst.InstanceType == "" {
			t.Errorf("instance %d has empty InstanceType", i)
		}
		if inst.Status == "" {
			t.Errorf("instance %d has empty Status", i)
		}

		t.Logf("Instance %d: ID=%s, Type=%s, Status=%s, Location=%s",
			i, inst.ID, inst.InstanceType, inst.Status, inst.Location)
	}
}

// TestRateLimiting tests how the service handles rate limiting
func TestRateLimiting_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	// Make multiple rapid requests to test rate limiting
	const numRequests = 10
	errors := make([]error, numRequests)

	for i := 0; i < numRequests; i++ {
		_, errors[i] = svc.ListInstances(nil)
		if i < numRequests-1 {
			time.Sleep(100 * time.Millisecond) // Small delay
		}
	}

	errorCount := 0
	for _, err := range errors {
		if err != nil {
			errorCount++
			t.Logf("Request error (expected for rate limiting): %v", err)
		}
	}

	// We expect some requests might fail due to rate limiting
	if errorCount == numRequests {
		t.Error("all requests failed, this might indicate a bigger issue")
	}

	t.Logf("Rate limiting test: %d/%d requests succeeded", numRequests-errorCount, numRequests)
}
