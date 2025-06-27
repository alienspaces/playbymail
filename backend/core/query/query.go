package query

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/convert"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/querier"
)

type Query struct {
	log    logger.Logger
	tx     pgx.Tx
	config Config
}

type Config struct {
	Name        string
	SQL         string
	ArrayFields set.Set[string]
}

var _ querier.Querier = &Query{}

func NewQuery(l logger.Logger, tx pgx.Tx, cfg Config) (*Query, error) {

	q := Query{
		log:    l,
		tx:     tx,
		config: cfg,
	}

	err := q.init()
	if err != nil {
		l.Error("failed new query >%v<", err)
		return nil, err
	}

	return &q, nil
}

func (q *Query) init() error {

	q.log.Debug("initialising query %s", q.Name())

	if q.tx == nil {
		return errors.New("query Tx is nil, cannot initialise")
	}

	if q.Name() == "" {
		return errors.New("query Name is empty, cannot initialise")
	}

	if q.SQL() == "" {
		return errors.New("query SQL is empty, cannot initialise")
	}

	if q.ArrayFields() == nil {
		return errors.New("repository ArrayFields is nil, cannot initialise")
	}

	return nil
}

func (q *Query) Name() string {
	return q.config.Name
}

func (q *Query) ArrayFields() set.Set[string] {
	return q.config.ArrayFields
}

func (q *Query) SQL() string {
	return q.config.SQL
}

func (q *Query) Exec(args pgx.NamedArgs) (pgconn.CommandTag, error) {
	l := q.log.WithFunctionContext("Exec")

	sql := q.SQL()

	l.Debug("SQL >%s<", sql)

	res, err := q.tx.Exec(context.Background(), sql, args)
	if err != nil {
		l.Warn("failed exec >%v<", err)
		return pgconn.CommandTag{}, err
	}

	return res, err
}

func (q *Query) Query(opts *coresql.Options) (pgx.Rows, error) {
	l := q.log.WithFunctionContext("NamedQuery")

	sql := q.SQL()

	opts, err := q.resolveOpts(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opts: sql >%s< opts >%#v< >%v<", sql, opts, err)
	}

	sql, args, err := coresql.From(sql, opts)
	if err != nil {
		q.log.Warn("failed generating query >%v<", err)
		return nil, err
	}

	l.Debug("Resulting SQL >%s< Params >%#v<", sql, args)

	rows, err := q.tx.Query(context.Background(), sql, args)
	if err != nil {
		l.Warn("failed query >%v<", err)
		return nil, err
	}

	return rows, err
}

func (q *Query) resolveOpts(opts *coresql.Options) (*coresql.Options, error) {
	if opts == nil {
		return opts, nil
	}

	for i, p := range opts.Params {
		if p.Op != "" {
			// if Op is specified, it is assumed you know what you're doing
			continue
		}

		switch t := p.Val.(type) {
		case []string:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case []int:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case time.Time:
			// Postgres stores the time in microseconds.
			// pgx requires the time to be in UTC.
			p.Val = t.Round(time.Microsecond).UTC()
		}

		isArrayField := q.ArrayFields().Has(p.Col)
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

		opts.Params[i] = p
	}

	return opts, nil
}
