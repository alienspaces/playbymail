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
	tx                pgx.Tx
	tableName         string
	attributes        []string
	setAttributes     []string
	attributeIndex    set.Set[string]
	arrayFields       set.Set[string]
	isRLSDisabled     bool
	rlsIdentifiers    map[string][]string
	rlsConstraints    []repositor.RLSConstraint
	hadRLSIdentifiers bool // Track if identifiers were ever set (to distinguish test harness from API mode)
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
		tx:             args.Tx,
		tableName:      args.TableName,
		attributes:     tag.GetFieldTagValues(args.Record, "db"),
		arrayFields:    tag.GetArrayFieldTagValues(args.Record, "db"),
		isRLSDisabled:  args.IsRLSDisabled,
		rlsIdentifiers: make(map[string][]string),
		rlsConstraints: []repositor.RLSConstraint{},
	}

	setAttributes := []string{}
	for _, attr := range r.attributes {
		if attr == "id" || attr == "created_at" || attr == "deleted_at" {
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

	// fmt.Printf("**** GetRows Pre-Resolve Params >%v< Attrs >%v<\n", opts.Params, r.attributeIndex)
	opts, err = r.resolveOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opts: sql >%s< opts >%#v< >%v<", sql, opts, err)
	}
	// fmt.Printf("**** GetRows Post-Resolve Params >%v<\n", opts.Params)

	sql, args, err := coresql.From(sql, opts)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("**** GetRows SQL >%q< Args >%v<\n", sql, args)

	// start := time.Now()
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

	sql := r.withRLS(fmt.Sprintf(`
UPDATE %s SET
%s
WHERE id = @id
AND   deleted_at IS NULL
`,
		r.TableName(),
		setAttributeAndValuePlaceholders(r.SetAttributes())))

	return sql + fmt.Sprintf("RETURNING %s\n", strings.Join(r.Attributes(), ", "))
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
	sql := r.withRLS(fmt.Sprintf(`
UPDATE %s SET deleted_at = @deleted_at WHERE id = @id AND deleted_at IS NULL
`, r.TableName()))
	return sql + fmt.Sprintf("RETURNING %s\n", strings.Join(r.Attributes(), ", "))
}

// RemoveOneSQL generates a physical delete SQL statement
func (r *Repository) RemoveOneSQL() string {
	return fmt.Sprintf(`
DELETE FROM %s WHERE id = @id
`, r.TableName())
}

func (r *Repository) SetRLSIdentifiers(identifiers map[string][]string) {
	if r.isRLSDisabled {
		return
	}

	// Store all provided identifiers - filtering for direct use happens in withRLS
	// This allows SetRLSIdentifiers to be called before SetRLSConstraints
	if len(identifiers) > 0 {
		r.rlsIdentifiers = identifiers
		r.hadRLSIdentifiers = true
	}
}

// SetRLSConstraints sets the RLS constraints that will be automatically applied
// when the repository record has matching column names.
func (r *Repository) SetRLSConstraints(constraints []repositor.RLSConstraint) {
	if r.isRLSDisabled {
		return
	}

	filtered := []repositor.RLSConstraint{}

	// Only include constraints for columns that exist in this repository's record
	attributeSet := set.FromSlice(r.Attributes())

	for _, constraint := range constraints {
		if attributeSet.Has(constraint.Column) {
			filtered = append(filtered, constraint)
		}
	}

	r.rlsConstraints = filtered
}

func (r *Repository) withRLS(sql string) string {

	// Early return if RLS is disabled or no RLS data is set
	if r.isRLSDisabled || (len(r.rlsIdentifiers) == 0 && len(r.rlsConstraints) == 0) {
		return sql
	}

	// Apply RLS identifier filters (direct column filters)
	// Only apply filters for identifiers that correspond to actual columns in the record
	attributeSet := set.FromSlice(r.Attributes())
	for attr, ids := range r.rlsIdentifiers {
		// Skip identifiers that don't correspond to record columns
		// (they may be kept for constraint substitution but shouldn't be used as direct filters)
		if !attributeSet.Has(attr) {
			continue
		}
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

	// Apply RLS constraint filters (derived column filters via SQL subqueries)
	for _, constraint := range r.rlsConstraints {
		// Substitute RLS constraint placeholders with RLS identifier values
		// Skip constraint if any required RLS identifier is missing or empty
		constraintSQL := constraint.SQLTemplate
		allRequiredIdentifiersSubstituted := true
		for _, requiredRLSIdentifier := range constraint.RequiredRLSIdentifiers {
			if paramValues, ok := r.rlsIdentifiers[requiredRLSIdentifier]; ok && len(paramValues) > 0 {
				// Substitute the first value from RLS identifiers
				placeholder := fmt.Sprintf(":%s", requiredRLSIdentifier)
				constraintSQL = strings.ReplaceAll(constraintSQL, placeholder, fmt.Sprintf("'%s'", paramValues[0]))
			} else {
				allRequiredIdentifiersSubstituted = false
				break
			}
		}

		if !allRequiredIdentifiersSubstituted {
			// Skip constraint if required RLS identifiers not available
			// Only apply fail-safe if we had identifiers provided (API mode)
			// This ensures we don't return data when we can't guarantee constraint safety
			// but allows test harness mode (no identifiers) to work without fail-safe
			if len(constraint.RequiredRLSIdentifiers) > 0 && r.hadRLSIdentifiers {
				// Fail-safe: required RLS identifiers not available
				// We can't guarantee they'll be available, so return no data
				sql += "AND 1=0\n"
			}
			continue
		}

		// Only apply constraint if all required RLS identifiers were successfully substituted
		if allRequiredIdentifiersSubstituted {
			// Verify no placeholders remain (safety check)
			if strings.Contains(constraintSQL, ":") {
				// Still has placeholders - skip to avoid SQL errors
				continue
			}
			// Append constraint SQL to WHERE clause
			sql += fmt.Sprintf("AND %s %s\n", constraint.Column, constraintSQL)
		}
	}

	return sql
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
