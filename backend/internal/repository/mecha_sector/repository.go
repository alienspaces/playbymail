package mecha_sector

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

const (
	TableName string = mecha_record.TableMechaSector
)

func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[mecha_record.MechaSector, *mecha_record.MechaSector](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    mecha_record.MechaSector{},
		},
	)
}
