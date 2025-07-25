package runner

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
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
func NewRunnerWithConfig(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], cfg config.Config) (*Runner, error) {

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

	// Add authentication function
	r.AuthenticateRequestFunc = r.authenticateRequestFunc

	// Add RLS function for row-level security
	l.Warn("(playbymail) setting RLS function")
	r.RLSFunc = r.rlsFunc

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		r.gameHandlerConfig,
		r.accountHandlerConfig,
		r.gameSubscriptionHandlerConfig,
		r.gameAdministrationHandlerConfig,
		// Adventure Game Handlers
		adventure_game.AdventureGameHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		r.HandlerConfig = server.MergeHandlerConfigs(r.HandlerConfig, cfg)
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

// rlsFunc determines what game resources the authenticated user has access to
func (rnr *Runner) rlsFunc(l logger.Logger, m domainer.Domainer, authedReq server.AuthenData) (server.RLS, error) {

	l.Info("(playbymail) rlsFunc called for account ID: %s", authedReq.Account.ID)

	mm := m.(*domain.Domain)

	// Get all games the user has access to through subscriptions
	gameSubscriptions, err := mm.GameSubscriptionRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: "account_id",
				Val: authedReq.Account.ID,
			},
		},
	})
	if err != nil {
		l.Warn("(playbymail) failed to get game subscriptions >%v<", err)
		return server.RLS{}, err
	}

	// Extract game IDs from subscriptions
	gameIDs := make([]string, 0, len(gameSubscriptions))
	for _, sub := range gameSubscriptions {
		gameIDs = append(gameIDs, sub.GameID)
	}

	// Get all games the user administers
	gameAdministrations, err := mm.GameAdministrationRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: "account_id",
				Val: authedReq.Account.ID,
			},
		},
	})
	if err != nil {
		l.Warn("(playbymail) failed to get game administrations >%v<", err)
		return server.RLS{}, err
	}

	// Add administered game IDs
	for _, admin := range gameAdministrations {
		gameIDs = append(gameIDs, admin.GameID)
	}

	// Get all games the user owns (if account_id is the owner field)
	userGames, err := mm.GameRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: "account_id",
				Val: authedReq.Account.ID,
			},
		},
	})
	if err != nil {
		l.Warn("(playbymail) failed to get user-owned games >%v<", err)
		return server.RLS{}, err
	}
	for _, game := range userGames {
		gameIDs = append(gameIDs, game.ID)
	}

	// Deduplicate gameIDs
	gameIDSet := make(map[string]struct{})
	for _, id := range gameIDs {
		gameIDSet[id] = struct{}{}
	}
	uniqueGameIDs := make([]string, 0, len(gameIDSet))
	for id := range gameIDSet {
		uniqueGameIDs = append(uniqueGameIDs, id)
	}

	// Create RLS identifiers map
	identifiers := map[string][]string{
		"account_id": {authedReq.Account.ID},
	}

	// Add game IDs if user has access to any games
	if len(uniqueGameIDs) > 0 {
		identifiers["game_id"] = uniqueGameIDs
	}

	l.Info("(playbymail) RLS identifiers: %+v", identifiers)

	return server.RLS{
		Identifiers: identifiers,
	}, nil
}

// loggerWithFunctionContext provides a logger with function context
func loggerWithFunctionContext(l logger.Logger, functionName string) logger.Logger {
	return logging.LoggerWithFunctionContext(l, "runner", functionName)
}

// GetHandlerConfig returns the HandlerConfig map (for test compatibility)
func (r *Runner) GetHandlerConfig() map[string]server.HandlerConfig {
	return r.HandlerConfig
}
