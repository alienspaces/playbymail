package handler_rls

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const (
	rlsIdentifierAccountID     = "account_id"
	rlsIdentifierAccountUserID = "account_user_id"
)

// gameRLSConstraints defines constraints that will be automatically applied
// when repositories have matching column names. Constraints use SQL subqueries
// to filter based on game_subscription relationships.
var gameRLSConstraints = []repositor.RLSConstraint{
	{
		Column: game_record.FieldGameSubscriptionInstanceGameInstanceID,
		SQLTemplate: fmt.Sprintf(
			"IN (SELECT %s FROM %s gsi INNER JOIN %s gs ON gsi.%s = gs.%s WHERE gsi.%s = :account_id AND gs.%s = '%s' AND gs.%s = '%s')",
			game_record.FieldGameSubscriptionInstanceGameInstanceID,
			game_record.TableGameSubscriptionInstance,
			game_record.TableGameSubscription,
			game_record.FieldGameSubscriptionInstanceGameSubscriptionID,
			game_record.FieldGameSubscriptionID,
			game_record.FieldGameSubscriptionInstanceAccountID,
			game_record.FieldGameSubscriptionSubscriptionType,
			game_record.GameSubscriptionTypeManager,
			game_record.FieldGameSubscriptionStatus,
			game_record.GameSubscriptionStatusActive,
		),
		RequiredRLSIdentifiers: []string{"account_id"},
	},
	{
		Column: game_record.FieldGameSubscriptionGameID,
		SQLTemplate: fmt.Sprintf("IN (SELECT %s FROM %s WHERE %s = :account_id AND %s IN ('%s', '%s') AND %s = '%s')",
			game_record.FieldGameSubscriptionGameID,
			game_record.TableGameSubscription,
			game_record.FieldGameSubscriptionAccountID,
			game_record.FieldGameSubscriptionSubscriptionType,
			game_record.GameSubscriptionTypeDesigner,
			game_record.GameSubscriptionTypeManager,
			game_record.FieldGameSubscriptionStatus,
			game_record.GameSubscriptionStatusActive,
		),
		RequiredRLSIdentifiers: []string{"account_id"},
	},
	// Allow access via account_user_id or account_id for tables with an account_id column.
	// Also allow access if the record ID corresponds to a game the user has a subscription to.
	// We allow IS NULL to support tables (like account_subscription) where this column might be null
	// but another column (account_user_id) provides access control.
	{
		Column: "account_id",
		SQLTemplate: fmt.Sprintf("IS NULL OR account_id IN (:account_user_id, :account_id) OR id IN (SELECT %s FROM %s WHERE ((%s = :account_id AND %s IN ('%s', '%s')) OR (%s = :account_user_id AND %s = '%s')) AND %s = '%s')",
			game_record.FieldGameSubscriptionGameID,
			game_record.TableGameSubscription,
			game_record.FieldGameSubscriptionAccountID,
			game_record.FieldGameSubscriptionSubscriptionType,
			game_record.GameSubscriptionTypeDesigner,
			game_record.GameSubscriptionTypeManager,
			"account_user_id", // FieldGameSubscriptionAccountUserID
			game_record.FieldGameSubscriptionSubscriptionType,
			game_record.GameSubscriptionTypePlayer,
			game_record.FieldGameSubscriptionStatus,
			game_record.GameSubscriptionStatusActive,
		),
		RequiredRLSIdentifiers: []string{"account_user_id", "account_id"},
		SkipSelfMapping:        true,
	},
	{
		Column:                 "account_user_id",
		SQLTemplate:            "IS NULL OR account_user_id = :account_user_id",
		RequiredRLSIdentifiers: []string{"account_user_id"},
		SkipSelfMapping:        true,
	},
}

// HandlerRLSFunc determines what game resources the authenticated user has access to.
// It sets RLS identifiers (account_id) and RLS constraints (game_id, game_instance_id)
// that are automatically applied by the repository layer based on column presence.
func HandlerRLSFunc(l logger.Logger, m domainer.Domainer, authedReq server.AuthenData) (server.RLS, error) {

	l.Info("(playbymail) rlsFunc called for account: ID=%s Email=%s",
		authedReq.AccountUser.ID, authedReq.AccountUser.Email)

	// Create RLS identifiers map
	// account_id is always set from authenticated account (direct relationship)
	// game_id and game_instance_id filtering is handled via RLS constraints
	identifiers := map[string][]string{
		rlsIdentifierAccountID:     {authedReq.AccountUser.AccountID},
		rlsIdentifierAccountUserID: {authedReq.AccountUser.ID},
	}

	return server.RLS{
		Identifiers: identifiers,
		Constraints: gameRLSConstraints,
	}, nil
}

// HandlerRLSAccountIdentifierValue returns the account identifier value for the current authenticated user
func HandlerRLSAccountIdentifierValue(l logger.Logger, r *http.Request) (string, error) {
	rlsData := server.GetRequestRLSData(l, r)
	if rlsData == nil {
		return "", nil
	}

	return rlsData.Identifiers[rlsIdentifierAccountID][0], nil
}

// HandlerRLSGameIdentifierValues returns the list of game identifier the current authenticated user has access to
// Note: game_id is no longer set as an RLS identifier - access is controlled via RLS constraints
func HandlerRLSGameIdentifierValues(l logger.Logger, r *http.Request) ([]string, error) {
	// Game access is now controlled via RLS constraints on game_id column
	// Return empty slice as identifiers are no longer used for games
	return []string{}, nil
}

// HandlerRLSGameSubscriptionIdentifierValues returns the list of game subscription identifier the current authenticated user has access to
// Note: game_subscription_id is no longer set as an RLS identifier - access is controlled via RLS constraints
func HandlerRLSGameSubscriptionIdentifierValues(l logger.Logger, r *http.Request) ([]string, error) {
	// Game subscription access is now controlled via RLS constraints
	// Return empty slice as identifiers are no longer used for game subscriptions
	return []string{}, nil
}
