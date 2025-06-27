package nulltime

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
)

func TestFromTime(t *testing.T) {

	tests := map[string]struct {
		time          time.Time
		expectedValid bool
	}{
		"valid time": {
			time:          time.Now(),
			expectedValid: true,
		},
		"invalid time": {
			time:          time.Time{},
			expectedValid: false,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := FromTime(tc.time)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.time, converted.Time)
		})
	}
}

func TestToTimePtr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullTime
		expected *time.Time
	}{
		"valid": {
			ns: sql.NullTime{
				Time:  time.Now().Round(time.Duration(1) * time.Hour),
				Valid: true,
			},
			expected: convert.Ptr(time.Now().Round(time.Duration(1) * time.Hour)),
		},
		"invalid": {
			ns: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := ToTimePtr(tc.ns)

			require.Equal(t, tc.expected, converted)
		})
	}
}
