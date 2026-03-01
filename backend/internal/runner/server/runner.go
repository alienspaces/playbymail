package runner

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/runnable"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobclient"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/account"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/catalog"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_rls"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/player"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
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
func NewRunner(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scnr turnsheet.TurnSheetScanner) (*Runner, error) {
	l = l.WithPackageContext("runner")

	cr, err := server.NewRunner(cfg.Config, l, s, j)
	if err != nil {
		err := fmt.Errorf("failed core runner >%v<", err)
		l.Warn(err.Error())
		return nil, err
	}

	r := Runner{
		Runner: cr,
		Config: cfg,
	}

	l.Warn("setting handler function")
	r.HandlerFunc = r.handlerFunc

	l.Warn("setting run daemon function")
	r.RunDaemonFunc = r.runDaemonFunc

	l.Warn("setting domain function on runner: %p", &r)
	r.DomainFunc = r.domainFunc

	l.Warn("setting job client function on runner: %p", &r)
	r.JobClientFunc = r.jobClientFunc

	l.Warn("setting authenticate request function")
	r.AuthenticateRequestFunc = r.authenticateRequestFunc

	l.Warn("setting RLS function")
	r.RLSFunc = handler_rls.HandlerRLSFunc

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(config.Config, logger.Logger, turnsheet.TurnSheetScanner) (map[string]server.HandlerConfig, error){
		// Account related handlers
		account.AccountHandlerConfig,
		// Adventure Game handlers
		adventure_game.AdventureGameHandlerConfig,
		// Catalog handlers (public)
		catalog.CatalogHandlerConfig,
		// Player handlers
		player.PlayerHandlerConfig,
		// Game handlers
		game.GameHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		handlerConfig, err := fn(cfg, l, scnr)
		if err != nil {
			return nil, err
		}
		r.HandlerConfig = server.MergeHandlerConfigs(r.HandlerConfig, handlerConfig)
	}

	return &r, nil
}

// GetHandlerConfig returns the HandlerConfig map
// TODO: The core runner equivalent method should work fine and this should be removed.
func (r *Runner) GetHandlerConfig() map[string]server.HandlerConfig {
	return r.HandlerConfig
}

// jobClientFunc returns a new job client instance.
func (rnr *Runner) jobClientFunc(l logger.Logger, s storer.Storer) (*river.Client[pgx.Tx], error) {
	l = l.WithFunctionContext("jobClientFunc")

	l.Info("JobClientFunc called on runner: %p", rnr)
	l.Info("calling jobclient.NewJobClient")

	// This job client is only used for registering jobs within the handler functions
	// and is not used for processing jobs so should not need an actual emailer.
	j, err := jobclient.NewJobClient(l, rnr.Config, s, nil, []string{jobqueue.QueueDefault})
	if err != nil {
		l.Warn("failed new job client >%v<", err)
		return nil, err
	}
	return j, nil
}

// domainFunc returns a new domain model instance.
func (rnr *Runner) domainFunc(l logger.Logger) (domainer.Domainer, error) {
	l = l.WithFunctionContext("domainFunc")

	l.Info("DomainFunc called on runner: %p", rnr)
	l.Info("calling domain.NewDomain")

	m, err := domain.NewDomain(l, rnr.Config)
	if err != nil {
		l.Warn("failed new domain >%v<", err)
		return nil, err
	}

	// Initialize turn sheet processors
	l.Info("initializing turn sheet processors")

	l.Info("successfully initialized turn sheet processors")

	return m, nil
}

// authenticateRequestFunc authenticates a request based on the authentication type
func (rnr *Runner) authenticateRequestFunc(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenData, error) {

	switch authType {
	case server.AuthenticationTypeToken:
		return handler_auth.AuthenticateRequestTokenFunc(rnr.Config, l, m, r)
	default:
		return server.AuthenData{}, coreerror.NewUnauthenticatedError("unsupported authentication type")
	}
}
