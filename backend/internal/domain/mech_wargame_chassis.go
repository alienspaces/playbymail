package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameChassisRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameChassis, error) {
	l := m.Logger("GetManyMechWargameChassisRecs")

	l.Debug("getting many mech_wargame_chassis records opts >%#v<", opts)

	r := m.MechWargameChassisRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameChassisRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameChassis, error) {
	l := m.Logger("GetMechWargameChassisRec")

	l.Debug("getting mech_wargame_chassis record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameChassisRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameChassis, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameChassisRec(rec *mech_wargame_record.MechWargameChassis) (*mech_wargame_record.MechWargameChassis, error) {
	l := m.Logger("CreateMechWargameChassisRec")

	l.Debug("creating mech_wargame_chassis record >%#v<", rec)

	if err := m.validateMechWargameChassisRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_chassis record >%v<", err)
		return rec, err
	}

	r := m.MechWargameChassisRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechWargameChassisRec(rec *mech_wargame_record.MechWargameChassis) (*mech_wargame_record.MechWargameChassis, error) {
	l := m.Logger("UpdateMechWargameChassisRec")

	currRec, err := m.GetMechWargameChassisRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mech_wargame_chassis record >%#v<", rec)

	if err := m.validateMechWargameChassisRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mech_wargame_chassis record >%v<", err)
		return rec, err
	}

	r := m.MechWargameChassisRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechWargameChassisRec(recID string) error {
	l := m.Logger("DeleteMechWargameChassisRec")

	l.Debug("deleting mech_wargame_chassis record ID >%s<", recID)

	_, err := m.GetMechWargameChassisRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechWargameChassisRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechWargameChassisRec(recID string) error {
	l := m.Logger("RemoveMechWargameChassisRec")

	l.Debug("removing mech_wargame_chassis record ID >%s<", recID)

	r := m.MechWargameChassisRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
