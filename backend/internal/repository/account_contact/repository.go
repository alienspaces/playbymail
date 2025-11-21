package account_contact

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

const (
	TableName string = account_record.TableAccountContact
)

// NewRepository -
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[account_record.AccountContact](
		repository.NewArgs{
			Tx:        tx,
			TableName: TableName,
			Record:    account_record.AccountContact{},
		},
	)
}
