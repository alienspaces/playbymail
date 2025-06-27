package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
)

type GenericView[Rec any] struct {
	Repository
}

var _ repositor.Repositor = &GenericView[struct{}]{}

// NewGenericView returns a view repository that supports only GetOne and GetMany.
//
// If there is an error initialising the repository, a non-zero but improperly
// uninitialised repository will be returned.
func NewGenericView[Rec any](args NewArgs) (*GenericView[Rec], error) {

	r, err := New(args)
	if err != nil {
		return nil, err
	}

	return &GenericView[Rec]{
		Repository: r,
	}, nil
}

func (r *GenericView[Rec]) NewRecord() *Rec {
	return new(Rec)
}

func (r *GenericView[Rec]) NewRecordSlice() []*Rec {
	return make([]*Rec, 0)
}

// GetMany -
func (r *GenericView[Rec]) GetMany(opts *coresql.Options) ([]*Rec, error) {

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
func (r *GenericView[Rec]) GetOne(id string, lock *coresql.Lock) (*Rec, error) {

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

// CreateOne -
func (r *GenericView[Rec]) CreateOne(rec *Rec) error {
	return fmt.Errorf("repository for view does not support create")
}

// UpdateOne -
func (r *GenericView[Rec]) UpdateOne(rec *Rec) error {
	return fmt.Errorf("repository for view does not support update")
}

// DeleteOne -
func (r *GenericView[Rec]) DeleteOne(id any) error {
	return fmt.Errorf("repository for view does not support delete")
}

// RemoveOne -
func (r *GenericView[Rec]) RemoveOne(id any) error {
	return fmt.Errorf("repository for view does not support remove")
}
