package deps

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/internal/jobclient"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

func Default(cfg config.Config) (*log.Log, *store.Store, *river.Client[pgx.Tx], error) {

	// Logger
	l, err := log.NewLogger(cfg.Config)
	if err != nil {
		fmt.Printf("failed new logger >%v<", err)
		return nil, nil, nil, err
	}

	// Storer
	s, err := store.NewStore(cfg.Config)
	if err != nil {
		fmt.Printf("failed new store >%v<", err)
		return nil, nil, nil, err
	}

	// River
	j, err := jobclient.NewJobClient(l, s, []string{jobqueue.QueueDefault})
	if err != nil {
		fmt.Printf("failed new job client >%v<", err)
		return nil, nil, nil, err
	}

	return l, s, j, nil
}
