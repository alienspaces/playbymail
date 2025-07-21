package adventure_game_location_link_requirement

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

const TableName = record.TableAdventureGameLocationLinkRequirement

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[record.AdventureGameLocationLinkRequirement](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    record.AdventureGameLocationLinkRequirement{},
		},
	)
}
