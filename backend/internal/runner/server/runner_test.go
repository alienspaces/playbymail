package runner

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func init() {
	// Set timezone to UTC for all tests
	time.Local = time.UTC
}

func newDefaultDependencies(t *testing.T) (*log.Log, *store.Store, *river.Client[pgx.Tx]) {
	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	return l, s, j
}

func newTestHarness(t *testing.T) *harness.Testing {

	// Default test harness data configuration
	config := harness.DefaultDataConfig

	// Default dependencies
	l, s, j := newDefaultDependencies(t)

	// Test harness
	h, err := harness.NewTesting(l, s, j, config)
	require.NoError(t, err, "NewTesting returns without error")

	// For handler tests the test harness needs to commit data as the handler
	// creates a new database transaction.
	h.ShouldCommitData = true

	// Ensure proper cleanup before creating new data
	err = h.Teardown()
	require.NoError(t, err, "Teardown returns without error")

	return h
}

func TestNewRunner(t *testing.T) {

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.Default(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	r, err := NewRunner(l, s, j, cfg)
	require.NoError(t, err, "NewRunner returns without error")

	err = r.Init(s)
	require.NoError(t, err, "Init returns without error")
}

func TestMergeHandlerConfigs(t *testing.T) {
	// Test case 1: Merge with nil first map
	hc2 := map[string]server.HandlerConfig{
		"test1": {Name: "test1", Method: "GET", Path: "/test1"},
		"test2": {Name: "test2", Method: "POST", Path: "/test2"},
	}

	result := mergeHandlerConfigs(nil, hc2)

	if len(result) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(result))
	}

	if result["test1"].Name != "test1" {
		t.Errorf("Expected test1, got %s", result["test1"].Name)
	}

	if result["test2"].Name != "test2" {
		t.Errorf("Expected test2, got %s", result["test2"].Name)
	}

	// Test case 2: Merge with existing map
	hc1 := map[string]server.HandlerConfig{
		"existing": {Name: "existing", Method: "GET", Path: "/existing"},
	}

	result = mergeHandlerConfigs(hc1, hc2)

	if len(result) != 3 {
		t.Errorf("Expected 3 handlers, got %d", len(result))
	}

	if result["existing"].Name != "existing" {
		t.Errorf("Expected existing, got %s", result["existing"].Name)
	}

	// Test case 3: Override existing key
	hc1Fresh := map[string]server.HandlerConfig{
		"existing": {Name: "existing", Method: "GET", Path: "/existing"},
	}
	hc3 := map[string]server.HandlerConfig{
		"existing": {Name: "overridden", Method: "PUT", Path: "/overridden"},
	}

	result = mergeHandlerConfigs(hc1Fresh, hc3)

	if len(result) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(result))
	}

	if result["existing"].Name != "overridden" {
		t.Errorf("Expected overridden, got %s", result["existing"].Name)
	}
}
