package domain

import (
	"regexp"
	"strings"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
)

var reDBConstraint = regexp.MustCompile(`\"(.*?)\"`)

func databaseError(err error) error {

	// Handle specific constraint errors
	rs := reDBConstraint.FindStringSubmatch(err.Error())
	if len(rs) > 0 {
		switch rs[1] {
		//
		}
	}

	if strings.Contains(err.Error(), "timeout") {
		return coreerror.NewUnavailableError()
	}

	return Internal("database error >%s<", err.Error())
}

func RequiredField(fieldName string) error {
	return coreerror.RequiredField(fieldName)
}

func InvalidField(fieldName string, fieldValue string, reason string) error {
	return coreerror.InvalidField(fieldName, fieldValue, reason)
}

func Internal(reason string, args ...any) error {
	return coreerror.NewInternalError(reason, args...)
}

func NotFound(recordName, id string) error {
	return coreerror.NewNotFoundError(recordName, id)
}
