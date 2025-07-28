package game_configuration

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const TableName = game_record.TableGameConfiguration

// NewRepository matches the RepositoryConstructor signature
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[game_record.GameConfiguration](repository.NewArgs{
		Tx:        tx,
		TableName: TableName,
		Record:    game_record.GameConfiguration{},
	})
}
