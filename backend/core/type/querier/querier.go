package querier

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
)

type Querier interface {
	Name() string
	Exec(args pgx.NamedArgs) (pgconn.CommandTag, error)
	Query(opts *coresql.Options) (pgx.Rows, error)
}
