package mecha_game_mech_instance

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

const (
	TableName string = mecha_game_record.TableMechaGameMechInstance
)

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[mecha_game_record.MechaGameMechInstance, *mecha_game_record.MechaGameMechInstance](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    mecha_game_record.MechaGameMechInstance{},
		},
	)
}
