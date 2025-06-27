package fake

import (
	"strings"

	"github.com/brianvoe/gofakeit/v6"

	"gitlab.com/alienspaces/playbymail/core/record"
)

func UniqueProgramNameStr() string {
	return gofakeit.Name() + " - " + record.NewRecordID()
}

func UniqueProjectAbbreviationStr() string {
	return record.NewRecordID()[0:3]
}

func UniqueSpherePayerCodeStr() string {
	return strings.Replace(record.NewRecordID(), "-", "", -1)
}

func UniqueSphereProcessorCodeStr() string {
	return strings.Replace(record.NewRecordID(), "-", "", -1)
}

func UniqueDomainNameStr() string {
	return gofakeit.DomainName() + "-" + record.NewRecordID()
}

func UniqueSubDomainNameStr() string {
	return gofakeit.DomainName() + "-" + record.NewRecordID()
}
