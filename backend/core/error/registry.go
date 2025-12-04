package error

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var (
	reArray = regexp.MustCompile(`(?m)\.(\d+)(\.)?`)
)

var registry = Registry{
	ErrorCodeSchemaValidation: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeSchemaValidation,
		Message:        "Request body failed JSON schema validation.",
	},
	ErrorCodeInvalidData: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeInvalidData,
		Message:        "Request body contains invalid JSON.",
	},
	ErrorCodeInvalidParam: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeInvalidParam,
		Message:        "The value for the parameter is invalid.",
	},
	ErrorCodeNotFound: Error{
		HttpStatusCode: http.StatusNotFound,
		ErrorCode:      ErrorCodeNotFound,
		Message:        "Resource not found.",
	},
	ErrorCodeUnauthorized: Error{
		HttpStatusCode: http.StatusForbidden,
		ErrorCode:      ErrorCodeUnauthorized,
		Message:        "Permission to the requested resource is denied.",
	},
	ErrorCodeUnauthenticated: Error{
		HttpStatusCode: http.StatusUnauthorized,
		ErrorCode:      ErrorCodeUnauthenticated,
		Message:        "Authentication information is missing or invalid.",
	},
	ErrorCodeUnavailable: Error{
		HttpStatusCode: http.StatusServiceUnavailable,
		ErrorCode:      ErrorCodeUnavailable,
		Message:        "Server overloaded: unable to process request",
	},
	ErrorCodeMalformedData: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeMalformedData,
		Message:        "Malformed data: unable to process the request",
	},
	ErrorCodeMissingData: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeMissingData,
		Message:        "Missing data: unable to process the request",
	},
	ErrorCodeInternal: Error{
		HttpStatusCode: http.StatusInternalServerError,
		ErrorCode:      ErrorCodeInternal,
		Message:        "An internal error has occurred.",
	},
}

func GetRegistryError(code Code) Error {
	return deepcopy(registry[code])
}

func deepcopy(e Error) Error {
	detail := e.SchemaValidationErrors

	if len(detail) > 0 {
		e.SchemaValidationErrors = make([]SchemaValidationError, len(detail))
		copy(e.SchemaValidationErrors, detail)
	}

	return e
}

func NewInternalError(reason string, args ...any) Error {
	err := GetRegistryError(ErrorCodeInternal)
	err.Reason = fmt.Sprintf(reason, args...)
	return err
}

func NewNotFoundError(recordName string, id string) Error {
	err := GetRegistryError(ErrorCodeNotFound)
	err.Message = fmt.Sprintf("%s with ID >%s< not found", recordName, id)
	return err
}

func NewUnavailableError() Error {
	return GetRegistryError(ErrorCodeUnavailable)
}

func NewMalformedDataError(message string, args ...any) Error {
	err := GetRegistryError(ErrorCodeMalformedData)
	if message != "" {
		err.Message = fmt.Sprintf(message, args...)
	}
	return err
}

func NewMissingDataError(message string, args ...any) Error {
	err := GetRegistryError(ErrorCodeMissingData)
	if message != "" {
		err.Message = fmt.Sprintf(message, args...)
	}
	return err
}

func NewUnauthorizedError() Error {
	return GetRegistryError(ErrorCodeUnauthorized)
}

func NewUnauthenticatedError(message string, args ...any) Error {
	err := GetRegistryError(ErrorCodeUnauthenticated)
	if message != "" {
		err.Message = fmt.Sprintf(message, args...)
	}
	return err
}

func NewParamError(message string, args ...any) Error {
	err := GetRegistryError(ErrorCodeInvalidParam)
	if message != "" {
		err.Message = fmt.Sprintf(message, args...)
	}
	return err
}

func NewHeaderError(message string, args ...any) Error {
	err := GetRegistryError(ErrorCodeInvalidHeader)
	if message != "" {
		err.Message = fmt.Sprintf(message, args...)
	}
	return err
}

func NewInvalidDataError(message string, args ...any) Error {
	err := GetRegistryError(ErrorCodeInvalidData)
	if message != "" {
		err.Message = fmt.Sprintf(message, args...)
	}
	return err
}

func NewInvalidError(errorCodeSuffix string, message string, args ...any) Error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      CreateErrorCode(ValidationErrorInvalid, errorCodeSuffix),
		Message:        fmt.Sprintf(message, args...),
	}
}

func NewInvalidActionError(action, message string, args ...any) error {
	return Error{
		HttpStatusCode: http.StatusConflict,
		ErrorCode:      CreateErrorCode(ValidationErrorInvalidAction, action),
		Message:        fmt.Sprintf(message, args...),
	}
}

func NewUnsupportedError(errorCodeSuffix string, message string, args ...any) Error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      CreateErrorCode(ValidationErrorUnsupported, errorCodeSuffix),
		Message:        fmt.Sprintf(message, args...),
	}
}

// InvalidUUID is a convenience function for invalid UUID errors that
// provides a standard formatted error message.
func InvalidUUID(field, value string) Error {
	return NewInvalidError(field, "%s >%s< is invalid, not a UUID", field, value)
}

// InvalidField is a convenience function for invalid field errors that
// provides a standard formatted error message.
func InvalidField(fieldName, fieldValue, reason string) error {
	if reason != "" {
		return NewInvalidError(fieldName, "%s >%s< is invalid, %s", fieldName, fieldValue, reason)
	}
	return NewInvalidError(fieldName, "%s >%s< is invalid", fieldName, fieldValue)
}

// InvalidAction is a convenience function for invalid action errors that
// provides a standard formatted error message.
func InvalidAction(actionType, reason string) error {
	if reason != "" {
		return NewInvalidActionError(actionType, "%s could not be completed, %s", actionType, reason)
	}
	return NewInvalidActionError(actionType, "%s could not be completed", actionType)
}

// RequiredField is a convenience function for empty field errors that
// provides a standard formatted error message.
func RequiredField(field string) Error {
	return NewInvalidError(field, "%s should not be empty", field)
}

// RequiredPathParameter is a convenience function for missing path parameter errors that
// provides a standard formatted error message.
func RequiredPathParameter(parameter string) Error {
	return NewInvalidError(parameter, "Path parameter '%s' is required", parameter)
}

// RequiredQueryParameter is a convenience function for missing query parameter errors that
// provides a standard formatted error message.
func RequiredQueryParameter(parameter string) Error {
	return NewInvalidError(parameter, "Query parameter '%s' is required", parameter)
}

// ImmutableField is a convenience function for immutable field errors that
// provides a standard formatted error message.
func ImmutableField(field string) Error {
	return NewInvalidError(field, "%s cannot be modified", field)
}

func CreateErrorCode(errorType ValidationErrorType, field string) Code {
	return Code(fmt.Sprintf("%s_%s", errorType, field))
}

func NewSchemaValidationError(resultErrors []gojsonschema.ResultError) Error {
	e := GetRegistryError(ErrorCodeSchemaValidation)

	resultErrors = filterNonUserFriendlyErrors(resultErrors)

	for _, re := range resultErrors {
		sve := setDataPath(SchemaValidationError{}, re)
		sve = setMessage(sve, re)
		e.SchemaValidationErrors = append(e.SchemaValidationErrors, sve)
	}

	return e
}

func filterNonUserFriendlyErrors(re []gojsonschema.ResultError) []gojsonschema.ResultError {
	var friendly []gojsonschema.ResultError
	var unfriendly []gojsonschema.ResultError

	// These errors refer to conditionals in the schema that may not be understood by end-users.
	for _, err := range re {
		errType := err.Type()
		switch errType {
		case "number_any_of", "number_one_of", "number_all_of", "number_not", "condition_then", "condition_else":
			unfriendly = append(unfriendly, err)
		default:
			friendly = append(friendly, err)
		}
	}

	// The non-user friendly errors are _usually_ accompanied by a more specific user-friendly error.
	if len(friendly) == 0 {
		return unfriendly
	}

	return friendly
}

func setDataPath(sve SchemaValidationError, re gojsonschema.ResultError) SchemaValidationError {
	var field string
	if re.Type() == "required" {
		field = re.Details()["property"].(string)
	} else {
		field = re.Field()
	}

	sve.DataPath = "$"

	// not sure if it is possible for the field to be empty, but to be safe the path is set to "$"
	switch field {
	case "", "(root)":
		return sve
	}

	// reformat fields with array index and prefix with "$." (e.g contacts.0.type -> $.contacts[0].type, contacts.0 -> $.contacts[0])
	sve.DataPath = sve.DataPath + "." + reArray.ReplaceAllString(field, "[$1]$2")

	return sve
}

// setMessage sets the detail of the validation error with the reformatted errors returned from the validation.
func setMessage(sve SchemaValidationError, re gojsonschema.ResultError) SchemaValidationError {
	switch re.Type() {
	case "number_gte", "number_gt", "number_lte", "number_lt", "format", "pattern", "array_min_items", "array_max_items":
		sve.Message = re.String()
		if strings.Contains(sve.Message, " 1 items") {
			sve.Message = strings.ReplaceAll(sve.Message, " 1 items", " 1 item")
		}
	default:
		sve.Message = re.Description()
	}

	// clean up message to avoid repeating the property
	if strings.Contains(sve.Message, re.Field()+": ") {
		sve.Message = strings.ReplaceAll(sve.Message, re.Field()+": ", "")
	}
	if strings.Contains(sve.Message, re.Field()+" must") {
		sve.Message = strings.ReplaceAll(sve.Message, re.Field()+" must", "Must")
	}
	if strings.Contains(sve.Message, re.Field()+" does") {
		sve.Message = strings.ReplaceAll(sve.Message, re.Field()+" does", "Does")
	}

	return sve
}
