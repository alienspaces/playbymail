package game_administration

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const TableName = game_record.TableGameAdministration

// NewRepository implements core domain.RepositoryConstructor
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[game_record.GameAdministration](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    game_record.GameAdministration{},
		},
	)
}
