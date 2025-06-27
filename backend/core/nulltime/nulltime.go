package nulltime

import (
	"database/sql"
	"time"
)

func FromTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func FromTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}

	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func ToTime(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

func ToTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func ToTimePtrStrict(nt sql.NullTime) *time.Time {
	if !IsValid(nt) {
		return nil
	}

	return &nt.Time
}

func IsValid(nt sql.NullTime) bool {
	return nt.Valid && !nt.Time.IsZero()
}
