// Package repository provides methods for interacting with the database
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/tag"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
)

// Repository -
type Repository struct {
	tx             pgx.Tx
	tableName      string
	attributes     []string
	setAttributes  []string
	attributeIndex set.Set[string]
	arrayFields    set.Set[string]
	isRLSDisabled  bool
	rlsIdentifiers map[string][]string
}

var _ repositor.Repositor = &Repository{}

type NewArgs struct {
	Tx            pgx.Tx
	TableName     string
	Record        any
	IsRLSDisabled bool
}

// If there is an error initialising the repository, the returned repository
// will be non-zero but improperly initialised.
func New(args NewArgs) (Repository, error) {

	if args.TableName == "" {
		return Repository{}, fmt.Errorf("failed new repository, missing arg TableName")
	}
	if args.Record == nil {
		return Repository{}, fmt.Errorf("failed new repository, missing arg Record")
	}

	r := Repository{
		tx:            args.Tx,
		tableName:     args.TableName,
		attributes:    tag.GetFieldTagValues(args.Record, "db"),
		arrayFields:   tag.GetArrayFieldTagValues(args.Record, "db"),
		isRLSDisabled: args.IsRLSDisabled,
	}

	setAttributes := []string{}
	for _, attr := range r.attributes {
		if attr == "created_at" || attr == "deleted_at" {
			continue
		}
		setAttributes = append(setAttributes, attr)
	}
	r.setAttributes = setAttributes

	if r.tx == nil {
		return Repository{}, errors.New("repository Tx is nil, cannot initialise")
	}

	if r.TableName() == "" {
		return Repository{}, errors.New("repository TableName is empty, cannot initialise")
	}

	if len(r.Attributes()) == 0 {
		return Repository{}, errors.New("repository Attributes are empty, cannot initialise")
	}

	if r.ArrayFields() == nil {
		return Repository{}, errors.New("repository ArrayFields is nil, cannot initialise")
	}

	attributeIndex := map[string]struct{}{}
	for _, attribute := range r.Attributes() {
		attributeIndex[attribute] = struct{}{}
	}
	r.attributeIndex = attributeIndex

	return r, nil
}

func (r *Repository) TableName() string {
	return r.tableName
}

func (r *Repository) Attributes() []string {
	return r.attributes
}

func (r *Repository) SetAttributes() []string {
	return r.setAttributes
}

func (r *Repository) ArrayFields() set.Set[string] {
	return r.arrayFields
}

func (r *Repository) Tx() pgx.Tx {
	return r.tx
}

// DeleteOne -
func (r *Repository) DeleteOne(id any) error {

	params := pgx.NamedArgs{
		"id":         id,
		"deleted_at": record.NewRecordNullTimestamp(),
	}

	res, err := r.tx.Exec(context.Background(), r.DeleteOneSQL(), params)
	if err != nil {
		return err
	}

	raf := res.RowsAffected()
	if raf != 1 {
		return fmt.Errorf("expecting to delete exactly one row but deleted >%d<", raf)
	}

	return nil
}

// RemoveOne -
func (r *Repository) RemoveOne(id any) error {

	params := pgx.NamedArgs{
		"id": id,
	}

	res, err := r.tx.Exec(context.Background(), r.RemoveOneSQL(), params)
	if err != nil {
		return err
	}

	raf := res.RowsAffected()
	if raf != 1 {
		return fmt.Errorf("expecting to remove exactly one row but removed >%d<", raf)
	}

	return nil
}

// GetRows returns the rows result from the provided SQL and options
func (r *Repository) GetRows(sql string, opts *coresql.Options) (rows pgx.Rows, err error) {

	opts, err = r.resolveOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opts: sql >%s< opts >%#v< >%v<", sql, opts, err)
	}

	sql, args, err := coresql.From(sql, opts)
	if err != nil {
		return nil, err
	}

	rows, err = r.tx.Query(context.Background(), sql, args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// GetOneSQL - This SQL statement ends with a newline so that any parameters
// can be easily appended.
func (r *Repository) GetOneSQL() string {
	return r.withRLS(fmt.Sprintf(`
SELECT %s FROM %s WHERE deleted_at IS NULL
`,
		strings.Join(r.Attributes(), ", "),
		r.TableName()))
}

// GetManySQL - This SQL statement ends with a newline so that any parameters
// can be easily appended.
func (r *Repository) GetManySQL() string {
	return r.withRLS(fmt.Sprintf(`
SELECT %s FROM %s WHERE deleted_at IS NULL
`,
		strings.Join(r.Attributes(), ", "),
		r.TableName()))
}

func (r *Repository) withRLS(sql string) string {

	for attr, ids := range r.rlsIdentifiers {
		var strBuilder strings.Builder
		for i, id := range ids {
			strBuilder.WriteString("'")
			strBuilder.WriteString(id)

			if i != len(ids)-1 {
				strBuilder.WriteString("',")
			} else {
				strBuilder.WriteString("'")
			}
		}
		sql += fmt.Sprintf("AND %s IN (%s)\n", attr, strBuilder.String())
	}

	return sql
}

// CreateOneSQL generates an insert SQL statements
func (r *Repository) CreateOneSQL() string {
	return fmt.Sprintf(`
INSERT INTO %s (
%s
) VALUES (
%s
)
RETURNING %s
`,
		r.TableName(),
		strings.Join(r.Attributes(), ",\n"),
		insertValuePlaceholders(r.Attributes()),
		strings.Join(r.Attributes(), ", "))
}

func insertValuePlaceholders(attributes []string) string {
	var strBuilder strings.Builder

	for i, attr := range attributes {
		strBuilder.WriteString(fmt.Sprintf("@%s", attr))
		if i != len(attributes)-1 {
			strBuilder.WriteString(",\n")
		}
	}

	return strBuilder.String()
}

// UpdateOneSQL generates an update SQL statement
func (r *Repository) UpdateOneSQL() string {

	return fmt.Sprintf(`
UPDATE %s SET
%s
WHERE id = @id
AND   deleted_at IS NULL
RETURNING %s
`,
		r.TableName(),
		setAttributeAndValuePlaceholders(r.SetAttributes()),
		strings.Join(r.Attributes(), ", "))
}

func setAttributeAndValuePlaceholders(attributes []string) string {
	var strBuilder strings.Builder

	for i, attr := range attributes {
		strBuilder.WriteString(attr)
		strBuilder.WriteString(" = ")
		strBuilder.WriteString(fmt.Sprintf("@%s", attr))

		if i != len(attributes)-1 {
			strBuilder.WriteString(",\n")
		}
	}

	return strBuilder.String()
}

// DeleteOneSQL generates a logical delete SQL statement
func (r *Repository) DeleteOneSQL() string {
	return fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL RETURNING %s
`, r.TableName(), strings.Join(r.Attributes(), ", "))
}

// RemoveOneSQL generates a physical delete SQL statement
func (r *Repository) RemoveOneSQL() string {
	return fmt.Sprintf(`
DELETE FROM %s WHERE id = @id
`, r.TableName())
}

func (r *Repository) SetRLS(identifiers map[string][]string) {
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

func (r *Repository) resolveOptions(opts *coresql.Options) (*coresql.Options, error) {
	if opts == nil {
		return opts, nil
	}

	params := []coresql.Param{}

	for _, p := range opts.Params {

		// Skip parameters that aren't valid attributes for the record
		if _, ok := r.attributeIndex[p.Col]; !ok {
			continue
		}

		switch t := p.Val.(type) {
		case []string:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case []int:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case []int16:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case []int32:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case []int64:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case time.Time:
			// Postgres stores the time in microseconds.
			// pgx requires the time to be in UTC.
			p.Val = t.Round(time.Microsecond).UTC()
		}

		// if Op is specified, it is assumed you know what you're doing
		if p.Op != "" {
			params = append(params, p)
			continue
		}

		isArrayField := r.ArrayFields().Has(p.Col)
		if isArrayField {
			if len(p.Array) > 0 {
				p.Op = coresql.OpContains
			} else {
				p.Op = coresql.OpAny
			}
		} else {
			if len(p.Array) > 0 {
				p.Op = coresql.OpIn
			} else {
				p.Op = coresql.OpEqual
			}
		}

		params = append(params, p)
	}

	opts.Params = params

	return opts, nil
}
