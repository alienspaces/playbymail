package domain

import (
	"regexp"
	"strings"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

var reDBConstraint = regexp.MustCompile(`\"(.*?)\"`)

func databaseError(err error) error {

	// TODO CX-??: Make this database error handling a mapping of database error strings to errors for simpler registration and lookup

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

const (
	// General
	ReasonReferenceSelf         = "referencing self not allowed"
	ReasonReferenceNotFound     = "referenced record not found"
	ReasonExpiryBeforeEffective = "expiry time before effective time not allowed"
	ReasonEffectiveBeforeNow    = "effective time of this record or its parent is before now"
	ReasonReferencesFound       = "referencing records found"
	ReasonPropertyModification  = "property modification not allowed"
)

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// TODO: Move/revise these into core domain

// Use for runtime errors where object properties are missing or function arguments are
// missing. The error will be logged at warning log level using the provided logger.
func InternalInvalidArgumentOrProperty(l logger.Logger, arg string) error {
	err := coreerror.NewInternalError("missing argument or property >%s<", arg)
	l.Warn(err.Error())
	return err
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

// Specific service errors
