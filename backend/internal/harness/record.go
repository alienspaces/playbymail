package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
)

// UniqueName appends a UUID4 to the end of the name to make it unique
// for parallel test execution.
func UniqueName(name string) string {
	if name == "" {
		name = gofakeit.Color()
	}
	return fmt.Sprintf("%s (%s)", name, corerecord.NewRecordID())
}

// NormalName removes the unique UUID4 from the end of the name to make it normal for
// test harness functions that return a record based on its non unique name.
func NormalName(name string) string {
	return name[:len(name)-39]
}

// UniqueEmail appends a UUID4 to the end of the email to make it unique
// for parallel test execution.
func UniqueEmail(email string) string {
	return fmt.Sprintf("%s-%s", email, corerecord.NewRecordID())
}
