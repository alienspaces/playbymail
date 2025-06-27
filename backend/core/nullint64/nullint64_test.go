package nullint64

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
)

func TestFromInt64(t *testing.T) {

	tests := map[string]struct {
		value         int64
		expectedValid bool
	}{
		"valid value": {
			value:         464,
			expectedValid: true,
		},
		"valid zero value": {
			value:         0,
			expectedValid: true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := FromInt64(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.Int64)
		})
	}
}

func TestFromInt64Ptr(t *testing.T) {

	tests := map[string]struct {
		value         *int64
		expectedValid bool
	}{
		"valid value": {
			value:         convert.Ptr(int64(654)),
			expectedValid: true,
		},
		"valid empty": {
			value:         convert.Ptr(int64(0)),
			expectedValid: true,
		},
		"invalid nil": {
			value:         nil,
			expectedValid: false,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := FromInt64Ptr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.Int64)
			}
		})
	}
}

func TestToInt64(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt64
		expected int64
		wantErr  bool
	}{
		"valid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: true,
			},
			expected: 654,
			wantErr:  false,
		},
		"invalid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: false,
			},
			expected: 0,
			wantErr:  true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted, err := ToInt64(tc.ns)
			if tc.wantErr {
				require.Error(t, err, "expect error")
				return
			}
			require.NoError(t, err, "expect no error")
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}

func TestToInt64Ptr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt64
		expected *int64
		wantErr  bool
	}{
		"valid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: true,
			},
			expected: convert.Ptr(int64(654)),
			wantErr:  false,
		},
		"invalid": {
			ns: sql.NullInt64{
				Int64: 654,
				Valid: false,
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted, err := ToInt64Ptr(tc.ns)
			if tc.wantErr {
				require.Error(t, err, "expect error")
				return
			}
			require.NoError(t, err, "expect no error")

			require.Equal(t, tc.expected, converted)
		})
	}
}
