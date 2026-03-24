package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaLanceRecs(opts *coresql.Options) ([]*mecha_record.MechaLance, error) {
	l := m.Logger("GetManyMechaLanceRecs")

	l.Debug("getting many mecha_lance records opts >%#v<", opts)

	r := m.MechaLanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaLanceRec(recID string, lock *coresql.Lock) (*mecha_record.MechaLance, error) {
	l := m.Logger("GetMechaLanceRec")

	l.Debug("getting mecha_lance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaLanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaLance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaLanceRec(rec *mecha_record.MechaLance) (*mecha_record.MechaLance, error) {
	l := m.Logger("CreateMechaLanceRec")

	l.Debug("creating mecha_lance record >%#v<", rec)

	if err := m.validateMechaLanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_lance record >%v<", err)
		return rec, err
	}

	r := m.MechaLanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaLanceRec(rec *mecha_record.MechaLance) (*mecha_record.MechaLance, error) {
	l := m.Logger("UpdateMechaLanceRec")

	currRec, err := m.GetMechaLanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_lance record >%#v<", rec)

	if err := m.validateMechaLanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_lance record >%v<", err)
		return rec, err
	}

	r := m.MechaLanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaLanceRec(recID string) error {
	l := m.Logger("DeleteMechaLanceRec")

	l.Debug("deleting mecha_lance record ID >%s<", recID)

	_, err := m.GetMechaLanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaLanceRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaLanceRec(recID string) error {
	l := m.Logger("RemoveMechaLanceRec")

	l.Debug("removing mecha_lance record ID >%s<", recID)

	r := m.MechaLanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
