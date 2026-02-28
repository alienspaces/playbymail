package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/repository"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

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

// GetGameRecByIDForJoinProcess retrieves a game record by ID without RLS filtering.
// This is intended for use when processing a join game turn sheet, where the turn sheet
// code itself serves as authorization and the caller needs access to the game record
// regardless of their own subscription status.
func (m *Domain) GetGameRecByIDForJoinProcess(recID string) (*game_record.Game, error) {
	l := m.Logger("GetGameRecByIDForJoinProcess")
	l.Debug("getting game record by ID for join process (no RLS) >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r, err := repository.NewGeneric[game_record.Game](repository.NewArgs{
		Tx:            m.Tx,
		TableName:     game_record.TableGame,
		Record:        game_record.Game{},
		IsRLSDisabled: true,
	})
	if err != nil {
		return nil, err
	}
	rec, err := r.GetOne(recID, nil)
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
