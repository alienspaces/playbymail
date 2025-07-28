package error

import (
	"strings"
)

type ValidationErrorType int

const (
	ValidationErrorUnsupported ValidationErrorType = iota
	ValidationErrorInvalid
	ValidationErrorInvalidAction
)

func (t ValidationErrorType) String() string {
	return [...]string{"unsupported", "invalid", "invalid_action"}[t]
}

type LinkedFields struct {
	LinkedField string
	Fields      []string
}

// CreateRegistry takes an error type and a list of subjects that may
// be a list of actions or fields that can be validated and returns a
// Registry of error codes and errors
func CreateRegistry(et ValidationErrorType, subjects ...string) Registry {
	errorCollection := Registry{}

	for _, subject := range subjects {
		errCode := CreateErrorCode(et, subject)

		var e Error
		switch et {
		case ValidationErrorInvalid:
			e, _ = ToError(NewInvalidError(subject, "The property '%s' is %s.", subject, et))
		case ValidationErrorInvalidAction:
			e, _ = ToError(NewInvalidActionError(subject, "The property '%s' is %s.", subject, et))
		case ValidationErrorUnsupported:
			e, _ = ToError(NewUnsupportedError(subject, "The property '%s' is %s.", subject, et))
		}
		errorCollection[errCode] = e
	}

	return errorCollection
}

func CreateLinkedRegistry(et ValidationErrorType, linkedFields []LinkedFields) Registry {
	errorCollection := Registry{}

	for _, f := range linkedFields {
		errCode := CreateErrorCode(et, f.LinkedField)
		combinationMsg := strings.Join(f.Fields, " & ")

		var e Error
		if et == ValidationErrorInvalid {
			e, _ = ToError(NewInvalidError(f.LinkedField, "The combination of %s is %s.", combinationMsg, et))
		} else {
			e, _ = ToError(NewUnsupportedError(f.LinkedField, "The combination of %s is %s.", combinationMsg, et))
		}
		errorCollection[errCode] = e
	}

	return errorCollection
}
