package fake

import (
	"github.com/brianvoe/gofakeit/v6"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
)

func EmailStrPtr() *string {
	return convert.Ptr(gofakeit.Email())
}

func FirstNameStrPtr() *string {
	return convert.Ptr(gofakeit.FirstName())
}

func LastNameStrPtr() *string {
	return convert.Ptr(gofakeit.LastName())
}

func UUIDStrPtr() *string {
	return convert.Ptr(record.NewRecordID())
}

func CompanyNameStrPtr() *string {
	return convert.Ptr(gofakeit.Company())
}
