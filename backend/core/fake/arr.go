package fake

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
)

// CurrenciesPqStrArr returns a slice of the specified length containing random unique ISO-4217 currencies.
//
// CurrenciesPqStrArr should only be used to generate small slices with a length much less than 163
// (the number of currencies specified in gofakeit).
func CurrenciesPqStrArr(length int) pq.StringArray {
	if length == 0 {
		return pq.StringArray{}
	}

	var currencies pq.StringArray
	hasCurrency := map[string]struct{}{}

	for i := 0; i < length; i++ {
		c := gofakeit.CurrencyShort()
		_, ok := hasCurrency[c]
		for ok {
			c = gofakeit.CurrencyShort()
			_, ok = hasCurrency[c]
		}

		hasCurrency[c] = struct{}{}
		currencies = append(currencies, c)
	}

	return currencies
}

// GENERICS-CANDIDATE
func CurrenciesStrSlice(length int) []string {
	if length == 0 {
		return []string{}
	}

	var currencies []string
	hasCurrency := map[string]struct{}{}

	for i := 0; i < length; i++ {
		c := gofakeit.CurrencyShort()
		_, ok := hasCurrency[c]
		for ok {
			c = gofakeit.CurrencyShort()
			_, ok = hasCurrency[c]
		}

		hasCurrency[c] = struct{}{}
		currencies = append(currencies, c)
	}

	return currencies
}

// GENERICS-CANDIDATE
func CountriesStrSlice(length int) []string {
	if length == 0 {
		return []string{}
	}

	var countries []string
	hasCountry := map[string]struct{}{}

	for i := 0; i < length; i++ {
		c := gofakeit.CountryAbr()
		_, ok := hasCountry[c]
		for ok {
			c = gofakeit.CountryAbr()
			_, ok = hasCountry[c]
		}

		hasCountry[c] = struct{}{}
		countries = append(countries, c)
	}

	return countries
}

// GENERICS-CANDIDATE
func LocaleStrSlice(length int) []string {
	if length == 0 {
		return []string{}
	}

	var locales []string
	hasLocale := map[string]struct{}{}

	for i := 0; i < length; i++ {
		c := gofakeit.LanguageBCP()
		_, ok := hasLocale[c]
		for ok {
			c = gofakeit.LanguageBCP()
			_, ok = hasLocale[c]
		}

		hasLocale[c] = struct{}{}
		locales = append(locales, c)
	}

	return locales
}
