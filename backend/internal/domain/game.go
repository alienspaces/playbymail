package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameRecs -
func (m *Domain) GetManyGameRecs(opts *coresql.Options) ([]*record.Game, error) {
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
func (m *Domain) GetGameRec(recID string, lock *coresql.Lock) (*record.Game, error) {
	l := m.Logger("GetGameRec")

	l.Debug("getting client record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGame, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameRec -
func (m *Domain) CreateGameRec(rec *record.Game) (*record.Game, error) {
	l := m.Logger("CreateGameRec")

	l.Debug("creating client record >%#v<", rec)

	r := m.GameRepository()

	if err := m.validateGameRecForCreate(rec); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return rec, err
	}

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateGameRec -
func (m *Domain) UpdateGameRec(next *record.Game) (*record.Game, error) {
	l := m.Logger("UpdateGameRec")

	curr, err := m.GetGameRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating client record >%#v<", next)

	if err := m.validateGameRecForUpdate(next, curr); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return next, err
	}

	r := m.GameRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}

	return next, nil
}

// DeleteGameRec -
func (m *Domain) DeleteGameRec(recID string) error {
	l := m.Logger("DeleteGameRec")

	l.Debug("deleting client record ID >%s<", recID)

	rec, err := m.GetGameRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameRepository()

	if err := m.validateGameRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
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

	l.Debug("removing client record ID >%s<", recID)

	rec, err := m.GetGameRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameRepository()

	if err := m.validateGameRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
