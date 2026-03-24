package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaLanceMechRecs(opts *coresql.Options) ([]*mecha_record.MechaLanceMech, error) {
	l := m.Logger("GetManyMechaLanceMechRecs")

	l.Debug("getting many mecha_lance_mech records opts >%#v<", opts)

	r := m.MechaLanceMechRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaLanceMechRec(recID string, lock *coresql.Lock) (*mecha_record.MechaLanceMech, error) {
	l := m.Logger("GetMechaLanceMechRec")

	l.Debug("getting mecha_lance_mech record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaLanceMechRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaLanceMech, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaLanceMechRec(rec *mecha_record.MechaLanceMech) (*mecha_record.MechaLanceMech, error) {
	l := m.Logger("CreateMechaLanceMechRec")

	l.Debug("creating mecha_lance_mech record >%#v<", rec)

	if err := m.validateMechaLanceMechRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_lance_mech record >%v<", err)
		return rec, err
	}

	r := m.MechaLanceMechRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaLanceMechRec(rec *mecha_record.MechaLanceMech) (*mecha_record.MechaLanceMech, error) {
	l := m.Logger("UpdateMechaLanceMechRec")

	currRec, err := m.GetMechaLanceMechRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_lance_mech record >%#v<", rec)

	if err := m.validateMechaLanceMechRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_lance_mech record >%v<", err)
		return rec, err
	}

	r := m.MechaLanceMechRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaLanceMechRec(recID string) error {
	l := m.Logger("DeleteMechaLanceMechRec")

	l.Debug("deleting mecha_lance_mech record ID >%s<", recID)

	_, err := m.GetMechaLanceMechRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaLanceMechRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaLanceMechRec(recID string) error {
	l := m.Logger("RemoveMechaLanceMechRec")

	l.Debug("removing mecha_lance_mech record ID >%s<", recID)

	r := m.MechaLanceMechRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
