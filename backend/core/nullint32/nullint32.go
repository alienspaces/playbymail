package nullint32

import (
	"database/sql"
	"fmt"
)

func FromInt32(s int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: s,
		Valid: true,
	}
}

func FromInt32Ptr(s *int32) sql.NullInt32 {
	if s == nil {
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	}

	return sql.NullInt32{
		Int32: *s,
		Valid: true,
	}
}

func ToInt32(ns sql.NullInt32) (int32, error) {
	if !ns.Valid {
		return 0, fmt.Errorf("type NullInt32 is not valid, cannot convert to int32")
	}
	return ns.Int32, nil
}

func ToInt32Ptr(ns sql.NullInt32) (*int32, error) {
	if !ns.Valid {
		return nil, fmt.Errorf("type NullInt32 is not valid, cannot convert to *int32")
	}
	return &ns.Int32, nil
}

func IsValid(ns sql.NullInt32) bool {
	return ns.Valid
}
