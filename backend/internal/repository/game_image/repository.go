package game_image

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const (
	TableName string = game_record.TableGameImage
)

// NewRepository -
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[game_record.GameImage](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    game_record.GameImage{},
		},
	)
}
