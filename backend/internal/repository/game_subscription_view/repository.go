package game_subscription_view

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const TableName = game_record.TableGameSubscriptionView

// NewRepository implements core domain.RepositoryConstructor
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGenericView[game_record.GameSubscriptionView](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    game_record.GameSubscriptionView{},
		},
	)
}
