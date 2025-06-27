package storer

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storer interface {
	Pool() (*pgxpool.Pool, error)
	ClosePool() error
	BeginTx() (pgx.Tx, error)
}
