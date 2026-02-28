package turnsheet_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// Common test utilities for adventure game turn sheet tests.
// This file provides shared functions and test scaffolding that can be used
// across all adventure game turn sheet test files (join game, location choice,
// combat, inventory, etc.). As new turn sheet types are added, common patterns
// should be extracted here to maintain consistency and reduce duplication.

// Test IDs for generating deterministic turn sheet codes
const (
	testGameID             = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	testAccountID          = "c3d4e5f6-a7b8-9012-cdef-123456789012"
	testGameTurnSheetID    = "d4e5f6a7-b8c9-0123-def0-234567890123"
	testGameSubscriptionID = "e5f6a7b8-c9d0-1234-ef01-345678901234"
)

// loadTestBackgroundImage loads a test background image and converts it to a base64 data URL.
// This is used by all adventure game turn sheet tests since background images are standard.
func loadTestBackgroundImage(t *testing.T, filename string) string {
	t.Helper()

	imageData, err := os.ReadFile(filename)
	require.NoError(t, err, "should load background image")

	// Detect MIME type
	mimeType := http.DetectContentType(imageData)

	// Create base64 data URL
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

	return dataURL
}

// generateTestJoinTurnSheetCode generates a realistic join game turn sheet code for testing.
// It uses the actual GenerateJoinTurnSheetCode function with deterministic test IDs.
func generateTestJoinTurnSheetCode(t *testing.T) string {
	t.Helper()

	code, err := turnsheetutil.GenerateJoinGameTurnSheetCode(testGameSubscriptionID)
	require.NoError(t, err, "should generate join turn sheet code")

	return code
}

// generateTestTurnSheetCode generates a realistic playing game turn sheet code for testing.
// It uses the actual GenerateTurnSheetCode function with deterministic test IDs.
// Note: This uses a placeholder game_subscription_instance_id since we don't have a real instance in unit tests
func generateTestTurnSheetCode(t *testing.T) string {
	t.Helper()

	code, err := turnsheetutil.GeneratePlayGameTurnSheetCode(testGameTurnSheetID)
	require.NoError(t, err, "should generate turn sheet code")

	return code
}
