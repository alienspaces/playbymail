package error

import (
	"errors"
	"fmt"
	"strings"
)

type Errors []Error

func (e Errors) Error() string {
	sb := strings.Builder{}

	for _, er := range e {
		sb.WriteString(er.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func ToErrors(errs ...error) (Errors, error) {
	var results Errors

	for _, e := range errs {
		result, err := ToError(e)
		if err == nil {
			results = append(results, result)
			continue
		}

		var errs Errors
		if errors.As(e, &errs) {
			results = append(results, errs...)
			continue
		}

		return nil, fmt.Errorf("failed to convert err to coreerr >%#v<", err)
	}

	return results, nil
}
