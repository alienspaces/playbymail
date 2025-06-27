package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

var (
	ErrMissingHeaders         = errors.New("missing headers")
	ErrMissingRequiredHeaders = errors.New("missing required header columns")
	ErrDuplicateHeaders       = errors.New("duplicate headers")
)

type ErrHeader struct {
	Fields []string
	Err    error
}

func (e ErrHeader) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf("%v: %v", e.Err, e.Fields)
	}

	return e.Err.Error()
}

func ValidateHeader(data []byte, requiredHeaders set.Set[string]) error {
	reader := csv.NewReader(bytes.NewReader(data))

	headers, err := reader.Read()
	if err != nil {
		return ErrHeader{Err: ErrMissingHeaders}
	}

	// Check whether headers are missing first to prevent ErrDuplicateHeaders
	// being returned when all headers are missing and the first row happens
	// to have duplicate data.
	headerSet := set.New[string](headers...)

	missingHeaders := set.Difference(requiredHeaders, headerSet)
	if len(missingHeaders) == len(requiredHeaders) {
		return ErrHeader{Err: ErrMissingHeaders}
	}

	dedupSet := make(set.Set[string])
	for _, h := range headers {
		if _, exists := dedupSet[h]; exists {
			return ErrHeader{
				Fields: []string{h},
				Err:    ErrDuplicateHeaders,
			}
		}
		dedupSet.Add(h)
	}

	if len(missingHeaders) > 0 {
		return ErrHeader{
			Fields: missingHeaders.ToSlice(),
			Err:    ErrMissingRequiredHeaders,
		}
	}

	return nil
}
