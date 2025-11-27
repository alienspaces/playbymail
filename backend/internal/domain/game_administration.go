package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetManyGameAdministrationRecs -
func (m *Domain) GetManyGameAdministrationRecs(opts *sql.Options) ([]*game_record.GameAdministration, error) {
	l := m.Logger("GetManyGameAdministrationRecs")

	l.Debug("getting many game_administration records opts >%#v<", opts)

	r := m.GameAdministrationRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameAdministrationRec -
func (m *Domain) GetGameAdministrationRec(recID string, lock *sql.Lock) (*game_record.GameAdministration, error) {
	l := m.Logger("GetGameAdministrationRec")

	l.Debug("getting game_administration record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameAdministrationRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameAdministration, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameAdministrationRec -
func (m *Domain) CreateGameAdministrationRec(rec *game_record.GameAdministration) (*game_record.GameAdministration, error) {
	l := m.Logger("CreateGameAdministrationRec")

	l.Debug("creating game_administration record >%#v<", rec)

	if err := m.validateGameAdministrationRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_administration record >%v<", err)
		return rec, err
	}

	r := m.GameAdministrationRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateGameAdministrationRec -
func (m *Domain) UpdateGameAdministrationRec(rec *game_record.GameAdministration) (*game_record.GameAdministration, error) {
	l := m.Logger("UpdateGameAdministrationRec")

	_, err := m.GetGameAdministrationRec(rec.ID, sql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game_administration record >%#v<", rec)

	if err := m.validateGameAdministrationRecForUpdate(rec); err != nil {
		l.Warn("failed to validate game_administration record >%v<", err)
		return rec, err
	}

	r := m.GameAdministrationRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteGameAdministrationRec -
func (m *Domain) DeleteGameAdministrationRec(recID string) error {
	l := m.Logger("DeleteGameAdministrationRec")

	l.Debug("deleting game_administration record ID >%s<", recID)

	rec, err := m.GetGameAdministrationRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameAdministrationRepository()

	if err := m.validateGameAdministrationRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveGameAdministrationRec -
func (m *Domain) RemoveGameAdministrationRec(recID string) error {
	l := m.Logger("RemoveGameAdministrationRec")

	l.Debug("removing game_administration record ID >%s<", recID)

	rec, err := m.GetGameAdministrationRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameAdministrationRepository()

	if err := m.validateGameAdministrationRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
