package domain

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/querier"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
)

// Domain -
type Domain struct {
	Log                    logger.Logger
	Repositories           map[string]repositor.Repositor
	RepositoryConstructors []RepositoryConstructor
	Queries                map[string]querier.Querier
	Tx                     pgx.Tx
	Err                    error

	// composable functions
	QueriesFunc func(tx pgx.Tx) ([]querier.Querier, error)
	RLSFunc     func(identifiers map[string][]string)
}

type RepositoryConstructor func(logger.Logger, pgx.Tx) (repositor.Repositor, error)

var _ domainer.Domainer = &Domain{}

// NewDomain - intended for testing only, maybe move into test files..
func NewDomain(l logger.Logger, repositoryConstructors []RepositoryConstructor) (m *Domain, err error) {

	if l == nil {
		return nil, fmt.Errorf("failed new domain, missing logger")
	}

	l, err = l.NewInstance()
	if err != nil {
		return nil, err
	}

	m = &Domain{
		Log: l.WithPackageContext("domain"),
		Err: nil,
	}

	// Repository constructors are used to create repositories for the domain
	m.RepositoryConstructors = repositoryConstructors

	// Default function for setting RLS constraints on repositories
	m.RLSFunc = m.SetRLS

	return m, nil
}

// Init -
func (m *Domain) Init(tx pgx.Tx) (err error) {

	if tx == nil {
		msg := "failed init, tx is required"
		m.Log.Warn(msg)
		return errors.New(msg)
	}

	if m.QueriesFunc == nil {
		m.QueriesFunc = m.NewQueries
	}

	if m.RLSFunc == nil {
		m.RLSFunc = m.SetRLS
	}

	m.Tx = tx

	// repositories
	m.Repositories, err = m.NewRepositories(tx)
	if err != nil {
		m.Log.Warn("(core) failed repositories func >%v<", err)
		return err
	}

	// queries
	queries, err := m.QueriesFunc(tx)
	if err != nil {
		m.Log.Warn("(core) failed queries func >%v<", err)
		return err
	}

	m.Queries = make(map[string]querier.Querier)
	for _, q := range queries {
		m.Queries[q.Name()] = q
	}

	return nil
}

// NewQueries returns a list of queries that will be used
func (m *Domain) NewQueries(tx pgx.Tx) ([]querier.Querier, error) {
	return nil, nil
}

// NewRepositories returns a list of repositories that will be used
func (m *Domain) NewRepositories(tx pgx.Tx) (map[string]repositor.Repositor, error) {

	repositories := map[string]repositor.Repositor{}
	for _, newRepo := range m.RepositoryConstructors {
		r, err := newRepo(m.Log, tx)
		if err != nil {
			m.Log.Warn("(core) failed new %s repository >%v<", r.TableName(), err)
			return nil, err
		}
		repositories[r.TableName()] = r
	}

	return repositories, nil
}

// SetRLS -
func (m *Domain) SetRLS(identifiers map[string][]string) {

	// We'll be resetting the "id" key when we use the map
	ri := maps.Clone(identifiers)

	for tableName := range m.Repositories {

		// When the repository table name matches an RLS identifier key, we apply the
		// RLS constraints to the "id" column to enforce any RLS constraints on itself!
		// Can this be done inside repository core code on itself? Absolutely... but it
		// would be making a naive assumption about conventions. This project's convention
		// is to name foreign key columns according to the table name it foreign keys to.
		// If that convention is not followed, then the following block would not work.
		if _, ok := ri[tableName+"_id"]; ok {
			ri["id"] = ri[tableName+"_id"]
			m.Repositories[tableName].SetRLS(ri)
			continue
		}
		m.Repositories[tableName].SetRLS(identifiers)
	}
}

// SetTxLockTimeout -
func (m *Domain) SetTxLockTimeout(timeoutSecs float64) error {
	if m.Tx == nil {
		err := fmt.Errorf("cannot set transaction lock timeout seconds, database Tx is nil")
		return err
	}

	// If we SET, instead of SET LOCAL, lock_timeout would be at the session-level.
	// Since we use connection pooling, this would mean that different sessions (and therefore requests)
	// would have different, unknown lock_timeout parameters.

	timeoutMs := timeoutSecs * 1000
	_, err := m.Tx.Exec(context.TODO(), fmt.Sprintf("SET LOCAL lock_timeout = %d", int(timeoutMs)))
	if err != nil {
		err = fmt.Errorf("failed setting transaction lock timeout seconds %w", err)
		return err
	}

	m.Log.Debug("lock timeout seconds set to >%fs<", timeoutSecs)

	return nil
}

// Commit commits the current model database transaction
func (m *Domain) Commit() error {
	if m.Tx != nil {
		err := m.Tx.Commit(context.TODO())
		if err != nil {
			m.Log.Warn("(core) failed committing transaction >%v<", err)
		}
		return nil
	}
	err := fmt.Errorf("cannot commit transaction, database Tx is nil")
	return err
}

// Rollback rolls back the current model database transaction
func (m *Domain) Rollback() error {
	if m.Tx != nil {
		return m.Tx.Rollback(context.TODO())
	}
	err := fmt.Errorf("cannot rollback transaction, database Tx is nil")
	return err
}
