package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

// TestRunner - allow Init and Run functions to be defined by tests
type TestRunner struct {
	Runner
	InitFunc func(l logger.Logger, s storer.Storer) error
}

func (rnr *TestRunner) Init(l logger.Logger, s storer.Storer) error {
	rnr.Log = l

	if rnr.InitFunc == nil {
		return rnr.Runner.Init(s)
	}
	return rnr.InitFunc(l, s)
}

func (rnr *TestRunner) mockAuthenticateRequestFunc(l logger.Logger, m domainer.Domainer, r *http.Request, authType AuthenticationType) (AuthenData, error) {
	return AuthenData{
		Type: AuthenticatedTypeToken,
		Account: AuthenticatedAccount{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
		},
	}, nil
}

func newTestRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], cfg config.Config) (*TestRunner, error) {
	cRnr, err := NewRunnerWithConfig(l, s, j, cfg)
	if err != nil {
		return nil, err
	}

	tr := TestRunner{
		Runner: *cRnr,
	}

	// tr.DomainFunc = func(l logger.Logger) (domainer.Domainer, error) {
	// 	return domain.NewDomain(l)
	// }

	tr.AuthenticateRequestFunc = tr.mockAuthenticateRequestFunc

	return &tr, nil
}

func TestRunnerInit(t *testing.T) {

	l, s, j := newDefaultDependencies(t)

	defer func() {
		err := s.ClosePool()
		require.NoError(t, err, "ClosePool should return no error")
	}()

	tr, err := newTestRunner(l, s, j, config.Config{})
	require.NoError(t, err, "newTestRunner returns without error")

	err = tr.Init(l, s)
	require.NoError(t, err, "Runner Init returns without error")

	tr.InitFunc = func(l logger.Logger, s storer.Storer) error {
		return errors.New("Init failed")
	}

	err = tr.Init(l, s)
	require.Error(t, err, "Runner Init returns with error")
}

func Test_RunnerServerError(t *testing.T) {
	t.Parallel()

	l, s, j := newDefaultDependencies(t)

	defer func() {
		err := s.ClosePool()
		require.NoError(t, err, "ClosePool should return no error")
	}()

	tr, err := newTestRunner(l, s, j, config.Config{})
	require.NoError(t, err, "newTestRunner returns without error")

	tr.RunHTTPFunc = func(args map[string]interface{}) error {
		return fmt.Errorf("Run server error")
	}

	tr.RunHTTPFunc = func(args map[string]interface{}) error {
		return fmt.Errorf("Run server error")
	}

	err = tr.Init(l, s)
	require.NoError(t, err, "Runner Init returns without error")

	err = tr.Run(nil)
	require.Error(t, err, "Runner Run returns with error")
}

func Test_RunnerDaemonError(t *testing.T) {
	t.Parallel()

	l, s, j := newDefaultDependencies(t)

	defer func() {
		err := s.ClosePool()
		require.NoError(t, err, "ClosePool should return no error")
	}()

	tr, err := newTestRunner(l, s, j, config.Config{})
	require.NoError(t, err, "newTestRunner returns without error")

	tr.RunHTTPFunc = func(args map[string]interface{}) error {
		return nil
	}

	tr.RunDaemonFunc = func(ctx context.Context, args map[string]interface{}) error {
		return fmt.Errorf("Run daemon error")
	}

	err = tr.Init(l, s)
	require.NoError(t, err, "Runner Init returns without error")

	err = tr.Run(nil)
	require.Error(t, err, "Runner Run returns with error")
}

func Test_registerRoutes(t *testing.T) {

	l, s, j := newDefaultDependencies(t)

	defer func() {
		err := s.ClosePool()
		require.NoError(t, err, "ClosePool should return no error")
	}()

	tr, err := newTestRunner(l, s, j, config.Config{})
	require.NoError(t, err, "newTestRunner returns without error")

	err = tr.Init(l, s)
	require.NoError(t, err, "Runner Init returns without error")

	r := httprouter.New()
	router, err := tr.registerRoutes(r)
	require.NoError(t, err, "Router returns without error")
	require.NotNil(t, router, "Router returns a router")

	// Test default configured routes
	handle, _, _ := router.Lookup(http.MethodGet, "/healthz")
	require.NotNil(t, handle, "Handle for /healthz is not nil")

	// Test custom routes
	tr.RouterFunc = func(r *httprouter.Router) (*httprouter.Router, error) {
		h, err := tr.ApplyMiddleware(HandlerConfig{Path: "/custom"}, tr.HandlerFunc)
		if err != nil {
			return nil, err
		}
		r.GET("/custom", h)
		return r, nil
	}

	r = httprouter.New()
	router, err = tr.registerRoutes(r)
	require.NoError(t, err, "Router returns without error")
	require.NotNil(t, router, "Router returns a router")

	// Test custom configured routes
	handle, _, _ = router.Lookup(http.MethodGet, "/custom")
	require.NotNil(t, handle, "Handle for /custom is not nil")

	// Test custom router error
	tr.RouterFunc = func(r *httprouter.Router) (*httprouter.Router, error) {
		return nil, errors.New("Failed router")
	}

	r = httprouter.New()
	router, err = tr.registerRoutes(r)
	require.Error(t, err, "Router returns with error")
	require.Nil(t, router, "Router returns nil")
}

func Test_ApplyMiddleware(t *testing.T) {

	l, s, j := newDefaultDependencies(t)

	defer func() {
		err := s.ClosePool()
		require.NoError(t, err, "ClosePool should return no error")
	}()

	tr, err := newTestRunner(l, s, j, config.Config{})
	require.NoError(t, err, "newTestRunner returns without error")

	err = tr.Init(l, s)
	require.NoError(t, err, "Runner Init returns without error")

	// Test default middleware
	handle, err := tr.ApplyMiddleware(HandlerConfig{Path: "/"}, tr.HandlerFunc)
	require.NoError(t, err, "Middleware returns without error")
	require.NotNil(t, handle, "Middleware returns a handle")

	// Test custom middleware
	tr.HandlerMiddlewareFuncs = func() []MiddlewareFunc {
		return []MiddlewareFunc{
			func(hc HandlerConfig, h Handle) (Handle, error) {
				return nil, errors.New("Failed middleware")
			},
		}
	}

	handle, err = tr.ApplyMiddleware(HandlerConfig{Path: "/"}, tr.HandlerFunc)
	require.Error(t, err, "Middleware returns with error")
	require.Nil(t, handle, "Middleware returns nil")
}
