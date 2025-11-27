package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetManyGameInstanceParameterRecs -
func (m *Domain) GetManyGameInstanceParameterRecs(opts *coresql.Options) ([]*game_record.GameInstanceParameter, error) {
	l := m.Logger("GetManyGameInstanceParameterRecs")

	l.Debug("getting many game_instance_parameter records opts >%#v<", opts)

	r := m.GameInstanceParameterRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameInstanceParameterRec -
func (m *Domain) GetGameInstanceParameterRec(recID string, lock *coresql.Lock) (*game_record.GameInstanceParameter, error) {
	l := m.Logger("GetGameInstanceParameterRec")

	l.Debug("getting game_instance_parameter record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameInstanceParameterRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameInstanceParameter, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameInstanceParameterRec -
func (m *Domain) CreateGameInstanceParameterRec(rec *game_record.GameInstanceParameter) (*game_record.GameInstanceParameter, error) {
	l := m.Logger("CreateGameInstanceParameterRec")

	l.Debug("creating game_instance_parameter record >%#v<", rec)

	if err := m.validateGameInstanceParameterRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_instance_parameter record >%v<", err)
		return rec, err
	}

	// Use comprehensive validation for business logic
	if err := m.ValidateGameInstanceParameter(rec); err != nil {
		l.Warn("failed comprehensive validation for game_instance_parameter record >%v<", err)
		return rec, err
	}

	r := m.GameInstanceParameterRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateGameInstanceParameterRec -
func (m *Domain) UpdateGameInstanceParameterRec(rec *game_record.GameInstanceParameter) (*game_record.GameInstanceParameter, error) {
	l := m.Logger("UpdateGameInstanceParameterRec")

	_, err := m.GetGameInstanceParameterRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game_instance_parameter record >%#v<", rec)

	if err := m.validateGameInstanceParameterRecForUpdate(rec); err != nil {
		l.Warn("failed to validate game_instance_parameter record >%v<", err)
		return rec, err
	}

	// Use comprehensive validation for business logic
	if err := m.ValidateGameInstanceParameter(rec); err != nil {
		l.Warn("failed comprehensive validation for game_instance_parameter record >%v<", err)
		return rec, err
	}

	r := m.GameInstanceParameterRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteGameInstanceParameterRec -
func (m *Domain) DeleteGameInstanceParameterRec(recID string) error {
	l := m.Logger("DeleteGameInstanceParameterRec")

	l.Debug("deleting game_instance_parameter record ID >%s<", recID)

	_, err := m.GetGameInstanceParameterRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameInstanceParameterRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveGameInstanceParameterRec -
func (m *Domain) RemoveGameInstanceParameterRec(recID string) error {
	l := m.Logger("RemoveGameInstanceParameterRec")

	l.Debug("removing game_instance_parameter record ID >%s<", recID)

	_, err := m.GetGameInstanceParameterRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameInstanceParameterRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// GetGameInstanceParameterRecsByGameInstanceID gets all parameters for a specific game instance
func (m *Domain) GetGameInstanceParameterRecsByGameInstanceID(gameInstanceID string) ([]*game_record.GameInstanceParameter, error) {
	l := m.Logger("GetGameInstanceParameterRecsByGameInstanceID")

	l.Debug("getting game_instance_parameter records for game_instance_id >%s<", gameInstanceID)

	r := m.GameInstanceParameterRepository()

	opts := &coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameInstanceParameterGameInstanceID,
				Val: gameInstanceID,
			},
		},
	}

	return r.GetMany(opts)
}
