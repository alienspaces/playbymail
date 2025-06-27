package store

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

type Store struct {
	log     logger.Logger
	pgxPool *pgxpool.Pool
	config  config.Config
}

var _ storer.Storer = &Store{}

func NewStore(cfg config.Config) (*Store, error) {

	l, err := log.NewLogger(cfg)
	if err != nil {
		return nil, err
	}

	s := Store{
		log:    l.WithPackageContext("store"),
		config: cfg,
	}

	// Check timezone
	if time.Local.String() != "UTC" {
		return nil, fmt.Errorf("failed timezone check, not UTC >%s<", time.Local.String())
	}

	return &s, nil
}

// Pool returns a pgx connection pool
func (s *Store) Pool() (*pgxpool.Pool, error) {
	l := s.log.WithFunctionContext("Pool")

	if s.pgxPool != nil {
		return s.pgxPool, nil
	}

	l.Info("(core) connecting to postgres pool >%s<", s.config.DatabaseURL)

	pool, err := connectPgx(s.log, s.config)
	if err != nil {
		l.Warn("(core) failed connecting to postgres pool >%v<", err)
		return nil, err
	}

	s.pgxPool = pool

	return pool, nil
}

// ClosePool closes the database connection pool
func (s *Store) ClosePool() error {
	l := s.log.WithFunctionContext("ClosePool")

	if s.pgxPool == nil {
		// Leaving this as a warning for now just to track what
		// scenarios this might happen.
		l.Warn("(core) not closing pool, pool is nil")
		return nil
	}

	l.Debug("(core) stat >%+v<", spew.Sdump(s.pgxPool.Stat()))

	l.Info("(core) closing postgres pool")

	s.pgxPool.Close()
	s.pgxPool = nil

	return nil
}

// BeginTx returns a pgx transaction. You are responsible for
// for handling commit or rollback.
func (s *Store) BeginTx() (pgx.Tx, error) {
	l := s.log.WithFunctionContext("BeginTx")

	var err error

	pool, err := s.Pool()
	if err != nil {
		l.Warn("(core) failed Pool >%#v<", err)
		return nil, err
	}

	tx, err := pool.Begin(context.TODO())
	if err != nil {
		l.Warn("(core) failed begin >%v<", err)
		return nil, err
	}

	return tx, nil
}
