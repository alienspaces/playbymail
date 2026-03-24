package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaSectorRecs(opts *coresql.Options) ([]*mecha_record.MechaSector, error) {
	l := m.Logger("GetManyMechaSectorRecs")

	l.Debug("getting many mecha_sector records opts >%#v<", opts)

	r := m.MechaSectorRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaSectorRec(recID string, lock *coresql.Lock) (*mecha_record.MechaSector, error) {
	l := m.Logger("GetMechaSectorRec")

	l.Debug("getting mecha_sector record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaSectorRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaSector, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaSectorRec(rec *mecha_record.MechaSector) (*mecha_record.MechaSector, error) {
	l := m.Logger("CreateMechaSectorRec")

	l.Debug("creating mecha_sector record >%#v<", rec)

	if err := m.validateMechaSectorRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_sector record >%v<", err)
		return rec, err
	}

	r := m.MechaSectorRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaSectorRec(rec *mecha_record.MechaSector) (*mecha_record.MechaSector, error) {
	l := m.Logger("UpdateMechaSectorRec")

	currRec, err := m.GetMechaSectorRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_sector record >%#v<", rec)

	if err := m.validateMechaSectorRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_sector record >%v<", err)
		return rec, err
	}

	r := m.MechaSectorRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaSectorRec(recID string) error {
	l := m.Logger("DeleteMechaSectorRec")

	l.Debug("deleting mecha_sector record ID >%s<", recID)

	_, err := m.GetMechaSectorRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaSectorRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaSectorRec(recID string) error {
	l := m.Logger("RemoveMechaSectorRec")

	l.Debug("removing mecha_sector record ID >%s<", recID)

	r := m.MechaSectorRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
