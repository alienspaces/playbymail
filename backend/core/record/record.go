package record

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/r3labs/diff/v3"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
)

const (
	FieldID        = "id"
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
	FieldDeletedAt = "deleted_at"
)

// NOTE:
// Use sql.NullXxx types when the underlying database column:
// - has a NOT NULL check constraint;
// - does not otherwise have any other CHECK constraints;
// - is not an eval type;
// - is not an uuid type;
// - does not have a foreign key constraint; and
// - you don't want it to default to Go's default value for the property type.

// Record is the base record struct most service records are composed from.
//
// ID is typically a UUID and is a required field.
//
// CreatedAt is the UTC timestamp for when the record was created.
//
// UpdatedAt is the UTC timestamp for when the record was last updated.
//
// DeletedAt is the UTC timestamp for when the record was logically deleted.
type Record struct {
	ID        string       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (r *Record) clearID() {
	r.ID = ""
}
func (r *Record) clearTimestamps() {
	r.CreatedAt = time.Time{}
	r.UpdatedAt = nulltime.FromTime(time.Time{})
	r.DeletedAt = nulltime.FromTime(time.Time{})
}

func (r *Record) ResolveID() *Record {
	if r.ID == "" {
		r.ID = NewRecordID()
	}

	return r
}

func (r *Record) ToNamedArgs() pgx.NamedArgs {
	args := pgx.NamedArgs{
		"id":         r.ID,
		"created_at": r.CreatedAt,
		"updated_at": r.UpdatedAt,
		"deleted_at": r.DeletedAt,
	}
	return args
}

func (r *Record) SetCreatedAt(t time.Time) *Record {
	r.CreatedAt = t
	return r
}

func (r *Record) GetUpdatedAt() sql.NullTime {
	return r.UpdatedAt
}

func (r *Record) SetUpdatedAt(t sql.NullTime) *Record {
	r.UpdatedAt = t
	return r
}

type EqualityFlag string

const EqualityFlagExcludeID EqualityFlag = "eo-exclude-id"
const EqualityFlagExcludeTimestamps EqualityFlag = "flag-exclude-timestamps"

type EquatableRecord interface {
	clearID()
	clearTimestamps()
}

func RecordEqual(p, pp EquatableRecord, flags ...EqualityFlag) (bool, error) {
	changelog, err := RecordDiff(p, pp, flags...)
	return len(changelog) == 0, err
}

func RecordDiff(p, pp EquatableRecord, flags ...EqualityFlag) (diff.Changelog, error) {
	for idx := range flags {
		switch flags[idx] {
		case EqualityFlagExcludeTimestamps:
			p.clearTimestamps()
			pp.clearTimestamps()
		case EqualityFlagExcludeID:
			p.clearID()
			pp.clearID()
		}
	}
	return diff.Diff(p, pp)
}

// NewRecordID -
func NewRecordID() string {
	return uuid.NewString()
}

func NewRecordIDPtr() *string {
	return convert.Ptr(uuid.NewString())
}

func NewRecordIDNullStr() sql.NullString {
	return nullstring.FromString(uuid.NewString())
}

// NewRecordTimestamp -
func NewRecordTimestamp() time.Time {
	return timestamp()
}

// RFC3339Microseconds may be used for time formatting since postgres stores time to the microsecond
const RFC3339Microseconds = "2006-01-02T15:04:05.000000Z07:00"

func NewRecordTimestampStr() string {
	return timestamp().Format(RFC3339Microseconds)
}

func NewRecordTimestampPtr() *time.Time {
	return convert.Ptr(timestamp())
}

// NewRecordNullTimestamp -
func NewRecordNullTimestamp() sql.NullTime {
	return nulltime.FromTime(timestamp())
}

func ToRecordNullTimestamp(t time.Time) sql.NullTime {
	return nulltime.FromTime(ToRecordTimestamp(t))
}

func ToRecordTimestamp(t time.Time) time.Time {
	// Why round to 1 microsecond? https://gitlab.com/msts-enterprise/ox/pricing/-/wikis/Engineering/Timestamp-Precision
	return t.UTC().Round(time.Microsecond)
}

func timestamp() time.Time {
	return ToRecordTimestamp(time.Now())
}
