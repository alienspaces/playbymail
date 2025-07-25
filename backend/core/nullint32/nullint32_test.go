package nullint32

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
)

func TestFromInt32(t *testing.T) {

	tests := map[string]struct {
		value         int32
		expectedValid bool
	}{
		"valid value": {
			value:         432,
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
			converted := FromInt32(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.value, converted.Int32)
		})
	}
}

func TestFromInt32Ptr(t *testing.T) {

	tests := map[string]struct {
		value         *int32
		expectedValid bool
	}{
		"valid value": {
			value:         convert.Ptr(int32(654)),
			expectedValid: true,
		},
		"valid empty": {
			value:         convert.Ptr(int32(0)),
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
			converted := FromInt32Ptr(tc.value)

			require.Equal(t, tc.expectedValid, converted.Valid)

			if tc.value != nil {
				require.Equal(t, *tc.value, converted.Int32)
			}
		})
	}
}

func TestToInt32(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt32
		expected int32
		wantErr  bool
	}{
		"valid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: true,
			},
			expected: 654,
			wantErr:  false,
		},
		"invalid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: false,
			},
			expected: 0,
			wantErr:  true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted, err := ToInt32(tc.ns)
			if tc.wantErr {
				require.Error(t, err, "expect error")
				return
			}
			require.NoError(t, err, "expect no error")
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}

func TestToInt32Ptr(t *testing.T) {

	tests := map[string]struct {
		ns       sql.NullInt32
		expected *int32
		wantErr  bool
	}{
		"valid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: true,
			},
			expected: convert.Ptr(int32(654)),
			wantErr:  false,
		},
		"invalid": {
			ns: sql.NullInt32{
				Int32: 654,
				Valid: false,
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted, err := ToInt32Ptr(tc.ns)
			if tc.wantErr {
				require.Error(t, err, "expect error")
				return
			}
			require.NoError(t, err, "expect no error")
			require.Equal(t, tc.expected, converted, "expect value")
		})
	}
}
