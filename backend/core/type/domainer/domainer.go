package domainer

import (
	"github.com/jackc/pgx/v5"
)

// Domainer -
type Domainer interface {
	GetTx() (pgx.Tx, error)
	Init(tx pgx.Tx) (err error)
	SetTxLockTimeout(timeoutSecs float64) error
	Commit() error
	Rollback() error
}
