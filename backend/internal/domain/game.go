package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetManyAccountGameViewRecs returns account game view records
func (m *Domain) GetManyAccountGameViewRecs(opts *coresql.Options) ([]*game_record.AccountGameView, error) {
	l := m.Logger("GetManyAccountGameViewRecs")

	l.Debug("getting many account_game_view records opts >%#v<", opts)

	r := m.AccountGameViewRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetManyManagerGameInstanceViewRecs returns manager game instance view records
func (m *Domain) GetManyManagerGameInstanceViewRecs(opts *coresql.Options) ([]*game_record.ManagerGameInstanceView, error) {
	l := m.Logger("GetManyManagerGameInstanceViewRecs")

	l.Debug("getting many manager_game_instance_view records opts >%#v<", opts)

	r := m.ManagerGameInstanceViewRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetManyCatalogGameInstanceViewRecs returns catalog game instance view records
func (m *Domain) GetManyCatalogGameInstanceViewRecs(opts *coresql.Options) ([]*game_record.CatalogGameInstanceView, error) {
	l := m.Logger("GetManyCatalogGameInstanceViewRecs")

	l.Debug("getting many catalog_game_instance_view records opts >%#v<", opts)

	r := m.CatalogGameInstanceViewRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetManyGameRecs -
func (m *Domain) GetManyGameRecs(opts *coresql.Options) ([]*game_record.Game, error) {
	l := m.Logger("GetManyGameRecs")

	l.Debug("getting many client records opts >%#v<", opts)

	r := m.GameRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameRec -
func (m *Domain) GetGameRec(recID string, lock *coresql.Lock) (*game_record.Game, error) {
	l := m.Logger("GetGameRec")

	l.Debug("getting client record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGame, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameRec -
func (m *Domain) CreateGameRec(rec *game_record.Game) (*game_record.Game, error) {
	l := m.Logger("CreateGameRec")

	l.Debug("creating client record >%#v<", rec)

	r := m.GameRepository()

	if err := m.validateGameRecForCreate(rec); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return rec, err
	}

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create client record >%v<", err)
		return rec, databaseError(err)
	}

	return createdRec, nil
}

// UpdateGameRec -
func (m *Domain) UpdateGameRec(rec *game_record.Game) (*game_record.Game, error) {
	l := m.Logger("UpdateGameRec")

	curr, err := m.GetGameRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game record ID >%s< >%#v<", rec.ID, rec)

	if err := m.validateGameRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate game record >%v<", err)
		return rec, err
	}

	r := m.GameRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteGameRec -
func (m *Domain) DeleteGameRec(recID string) error {
	l := m.Logger("DeleteGameRec")

	l.Debug("deleting game record ID >%s<", recID)

	rec, err := m.GetGameRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameRepository()

	if err := m.validateGameRecForDelete(rec); err != nil {
		l.Warn("failed to validate game record >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveGameRec -
func (m *Domain) RemoveGameRec(recID string) error {
	l := m.Logger("RemoveGameRec")

	l.Debug("removing game record ID >%s<", recID)

	rec, err := m.GetGameRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameRepository()

	if err := m.validateGameRecForDelete(rec); err != nil {
		l.Warn("failed to validate game record >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// CreateGame -
// Note: Games no longer have account_id. Designer subscriptions must be created separately
// via the game_subscription handler after game creation.
func (m *Domain) CreateGame(rec *game_record.Game) (*game_record.Game, *game_record.GameSubscription, error) {
	l := m.Logger("CreateGame")

	l.Debug("creating game >%#v<", rec)

	createdRec, err := m.CreateGameRec(rec)
	if err != nil {
		l.Warn("failed to create game record >%v<", err)
		return nil, nil, err
	}

	// Games no longer have owners - designer subscriptions are created separately
	// Return nil for subscription to indicate it should be created via handler
	return createdRec, nil, nil
}
