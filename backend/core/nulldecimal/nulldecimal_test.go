package nulldecimal

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
)

func TestFromDecimal(t *testing.T) {

	tests := map[string]struct {
		decimal       decimal.Decimal
		expectedValid bool
	}{
		"valid decimal": {
			decimal:       decimal.NewFromInt(10),
			expectedValid: true,
		},
		"invalid decimal": {
			decimal:       decimal.Decimal{},
			expectedValid: false,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := FromDecimal(tc.decimal)

			require.Equal(t, tc.expectedValid, converted.Valid)
			require.Equal(t, tc.decimal, converted.Decimal)
		})
	}
}

func TestToDecimalPtr(t *testing.T) {

	tests := map[string]struct {
		nd       decimal.NullDecimal
		expected *decimal.Decimal
	}{
		"valid": {
			nd: decimal.NullDecimal{
				Decimal: decimal.NewFromInt(10),
				Valid:   true,
			},
			expected: convert.Ptr(decimal.NewFromInt(10)),
		},
		"invalid": {
			nd: decimal.NullDecimal{
				Decimal: decimal.Decimal{},
				Valid:   false,
			},
			expected: nil,
		},
	}

	for tcName, tc := range tests {
		t.Run(tcName, func(t *testing.T) {
			t.Logf("Running test >%s<", tcName)
			converted := ToDecimalPtr(tc.nd)

			require.Equal(t, tc.expected, converted)
		})
	}
}
