package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

type TestRecord struct {
	Name sql.NullString `db:"name"`
	record.Record
}

func (r *TestRecord) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args["name"] = r.Name
	return args
}

func setup(s storer.Storer) (func() error, error) {

	sql := `
CREATE TABLE test
(
	id           UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	name         TEXT NOT NULL,
	created_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at   TIMESTAMPTZ,
	deleted_at   TIMESTAMPTZ
);	
`

	tx, err := s.BeginTx()
	defer func() {
		tx.Rollback(context.Background())
	}()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(context.Background(), sql)
	if err != nil {
		return nil, err
	}

	tx.Commit(context.Background())

	teardown := func() error {
		sql := `
DROP TABLE test		
`
		tx, err := s.BeginTx()
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), sql)
		defer func() {
			tx.Rollback(context.Background())
		}()
		if err != nil {
			return err
		}

		tx.Commit(context.Background())

		return nil
	}

	return teardown, nil
}

func Test_Generic(t *testing.T) {
	_, s, err := newDependencies()
	require.NoError(t, err, "NewDependencies returns without error")
	defer func() {
		err = s.ClosePool()
		require.NoError(t, err, "ClosePool returns without error")
	}()

	teardown, err := setup(s)
	defer func() {
		err = teardown()
		require.NoError(t, err, "setup returns without error")
	}()
	require.NoError(t, err, "setup returns without error")

	tx, err := s.BeginTx()
	require.NoError(t, err, "BeginTx returns without error")
	defer func() {
		tx.Rollback(context.Background())
	}()

	r, err := NewGeneric[TestRecord](
		NewArgs{
			Tx:        tx,
			TableName: "test",
			Record:    TestRecord{},
		},
	)
	require.NoError(t, err, "NewGeneric returns without error")

	testRec := &TestRecord{
		Name: record.NewRecordIDNullStr(),
	}

	t.Run("test record named args", func(t *testing.T) {
		err := TestToNamedArgs(testRec)
		require.NoError(t, err, "testToNamedArgs returns without error with record >%v<", testRec)
	})

	var cRec, gRec, uRec *TestRecord
	var gRecs []*TestRecord

	tname := "generic when create one then ok"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)
		cRec, err = r.CreateOne(testRec)
		require.NoError(t, err, "CreateOne returns without error")
		require.NotNil(t, cRec, "CreateOne returns a record")
		require.NotEmpty(t, cRec.ID, "CreateOne record ID is not empty")
		require.NotEmpty(t, cRec.Name, "CreateOne record Name is not empty")
		require.NotEmpty(t, cRec.CreatedAt, "CreateOne record CreatedAt is not empty")
		zone, _ := cRec.CreatedAt.Zone()
		require.Equal(t, "UTC", zone, "CreateOne record CreatedAt timezone equals expected")
		require.False(t, nulltime.IsValid(cRec.UpdatedAt), "CreateOne record UpdatedAt is not valid")
		require.False(t, nulltime.IsValid(cRec.DeletedAt), "CreateOne record DeletedAt is not valid")
	})

	tname = "generic when get one then ok"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)
		gRec, err = r.GetOne(cRec.ID, nil)
		require.NoError(t, err, "GetOne returns without error")
		require.NotNil(t, gRec, "GetOne returns a record")
		require.NotEmpty(t, gRec.ID, "GetOne record ID is not empty")
		require.Equal(t, cRec.Name, gRec.Name, "GetOne record Name equals expected")
		require.Equal(t, cRec.CreatedAt, gRec.CreatedAt, "GetOne record CreatedAt is equals expected")
		zone, _ := gRec.CreatedAt.Zone()
		require.Equal(t, "UTC", zone, "GetOne record CreatedAt timezone equals expected")
		require.False(t, nulltime.IsValid(gRec.UpdatedAt), "GetOne record UpdatedAt is not valid")
		require.False(t, nulltime.IsValid(gRec.DeletedAt), "GetOne record DeletedAt is not valid")
	})

	tname = "generic when get many then ok"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)
		gRecs, err = r.GetMany(&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "id",
					Val: gRec.ID,
				},
			},
		})
		require.NoError(t, err, "GetMany returns without error")
		require.Equal(t, 1, len(gRecs), "GetMany returns expected number of records")
		require.NotEmpty(t, gRecs[0].ID, "GetMany record ID is not empty")
		require.Equal(t, cRec.Name, gRecs[0].Name, "GetMany record Name equals expected")
		require.Equal(t, cRec.CreatedAt, gRecs[0].CreatedAt, "GetMany record CreatedAt is equals expected")
		zone, _ := gRecs[0].CreatedAt.Zone()
		require.Equal(t, "UTC", zone, "GetMany record CreatedAt timezone equals expected")
		require.False(t, nulltime.IsValid(gRecs[0].UpdatedAt), "GetMany record UpdatedAt is not valid")
		require.False(t, nulltime.IsValid(gRecs[0].DeletedAt), "GetMany record DeletedAt is not valid")
	})

	tname = "generic when update one then ok"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)
		uRec, err = r.UpdateOne(gRec)
		require.NoError(t, err, "UpdateOne returns without error")
		require.NotNil(t, uRec, "UpdateOne returns a record")
		require.NotEmpty(t, uRec.ID, "UpdateOne record ID is not empty")
		require.Equal(t, gRec.Name, uRec.Name, "UpdateOne record Name equals expected")
		require.Equal(t, gRec.CreatedAt, uRec.CreatedAt, "UpdateOne record CreatedAt equals expected")
		zone, _ := uRec.CreatedAt.Zone()
		require.Equal(t, "UTC", zone, "UpdateOne record CreatedAt timezone equals expected")
		require.True(t, nulltime.IsValid(uRec.UpdatedAt), "UpdateOne record UpdatedAt is valid")
		zone, _ = uRec.UpdatedAt.Time.Zone()
		require.Equal(t, "UTC", zone, "UpdateOne record UpdatedAt timezone equals expected")
		require.False(t, nulltime.IsValid(uRec.DeletedAt), "UpdateOne record DeletedAt is not valid")

		currUpdatedAt := uRec.UpdatedAt
		uRec, err = r.UpdateOne(gRec)
		require.NotEqual(t, currUpdatedAt, uRec.UpdatedAt, "UpdateOne returns modified UpdatedAt")
	})

	tname = "generic when delete one then ok"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)
		err = r.DeleteOne(uRec.ID)
		require.NoError(t, err, "DeleteOne returns without error")
	})

	tname = "generic when get one deleted then error"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)
		_, err = r.GetOne(uRec.ID, nil)
		require.Error(t, err, "GetOne after delete returns with error")
	})

	tname = "generic when remove one then ok"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)
		err = r.RemoveOne(uRec.ID)
		require.NoError(t, err, "RemoveOne returns without error")
	})

	tname = "generic when create one not null constraint then error"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)

		tx, err := s.BeginTx()
		require.NoError(t, err, "BeginTx returns without error")
		defer func() {
			tx.Rollback(context.Background())
		}()

		r, err := NewGeneric[TestRecord](
			NewArgs{
				Tx:        tx,
				TableName: "test",
				Record:    TestRecord{},
			},
		)
		require.NoError(t, err, "NewGeneric returns without error")

		testRecWithErr := &TestRecord{}

		cErrRec, err := r.CreateOne(testRecWithErr)
		require.Error(t, err, "CreateOne returns with error")
		require.Nil(t, cErrRec, "CreateOne with error does not return a record")
	})

	tname = "generic when get one unknown ID then not found error"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)

		tx, err := s.BeginTx()
		require.NoError(t, err, "BeginTx returns without error")
		defer func() {
			tx.Rollback(context.Background())
		}()

		r, err := NewGeneric[TestRecord](
			NewArgs{
				Tx:        tx,
				TableName: "test",
				Record:    TestRecord{},
			},
		)
		require.NoError(t, err, "NewGeneric returns without error")

		gErrRec, err := r.GetOne(record.NewRecordID(), nil)
		require.Equal(t, pgx.ErrNoRows.Error(), err.Error(), "GetOne returns expected error")
		require.Nil(t, gErrRec, "GetOne does not return a record")
	})

	tname = "generic when get one invalid ID then invalid input error"
	t.Run(tname, func(t *testing.T) {
		t.Logf("Running %s", tname)

		tx, err := s.BeginTx()
		require.NoError(t, err, "BeginTx returns without error")
		defer func() {
			tx.Rollback(context.Background())
		}()

		r, err := NewGeneric[TestRecord](
			NewArgs{
				Tx:        tx,
				TableName: "test",
				Record:    TestRecord{},
			},
		)
		require.NoError(t, err, "NewGeneric returns without error")

		gErrRec, err := r.GetOne("notauuid", nil)
		require.True(t, strings.Contains(err.Error(), "invalid input syntax for type uuid"), "GetOne returns expected error")
		require.Nil(t, gErrRec, "GetOne does not return a record")
	})
}
