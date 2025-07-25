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
	"gitlab.com/alienspaces/playbymail/internal/runner/server/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// Runner -
type Runner struct {
	// Embed *server.Runner as a pointer so field changes are shared.
	// Value embedding would copy fields, causing changes to be lost.
	// Pointer embedding ensures all references use the same data.
	*server.Runner
	Config config.Config
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
		Config: cfg,
	}

	l.Warn("(playbymail) setting handler function")
	r.HandlerFunc = r.handlerFunc

	l.Warn("(playbymail) setting run daemon function")
	r.RunDaemonFunc = r.runDaemonFunc

	l.Warn("(playbymail) setting domain function on runner: %p", &r)
	r.DomainFunc = r.domainFunc

	l.Warn("(playbymail) setting authenticate request function")

	// Add mock authentication function for testing
	r.AuthenticateRequestFunc = r.authenticateRequestFunc

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		r.gameHandlerConfig,
		r.accountHandlerConfig,
		// Adventure Game Handlers
		adventure_game.AdventureGameHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		r.HandlerConfig = mergeHandlerConfigs(r.HandlerConfig, cfg)
	}

	return &r, nil
}

// DomainFunc -
func (rnr *Runner) domainFunc(l logger.Logger) (domainer.Domainer, error) {
	l.Info("(playbymail) DomainFunc called on runner: %p", rnr)
	l.Info("(playbymail) calling domain.NewDomain")
	m, err := domain.NewDomain(l, rnr.JobClient, rnr.Config)
	if err != nil {
		l.Warn("(playbymail) failed new domain >%v<", err)
		return nil, err
	}
	return m, nil
}

// authenticateRequestFunc authenticates a request based on the authentication type
func (rnr *Runner) authenticateRequestFunc(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenData, error) {

	switch authType {
	case server.AuthenticationTypeToken:
		return rnr.authenticateRequestTokenFunc(l, m, r)
	default:
		return server.AuthenData{}, nil
	}
}

// authenticateRequestTokenFunc authenticates a request based on a session token. Returning anything
// other than an AuthenData{} with a valid Typewill result in a 401 Unauthorized response.
func (rnr *Runner) authenticateRequestTokenFunc(l logger.Logger, m domainer.Domainer, r *http.Request) (server.AuthenData, error) {

	l.Info("(playbymail) authenticateRequestTokenFunc called")

	mm := m.(*domain.Domain)

	accountRec, err := mm.VerifyAccountSessionToken(r.Header.Get("Authorization"))
	if err != nil {
		l.Warn("(playbymail) failed to verify account session token >%v<", err)
		return server.AuthenData{}, err
	}

	if accountRec == nil {
		l.Warn("(playbymail) no account found for session token")
		return server.AuthenData{}, nil
	}

	return server.AuthenData{
		Type: server.AuthenticatedTypeToken,
		Account: server.AuthenticatedAccount{
			ID:    accountRec.ID,
			Name:  accountRec.Name,
			Email: accountRec.Email,
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
