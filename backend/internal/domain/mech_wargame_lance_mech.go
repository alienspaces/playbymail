package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameLanceMechRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameLanceMech, error) {
	l := m.Logger("GetManyMechWargameLanceMechRecs")

	l.Debug("getting many mech_wargame_lance_mech records opts >%#v<", opts)

	r := m.MechWargameLanceMechRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameLanceMechRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameLanceMech, error) {
	l := m.Logger("GetMechWargameLanceMechRec")

	l.Debug("getting mech_wargame_lance_mech record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameLanceMechRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameLanceMech, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameLanceMechRec(rec *mech_wargame_record.MechWargameLanceMech) (*mech_wargame_record.MechWargameLanceMech, error) {
	l := m.Logger("CreateMechWargameLanceMechRec")

	l.Debug("creating mech_wargame_lance_mech record >%#v<", rec)

	if err := m.validateMechWargameLanceMechRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_lance_mech record >%v<", err)
		return rec, err
	}

	r := m.MechWargameLanceMechRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechWargameLanceMechRec(rec *mech_wargame_record.MechWargameLanceMech) (*mech_wargame_record.MechWargameLanceMech, error) {
	l := m.Logger("UpdateMechWargameLanceMechRec")

	currRec, err := m.GetMechWargameLanceMechRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mech_wargame_lance_mech record >%#v<", rec)

	if err := m.validateMechWargameLanceMechRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mech_wargame_lance_mech record >%v<", err)
		return rec, err
	}

	r := m.MechWargameLanceMechRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechWargameLanceMechRec(recID string) error {
	l := m.Logger("DeleteMechWargameLanceMechRec")

	l.Debug("deleting mech_wargame_lance_mech record ID >%s<", recID)

	_, err := m.GetMechWargameLanceMechRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechWargameLanceMechRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechWargameLanceMechRec(recID string) error {
	l := m.Logger("RemoveMechWargameLanceMechRec")

	l.Debug("removing mech_wargame_lance_mech record ID >%s<", recID)

	r := m.MechWargameLanceMechRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
