package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
)

const (
	// UniqueNameIDLength is the length of the unique ID used in UniqueName
	UniqueNameIDLength = 9
	// UniqueEmailIDLength is the length of the unique ID used in UniqueEmail
	UniqueEmailIDLength = 8
)

// UniqueName appends a pseudo-unique short ID to the front of name, for parallel test execution.
// Format: "(9chars) name"
func UniqueName(name string) string {
	if name == "" {
		name = gofakeit.Color()
	}
	return fmt.Sprintf("(%s) %s", corerecord.NewRecordID()[:UniqueNameIDLength], name)
}

// NormalName removes the bracketed ID prefix (added by UniqueName) from the start of name.
// Format: "(9chars) name" -> "name"
func NormalName(name string) string {
	// Skip past "(9chars) " which is 1 + 9 + 1 + 1 = 12 characters
	prefixLen := 1 + UniqueNameIDLength + 1 + 1 // "(" + ID + ")" + " "
	if len(name) <= prefixLen {
		return name
	}
	return name[prefixLen:]
}

// UniqueEmail appends a UUID4 prefix to the email to make it unique
// for parallel test execution.
// Format: "8chars-email"
func UniqueEmail(email string) string {
	return fmt.Sprintf("%s-%s", corerecord.NewRecordID()[:UniqueEmailIDLength], email)
}

// NormalEmail removes the unique ID prefix (added by UniqueEmail) from the start of email.
// Format: "8chars-email" -> "email"
func NormalEmail(email string) string {
	// Skip past "8chars-" which is 8 + 1 = 9 characters
	prefixLen := UniqueEmailIDLength + 1 // ID + "-"
	if len(email) <= prefixLen {
		return email
	}
	return email[prefixLen:]
}
