package harness_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/harness"
)

func TestUniqueName_NormalName(t *testing.T) {
	input := harness.UniqueName("Test Normal Name")
	actual := harness.NormalName(input)
	require.Equal(t, "Test Normal Name", actual)
}

func TestUniqueEmail_NormalEmail(t *testing.T) {
	input := harness.UniqueEmail("test@example.com")
	actual := harness.NormalEmail(input)
	require.Equal(t, "test@example.com", actual)
}
