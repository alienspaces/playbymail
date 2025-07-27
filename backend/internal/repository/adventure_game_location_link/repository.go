package adventure_game_location_link

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

const TableName string = adventure_game_record.TableAdventureGameLocationLink

// NewRepository -
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[adventure_game_record.AdventureGameLocationLink](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    adventure_game_record.AdventureGameLocationLink{},
		},
	)
}
