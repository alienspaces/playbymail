package harness

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

// CreateDataFunc - callback function that creates test data
type CreateDataFunc func() error

// RemoveDataFunc - callback function that removes test data
type RemoveDataFunc func() error

// Testing -
type Testing struct {
	Log       logger.Logger
	Store     storer.Storer
	JobClient *river.Client[pgx.Tx]
	Domain    domainer.Domainer

	// ShouldCommitData is used to determine whether Setup and Teardown should commit data to the DB.
	// This should only be true if changes in one transaction must be visible in another (e.g., handler tests).
	ShouldCommitData bool

	// domainer function
	DomainFunc func() (domainer.Domainer, error)

	// Composable functions
	CreateDataFunc CreateDataFunc
	RemoveDataFunc RemoveDataFunc

	tx pgx.Tx
}

// NewTesting -
func NewTesting(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (t *Testing, err error) {

	// Require logger and store
	if l == nil || s == nil {
		return nil, fmt.Errorf("missing logger >%v< or storer >%v<, cannot create new test harness", l, s)
	}

	t = &Testing{
		Log:       l,
		Store:     s,
		JobClient: j,
	}
	return t, nil
}

// Init -
func (t *Testing) Init() (err error) {
	cfg := config.Config{}
	err = config.Parse(&cfg)
	if err != nil {
		return err
	}

	// logger
	if t.Log == nil {
		t.Log, err = log.NewLogger(cfg)
		if err != nil {
			return err
		}
	}
	t.Log = t.Log.WithApplicationContext("harness")

	// storer
	if t.Store == nil {
		t.Store, err = store.NewStore(cfg)
		if err != nil {
			return err
		}
	}

	// domainer
	if t.DomainFunc != nil {
		t.Domain, err = t.DomainFunc()
		if err != nil {
			t.Log.Warn("failed new domainer >%v<", err)
			return err
		}
	}

	t.Log.Debug("test harness ready")

	return nil
}

// InitTx -
func (t *Testing) InitTx() (pgx.Tx, error) {
	l := t.Log.WithFunctionContext("InitTx")

	if t.tx != nil {
		l.Debug("Skipping Tx initialisation")
		return t.tx, nil
	}

	l.Debug("Initialising Tx")

	tx, err := t.Store.BeginTx()
	if err != nil {
		l.Warn("failed getting database tx >%v<", err)
		return nil, err
	}

	t.tx = tx

	err = t.Domain.Init(t.tx)
	if err != nil {
		l.Warn("failed model init >%v<", err)
		return nil, err
	}

	return t.tx, nil
}

// CommitTx -
func (t *Testing) CommitTx() (err error) {
	l := t.Log.WithFunctionContext("CommitTx")

	l.Debug("Committing Tx")

	err = t.tx.Commit(context.TODO())
	if err != nil {
		return err
	}
	t.tx = nil

	return nil
}

// RollbackTx -
func (t *Testing) RollbackTx() (err error) {
	l := t.Log.WithFunctionContext("RollbackTx")

	l.Debug("Rolling back Tx")

	err = t.tx.Rollback(context.TODO())
	if err != nil {
		return err
	}
	t.tx = nil

	return nil
}

// Setup -
//
// If ShouldCommitData is false, the tx is returned. The caller can perform other queries, but must commit or rollback the tx.
func (t *Testing) Setup() (pgx.Tx, error) {
	l := t.Log.WithFunctionContext("Setup")

	_, err := t.InitTx()
	if err != nil {
		l.Warn("failed init >%v<", err)
		return nil, err
	}

	// Create data function is expected to create and manage its own store
	if t.CreateDataFunc != nil {
		l.Debug("Creating test data")
		err := t.CreateDataFunc()
		if err != nil {
			l.Warn("failed creating data >%v<", err)
			return nil, err
		}
	}

	// Commit data when configured, otherwise we are leaving it up to tests
	// to explicitly commit or rollback.
	if t.ShouldCommitData {
		err = t.CommitTx()
		if err != nil {
			l.Warn("failed committing data >%v<", err)
			return nil, err
		}
		return nil, nil
	}

	return t.tx, nil
}

// Teardown -
func (t *Testing) Teardown() error {
	l := t.Log.WithFunctionContext("Teardown")

	// If we are not committing data, we need to rollback the tx.
	if !t.ShouldCommitData {
		l.Info("Rolling back Tx")
		err := t.RollbackTx()
		if err != nil {
			l.Warn("failed rolling back data >%v<", err)
			return err
		}
		return nil
	}

	_, err := t.InitTx()
	if err != nil {
		l.Warn("failed init >%v<", err)
		return err
	}

	// Remove data function is expected to create and manage its own store
	if t.RemoveDataFunc != nil {
		l.Debug("Removing test data")
		err := t.RemoveDataFunc()
		if err != nil {
			l.Warn("failed removing data >%v<", err)
			return err
		}
	}

	if t.ShouldCommitData {
		l.Debug("Committing database tx")
		err := t.CommitTx()
		if err != nil {
			l.Warn("failed committing data >%v<", err)
			return err
		}
	} else {
		l.Debug("Rollback database tx")
		err := t.RollbackTx()
		if err != nil {
			l.Warn("failed rolling back data >%v<", err)
			return err
		}
	}

	return nil
}

// Shutdown -
func (t *Testing) Shutdown(test *testing.T) {
	err := t.Store.ClosePool()
	require.NoError(test, err, "ClosePool should return no error")
}

// UpdateRecordCreatedAt bypasses default repository behaviour of not
// allowing created_at timestamps to be modified. This method should
// only be used under testing scenarios to create historical data.
func (t *Testing) UpdateRecordCreatedAt(tx pgx.Tx, tableName, recordID string, createdAt time.Time) error {

	args := pgx.NamedArgs{
		"id":         recordID,
		"created_at": createdAt,
	}
	sql := fmt.Sprintf("UPDATE %s SET created_at = @created_at WHERE id = @id", tableName)

	rows, err := tx.Query(context.Background(), sql, args)
	if err != nil {
		err = fmt.Errorf("failed updating record created_at: %w", err)
		return err
	}
	defer rows.Close()

	t.Log.Debug("Updated table name >%s< record ID >%s< createdAt >%s<",
		tableName,
		recordID,
		createdAt.String(),
	)

	return nil
}

// Logger - Returns a logger with package context and provided function context.
// This is the preferred way to access logging functionality.
func (t *Testing) Logger(functionName string) logger.Logger {
	return t.Log.WithPackageContext("harness").WithFunctionContext(functionName)
}
