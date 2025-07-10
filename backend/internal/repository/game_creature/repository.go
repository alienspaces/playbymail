package game_creature

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

const TableName = record.TableGameCreature

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[record.GameCreature](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    record.GameCreature{},
		},
	)
}
