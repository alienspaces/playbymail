package harness

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/harness"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// Testing -
type Testing struct {
	*harness.Testing
	Config config.Config
	// Scanner is the implementation of the turn sheet scanner interface
	Scanner      turnsheet.TurnSheetScanner
	DataConfig   DataConfig
	Data         Data
	teardownData Data
}

// NewTesting -
func NewTesting(cfg config.Config, l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], scanner turnsheet.TurnSheetScanner, dcfg DataConfig) (t *Testing, err error) {

	h, err := harness.NewTesting(l, s, j)
	if err != nil {
		return nil, err
	}

	t = &Testing{
		Testing: h,
		Config:  cfg,
		Scanner: scanner,
	}

	// domainer
	t.DomainFunc = t.domainer

	// data
	t.CreateDataFunc = t.CreateData
	t.RemoveDataFunc = t.RemoveData

	t.DataConfig = dcfg
	t.Data = Data{}
	t.teardownData = Data{}

	err = t.Init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// domainer -
func (t *Testing) domainer() (domainer.Domainer, error) {
	l := t.Logger("domainer")

	m, err := domain.NewDomain(t.Log, t.Config)
	if err != nil {
		l.Warn("failed new domain >%v<", err)
		return nil, err
	}

	return m, nil
}

// CreateData - Custom data
func (t *Testing) CreateData() error {
	l := t.Logger("CreateData")

	t.Data = initialiseDataStores()
	t.teardownData = initialiseTeardownDataStores()

	l.Debug("creating test data")

	// First pass: Create all account users
	for _, accountConfig := range t.DataConfig.AccountConfigs {
		accountRec, err := t.createAccountUserRec(accountConfig)
		if err != nil {
			l.Warn("failed creating account record >%v<", err)
			return err
		}
		l.Debug("created account record ID >%s< Email >%s<", accountRec.ID, accountRec.Email)

		// Create account contact for the account (required for player subscriptions)
		_, err = t.createAccountUserContactRec(accountRec.ID)
		if err != nil {
			l.Warn("failed creating account contact record >%v<", err)
			return err
		}
		l.Debug("created account contact record for account ID >%s<", accountRec.ID)
	}

	// Second pass: Create all games (without instances yet)
	for _, gameConfig := range t.DataConfig.GameConfigs {

		// Create game record
		gameRec, err := t.createGameRec(gameConfig)
		if err != nil {
			l.Warn("failed creating game record >%v<", err)
			return err
		}
		l.Debug("created game record ID >%s< Name >%s<", gameRec.ID, gameRec.Name)

		// ------------------------------------------------------------
		// Adventure game specific records for this game
		// ------------------------------------------------------------

		err = t.createAdventureGameRecords(gameConfig, gameRec)
		if err != nil {
			l.Warn("failed creating adventure game records >%v<", err)
			return err
		}
		l.Debug("created adventure game records for game >%s<", gameRec.ID)
	}

	// Third pass: Create designer game subscriptions first (they don't require game instances)
	// These are needed to create game instances which require a subscription
	for _, accountConfig := range t.DataConfig.AccountConfigs {
		accountRec, err := t.Data.GetAccountUserRecByRef(accountConfig.Reference)
		if err != nil {
			l.Warn("failed getting account user record for reference >%s<: %v", accountConfig.Reference, err)
			return err
		}

		// Create designer subscriptions first (they don't need game instances)
		for _, subscriptionConfig := range accountConfig.GameSubscriptionConfigs {
			if subscriptionConfig.SubscriptionType == game_record.GameSubscriptionTypeDesigner {
				_, err = t.createGameSubscriptionRec(subscriptionConfig, accountRec)
				if err != nil {
					l.Warn("failed creating designer game_subscription record >%v<", err)
					return err
				}
				l.Debug("created designer game_subscription record for account >%s<", accountRec.ID)
			}
		}
	}

	// Fourth pass: Create all game instances (now that designer subscriptions exist)
	for _, gameConfig := range t.DataConfig.GameConfigs {
		gameRec, err := t.Data.GetGameRecByRef(gameConfig.Reference)
		if err != nil {
			l.Warn("failed getting game record for reference >%s<: %v", gameConfig.Reference, err)
			return err
		}

		// Create game instance records for this game
		for _, gameInstanceConfig := range gameConfig.GameInstanceConfigs {
			gameInstanceRec, err := t.createGameInstanceRec(gameInstanceConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game_instance record >%v<", err)
				return err
			}
			l.Debug("created game_instance record ID >%s<", gameInstanceRec.ID)

			// Create game instance parameter records for this game instance
			for _, instanceParameterConfig := range gameInstanceConfig.GameInstanceParameterConfigs {
				instanceParameterRec, err := t.createGameInstanceParameterRec(instanceParameterConfig, gameInstanceRec)
				if err != nil {
					l.Warn("failed creating game_instance_parameter record >%v<", err)
					return err
				}
				l.Debug("created game_instance_parameter record ID >%s<", instanceParameterRec.ID)
			}

			// ------------------------------------------------------------
			// Adventure game specific instance records
			// ------------------------------------------------------------

			err = t.createAdventureGameInstanceRecords(gameConfig, gameInstanceConfig, gameInstanceRec)
			if err != nil {
				l.Warn("failed creating adventure game instance records >%v<", err)
				return err
			}
			l.Debug("created adventure game instance records for game instance >%s<", gameInstanceRec.ID)
		}
	}

	// Fifth pass: Create manager and player subscriptions (now that game instances exist)
	// These subscriptions require game instances to be created first
	for _, accountConfig := range t.DataConfig.AccountConfigs {
		accountRec, err := t.Data.GetAccountUserRecByRef(accountConfig.Reference)
		if err != nil {
			l.Warn("failed getting account user record for reference >%s<: %v", accountConfig.Reference, err)
			return err
		}

		// Create manager and player subscriptions (they require game instances)
		for _, subscriptionConfig := range accountConfig.GameSubscriptionConfigs {
			if subscriptionConfig.SubscriptionType == game_record.GameSubscriptionTypeManager ||
				subscriptionConfig.SubscriptionType == game_record.GameSubscriptionTypePlayer {
				_, err = t.createGameSubscriptionRec(subscriptionConfig, accountRec)
				if err != nil {
					l.Warn("failed creating manager/player game_subscription record >%v<", err)
					return err
				}
				l.Debug("created %s game_subscription record for account >%s<", subscriptionConfig.SubscriptionType, accountRec.ID)
			}
		}
	}

	// TODO: Remove me
	for _, rec := range t.Data.AccountUserRecs {
		l.Info("created account user record ID >%s<", rec.ID)
	}

	l.Debug("created test data")

	// Force commit to ensure data is visible to handlers running in separate transactions
	if t.ShouldCommitData {
		l.Info("forcing commit of harness data")
		if err := t.CommitTx(); err != nil {
			l.Warn("failed force commit >%v<", err)
			return err
		}

		// Re-open transaction to satisfy callers expecting an open transaction (like Setup wrapper)
		tx, err := t.InitTx()
		if err != nil {
			l.Warn("failed to re-init tx >%v<", err)
			return err
		}
		l.Info("re-opened transaction for test harness")

		// Verify game visibility
		for _, gameID := range t.Data.Refs.GameRefs {
			var count int
			err := tx.QueryRow(context.Background(), "SELECT count(*) FROM game WHERE id = $1", gameID).Scan(&count)
			if err != nil {
				l.Warn("failed to query game count >%v<", err)
			} else {
				l.Info("VERIFICATION: Game ID >%s< Count >%d< (Should be 1)", gameID, count)
			}
		}
	}

	return nil
}

// RemoveData - Uses the domain remove methods to physically remove data
// from the database as opposed to the delete methods which only mark the
// record as deleted.
func (t *Testing) RemoveData() error {
	l := t.Logger("RemoveData")

	// ------------------------------------------------------------
	// Adventure game specific instance records
	// ------------------------------------------------------------

	err := t.removeAdventureGameInstanceRecords()
	if err != nil {
		l.Warn("failed removing adventure game instance records >%v<", err)
		return err
	}
	l.Debug("removed adventure game instance records")

	// Remove game turn sheet records first (they reference game instances)
	l.Debug("removing >%d< game turn sheet records", len(t.teardownData.GameTurnSheetRecs))
	for _, turnSheetRec := range t.teardownData.GameTurnSheetRecs {
		l.Debug("[teardown] game turn sheet ID: >%s<", turnSheetRec.ID)
		if turnSheetRec.ID == "" {
			l.Warn("[teardown] skipping game turn sheet with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameTurnSheetRec(turnSheetRec.ID)
		if err != nil {
			l.Warn("failed removing game turn sheet record >%v<", err)
			return err
		}
	}

	// Remove game instance parameter records
	l.Debug("removing >%d< game instance parameter records", len(t.teardownData.GameInstanceParameterRecs))
	for _, parameterRec := range t.teardownData.GameInstanceParameterRecs {
		l.Debug("[teardown] game instance parameter ID: >%s<", parameterRec.ID)
		if parameterRec.ID == "" {
			l.Warn("[teardown] skipping game instance parameter with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameInstanceParameterRec(parameterRec.ID)
		if err != nil {
			l.Warn("failed removing game instance parameter record >%v<", err)
			return err
		}
	}

	// Remove game subscription records BEFORE game instances
	// (game_subscription_instance links subscriptions to instances)
	// Remove game subscription instance records (before subscriptions)
	l.Debug("removing >%d< game subscription instance records", len(t.teardownData.GameSubscriptionInstanceRecs))
	for _, instanceLinkRec := range t.teardownData.GameSubscriptionInstanceRecs {
		l.Debug("[teardown] game subscription instance ID: >%s<", instanceLinkRec.ID)
		if instanceLinkRec.ID == "" {
			l.Warn("[teardown] skipping game subscription instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameSubscriptionInstanceRec(instanceLinkRec.ID)
		if err != nil {
			l.Warn("failed removing game subscription instance record >%v<", err)
			// Don't return error - continue with other records
		}
	}

	l.Debug("removing >%d< game subscription records", len(t.teardownData.GameSubscriptionRecs))
	for _, subscriptionRec := range t.teardownData.GameSubscriptionRecs {
		l.Debug("[teardown] game subscription ID: >%s<", subscriptionRec.ID)
		if subscriptionRec.ID == "" {
			l.Warn("[teardown] skipping game subscription with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameSubscriptionRec(subscriptionRec.ID)
		if err != nil {
			l.Warn("failed removing game subscription record >%v<", err)
			return err
		}
	}

	// ------------------------------------------------------------
	// Adventure game specific records
	// ------------------------------------------------------------

	err = t.removeAdventureGameRecords()
	if err != nil {
		l.Warn("failed removing adventure game records >%v<", err)
		return err
	}
	l.Debug("removed adventure game records")

	// Remove game instance records
	l.Debug("removing >%d< game instance records", len(t.teardownData.GameInstanceRecs))
	for _, instanceRec := range t.teardownData.GameInstanceRecs {
		l.Debug("[teardown] game instance ID: >%s<", instanceRec.ID)
		if instanceRec.ID == "" {
			l.Warn("[teardown] skipping game instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameInstanceRec(instanceRec.ID)
		if err != nil {
			l.Warn("failed removing game instance record >%v<", err)
			return err
		}
	}

	// Remove game image records
	l.Debug("removing >%d< game image records", len(t.teardownData.GameImageRecs))
	for _, imageRec := range t.teardownData.GameImageRecs {
		l.Debug("[teardown] game image ID: >%s<", imageRec.ID)
		if imageRec.ID == "" {
			l.Warn("[teardown] skipping game image with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameImageRec(imageRec.ID)
		if err != nil {
			l.Warn("failed removing game image record >%v<", err)
			// Don't return error - continue with other images and games
		}
	}

	// Remove games
	l.Debug("removing >%d< game records", len(t.teardownData.GameRecs))
	for _, gameRec := range t.teardownData.GameRecs {
		l.Debug("[teardown] game ID: >%s<", gameRec.ID)
		if gameRec.ID == "" {
			l.Warn("[teardown] skipping game with empty ID")
			continue
		}
		l.Debug("removing game record ID >%s<", gameRec.ID)
		err := t.Domain.(*domain.Domain).RemoveGameRec(gameRec.ID)
		if err != nil {
			l.Warn("failed removing game record >%v<", err)
			return err
		}
	}

	// Remove account subscriptions before accounts
	l.Debug("removing >%d< account subscription records", len(t.teardownData.AccountSubscriptionRecs))
	for _, rec := range t.teardownData.AccountSubscriptionRecs {
		l.Debug("[teardown] account subscription ID: >%s<", rec.ID)
		if rec.ID == "" {
			l.Warn("[teardown] skipping account subscription with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAccountSubscriptionRec(rec.ID)
		if err != nil {
			l.Warn("failed removing account subscription record >%v<", err)
			return err
		}
	}

	// Remove account user contacts before account users
	l.Debug("removing >%d< account user contact records", len(t.teardownData.AccountUserContactRecs))
	for _, rec := range t.teardownData.AccountUserContactRecs {
		l.Debug("[teardown] account user contact ID: >%s<", rec.ID)
		if rec.ID == "" {
			l.Warn("[teardown] skipping account user contact with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAccountUserContactRec(rec.ID)
		if err != nil {
			l.Warn("failed removing account user contact record >%v<", err)
			return err
		}
	}

	// Remove account users
	l.Debug("removing >%d< account user records", len(t.teardownData.AccountUserRecs))

	for _, rec := range t.teardownData.AccountUserRecs {
		l.Debug("[teardown] account user ID: >%s<", rec.ID)
		if rec.ID == "" {
			l.Warn("[teardown] skipping account user with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAccountUserRec(rec.ID)
		if err != nil {
			l.Warn("failed removing account user record >%v<", err)
			return err
		}
	}

	// Remove account records last — account_user rows reference them via account_id FK
	l.Debug("removing >%d< account records", len(t.teardownData.AccountRecs))

	for _, rec := range t.teardownData.AccountRecs {
		l.Debug("[teardown] account ID: >%s<", rec.ID)
		if rec.ID == "" {
			l.Warn("[teardown] skipping account with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveAccountRec(rec.ID)
		if err != nil {
			l.Warn("failed removing account record >%v<", err)
			return err
		}
	}

	l.Debug("removed test data")

	return nil
}

// AddGameImageRecToTeardown adds a game image record to the teardown data store
// so it will be cleaned up during teardown. This is useful for test cases that
// create images in separate transactions.
func (t *Testing) AddGameImageRecToTeardown(rec *game_record.GameImage) {
	t.teardownData.AddGameImageRec(rec)
}
