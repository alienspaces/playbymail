package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullbool"
	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullint64"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
)

// IsUUID - tests whether provided string is a valid UUID
func IsUUID(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		return false
	}
	return true
}

func ValidateUUIDField(fieldName, fieldValue string) error {
	if !IsUUID(fieldValue) {
		return coreerror.InvalidUUID(fieldName, fieldValue)
	}
	return nil
}

func ValidateNullUUIDField(fieldName string, fieldValue sql.NullString) error {
	if !nullstring.IsValid(fieldValue) || !IsUUID(fieldValue.String) {
		return coreerror.InvalidUUID(fieldName, fieldValue.String)
	}
	return nil
}

func ValidateByteSliceField(fieldName string, fieldValue []byte) error {
	if len(fieldValue) == 0 {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateStringField(fieldName, fieldValue string) error {
	if fieldValue == "" {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateNullStringField(fieldName string, fieldValue sql.NullString) error {
	if !nullstring.IsValid(fieldValue) {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateNullBoolField(fieldName string, fieldValue sql.NullBool) error {
	if !nullbool.IsValid(fieldValue) {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateTimeField(fieldName string, fieldValue time.Time) error {
	if fieldValue.IsZero() {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateNullTimeField(fieldName string, fieldValue sql.NullTime) error {
	if !nulltime.IsValid(fieldValue) {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateIntField(fieldName string, fieldValue int) error {
	if fieldValue == 0 {
		return coreerror.NewInvalidError(fieldName, "%s should not be zero", fieldName)
	}
	return nil
}

func ValidateNullInt32Field(fieldName string, fieldValue sql.NullInt32) error {
	if !nullint32.IsValid(fieldValue) {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateNullInt64Field(fieldName string, fieldValue sql.NullInt64) error {
	if !nullint64.IsValid(fieldValue) {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateStringArrayField(fieldName string, fieldValue pq.StringArray) error {
	if len(fieldValue) == 0 {
		return coreerror.RequiredField(fieldName)
	}
	return nil
}

func ValidateNullDecimalField(fieldName string, fieldValue decimal.NullDecimal) error {
	if !fieldValue.Valid {
		return coreerror.RequiredField(fieldName)
	}

	return nil
}

func ValidateEnumField(fieldName string, fieldValue string, enumSet set.Set[string]) error {
	if !enumSet.Has(fieldValue) {
		return coreerror.NewInvalidError(fieldName, "%s is not one of >%v<, but is >%s<", fieldName, enumSet, fieldValue)
	}

	return nil
}
