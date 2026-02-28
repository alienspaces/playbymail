package error

import (
	"errors"
	"fmt"
	"strings"
)

type Code string

const (
	ErrorCodeSchemaValidation Code = "body_not_matching_json_schema"
	ErrorCodeInvalidAction    Code = "invalid_action"
	ErrorCodeInvalidData      Code = "invalid_data"
	ErrorCodeInvalidHeader    Code = "invalid_header"
	ErrorCodeInvalidParam     Code = "invalid_parameter"
	ErrorCodeMalformedData    Code = "malformed_data"
	ErrorCodeMissingData      Code = "missing_data"
	ErrorCodeNotFound         Code = "resource_not_found"
	ErrorCodeUnauthorized     Code = "unauthorized"
	ErrorCodeUnauthenticated  Code = "unauthenticated"
	ErrorCodeUnavailable      Code = "unavailable"
	ErrorCodeInternal         Code = "internal_error"
)

// Error -
//
// Context should either be nil, or a struct with its own JSON tags.
type Error struct {
	HttpStatusCode         int                     `json:"-"`
	ErrorCode              Code                    `json:"code"`
	Message                string                  `json:"message"`
	Reason                 string                  `json:"-"` // Not exposed through JSON APIs
	SchemaValidationErrors []SchemaValidationError `json:"validationErrors,omitempty"`
	Context                any                     `json:"context,omitempty"`
}

func (e Error) Error() string {
	if e.Reason != "" {
		reason := fmt.Sprintf(">%s<", e.Reason)
		return fmt.Sprintf("%s: %s %s", e.ErrorCode, e.Message, reason)
	}
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

// WithContext sets context on the Error.
//
// context should be a struct with its own JSON tags.
func (e Error) WithContext(context any) Error {
	e.Context = context
	return e
}

// Registry is a map of error codes to errors
type Registry map[Code]Error

// Merge merges another error collection with this error collection returning a
// new error collection
func (c Registry) Merge(a Registry) Registry {
	for k, v := range c {
		a[k] = v
	}
	return a
}

type SchemaValidationError struct {
	DataPath string `json:"dataPath"`
	Message  string `json:"message"`
}

func (sve SchemaValidationError) GetField() string {
	field := strings.Split(sve.DataPath, ".")
	lastField := field[len(field)-1]
	return lastField
}

func IsError(e error) bool {
	var errorPtr Error
	return errors.As(e, &errorPtr)
}

func HasErrorCode(err error, c Code) bool {
	e, err := ToError(err)
	if err != nil {
		return false
	}

	return e.ErrorCode == c
}

func ToError(e error) (Error, error) {
	if e == nil {
		return Error{}, fmt.Errorf("err is nil when converting to coreerror.Error type")
	}

	var err Error
	if !errors.As(e, &err) {
		return Error{}, fmt.Errorf("failed to convert to coreerror.Error type >%v<", e)
	}

	if len(err.SchemaValidationErrors) == 0 {
		err.SchemaValidationErrors = nil
	}

	return err, nil
}

func ProcessParamError(err error) error {
	e, conversionErr := ToError(err)
	if conversionErr != nil {
		return err
	}

	if len(e.SchemaValidationErrors) == 0 {
		return NewParamError("%s", e.Error())
	}

	errStr := strings.Builder{}
	errStr.WriteString("Invalid parameter(s): ")
	for i, sve := range e.SchemaValidationErrors {
		if sve.GetField() == "$" {
			fmt.Fprintf(&errStr, "(%d) %s; ", i+1, sve.Message)
		} else {
			fmt.Fprintf(&errStr, "(%d) %s: %s; ", i+1, sve.GetField(), sve.Message)
		}
	}

	formattedErrString := errStr.String()
	formattedErrString = formattedErrString[0 : len(formattedErrString)-2] // remove extra space and semicolon
	formattedErrString += "."
	return NewParamError("%s", formattedErrString)
}
