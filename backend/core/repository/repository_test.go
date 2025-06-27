package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

type testNestedMultiTag struct {
	A          string   `db:"db_A" json:"json_A"`
	B          []string `db:"db_B" json:"json_B"`
	testNested `db:"nested" json:"NESTED"`
	//lint:ignore U1000 intended to be unused?
	testEmptyNested        `db:"empty_nested" json:"EMPTY_NESTED"`
	testSingleNestedDB     `db:"single_nested_db" json:"SINGLE_NESTED_DB"`
	testSingleNestedJSON   `db:"single_nested_json" json:"SINGLE_NESTED_JSON"`
	singleNestedNullTime   `db:"single_nested_null_time" json:"SINGLE_NESTED_NULL_TIME"`
	singleNestedNullString `db:"single_nested_null_string" json:"SINGLE_NESTED_NULL_STRING"`
	//lint:ignore U1000 intended to be unused?
	f int `db:"db_f"`
	G int `json:"json_G"`
	//lint:ignore U1000 intended to be unused?
	h int
	//lint:ignore U1000 intended to be unused?
	time time.Time `db:"db_time"`
	record.Record
}

type testNested struct {
	C string `json:"json_C"`
	//lint:ignore U1000 intended to be unused?
	d  []string `db:"db_d"`
	E  string   `db:"db_E" json:"json_E"`
	E2 int
}

type testEmptyNested struct{}

type testSingleNestedDB struct {
	E3 int `db:"db_E3"`
}

type testSingleNestedJSON struct {
	E4 int `json:"json_E4"`
}

type singleNestedNullTime struct {
	E5 sql.NullTime `db:"db_E5"`
}

type singleNestedNullString struct {
	E6 sql.NullString `db:"db_E6"`
}

func newDependencies() (logger.Logger, storer.Storer, error) {

	cfg := config.Config{}
	err := config.Parse(&cfg)
	if err != nil {
		return nil, nil, err
	}

	// logger
	l, err := log.NewLogger(cfg)
	if err != nil {
		return nil, nil, err
	}

	// storer
	s, err := store.NewStore(cfg)
	if err != nil {
		return nil, nil, err
	}

	return l, s, nil
}

func Test_GetOneSQL(t *testing.T) {
	_, s, err := newDependencies()
	require.NoError(t, err, "NewDependencies returns without error")
	defer func() {
		err = s.ClosePool()
		require.NoError(t, err, "ClosePool returns without error")
	}()

	tx, err := s.BeginTx()
	require.NoError(t, err, "BeginTx returns without error")
	defer func() {
		tx.Rollback(context.Background())
	}()

	r, err := New(NewArgs{
		Tx:        tx,
		TableName: "test",
		Record:    testNestedMultiTag{},
	})

	require.NoError(t, err, "Repository Init returns without error")

	tests := []struct {
		name          string
		identifiers   map[string][]string
		IsRLSDisabled bool
		want          set.Set[string]
	}{
		{
			name:          "RLS disabled, no identifiers",
			IsRLSDisabled: true,
			want:          set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName)}),
		},
		{
			name:          "RLS disabled, has identifiers",
			IsRLSDisabled: true,
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName)}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName)}),
		},
		{
			name: "RLS enabled, has identifiers",
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			// appending of the RLS constraints to the SQL query is not ordered because of map iteration
			want: set.FromSlice([]string{
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND db_A IN ('1','2')\nAND db_B IN ('a','b')",
					r.tableName),
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND db_B IN ('a','b')\nAND db_A IN ('1','2')",
					r.tableName),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.isRLSDisabled = tt.IsRLSDisabled
			r.SetRLS(tt.identifiers)
			require.Contains(t, tt.want, strings.TrimSpace(r.GetOneSQL()))
		})
	}
}

func Test_GetManySQL(t *testing.T) {
	_, s, err := newDependencies()
	require.NoError(t, err, "NewDependencies returns without error")
	defer func() {
		err = s.ClosePool()
		require.NoError(t, err, "ClosePool returns without error")
	}()

	tx, err := s.BeginTx()
	require.NoError(t, err, "BeginTx returns without error")
	defer func() {
		tx.Rollback(context.Background())
	}()

	r, err := New(NewArgs{
		Tx:        tx,
		TableName: "test",
		Record:    testNestedMultiTag{},
	})

	require.NoError(t, err, "Repository New returns without error")

	sqlTests := []struct {
		name          string
		identifiers   map[string][]string
		IsRLSDisabled bool
		want          set.Set[string]
	}{
		{
			name:          "RLS disabled, no identifiers",
			IsRLSDisabled: true,
			want:          set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName)}),
		},
		{
			name:          "RLS disabled, has identifiers",
			IsRLSDisabled: true,
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName)}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName)}),
		},
		{
			name: "RLS enabled, has identifiers",
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			// appending of the RLS constraints to the SQL query is not ordered because of map iteration
			want: set.FromSlice([]string{
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND db_A IN ('1','2')\nAND db_B IN ('a','b')",
					r.tableName),
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND db_B IN ('a','b')\nAND db_A IN ('1','2')",
					r.tableName),
			}),
		},
	}

	for _, tt := range sqlTests {
		t.Run(tt.name, func(t *testing.T) {
			r.isRLSDisabled = tt.IsRLSDisabled
			r.SetRLS(tt.identifiers)
			require.Contains(t, tt.want, strings.TrimSpace(r.GetManySQL()))
		})
	}

	type args struct {
		opts *coresql.Options
	}
	optTests := []struct {
		name string
		args args
		want *coresql.Options
	}{
		{
			name: "@> for multi value, @> for nested slice field for single value, IN with multi value int slice, = for string",
			args: args{
				opts: &coresql.Options{
					Params: []coresql.Param{
						{
							Col: "db_B",
							Val: []string{"a", "b"},
						},
						{
							Col: "db_d",
							Val: []string{"c"},
						},
						{
							Col: "db_f",
							Val: []int{1, 2},
						},
						{
							Col: "db_E6",
							Val: "a",
						},
					},
				},
			},
			want: &coresql.Options{
				Params: []coresql.Param{
					{
						Col:   "db_B",
						Op:    coresql.OpContains,
						Array: convert.GenericSlice([]string{"a", "b"}),
					},
					{
						Col:   "db_d",
						Op:    coresql.OpContains,
						Array: convert.GenericSlice([]string{"c"}),
					},
					{
						Col:   "db_f",
						Op:    coresql.OpIn,
						Array: convert.GenericSlice([]int{1, 2}),
					},
					{
						Col: "db_E6",
						Op:  coresql.OpEqual,
						Val: "a",
					},
				},
			},
		},
		{
			name: "ANY, ANY for nested slice field, = for int, IN with single value string slice",
			args: args{
				opts: &coresql.Options{
					Params: []coresql.Param{
						{
							Col: "db_B",
							Val: "a",
						},
						{
							Col: "db_d",
							Val: "c",
						},
						{
							Col: "db_f",
							Val: 1,
						},
						{
							Col: "db_E6",
							Val: []string{"a"},
						},
					},
				},
			},
			want: &coresql.Options{
				Params: []coresql.Param{
					{
						Col: "db_B",
						Op:  coresql.OpAny,
						Val: "a",
					},
					{
						Col: "db_d",
						Op:  coresql.OpAny,
						Val: "c",
					},
					{
						Col: "db_f",
						Op:  coresql.OpEqual,
						Val: 1,
					},
					{
						Col:   "db_E6",
						Op:    coresql.OpIn,
						Array: convert.GenericSlice([]string{"a"}),
					},
				},
			},
		},
	}
	for _, tt := range optTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.resolveOptions(tt.args.opts)
			require.NoError(t, err, "resolveOptions should not err")
			require.Equal(t, tt.want, got, "resolveOptions should return coresql.Options as expected")
		})
	}
}

func Test_CreateOneSQL(t *testing.T) {
	_, s, err := newDependencies()
	require.NoError(t, err, "NewDependencies returns without error")
	defer func() {
		err = s.ClosePool()
		require.NoError(t, err, "ClosePool returns without error")
	}()

	tx, err := s.BeginTx()
	require.NoError(t, err, "BeginTx returns without error")
	defer func() {
		tx.Rollback(context.Background())
	}()

	r, err := New(NewArgs{
		Tx:        tx,
		TableName: "test",
		Record:    testNestedMultiTag{},
	})

	require.NoError(t, err, "Repository New returns without error")

	require.Equal(t, fmt.Sprintf(`
INSERT INTO %s (
db_A,
db_B,
db_d,
db_E,
db_E3,
db_E5,
db_E6,
db_f,
db_time,
id,
created_at,
updated_at,
deleted_at
) VALUES (
@db_A,
@db_B,
@db_d,
@db_E,
@db_E3,
@db_E5,
@db_E6,
@db_f,
@db_time,
@id,
@created_at,
@updated_at,
@deleted_at
)
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName), r.CreateOneSQL())
}

func Test_UpdateOneSQL(t *testing.T) {
	_, s, err := newDependencies()
	require.NoError(t, err, "NewDependencies returns without error")
	defer func() {
		err = s.ClosePool()
		require.NoError(t, err, "ClosePool returns without error")
	}()

	tx, err := s.BeginTx()
	require.NoError(t, err, "BeginTx returns without error")
	defer func() {
		tx.Rollback(context.Background())
	}()

	r, err := New(NewArgs{
		Tx:        tx,
		TableName: "test",
		Record:    testNestedMultiTag{},
	})

	require.NoError(t, err, "Repository New returns without error")

	require.Equal(t, fmt.Sprintf(`
UPDATE %s SET
db_A = @db_A,
db_B = @db_B,
db_d = @db_d,
db_E = @db_E,
db_E3 = @db_E3,
db_E5 = @db_E5,
db_E6 = @db_E6,
db_f = @db_f,
db_time = @db_time,
id = @id,
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName), r.UpdateOneSQL())
}

func Test_DeleteOneSQL(t *testing.T) {
	_, s, err := newDependencies()
	require.NoError(t, err, "NewDependencies returns without error")
	defer func() {
		err = s.ClosePool()
		require.NoError(t, err, "ClosePool returns without error")
	}()

	tx, err := s.BeginTx()
	require.NoError(t, err, "BeginTx returns without error")
	defer func() {
		tx.Rollback(context.Background())
	}()

	r, err := New(NewArgs{
		Tx:        tx,
		TableName: "test",
		Record:    testNestedMultiTag{},
	})

	require.NoError(t, err, "Repository New returns without error")

	require.Equal(t, fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName), r.DeleteOneSQL())
}
