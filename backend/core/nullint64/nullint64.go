package nullint64

import (
	"database/sql"
	"fmt"
)

func FromInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func FromInt64Ptr(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}

	return sql.NullInt64{
		Int64: *i,
		Valid: true,
	}
}

func ToInt64(ni sql.NullInt64) (int64, error) {
	if !ni.Valid {
		return 0, fmt.Errorf("type NullInt64 is not valid, cannot convert to int64")
	}
	return ni.Int64, nil
}

func ToInt64Ptr(ni sql.NullInt64) (*int64, error) {
	if !ni.Valid {
		return nil, fmt.Errorf("type NullInt64 is not valid, cannot convert to *int64")
	}
	return &ni.Int64, nil
}

func IsValid(ni sql.NullInt64) bool {
	return ni.Valid
}
