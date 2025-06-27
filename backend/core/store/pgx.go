package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // blank import intended

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// connectPgx -
func connectPgx(l logger.Logger, cfg config.Config) (*pgxpool.Pool, error) {

	if l == nil {
		return nil, fmt.Errorf("missing logger, cannot connect to postgres database")
	}

	cs := cfg.DatabaseURL
	if cs == "" {
		err := fmt.Errorf("failed connectPgx, DATABASE_URL environment variable is not set")
		l.Warn(err.Error())
		return nil, err
	}
	l.Info("(core) connecting to postgres using DATABASE_URL")

	pgxcfg, err := pgxpool.ParseConfig(cs)
	if err != nil {
		l.Warn("(core) failed parse config >%v<", err)
		return nil, err
	}

	pgxcfg.MaxConns = int32(cfg.DatabaseMaxOpenConnections)
	pgxcfg.MaxConnIdleTime = time.Duration(cfg.DatabaseMaxIdleTimeMins) * time.Minute

	// All database timestamp columns in applications built using this core library
	// are assumed to be `TIMESTAMP WITH TIMEZONE` and are stored in UTC. The
	// following ensures that timestamps values scanned back into time.Time values
	// have the location set correctly to UTC and not Local.
	// https://github.com/jackc/pgx/pull/1948
	pgxcfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "timestamptz",
			OID:   pgtype.TimestamptzOID,
			Codec: &pgtype.TimestampCodec{ScanLocation: time.UTC},
		})
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), pgxcfg)
	if err != nil {
		l.Warn("(core) failed to new pool >%v<", err)
		return nil, err
	}

	l.Info("(core) pool config HealthCheckPeriod >%d<", pool.Config().HealthCheckPeriod)
	l.Info("(core) pool config MinConns >%d<", pool.Config().MinConns)
	l.Info("(core) pool config MaxConns >%d<", pool.Config().MaxConns)
	l.Info("(core) pool config MaxConnIdleTime >%d<", pool.Config().MaxConnIdleTime)
	l.Info("(core) pool config MaxConnLifetime >%d<", pool.Config().MaxConnLifetime)
	l.Info("(core) pool config MaxConnLifetimeJitter >%d<", pool.Config().MaxConnLifetimeJitter)

	return pool, nil
}
