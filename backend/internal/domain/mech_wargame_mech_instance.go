package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameMechInstanceRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameMechInstance, error) {
	l := m.Logger("GetManyMechWargameMechInstanceRecs")

	l.Debug("getting many mech_wargame_mech_instance records opts >%#v<", opts)

	r := m.MechWargameMechInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameMechInstanceRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameMechInstance, error) {
	l := m.Logger("GetMechWargameMechInstanceRec")

	l.Debug("getting mech_wargame_mech_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameMechInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameMechInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameMechInstanceRec(rec *mech_wargame_record.MechWargameMechInstance) (*mech_wargame_record.MechWargameMechInstance, error) {
	l := m.Logger("CreateMechWargameMechInstanceRec")

	l.Debug("creating mech_wargame_mech_instance record >%#v<", rec)

	if err := m.validateMechWargameMechInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_mech_instance record >%v<", err)
		return rec, err
	}

	r := m.MechWargameMechInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechWargameMechInstanceRec(rec *mech_wargame_record.MechWargameMechInstance) (*mech_wargame_record.MechWargameMechInstance, error) {
	l := m.Logger("UpdateMechWargameMechInstanceRec")

	currRec, err := m.GetMechWargameMechInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mech_wargame_mech_instance record >%#v<", rec)

	if err := m.validateMechWargameMechInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mech_wargame_mech_instance record >%v<", err)
		return rec, err
	}

	r := m.MechWargameMechInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) RemoveMechWargameMechInstanceRec(recID string) error {
	l := m.Logger("RemoveMechWargameMechInstanceRec")

	l.Debug("removing mech_wargame_mech_instance record ID >%s<", recID)

	r := m.MechWargameMechInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
