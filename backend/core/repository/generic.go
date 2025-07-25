package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
)

type Recorder[Rec any] interface {
	*Rec
	ResolveID() *record.Record
	ToNamedArgs() pgx.NamedArgs
	SetCreatedAt(time.Time) *record.Record
	GetUpdatedAt() sql.NullTime
	SetUpdatedAt(sql.NullTime) *record.Record
}

type Generic[Rec any, RecPtr Recorder[Rec]] struct {
	Repository
}

var _ repositor.Repositor = &Generic[record.Record, *record.Record]{}

// NewGeneric returns a repository that supports the One(O) and Many(M) variants
// of Get(O/M), Create(O), Update(O), Delete(O/M), and Remove(O/M).
//
// If there is an error initialising the repository, a non-zero but improperly
// initialised repository will be returned.
//
// The generic parameter type RecPtr must be a pointer for the Recorder interface
// to mutate record properties in-place.
func NewGeneric[Rec any, RecPtr Recorder[Rec]](args NewArgs) (*Generic[Rec, RecPtr], error) {

	r, err := New(args)
	if err != nil {
		return &Generic[Rec, RecPtr]{Repository: r}, err
	}

	return &Generic[Rec, RecPtr]{
		Repository: r,
	}, nil
}

func (r *Generic[Rec, RecPtr]) NewRecord() *Rec {
	return new(Rec)
}

func (r *Generic[Rec, RecPtr]) NewRecordSlice() []*Rec {
	return make([]*Rec, 0)
}

// GetMany -
func (r *Generic[Rec, RecPtr]) GetMany(opts *coresql.Options) ([]*Rec, error) {

	recs := r.NewRecordSlice()

	rows, err := r.GetRows(r.GetManySQL(), opts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rec, err := pgx.RowToStructByName[Rec](rows)
		if err != nil {
			return nil, err
		}
		recs = append(recs, &rec)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return recs, nil
}

// GetOne -
func (r *Generic[Rec, RecPtr]) GetOne(id string, lock *coresql.Lock) (*Rec, error) {

	opts := &coresql.Options{
		Params: []coresql.Param{
			{
				Col: "id",
				Val: id,
			},
		},
		Lock: lock,
	}

	rows, err := r.GetRows(r.GetOneSQL(), opts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		err := pgx.ErrNoRows
		return nil, err
	}

	rec, err := pgx.RowToStructByName[Rec](rows)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}

// We could also set the `rec` parameter as `RecPtr`, then we would not need to
// do the type assertion, but the return type would be an interface, resulting
// in a loss of type information for the caller.

// CreateOne -
func (r *Generic[Rec, RecPtr]) CreateOne(rec *Rec) (*Rec, error) {

	RecPtr(rec).ResolveID()
	RecPtr(rec).SetCreatedAt(record.NewRecordTimestamp())

	rows, err := r.tx.Query(context.Background(), r.CreateOneSQL(), RecPtr(rec).ToNamedArgs())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		err := fmt.Errorf("failed to create >%s< record, create returned no rows", r.TableName())
		return nil, err
	}

	newRec, err := pgx.RowToStructByName[Rec](rows)
	if err != nil {
		return nil, err
	}

	return &newRec, nil
}

// UpdateOne -
func (r *Generic[Rec, RecPtr]) UpdateOne(rec *Rec) (*Rec, error) {

	origUpdatedAt := RecPtr(rec).GetUpdatedAt()
	RecPtr(rec).SetUpdatedAt(record.NewRecordNullTimestamp())

	rows, err := r.tx.Query(context.Background(), r.UpdateOneSQL(), RecPtr(rec).ToNamedArgs())
	if err != nil {
		RecPtr(rec).SetUpdatedAt(origUpdatedAt)
		return rec, err
	}
	defer rows.Close()

	if !rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		err := fmt.Errorf("failed to update >%s< record, updated returned no rows", r.TableName())
		return nil, err
	}

	modRec, err := pgx.RowToStructByName[Rec](rows)
	if err != nil {
		return nil, err
	}

	return &modRec, nil
}

func (r *Generic[Rec, RecPtr]) SetRLS(identifiers map[string][]string) {
	if r.isRLSDisabled {
		return
	}

	filtered := map[string][]string{}

	// SELECT queries should only filter on rows with identifiers for resources matching
	// the attributes of the record. For example; a record with only `program_id`, but
	// not `client_id`, should not be filtered on `client_id`.
	for _, attr := range r.Attributes() {
		if ids, ok := identifiers[attr]; ok {
			filtered[attr] = ids
		}
	}

	if len(filtered) > 0 {
		r.rlsIdentifiers = filtered
	}
}
