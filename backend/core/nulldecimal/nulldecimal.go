package nulldecimal

import (
	"github.com/shopspring/decimal"
)

func FromDecimal(d decimal.Decimal) decimal.NullDecimal {
	if d.IsZero() {
		return decimal.NullDecimal{}
	}
	return decimal.NullDecimal{
		Decimal: d,
		Valid:   true,
	}
}

func FromDecimalPtr(d *decimal.Decimal) decimal.NullDecimal {
	if d == nil {
		return decimal.NullDecimal{
			Decimal: decimal.Zero,
			Valid:   false,
		}
	}

	return decimal.NullDecimal{
		Decimal: *d,
		Valid:   true,
	}
}

func ToDecimal(nd decimal.NullDecimal) decimal.Decimal {
	if !nd.Valid {
		return decimal.Zero
	}
	return nd.Decimal
}

func ToDecimalPtr(nd decimal.NullDecimal) *decimal.Decimal {
	if !nd.Valid {
		return nil
	}
	return &nd.Decimal
}

func IsValid(nd decimal.NullDecimal) bool {
	return nd.Valid && !nd.Decimal.IsZero()
}
