package adventure_game_location_instance

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

const TableName = record.TableAdventureGameLocationInstance

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[record.AdventureGameLocationInstance](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    record.AdventureGameLocationInstance{},
		},
	)
}
