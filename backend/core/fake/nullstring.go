package fake

import (
	"database/sql"

	"github.com/brianvoe/gofakeit/v6"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/record"
)

func EmailNullStr() sql.NullString {
	return nullstring.FromString(gofakeit.Email())
}

func UniqueProgramNameNullStr() sql.NullString {
	return nullstring.FromString(UniqueProgramNameStr())
}

func NameNullStr() sql.NullString {
	return nullstring.FromString(gofakeit.Name())
}

func FirstNameNullStr() sql.NullString {
	return nullstring.FromString(gofakeit.FirstName())
}

func LastNameNullStr() sql.NullString {
	return nullstring.FromString(gofakeit.LastName())
}

func CompanyNameNullStr() sql.NullString {
	return nullstring.FromString(gofakeit.Company())
}

func CountryNullStr() sql.NullString {
	return nullstring.FromString(gofakeit.Country())
}

func CityNullStr() sql.NullString {
	return nullstring.FromString(gofakeit.City())
}

func UUIDNullStr() sql.NullString {
	return nullstring.FromString(record.NewRecordID())
}

func ServiceCloudIDNullStr() sql.NullString {
	return nullstring.FromString(record.NewRecordID())
}
