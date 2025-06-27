package game

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

const (
	// TableName - underlying database table name used for configuration
	TableName string = record.TableGame
)

// NewRepository -
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[record.Game](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    record.Game{},
		},
	)
}
