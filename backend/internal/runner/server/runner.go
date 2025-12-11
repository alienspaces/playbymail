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
	"gitlab.com/alienspaces/playbymail/internal/jobclient"
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/account"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/adventure_game"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/game"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_rls"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/player"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
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
func NewRunner(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scnr turn_sheet.TurnSheetScanner) (*Runner, error) {
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

	l.Warn("(playbymail) setting handler function")
	r.HandlerFunc = r.handlerFunc

	l.Warn("(playbymail) setting run daemon function")
	r.RunDaemonFunc = r.runDaemonFunc

	l.Warn("(playbymail) setting domain function on runner: %p", &r)
	r.DomainFunc = r.domainFunc

	l.Warn("(playbymail) setting job client function on runner: %p", &r)
	r.JobClientFunc = r.jobClientFunc

	l.Warn("(playbymail) setting authenticate request function")
	r.AuthenticateRequestFunc = r.authenticateRequestFunc

	l.Warn("(playbymail) setting RLS function")
	r.RLSFunc = handler_rls.HandlerRLSFunc

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(config.Config, logger.Logger, turn_sheet.TurnSheetScanner) (map[string]server.HandlerConfig, error){
		// Account related handlers
		account.AccountHandlerConfig,
		// Adventure Game handlers
		adventure_game.AdventureGameHandlerConfig,
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

	l.Info("(playbymail) JobClientFunc called on runner: %p", rnr)
	l.Info("(playbymail) calling jobclient.NewJobClient")

	// This job client is only used for registering jobs within the handler functions
	// and is not used for processing jobs so should not need an actual emailer.
	j, err := jobclient.NewJobClient(l, rnr.Config, s, nil, []string{jobqueue.QueueDefault})
	if err != nil {
		l.Warn("(playbymail) failed new job client >%v<", err)
		return nil, err
	}
	return j, nil
}

// domainFunc returns a new domain model instance.
func (rnr *Runner) domainFunc(l logger.Logger) (domainer.Domainer, error) {
	l = l.WithFunctionContext("domainFunc")

	l.Info("(playbymail) DomainFunc called on runner: %p", rnr)
	l.Info("(playbymail) calling domain.NewDomain")

	m, err := domain.NewDomain(l, rnr.Config)
	if err != nil {
		l.Warn("(playbymail) failed new domain >%v<", err)
		return nil, err
	}

	// Initialize turn sheet processors
	l.Info("(playbymail) initializing turn sheet processors")

	l.Info("(playbymail) successfully initialized turn sheet processors")

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

	// Check for development mode authentication bypass
	if rnr.Config.AppEnv == "develop" {
		if bypassEmail := r.Header.Get("X-Bypass-Authentication"); bypassEmail != "" {
			l.Info("(playbymail) development mode: using bypass authentication for email >%s<", bypassEmail)

			// In development mode, query the actual account record by email
			// This bypasses the need for actual session tokens but uses real account data
			accountRecs, err := mm.GetManyAccountRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountEmail, Val: bypassEmail},
				},
				Limit: 1,
			})
			if err != nil {
				l.Warn("(playbymail) development mode: failed to get account by email >%s< >%v<", bypassEmail, err)
				return server.AuthenData{}, err
			}

			if len(accountRecs) == 0 {
				l.Warn("(playbymail) development mode: no account found for email >%s<", bypassEmail)
				return server.AuthenData{}, nil
			}

			accountRec := accountRecs[0]
			l.Info("(playbymail) development mode: found account ID >%s< for email >%s<", accountRec.ID, bypassEmail)

			// Get account contact name if available
			accountName := ""
			contactRecs, err := mm.GetManyAccountContactRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountContactAccountID, Val: accountRec.ID},
				},
				Limit: 1,
				OrderBy: []coresql.OrderBy{
					{Col: account_record.FieldAccountContactCreatedAt, Direction: coresql.OrderDirectionASC},
				},
			})
			if err == nil && len(contactRecs) > 0 {
				accountName = contactRecs[0].Name
			}

			return server.AuthenData{
				Type:    server.AuthenticatedTypeToken,
				RLSType: server.RLSTypeRestricted,
				Account: server.AuthenticatedAccount{
					ID:    accountRec.ID,
					Name:  accountName,
					Email: accountRec.Email,
				},
			}, nil
		}
	}

	// Try to extract token from multiple sources:
	// 1. Authorization header (preferred)
	// 2. Query parameter 'token' (for iframe/embed scenarios)
	var token string

	// First, try Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) >= 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
		l.Debug("(playbymail) token extracted from Authorization header")
	}

	// If no header token, try query parameter (for iframe usage)
	if token == "" {
		token = r.URL.Query().Get("token")
		if token != "" {
			l.Debug("(playbymail) token extracted from query parameter")
		}
	}

	// No token found
	if token == "" {
		l.Warn("(playbymail) no token found in Authorization header or query parameter")
		return server.AuthenData{}, nil
	}

	accountRec, err := mm.VerifyAccountSessionToken(token)
	if err != nil {
		l.Warn("(playbymail) failed to verify account session token >%v<", err)
		return server.AuthenData{}, err
	}

	if accountRec == nil {
		l.Warn("(playbymail) no account found for session token")
		return server.AuthenData{}, nil
	}

	// Get account contact name if available
	accountName := ""
	contactRecs, err := mm.GetManyAccountContactRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountContactAccountID, Val: accountRec.ID},
		},
		Limit: 1,
		OrderBy: []coresql.OrderBy{
			{Col: account_record.FieldAccountContactCreatedAt, Direction: coresql.OrderDirectionASC},
		},
	})
	if err == nil && len(contactRecs) > 0 {
		accountName = contactRecs[0].Name
	}

	authenData := server.AuthenData{
		Type:    server.AuthenticatedTypeToken,
		RLSType: server.RLSTypeRestricted,
		Account: server.AuthenticatedAccount{
			ID:    accountRec.ID,
			Name:  accountName,
			Email: accountRec.Email,
		},
	}

	l.Info("(playbymail) authenticated account: ID=%s Email=%s Name=%s",
		authenData.Account.ID, authenData.Account.Email, authenData.Account.Name)

	return authenData, nil
}
