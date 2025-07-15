package deps

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/email/fake"
	"gitlab.com/alienspaces/playbymail/core/email/forwardemail"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/core/type/emailer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/jobclient"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

func NewHarness(t *testing.T) *harness.Testing {
	dcfg := harness.DefaultDataConfig()
	cfg, err := config.Parse()
	require.NoError(t, err)
	l, s, j, err := Default(cfg)
	require.NoError(t, err)
	h, err := harness.NewTesting(l, s, j, cfg, dcfg)
	require.NoError(t, err)

	// We setup and teardown within the context of the test
	// so we don't need to commit the data to the database.
	h.ShouldCommitData = false

	return h
}

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

	// Emailer
	var e emailer.Emailer
	if cfg.EmailerFaked {
		fmt.Println("using fake emailer")
		e, err = fake.New(l, cfg.Config)
	} else {
		fmt.Println("using forward emailer")
		e, err = forwardemail.New(l, cfg.Config)
	}
	if err != nil {
		fmt.Printf("failed new emailer >%v<", err)
		return nil, nil, nil, err
	}

	// River
	j, err := jobclient.NewJobClient(l, cfg, s, e, []string{jobqueue.QueueDefault})
	if err != nil {
		fmt.Printf("failed new job client >%v<", err)
		return nil, nil, nil, err
	}

	return l, s, j, nil
}
