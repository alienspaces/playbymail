package mecha_game_squad_mech

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

const (
	TableName string = mecha_game_record.TableMechaGameSquadMech
)

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[mecha_game_record.MechaGameSquadMech](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    mecha_game_record.MechaGameSquadMech{},
		},
	)
}
