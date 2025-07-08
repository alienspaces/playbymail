package runner

import (
	"fmt"
	"maps"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/runnable"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// Runner -
type Runner struct {
	// Embed *server.Runner as a pointer so field changes are shared.
	// Value embedding would copy fields, causing changes to be lost.
	// Pointer embedding ensures all references use the same data.
	*server.Runner
}

// ensure we comply with the Runnerer interface
var _ runnable.Runnable = &Runner{}

// NewRunner -
func NewRunner(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], cfg config.Config) (*Runner, error) {

	l = l.WithPackageContext("runner")

	cr, err := server.NewRunnerWithConfig(l, s, j, cfg.Config)
	if err != nil {
		err := fmt.Errorf("failed core runner >%v<", err)
		l.Warn(err.Error())
		return nil, err
	}

	r := Runner{
		Runner: cr,
	}

	l.Warn("(runner) setting handler function")
	r.HandlerFunc = r.handlerFunc

	l.Warn("(runner) setting run daemon function")
	r.RunDaemonFunc = r.runDaemonFunc

	l.Warn("(runner) setting domain function on runner: %p", &r)
	r.DomainFunc = r.domainFunc

	l.Warn("(runner) setting authenticate request function")

	// Add mock authentication function for testing
	r.AuthenticateRequestFunc = r.mockAuthenticateRequest

	// TODO: Additional handler configs can be added here
	gameConfig, err := r.gameHandlerConfig(l)
	if err != nil {
		return nil, err
	}

	// Add handler configs using mergeHandlerConfigs for proper configuration merging
	r.HandlerConfig = mergeHandlerConfigs(r.HandlerConfig, gameConfig)

	accountConfig, err := r.accountHandlerConfig(l)
	if err != nil {
		return nil, err
	}
	r.HandlerConfig = mergeHandlerConfigs(r.HandlerConfig, accountConfig)

	locationConfig, err := r.locationHandlerConfig(l)
	if err != nil {
		return nil, err
	}
	r.HandlerConfig = mergeHandlerConfigs(r.HandlerConfig, locationConfig)

	locationLinkConfig, err := r.locationLinkHandlerConfig(l)
	if err != nil {
		return nil, err
	}
	r.HandlerConfig = mergeHandlerConfigs(r.HandlerConfig, locationLinkConfig)

	return &r, nil
}

// DomainFunc -
func (rnr *Runner) domainFunc(l logger.Logger) (domainer.Domainer, error) {
	l.Info("(runner) DomainFunc called on runner: %p", rnr)
	l.Info("(runner) calling domain.NewDomain")
	m, err := domain.NewDomain(l, rnr.JobClient)
	if err != nil {
		l.Warn("(runner) failed new domain >%v<", err)
		return nil, err
	}
	return m, nil
}

// mockAuthenticateRequest provides a mock authentication function for testing
func (rnr *Runner) mockAuthenticateRequest(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenticatedRequest, error) {
	// For testing, always return a mock authenticated request with admin permissions
	return server.AuthenticatedRequest{
		Type: server.AuthenticatedTypeAPIKey,
		User: server.AuthenticatedUser{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
		},
		Permissions: []server.AuthorizedPermission{
			"admin",
			"games:read",
			"games:write",
			"games:delete",
		},
	}, nil
}

func mergeHandlerConfigs(hc1 map[string]server.HandlerConfig, hc2 map[string]server.HandlerConfig) map[string]server.HandlerConfig {
	if hc1 == nil {
		hc1 = map[string]server.HandlerConfig{}
	}
	maps.Copy(hc1, hc2)
	return hc1
}

// loggerWithFunctionContext provides a logger with function context
func loggerWithFunctionContext(l logger.Logger, functionName string) logger.Logger {
	return logging.LoggerWithFunctionContext(l, "runner", functionName)
}
