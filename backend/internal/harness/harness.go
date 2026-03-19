package harness

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/harness"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
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

	if err := t.createAccounts(); err != nil {
		l.Warn("failed creating accounts >%v<", err)
		return err
	}

	if err := t.createGames(); err != nil {
		l.Warn("failed creating games >%v<", err)
		return err
	}

	if err := t.createGameSubscriptions(); err != nil {
		l.Warn("failed creating game subscriptions >%v<", err)
		return err
	}

	if err := t.createGameInstances(); err != nil {
		l.Warn("failed creating game instances >%v<", err)
		return err
	}

	l.Debug("created test data")

	return nil
}

func (t *Testing) createAccounts() error {
	l := t.Logger("createAccounts")

	// Create all account, account user, account user contact, and account subscription records
	var allAccountRecs []*account_record.Account
	var allAccountUserRecs []*account_record.AccountUser
	var allAccountUserContactRecs []*account_record.AccountUserContact
	var allAccountSubscriptionRecs []*account_record.AccountSubscription

	for _, accountConfig := range t.DataConfig.AccountConfigs {

		// Create account, account user, account user contact, and account subscription records
		accountRec, accountUserRecs, accountUserContactRecs, accountSubscriptionRecs, err := t.processAccountConfig(accountConfig)
		if err != nil {
			l.Warn("failed processing account config >%v<", err)
			return err
		}

		allAccountRecs = append(allAccountRecs, accountRec)
		allAccountUserRecs = append(allAccountUserRecs, accountUserRecs...)
		allAccountUserContactRecs = append(allAccountUserContactRecs, accountUserContactRecs...)
		allAccountSubscriptionRecs = append(allAccountSubscriptionRecs, accountSubscriptionRecs...)
	}

	l.Debug("created >%d< account records", len(allAccountRecs))
	l.Debug("created >%d< account user records", len(allAccountUserRecs))
	l.Debug("created >%d< account user contact records", len(allAccountUserContactRecs))
	l.Debug("created >%d< account subscription records", len(allAccountSubscriptionRecs))

	return nil
}

func (t *Testing) createGames() error {
	l := t.Logger("createGames")

	var allGameRecs []*game_record.Game
	var allGameImageRecs []*game_record.GameImage

	for i := range t.DataConfig.GameConfigs {
		gameRec, gameImageRecs, err := t.processGameConfig(t.DataConfig.GameConfigs[i])
		if err != nil {
			l.Warn("failed processing game config >%v<", err)
			return err
		}

		allGameRecs = append(allGameRecs, gameRec)
		allGameImageRecs = append(allGameImageRecs, gameImageRecs...)

		// Process adventure game config
		adventureGameRecs, err := t.processAdventureGameConfig(t.DataConfig.GameConfigs[i], gameRec)
		if err != nil {
			l.Warn("failed processing adventure game config >%v<", err)
			return err
		}

		l.Debug("created >%d< adventure game item records", len(adventureGameRecs.Items))
		l.Debug("created >%d< adventure game location records", len(adventureGameRecs.Locations))
		l.Debug("created >%d< adventure game creature records", len(adventureGameRecs.Creatures))
		l.Debug("created >%d< adventure game location link records", len(adventureGameRecs.LocationLinks))
		l.Debug("created >%d< adventure game location link requirement records", len(adventureGameRecs.LocationLinkRequirements))
		l.Debug("created >%d< adventure game character records", len(adventureGameRecs.Characters))
	}

	l.Debug("created >%d< game records", len(allGameRecs))
	l.Debug("created >%d< game image records", len(allGameImageRecs))

	return nil
}

func (t *Testing) createGameInstances() error {
	l := t.Logger("createGameInstances")

	var count int
	for i := range t.DataConfig.AccountUserGameSubscriptionConfigs {
		subConfig := &t.DataConfig.AccountUserGameSubscriptionConfigs[i]
		if subConfig.SubscriptionType != game_record.GameSubscriptionTypeManager || len(subConfig.GameInstanceConfigs) == 0 {
			continue
		}

		managerSubRec, err := t.Data.GetGameSubscriptionRecByRef(subConfig.Reference)
		if err != nil {
			l.Warn("failed getting manager subscription by ref >%s< >%v<", subConfig.Reference, err)
			return err
		}

		gameRec, err := t.Data.GetGameRecByRef(subConfig.GameRef)
		if err != nil {
			l.Warn("failed getting game record by reference >%s< >%v<", subConfig.GameRef, err)
			return err
		}

		gameConfig := t.findGameConfigByRef(subConfig.GameRef)
		if gameConfig == nil {
			l.Warn("game config for ref >%s< not found", subConfig.GameRef)
			return fmt.Errorf("game config for ref >%s< not found", subConfig.GameRef)
		}

		for _, gameInstanceConfig := range subConfig.GameInstanceConfigs {

			gameInstanceRec, err := t.createGameInstanceRec(gameInstanceConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game_instance record >%v<", err)
				return err
			}

			_, err = t.createGameSubscriptionInstanceRec(managerSubRec, gameInstanceRec)
			if err != nil {
				l.Warn("failed creating game_subscription_instance link >%v<", err)
				return err
			}

			for _, playerSubRef := range gameInstanceConfig.PlayerSubscriptionRefs {

				playerSubRec, err := t.Data.GetGameSubscriptionRecByRef(playerSubRef)
				if err != nil {
					l.Warn("failed getting player subscription by ref >%s< >%v<", playerSubRef, err)
					return err
				}

				_, err = t.createGameSubscriptionInstanceRec(playerSubRec, gameInstanceRec)
				if err != nil {
					l.Warn("failed creating game_subscription_instance link for player >%s< >%v<", playerSubRef, err)
					return err
				}
			}

			l.Debug("created game_instance record ID >%s< linked to manager subscription", gameInstanceRec.ID)

			for _, ipc := range gameInstanceConfig.GameInstanceParameterConfigs {
				if _, err := t.createGameInstanceParameterRec(ipc, gameInstanceRec); err != nil {
					l.Warn("failed creating game_instance_parameter record >%v<", err)
					return err
				}
			}

			if err = t.createAdventureGameInstanceRecords(gameInstanceConfig, gameInstanceRec); err != nil {
				l.Warn("failed creating adventure game instance records >%v<", err)
				return err
			}

			count++
		}
	}

	l.Debug("created >%d< game instance records", count)

	return nil
}

// findGameConfigByRef returns the GameConfig with the given reference, or nil if not found.
func (t *Testing) findGameConfigByRef(gameRef string) *GameConfig {
	for i := range t.DataConfig.GameConfigs {
		if t.DataConfig.GameConfigs[i].Reference == gameRef {
			return &t.DataConfig.GameConfigs[i]
		}
	}
	return nil
}

// RemoveData - Uses the domain remove methods to physically remove data
// from the database as opposed to the delete methods which only mark the
// record as deleted.
func (t *Testing) RemoveData() error {
	l := t.Logger("RemoveData")

	dom := t.Domain.(*domain.Domain)

	// Remove each game instance and all its associated data (turn sheets, instance records,
	// parameters, subscription instance links) via RemoveGameInstance.
	l.Debug("removing >%d< game instance records", len(t.teardownData.GameInstanceRecs))
	for _, instanceRec := range t.teardownData.GameInstanceRecs {
		l.Debug("[teardown] game instance ID: >%s<", instanceRec.ID)
		if instanceRec.ID == "" {
			l.Warn("[teardown] skipping game instance with empty ID")
			continue
		}
		if err := dom.RemoveGameInstance(instanceRec.ID); err != nil {
			l.Warn("failed removing game instance >%s< >%v<", instanceRec.ID, err)
			return err
		}
	}

	// Remove game subscription records (subscription_instance links are removed by RemoveGameInstance above)
	l.Debug("removing >%d< game subscription records", len(t.teardownData.GameSubscriptionRecs))
	for _, subscriptionRec := range t.teardownData.GameSubscriptionRecs {
		l.Debug("[teardown] game subscription ID: >%s<", subscriptionRec.ID)
		if subscriptionRec.ID == "" {
			l.Warn("[teardown] skipping game subscription with empty ID")
			continue
		}
		if err := dom.RemoveGameSubscriptionRec(subscriptionRec.ID); err != nil {
			l.Warn("failed removing game subscription record >%v<", err)
			return err
		}
	}

	// ------------------------------------------------------------
	// Adventure game specific records
	// ------------------------------------------------------------

	if err := t.removeAdventureGameRecords(); err != nil {
		l.Warn("failed removing adventure game records >%v<", err)
		return err
	}
	l.Debug("removed adventure game records")

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
