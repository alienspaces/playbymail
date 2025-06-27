package server

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/core/jobclient"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/store"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

// newDefaultDependencies -
func newDefaultDependencies(t *testing.T) (logger.Logger, storer.Storer, *river.Client[pgx.Tx]) {
	cfg := config.Config{}
	err := config.Parse(&cfg)
	require.NoError(t, err, "config.Parse returns without error")

	l, err := log.NewLogger(cfg)
	require.NoError(t, err, "NewLogger returns without error")

	s, err := store.NewStore(cfg)
	require.NoError(t, err, "NewStore returns without error")

	j, err := jobclient.NewJobClient(s, &river.Config{})
	require.NoError(t, err, "NewJobClient returns without error")

	return l, s, j
}

func TestNewServer(t *testing.T) {

	l, s, _ := newDefaultDependencies(t)

	defer func() {
		err := s.ClosePool()
		require.NoError(t, err, "ClosePool should return no error")
	}()

	tr := Runner{
		Log: l,
	}
	tr.DomainFunc = func(l logger.Logger) (domainer.Domainer, error) {
		return domain.NewDomain(l)
	}

	ts, err := NewServer(l, s, &tr)
	require.NoError(t, err, "NewServer returns without error")
	require.NotNil(t, ts, "Test server is not nil")
}
