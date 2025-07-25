package runner_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

// func NewTestRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (*runner.Runner, error) {

// 	cfg, err := config.Parse()
// 	if err != nil {
// 		return nil, err
// 	}
// 	rnr, err := runner.NewRunnerWithConfig(l, s, j, cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rnr.AuthenticateRequestFunc = func(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenData, error) {
// 		return server.AuthenData{
// 			Type: server.AuthenticatedTypeToken,
// 		}, nil
// 	}

// 	rnr.RLSFunc = func(l logger.Logger, m domainer.Domainer, authedReq server.AuthenData) (server.RLS, error) {
// 		return server.RLS{
// 			Identifiers: map[string][]string{},
// 		}, nil
// 	}

// 	return rnr, nil
// }

func TestNewRunner(t *testing.T) {

	cfg, err := config.Parse()
	require.NoError(t, err, "Parse returns without error")

	l, s, j, err := deps.NewDefaultDependencies(cfg)
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	r, err := testutil.NewTestRunner(l, s, j)
	require.NoError(t, err, "NewTestRunner returns without error")

	err = r.Init(s)
	require.NoError(t, err, "Init returns without error")
}
