package adventure_game_item_effect

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

const TableName = adventure_game_record.TableAdventureGameItemEffect

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[adventure_game_record.AdventureGameItemEffect](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    adventure_game_record.AdventureGameItemEffect{},
		},
	)
}
