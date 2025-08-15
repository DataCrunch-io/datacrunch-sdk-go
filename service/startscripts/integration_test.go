package startscripts_test

import (
	"testing"
	"time"

	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/credentials"
	"github.com/datacrunch-io/datacrunch-sdk-go/datacrunch/session"
	"github.com/datacrunch-io/datacrunch-sdk-go/service/startscripts"
)

func setupIntegrationTest(t *testing.T) *startscripts.StartScripts {
	t.Helper()

	sess := session.New(
		session.WithCredentialsProvider(credentials.NewSharedCredentials("", "testing")),
		session.WithTimeout(30*time.Second),
		session.WithDebug(false),
	)

	return startscripts.New(sess)
}

func TestCreateStartScript_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	scriptID, err := svc.CreateStartScript(&startscripts.CreateStartScriptInput{
		Name:   "My startup scrip",
		Script: "#!/bin/bash\n\necho hello world",
	})

	if err != nil {
		t.Fatalf("failed to create start script: %v", err)
	}

	t.Logf("Created start script with ID: %s", scriptID)

	// Cleanup: Delete the start script
	defer func() {
		t.Log("Cleaning up test start script...")
		err := svc.DeleteStartScript(scriptID)
		if err != nil {
			t.Errorf("failed to delete test start script %s: %v", scriptID, err)
		} else {
			t.Log("Successfully cleaned up test start script")
		}
	}()
}

func TestListStartScripts_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	svc := setupIntegrationTest(t)

	// create a start script
	scriptID, err := svc.CreateStartScript(&startscripts.CreateStartScriptInput{
		Name:   "My startup scrip",
		Script: "#!/bin/bash\n\necho hello world",
	})
	if err != nil {
		t.Fatalf("failed to create start script: %v", err)
	}
	t.Logf("Created start script with ID: %s", scriptID)

	startScripts, err := svc.ListStartScripts()
	if err != nil {
		t.Fatalf("failed to list start scripts: %v", err)
	}
	var found bool
	// look for scriptId in the list
	for _, script := range startScripts {
		if script.ID == scriptID {
			found = true
			t.Logf("Found start script with ID: %s", scriptID)
			break
		}
	}

	if !found {
		t.Fatalf("start script with ID %s not found", scriptID)
	}
	t.Logf("Found start script with ID: %s", scriptID)

	// cleanup
	defer func() {
		t.Log("Cleaning up test start script...")
		err := svc.DeleteStartScript(scriptID)
		if err != nil {
			t.Errorf("failed to delete test start script %s: %v", scriptID, err)
		} else {
			t.Log("Successfully cleaned up test start script")
		}
	}()
}
