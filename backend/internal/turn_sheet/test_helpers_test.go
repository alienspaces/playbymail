package turn_sheet_test

import (
	"os"
	"testing"
)

func requireOpenAIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Fatalf("OPENAI_API_KEY must be set to run scanner integration tests")
	}
}
