package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaMechInstanceRecs(opts *coresql.Options) ([]*mecha_record.MechaMechInstance, error) {
	l := m.Logger("GetManyMechaMechInstanceRecs")

	l.Debug("getting many mecha_mech_instance records opts >%#v<", opts)

	r := m.MechaMechInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaMechInstanceRec(recID string, lock *coresql.Lock) (*mecha_record.MechaMechInstance, error) {
	l := m.Logger("GetMechaMechInstanceRec")

	l.Debug("getting mecha_mech_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaMechInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaMechInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaMechInstanceRec(rec *mecha_record.MechaMechInstance) (*mecha_record.MechaMechInstance, error) {
	l := m.Logger("CreateMechaMechInstanceRec")

	l.Debug("creating mecha_mech_instance record >%#v<", rec)

	if err := m.validateMechaMechInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_mech_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaMechInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaMechInstanceRec(rec *mecha_record.MechaMechInstance) (*mecha_record.MechaMechInstance, error) {
	l := m.Logger("UpdateMechaMechInstanceRec")

	currRec, err := m.GetMechaMechInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_mech_instance record >%#v<", rec)

	if err := m.validateMechaMechInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_mech_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaMechInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) RemoveMechaMechInstanceRec(recID string) error {
	l := m.Logger("RemoveMechaMechInstanceRec")

	l.Debug("removing mecha_mech_instance record ID >%s<", recID)

	r := m.MechaMechInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
