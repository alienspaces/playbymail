package game_item

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

const TableName = record.TableGameItem

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[record.GameItem](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    record.GameItem{},
		},
	)
}
