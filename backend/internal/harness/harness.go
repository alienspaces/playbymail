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
)

// Testing -
type Testing struct {
	harness.Testing
	Data         Data
	teardownData Data
	DataConfig   DataConfig
}

// NewTesting -
func NewTesting(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx], config DataConfig) (t *Testing, err error) {

	t = &Testing{
		Testing: harness.Testing{
			Log:       l,
			Store:     s,
			JobClient: j,
		},
	}

	// Require service config, logger and store
	if t.Log == nil || t.Store == nil {
		return nil, fmt.Errorf("missing logger >%v< or storer >%v<, cannot create new test harness", t.Log, t.Store)
	}

	// domainer
	t.DomainFunc = t.domainer

	// data
	t.CreateDataFunc = t.CreateData
	t.RemoveDataFunc = t.RemoveData

	t.DataConfig = config
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

	m, err := domain.NewDomain(t.Log, t.JobClient)
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
	}

	for _, gameConfig := range t.DataConfig.GameConfigs {
		gameRec, err := t.createGameRec(gameConfig)
		if err != nil {
			l.Warn("failed creating game record >%v<", err)
			return err
		}
		l.Debug("created game record ID >%s< Name >%s<", gameRec.ID, gameRec.Name)

		for _, itemConfig := range gameConfig.GameItemConfigs {
			_, err = t.createGameItemRec(itemConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game_item record >%v<", err)
				return err
			}
			l.Debug("created game_item record for game >%s<", gameRec.ID)
		}

		for _, gameLocationConfig := range gameConfig.GameLocationConfigs {
			gameLocationRec, err := t.createGameLocationRec(gameLocationConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game location record >%v<", err)
				return err
			}
			l.Debug("created game location record ID >%s< Name >%s<", gameLocationRec.ID, gameLocationRec.Name)
		}

		for _, creatureConfig := range gameConfig.GameCreatureConfigs {
			creatureRec, err := t.createGameCreatureRec(creatureConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game creature record >%v<", err)
				return err
			}
			l.Debug("created game creature record ID >%s< Name >%s<", creatureRec.ID, creatureRec.Name)
		}

		for _, linkConfig := range gameConfig.GameLocationLinkConfigs {
			gameLocationLinkRec, err := t.createGameLocationLinkRec(linkConfig, gameRec)
			if err != nil {
				l.Warn("failed creating location link record >%v<", err)
				return err
			}
			l.Debug("created location link record ID >%s<", gameLocationLinkRec.ID)

			for _, reqConfig := range linkConfig.GameLocationLinkRequirementConfigs {
				_, err = t.createGameLocationLinkRequirementRec(reqConfig, gameLocationLinkRec)
				if err != nil {
					l.Warn("failed creating game_location_link_requirement record >%v<", err)
					return err
				}
				l.Debug("created game_location_link_requirement record for game >%s<", gameRec.ID)
			}
		}

		for _, charConfig := range gameConfig.GameCharacterConfigs {
			_, err = t.createGameCharacterRec(charConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game_character record >%v<", err)
				return err
			}
			l.Debug("created game_character record for game >%s<", gameRec.ID)
		}

		// Instance records

		// Create game instances for this game
		for _, gameInstanceConfig := range gameConfig.GameInstanceConfigs {
			gameInstanceRec, err := t.createGameInstanceRec(gameInstanceConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game_instance record >%v<", err)
				return err
			}
			l.Debug("created game_instance record ID >%s<", gameInstanceRec.ID)

			// Create location instances for this game instance
			for _, locationInstanceConfig := range gameInstanceConfig.GameLocationInstanceConfigs {
				locationInstanceRec, err := t.createGameLocationInstanceRec(locationInstanceConfig, gameInstanceRec)
				if err != nil {
					l.Warn("failed creating game_location_instance record >%v<", err)
					return err
				}
				l.Debug("created game_location_instance record ID >%s<", locationInstanceRec.ID)
			}

			// Create creature instances for this game instance
			for _, creatureInstanceConfig := range gameInstanceConfig.GameCreatureInstanceConfigs {
				creatureInstanceRec, err := t.createGameCreatureInstanceRec(creatureInstanceConfig, gameInstanceRec)
				if err != nil {
					l.Warn("failed creating game_creature_instance record >%v<", err)
					return err
				}
				l.Debug("created game_creature_instance record ID >%s<", creatureInstanceRec.ID)
			}

			// Create item instances for this game instance
			for _, itemInstanceConfig := range gameInstanceConfig.GameItemInstanceConfigs {
				itemInstanceRec, err := t.createGameItemInstanceRec(itemInstanceConfig, gameInstanceRec)
				if err != nil {
					l.Warn("failed creating game_item_instance record >%v<", err)
					return err
				}
				l.Debug("created game_item_instance record ID >%s<", itemInstanceRec.ID)
			}
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

	// Remove instance records

	// Remove game creature instances
	l.Debug("removing >%d< game creature instance records", len(t.teardownData.GameCreatureInstanceRecs))
	for _, creatureInstanceRec := range t.teardownData.GameCreatureInstanceRecs {
		l.Debug("[teardown] game creature instance ID: >%s<", creatureInstanceRec.ID)
		if creatureInstanceRec.ID == "" {
			l.Warn("[teardown] skipping game creature instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameCreatureInstanceRec(creatureInstanceRec.ID)
		if err != nil {
			l.Warn("failed removing game creature instance record >%v<", err)
			return err
		}
	}

	// Remove game item instances
	l.Debug("removing >%d< game item instance records", len(t.teardownData.GameItemInstanceRecs))
	for _, itemInstanceRec := range t.teardownData.GameItemInstanceRecs {
		l.Debug("[teardown] game item instance ID: >%s<", itemInstanceRec.ID)
		if itemInstanceRec.ID == "" {
			l.Warn("[teardown] skipping game item instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameItemInstanceRec(itemInstanceRec.ID)
		if err != nil {
			l.Warn("failed removing game item instance record >%v<", err)
			return err
		}
	}

	// Remove game location instances before game instances to avoid FK errors
	l.Debug("removing >%d< game location instance records", len(t.teardownData.GameLocationInstanceRecs))
	for _, locationInstanceRec := range t.teardownData.GameLocationInstanceRecs {
		l.Debug("[teardown] game location instance ID: >%s<", locationInstanceRec.ID)
		if locationInstanceRec.ID == "" {
			l.Warn("[teardown] skipping game location instance with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameLocationInstanceRec(locationInstanceRec.ID)
		if err != nil {
			l.Warn("failed removing game location instance record >%v<", err)
			return err
		}
	}

	// Remove game instance records before games to avoid FK errors
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

	// Remove game creature records
	l.Debug("removing >%d< game creature records", len(t.teardownData.GameCreatureRecs))
	for _, creatureRec := range t.teardownData.GameCreatureRecs {
		l.Debug("[teardown] game creature ID: >%s<", creatureRec.ID)
		if creatureRec.ID == "" {
			l.Warn("[teardown] skipping game creature with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameCreatureRec(creatureRec.ID)
		if err != nil {
			l.Warn("failed removing game creature record >%v<", err)
			return err
		}
	}

	// Remove game location link requirements
	l.Debug("removing >%d< game location link requirement records", len(t.teardownData.GameLocationLinkRequirementRecs))
	for _, reqRec := range t.teardownData.GameLocationLinkRequirementRecs {
		l.Debug("[teardown] game location link requirement ID: >%s<", reqRec.ID)
		if reqRec.ID == "" {
			l.Warn("[teardown] skipping game location link requirement with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameLocationLinkRequirementRec(reqRec.ID)
		if err != nil {
			l.Warn("failed removing game location link requirement record >%v<", err)
			return err
		}
	}

	// Remove game location links
	l.Debug("removing >%d< game location link records", len(t.teardownData.GameLocationLinkRecs))
	for _, linkRec := range t.teardownData.GameLocationLinkRecs {
		l.Debug("[teardown] game location link ID: >%s<", linkRec.ID)
		if linkRec.ID == "" {
			l.Warn("[teardown] skipping game location link with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameLocationLinkRec(linkRec.ID)
		if err != nil {
			l.Warn("failed removing game location link record >%v<", err)
			return err
		}
	}

	// Remove game location records before games to avoid FK errors
	l.Debug("removing >%d< game location records", len(t.teardownData.GameLocationRecs))
	for _, gameLocationRec := range t.teardownData.GameLocationRecs {
		l.Debug("[teardown] game location ID: >%s<", gameLocationRec.ID)
		if gameLocationRec.ID == "" {
			l.Warn("[teardown] skipping game location with empty ID")
			continue
		}
		l.Debug("removing game location record ID >%s<", gameLocationRec.ID)
		err := t.Domain.(*domain.Domain).RemoveGameLocationRec(gameLocationRec.ID)
		if err != nil {
			l.Warn("failed removing game location record >%v<", err)
			return err
		}
	}

	// Remove game item records before games to avoid FK errors
	l.Debug("removing >%d< game item records", len(t.teardownData.GameItemRecs))
	for _, itemRec := range t.teardownData.GameItemRecs {
		l.Debug("[teardown] game item ID: >%s<", itemRec.ID)
		if itemRec.ID == "" {
			l.Warn("[teardown] skipping game item with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameItemRec(itemRec.ID)
		if err != nil {
			l.Warn("failed removing game item record >%v<", err)
			return err
		}
	}

	// Remove game character records before games to avoid FK errors
	l.Debug("removing >%d< game character records", len(t.teardownData.GameCharacterRecs))
	for _, charRec := range t.teardownData.GameCharacterRecs {
		l.Debug("[teardown] game character ID: >%s<", charRec.ID)
		if charRec.ID == "" {
			l.Warn("[teardown] skipping game character with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameCharacterRec(charRec.ID)
		if err != nil {
			l.Warn("failed removing game character record >%v<", err)
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
	return nil
}

// Logger - Returns a logger with package context and provided function context
func (t *Testing) Logger(functionName string) logger.Logger {
	return t.Log.WithPackageContext("harness").WithFunctionContext(functionName)
}
