package harness

import (
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/harness"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// Testing -
type Testing struct {
	*harness.Testing
	Data         Data
	teardownData Data
	DataConfig   DataConfig
	Config       config.Config
}

// NewTesting -
func NewTesting(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], cfg config.Config, dcfg DataConfig) (t *Testing, err error) {

	h, err := harness.NewTesting(l, s, j)
	if err != nil {
		return nil, err
	}

	t = &Testing{
		Testing: h,
		Config:  cfg,
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

	for _, accountConfig := range t.DataConfig.AccountConfigs {
		accountRec, err := t.createAccountRec(accountConfig)
		if err != nil {
			l.Warn("failed creating account record >%v<", err)
			return err
		}
		l.Debug("created account record ID >%s< Email >%s<", accountRec.ID, accountRec.Email)

		// Create account contact for the account (required for player subscriptions)
		_, err = t.createAccountContactRec(accountRec.ID)
		if err != nil {
			l.Warn("failed creating account contact record >%v<", err)
			return err
		}
		l.Debug("created account contact record for account ID >%s<", accountRec.ID)
	}

	for _, gameConfig := range t.DataConfig.GameConfigs {

		// Create game record
		gameRec, err := t.createGameRec(gameConfig)
		if err != nil {
			l.Warn("failed creating game record >%v<", err)
			return err
		}
		l.Debug("created game record ID >%s< Name >%s<", gameRec.ID, gameRec.Name)

		// Create game subscription records for this game
		for _, subscriptionConfig := range gameConfig.GameSubscriptionConfigs {
			_, err = t.createGameSubscriptionRec(subscriptionConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game_subscription record >%v<", err)
				return err
			}
			l.Debug("created game_subscription record for game >%s<", gameRec.ID)
		}

		// ------------------------------------------------------------
		// Adventure game specific records for this game
		// ------------------------------------------------------------

		err = t.createAdventureGameRecords(gameConfig, gameRec)
		if err != nil {
			l.Warn("failed creating adventure game records >%v<", err)
			return err
		}
		l.Debug("created adventure game records for game >%s<", gameRec.ID)

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

			err = t.createAdventureGameInstanceRecords(gameInstanceConfig, gameInstanceRec)
			if err != nil {
				l.Warn("failed creating adventure game instance records >%v<", err)
				return err
			}
			l.Debug("created adventure game instance records for game instance >%s<", gameInstanceRec.ID)
		}
	}

	l.Debug("created test data")

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

	// ------------------------------------------------------------
	// Adventure game specific records
	// ------------------------------------------------------------

	err = t.removeAdventureGameRecords()
	if err != nil {
		l.Warn("failed removing adventure game records >%v<", err)
		return err
	}
	l.Debug("removed adventure game records")

	// Remove game subscription records
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

	// Remove game image records
	l.Debug("removing >%d< game image records", len(t.teardownData.GameImageRecs))
	for _, imageRec := range t.teardownData.GameImageRecs {
		l.Debug("[teardown] game image ID: >%s<", imageRec.ID)
		if imageRec.ID == "" {
			l.Warn("[teardown] skipping game image with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).DeleteGameImageRec(imageRec.ID)
		if err != nil {
			l.Warn("failed removing game image record >%v<", err)
			return err
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

	// Remove accounts
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
