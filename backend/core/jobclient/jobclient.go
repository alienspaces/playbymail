// Implements the Jobber interface
package jobclient

import (
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

func NewJobClient(s storer.Storer, riverConfig *river.Config) (*river.Client[pgx.Tx], error) {

	pool, err := s.Pool()
	if err != nil {
		return nil, err
	}

	riverClient, err := river.NewClient(riverpgxv5.New(pool), riverConfig)
	if err != nil {
		return nil, err
	}

	return riverClient, nil
}
