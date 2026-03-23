package mech_wargame_lance_mech

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

const (
	TableName string = mech_wargame_record.TableMechWargameLanceMech
)

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[mech_wargame_record.MechWargameLanceMech, *mech_wargame_record.MechWargameLanceMech](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    mech_wargame_record.MechWargameLanceMech{},
		},
	)
}
