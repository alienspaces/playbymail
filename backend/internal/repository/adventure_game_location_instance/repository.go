package adventure_game_location_instance

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	adventure_game_record "gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

const TableName = adventure_game_record.TableAdventureGameLocationInstance

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[adventure_game_record.AdventureGameLocationInstance](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    adventure_game_record.AdventureGameLocationInstance{},
		},
	)
}
