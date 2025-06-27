package currency

import (
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

// SubdenominationMonetaryUnitToLowestDenomination converts a monetary unit with decimal places
// beyond the lowest denomination to the lowest denomination unit using half-away rounding.
//
// SubdenominationMonetaryUnitToLowestDenomination handles currencies with decimal places where the value has less
// than the maximum number of decimal places the currency may have (e.g., 2.0 AUD or 2 AUD).
func SubdenominationMonetaryUnitToLowestDenomination(value string, currency string) (int64, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return 0, err
	}

	lc := accounting.LocaleInfo[currency]
	newRate := dec.Round(int32(lc.FractionLength))
	newRateCents := newRate.Shift(int32(lc.FractionLength))

	return newRateCents.IntPart(), nil
}

func LowestDenominationMonetaryUnitToFormatted(value string, currency string) (string, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return "", err
	}

	lc := LocaleInfo[currency]
	shifted := dec.Shift(int32(-lc.FractionLength))

	// TODO merge locale/currency data or use https://github.com/bojanz/currency

	// NOTE: this package is no longer actively maintained.
	// We only use the precision, thousand separator, and decimal separator to format the currency value using the pkg,
	// as there is no way to format with a custom `Accounting.Pre` value: https://github.com/leekchan/accounting/issues/25
	ac := accounting.Accounting{
		Precision: lc.FractionLength,
		Thousand:  lc.ThouSep,
		Decimal:   lc.DecSep,
	}

	// Consider extracting and adapting the FormatMoneyDecimal if bojanz/currency is inadequate.
	// The rest of the accounting package is not needed.
	formatted := ac.FormatMoneyDecimal(shifted)
	if lc.Pre {
		return lc.ComSymbol + formatted, nil
	}

	return formatted + lc.ComSymbol, nil
}
