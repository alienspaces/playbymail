package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameChassisRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameChassis, error) {
	l := m.Logger("GetManyMechaGameChassisRecs")

	l.Debug("getting many mecha_game_chassis records opts >%#v<", opts)

	r := m.MechaGameChassisRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameChassisRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameChassis, error) {
	l := m.Logger("GetMechaGameChassisRec")

	l.Debug("getting mecha_game_chassis record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameChassisRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameChassis, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameChassisRec(rec *mecha_game_record.MechaGameChassis) (*mecha_game_record.MechaGameChassis, error) {
	l := m.Logger("CreateMechaGameChassisRec")

	l.Debug("creating mecha_game_chassis record >%#v<", rec)

	if err := m.validateMechaGameChassisRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_chassis record >%v<", err)
		return rec, err
	}

	r := m.MechaGameChassisRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameChassisRec(rec *mecha_game_record.MechaGameChassis) (*mecha_game_record.MechaGameChassis, error) {
	l := m.Logger("UpdateMechaGameChassisRec")

	currRec, err := m.GetMechaGameChassisRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_chassis record >%#v<", rec)

	if err := m.validateMechaGameChassisRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_chassis record >%v<", err)
		return rec, err
	}

	r := m.MechaGameChassisRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameChassisRec(recID string) error {
	l := m.Logger("DeleteMechaGameChassisRec")

	l.Debug("deleting mecha_game_chassis record ID >%s<", recID)

	_, err := m.GetMechaGameChassisRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameChassisRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameChassisRec(recID string) error {
	l := m.Logger("RemoveMechaGameChassisRec")

	l.Debug("removing mecha_game_chassis record ID >%s<", recID)

	r := m.MechaGameChassisRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
