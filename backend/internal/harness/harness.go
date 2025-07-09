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

	l.Info("creating test data")

	for _, accountConfig := range t.DataConfig.AccountConfigs {
		accountRec, err := t.createAccountRec(accountConfig)
		if err != nil {
			l.Warn("failed creating account record >%v<", err)
			return err
		}
		l.Info("created account record ID >%s< Email >%s<", accountRec.ID, accountRec.Email)
	}

	for _, gameConfig := range t.DataConfig.GameConfigs {
		gameRec, err := t.createGameRec(gameConfig)
		if err != nil {
			l.Warn("failed creating game record >%v<", err)
			return err
		}
		l.Info("created game record ID >%s< Name >%s<", gameRec.ID, gameRec.Name)

		for _, gameLocationConfig := range gameConfig.GameLocationConfigs {
			gameLocationRec, err := t.createGameLocationRec(gameLocationConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game location record >%v<", err)
				return err
			}
			l.Info("created game location record ID >%s< Name >%s<", gameLocationRec.ID, gameLocationRec.Name)
		}

		for _, linkConfig := range gameConfig.LocationLinkConfigs {
			_, err := t.createLocationLinkRec(linkConfig)
			if err != nil {
				l.Warn("failed creating location link record >%v<", err)
				return err
			}
			l.Info("created location link record for game >%s<", gameRec.ID)
		}

		// Create game_character records for this game
		for _, charConfig := range gameConfig.GameCharacterConfigs {
			_, err = t.createGameCharacterRec(charConfig, gameRec)
			if err != nil {
				l.Warn("failed creating game_character record >%v<", err)
				return err
			}
			l.Info("created game_character record for game >%s<", gameRec.ID)
		}
	}

	l.Info("created test data")

	return nil
}

// RemoveData - Uses the domain remove methods to physically remove data
// from the database as opposed to the delete methods which only mark the
// record as deleted.
func (t *Testing) RemoveData() error {
	l := t.Logger("RemoveData")

	// Remove location links first to avoid foreign key constraint errors
	l.Info("removing >%d< location link records", len(t.teardownData.LocationLinkRecs))
	for _, linkRec := range t.teardownData.LocationLinkRecs {
		l.Info("[teardown] location link ID: >%s<", linkRec.ID)
		if linkRec.ID == "" {
			l.Warn("[teardown] skipping location link with empty ID")
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveLocationLinkRec(linkRec.ID)
		if err != nil {
			l.Warn("failed removing location link record >%v<", err)
			return err
		}
	}

	// Remove locations first to avoid foreign key constraint errors
	l.Info("removing >%d< game location records", len(t.teardownData.GameLocationRecs))

	for _, gameLocationRec := range t.teardownData.GameLocationRecs {
		l.Info("[teardown] game location ID: >%s<", gameLocationRec.ID)
		if gameLocationRec.ID == "" {
			l.Warn("[teardown] skipping game location with empty ID")
			continue
		}
		l.Info("removing game location record ID >%s<", gameLocationRec.ID)
		err := t.Domain.(*domain.Domain).RemoveGameLocationRec(gameLocationRec.ID)
		if err != nil {
			l.Warn("failed removing game location record >%v<", err)
			return err
		}
	}

	// Remove game characters
	l.Info("removing >%d< game character records", len(t.teardownData.GameCharacterRecs))
	for _, charRec := range t.teardownData.GameCharacterRecs {
		l.Info("[teardown] game character ID: >%s<", charRec.ID)
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
	l.Info("removing >%d< game records", len(t.teardownData.GameRecs))

	for _, gameRec := range t.teardownData.GameRecs {
		l.Info("[teardown] game ID: >%s<", gameRec.ID)
		if gameRec.ID == "" {
			l.Warn("[teardown] skipping game with empty ID")
			continue
		}
		l.Info("removing game record ID >%s<", gameRec.ID)
		err := t.Domain.(*domain.Domain).RemoveGameRec(gameRec.ID)
		if err != nil {
			l.Warn("failed removing game record >%v<", err)
			return err
		}
	}

	// Remove accounts
	l.Info("removing >%d< account records", len(t.teardownData.AccountRecs))

	for _, rec := range t.teardownData.AccountRecs {
		l.Info("[teardown] account ID: >%s<", rec.ID)
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
