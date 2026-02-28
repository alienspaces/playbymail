package turnsheet_test

import (
	"os"
	"testing"
)

func requireOpenAIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("OPENAI_API_KEY") == "" {
		ciHint := ""
		if os.Getenv("CI") != "" {
			ciHint = " In CI, ensure OPENAI_API_KEY is not marked as 'Protected' or the branch is protected."
		}
		t.Fatalf("OPENAI_API_KEY must be set to run scanner integration tests.%s", ciHint)
	}
}
