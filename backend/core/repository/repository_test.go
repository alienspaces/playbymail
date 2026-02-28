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
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

type testRecord struct {
	A                string   `db:"db_A" json:"json_A"`
	B                []string `db:"db_B" json:"json_B"`
	testNestedRecord `db:"nested" json:"NESTED"`
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

type testNestedRecord struct {
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
		Record:    testRecord{},
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
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName))}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName))}),
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
			r.SetRLSIdentifiers(tt.identifiers)
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
		Record:    testRecord{},
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
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName))}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName))}),
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
			r.SetRLSIdentifiers(tt.identifiers)
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
		Record:    testRecord{},
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
		Record:    testRecord{},
	})

	require.NoError(t, err, "Repository New returns without error")

	tests := []struct {
		name          string
		identifiers   map[string][]string
		IsRLSDisabled bool
		want          set.Set[string]
	}{
		{
			name:          "RLS disabled, no identifiers",
			IsRLSDisabled: true,
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
		{
			name:          "RLS disabled, has identifiers",
			IsRLSDisabled: true,
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
		{
			name: "RLS enabled, has identifiers",
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			// appending of the RLS constraints to the SQL query is not ordered because of map iteration
			want: set.FromSlice([]string{
				strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
AND db_A IN ('1','2')
AND db_B IN ('a','b')
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName)),
				strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
AND db_B IN ('a','b')
AND db_A IN ('1','2')
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName)),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.isRLSDisabled = tt.IsRLSDisabled
			r.SetRLSIdentifiers(tt.identifiers)
			sql := r.UpdateOneSQL()
			t.Logf("SQL: %s", sql)
			require.Contains(t, tt.want, strings.TrimSpace(sql))
		})
	}
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
		Record:    testRecord{},
	})

	require.NoError(t, err, "Repository New returns without error")

	tests := []struct {
		name          string
		identifiers   map[string][]string
		IsRLSDisabled bool
		want          set.Set[string]
	}{
		{
			name:          "RLS disabled, no identifiers",
			IsRLSDisabled: true,
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
		{
			name:          "RLS disabled, has identifiers",
			IsRLSDisabled: true,
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
		{
			name: "RLS enabled, has identifiers",
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			// appending of the RLS constraints to the SQL query is not ordered because of map iteration
			want: set.FromSlice([]string{
				strings.TrimSpace(fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL
AND db_A IN ('1','2')
AND db_B IN ('a','b')
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName)),
				strings.TrimSpace(fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL
AND db_B IN ('a','b')
AND db_A IN ('1','2')
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName)),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.isRLSDisabled = tt.IsRLSDisabled
			r.SetRLSIdentifiers(tt.identifiers)
			require.Contains(t, tt.want, strings.TrimSpace(r.DeleteOneSQL()))
		})
	}
}

func Test_SetRLSConstraints_GetOneSQL(t *testing.T) {
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
		Record:    testRecord{},
	})

	require.NoError(t, err, "Repository New returns without error")

	tests := []struct {
		name          string
		constraints   []repositor.RLSConstraint
		identifiers   map[string][]string
		IsRLSDisabled bool
		want          set.Set[string]
	}{
		{
			name:          "RLS disabled, constraints ignored",
			IsRLSDisabled: true,
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName))}),
		},
		{
			name: "constraint with matching column, all params available",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND (db_A IN (SELECT id FROM other_table WHERE account_id = '123'))", r.tableName))}),
		},
		{
			name: "constraint with non-matching column, skipped",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "non_existent_column",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName))}),
		},
		{
			name: "constraint with missing required params, fail-safe applied",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"other_id": {"456"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND 1=0", r.tableName))}),
		},
		{
			name: "constraint with missing required params, no identifiers, skipped",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{},
			want:        set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL", r.tableName))}),
		},
		{
			name: "multiple constraints, all applied",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM table_a WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
				{
					Column:                 "db_B",
					SQLTemplate:            "IN (SELECT id FROM table_b WHERE game_id = :game_id)",
					RequiredRLSIdentifiers: []string{"game_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
				"game_id":    {"456"},
			},
			// Order may vary due to slice iteration
			want: set.FromSlice([]string{
				strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND (db_A IN (SELECT id FROM table_a WHERE account_id = '123'))\nAND (db_B IN (SELECT id FROM table_b WHERE game_id = '456'))", r.tableName)),
				strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND (db_B IN (SELECT id FROM table_b WHERE game_id = '456'))\nAND (db_A IN (SELECT id FROM table_a WHERE account_id = '123'))", r.tableName)),
			}),
		},
		{
			name: "constraint with multiple required params",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id AND game_id = :game_id)",
					RequiredRLSIdentifiers: []string{"account_id", "game_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
				"game_id":    {"456"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND (db_A IN (SELECT id FROM other_table WHERE account_id = '123' AND game_id = '456'))", r.tableName))}),
		},
		{
			name: "constraint with multiple required params, one missing, fail-safe",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id AND game_id = :game_id)",
					RequiredRLSIdentifiers: []string{"account_id", "game_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND 1=0", r.tableName))}),
		},
		{
			name: "NULL-handling constraint skips identifier IN filter for same column",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IS NULL OR db_A IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
				"db_A":       {"val1"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND (db_A IS NULL OR db_A IN (SELECT id FROM other_table WHERE account_id = '123'))", r.tableName))}),
		},
		{
			name: "non-NULL constraint does not skip identifier IN filter for same column",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
				"db_A":       {"val1"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL\nAND db_A IN ('val1')\nAND (db_A IN (SELECT id FROM other_table WHERE account_id = '123'))", r.tableName))}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset repository state for each test
			r.isRLSDisabled = tt.IsRLSDisabled
			// Clear constraints - SetRLSConstraints will set them if they match
			r.rlsConstraints = nil
			r.rlsIdentifiers = map[string][]string{}
			r.hadRLSIdentifiers = false

			// Set constraints and identifiers (order no longer matters)
			if len(tt.constraints) > 0 {
				r.SetRLSConstraints(tt.constraints)
			} else {
				r.rlsConstraints = []repositor.RLSConstraint{}
			}
			r.SetRLSIdentifiers(tt.identifiers)

			t.Logf("Test: %s", tt.name)
			t.Logf("Constraints set: %d", len(tt.constraints))
			t.Logf("Constraints stored: %d", len(r.rlsConstraints))
			t.Logf("Identifiers input: %v", tt.identifiers)
			t.Logf("Identifiers stored: %v", r.rlsIdentifiers)
			t.Logf("Repository attributes: %v", r.Attributes())

			actualSQL := strings.TrimSpace(r.GetOneSQL())
			t.Logf("Actual SQL: %q", actualSQL)
			t.Logf("Expected set: %v", tt.want)
			require.Contains(t, tt.want, actualSQL)
		})
	}
}

func Test_SetRLSConstraints_UpdateOneSQL(t *testing.T) {
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
		Record:    testRecord{},
	})

	require.NoError(t, err, "Repository New returns without error")

	tests := []struct {
		name        string
		constraints []repositor.RLSConstraint
		identifiers map[string][]string
		want        set.Set[string]
	}{
		{
			name: "constraint applied before RETURNING",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
AND (db_A IN (SELECT id FROM other_table WHERE account_id = '123'))
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
		{
			name: "multiple constraints before RETURNING",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "db_A",
					SQLTemplate:            "IN (SELECT id FROM table_a WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
				{
					Column:                 "db_B",
					SQLTemplate:            "IN (SELECT id FROM table_b WHERE game_id = :game_id)",
					RequiredRLSIdentifiers: []string{"game_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
				"game_id":    {"456"},
			},
			// Order may vary due to slice iteration
			want: set.FromSlice([]string{
				strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
AND (db_A IN (SELECT id FROM table_a WHERE account_id = '123'))
AND (db_B IN (SELECT id FROM table_b WHERE game_id = '456'))
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName)),
				strings.TrimSpace(fmt.Sprintf(`
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
updated_at = @updated_at
WHERE id = @id
AND   deleted_at IS NULL
AND (db_B IN (SELECT id FROM table_b WHERE game_id = '456'))
AND (db_A IN (SELECT id FROM table_a WHERE account_id = '123'))
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName)),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.SetRLSConstraints(tt.constraints)
			r.SetRLSIdentifiers(tt.identifiers)
			sql := r.UpdateOneSQL()
			t.Logf("SQL: %s", sql)
			require.Contains(t, tt.want, strings.TrimSpace(sql))
		})
	}
}

func Test_SetRLSConstraints_DeleteOneSQL(t *testing.T) {
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
		Record:    testRecord{},
	})

	require.NoError(t, err, "Repository New returns without error")

	tests := []struct {
		name        string
		constraints []repositor.RLSConstraint
		identifiers map[string][]string
		want        set.Set[string]
	}{
		{
			name: "constraint applied before RETURNING",
			constraints: []repositor.RLSConstraint{
				{
					Column:                 "id",
					SQLTemplate:            "IN (SELECT id FROM other_table WHERE account_id = :account_id)",
					RequiredRLSIdentifiers: []string{"account_id"},
				},
			},
			identifiers: map[string][]string{
				"account_id": {"123"},
			},
			want: set.FromSlice([]string{strings.TrimSpace(fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL
AND (id IN (SELECT id FROM other_table WHERE account_id = '123'))
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time, id, created_at, updated_at, deleted_at
`, r.tableName))}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.SetRLSConstraints(tt.constraints)
			r.SetRLSIdentifiers(tt.identifiers)
			sql := r.DeleteOneSQL()
			t.Logf("SQL: %s", sql)
			require.Contains(t, tt.want, strings.TrimSpace(sql))
		})
	}
}
