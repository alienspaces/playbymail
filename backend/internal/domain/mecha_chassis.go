package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaChassisRecs(opts *coresql.Options) ([]*mecha_record.MechaChassis, error) {
	l := m.Logger("GetManyMechaChassisRecs")

	l.Debug("getting many mecha_chassis records opts >%#v<", opts)

	r := m.MechaChassisRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaChassisRec(recID string, lock *coresql.Lock) (*mecha_record.MechaChassis, error) {
	l := m.Logger("GetMechaChassisRec")

	l.Debug("getting mecha_chassis record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaChassisRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaChassis, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaChassisRec(rec *mecha_record.MechaChassis) (*mecha_record.MechaChassis, error) {
	l := m.Logger("CreateMechaChassisRec")

	l.Debug("creating mecha_chassis record >%#v<", rec)

	if err := m.validateMechaChassisRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_chassis record >%v<", err)
		return rec, err
	}

	r := m.MechaChassisRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaChassisRec(rec *mecha_record.MechaChassis) (*mecha_record.MechaChassis, error) {
	l := m.Logger("UpdateMechaChassisRec")

	currRec, err := m.GetMechaChassisRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_chassis record >%#v<", rec)

	if err := m.validateMechaChassisRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_chassis record >%v<", err)
		return rec, err
	}

	r := m.MechaChassisRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaChassisRec(recID string) error {
	l := m.Logger("DeleteMechaChassisRec")

	l.Debug("deleting mecha_chassis record ID >%s<", recID)

	_, err := m.GetMechaChassisRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaChassisRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaChassisRec(recID string) error {
	l := m.Logger("RemoveMechaChassisRec")

	l.Debug("removing mecha_chassis record ID >%s<", recID)

	r := m.MechaChassisRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
