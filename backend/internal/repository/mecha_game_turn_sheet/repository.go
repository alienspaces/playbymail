package mecha_game_turn_sheet

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

const (
	TableName string = mecha_game_record.TableMechaGameTurnSheet
)

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[mecha_game_record.MechaGameTurnSheet, *mecha_game_record.MechaGameTurnSheet](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    mecha_game_record.MechaGameTurnSheet{},
		},
	)
}
