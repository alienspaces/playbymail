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
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data         Data
	DataConfig   DataConfig
	teardownData teardownData
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
	t.teardownData = teardownData{}

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
	t.teardownData = teardownData{}

	l.Info("creating test data")

	for _, gameConfig := range t.DataConfig.GameConfig {
		gameRec, err := t.createGameRec(gameConfig)
		if err != nil {
			l.Warn("failed creating game record >%v<", err)
			return err
		}
		l.Debug("created game record ID >%s< Name >%s<", gameRec.ID, gameRec.Name)
		t.Data.AddGameRec(gameRec)
		t.teardownData.AddGameRec(gameRec)
	}

	l.Info("created test data")

	return nil
}

type teardownData struct {
	GameRecs []*record.Game
}

func (t *teardownData) AddGameRec(rec *record.Game) {
	for idx := range t.GameRecs {
		if t.GameRecs[idx].ID == rec.ID {
			t.GameRecs[idx] = rec
			return
		}
	}
	t.GameRecs = append(t.GameRecs, rec)
}

// RemoveData -
func (t *Testing) RemoveData() error {
	l := t.Logger("RemoveData")

	// Quick cleanup when data is not committed
	if !t.ShouldCommitData {
		t.Data = Data{}
		t.teardownData = teardownData{}
		return nil
	}

	l.Info("Removing test data")

	seen := map[string]bool{}

	l.Debug("Removing >%d< game records", len(t.teardownData.GameRecs))

	for _, rec := range t.teardownData.GameRecs {
		if seen[rec.ID] {
			continue
		}
		err := t.Domain.(*domain.Domain).RemoveGameRec(rec.ID)
		if err != nil {
			l.Warn("failed removing game record >%v<", err)
			return err
		}
		seen[rec.ID] = true
	}

	l.Debug("Removing >%d< game records", len(t.teardownData.GameRecs))

	t.Data = Data{}

	l.Info("Removed test data")

	return nil
}

// Logger - Returns a logger with package context and provided function context
func (t *Testing) Logger(functionName string) logger.Logger {
	return t.Log.WithPackageContext("harness").WithFunctionContext(functionName)
}
