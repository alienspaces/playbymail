package account_user

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

const (
	TableName string = account_record.TableAccountUser
)

// NewRepository -
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[account_record.AccountUser](
		repository.NewArgs{
			Tx:            tx,
			TableName:     TableName,
			Record:        account_record.AccountUser{},
			IsRLSDisabled: true,
		},
	)
}
