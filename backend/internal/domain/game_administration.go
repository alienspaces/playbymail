package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/sql"
	game_record "gitlab.com/alienspaces/playbymail/internal/record/game"
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
	r := m.GameAdministrationRepository()
	if err := m.validateGameAdministrationRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_administration record >%v<", err)
		return rec, err
	}
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameAdministrationRec -
func (m *Domain) UpdateGameAdministrationRec(next *game_record.GameAdministration) (*game_record.GameAdministration, error) {
	l := m.Logger("UpdateGameAdministrationRec")
	curr, err := m.GetGameAdministrationRec(next.ID, sql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_administration record >%#v<", next)
	if err := m.validateGameAdministrationRecForUpdate(next, curr); err != nil {
		l.Warn("failed to validate game_administration record >%v<", err)
		return next, err
	}
	r := m.GameAdministrationRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
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

// Validation stubs
func (m *Domain) validateGameAdministrationRecForCreate(rec *game_record.GameAdministration) error {
	l := m.Logger("validateGameAdministrationRecForCreate")
	l.Debug("validating game_administration record >%#v<", rec)
	// TODO: Add validation logic
	return nil
}
func (m *Domain) validateGameAdministrationRecForUpdate(next, curr *game_record.GameAdministration) error {
	l := m.Logger("validateGameAdministrationRecForUpdate")
	l.Debug("validating current game_administration record >%#v< against next >%#v<", curr, next)

	// TODO: Add validation logic
	return nil
}
func (m *Domain) validateGameAdministrationRecForDelete(rec *game_record.GameAdministration) error {
	l := m.Logger("validateGameAdministrationRecForDelete")
	l.Debug("validating game_administration record >%#v<", rec)
	// TODO: Add validation logic
	return nil
}
