package repositor

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
)

// Repositor -
type Repositor interface {
	TableName() string
	Attributes() []string
	ArrayFields() set.Set[string]
	Tx() pgx.Tx
	SetRLS(identifiers map[string][]string)
}
