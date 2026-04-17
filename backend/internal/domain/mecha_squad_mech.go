package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaSquadMechRecs(opts *coresql.Options) ([]*mecha_record.MechaSquadMech, error) {
	l := m.Logger("GetManyMechaSquadMechRecs")

	l.Debug("getting many mecha_squad_mech records opts >%#v<", opts)

	r := m.MechaSquadMechRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaSquadMechRec(recID string, lock *coresql.Lock) (*mecha_record.MechaSquadMech, error) {
	l := m.Logger("GetMechaSquadMechRec")

	l.Debug("getting mecha_squad_mech record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaSquadMechRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaSquadMech, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaSquadMechRec(rec *mecha_record.MechaSquadMech) (*mecha_record.MechaSquadMech, error) {
	l := m.Logger("CreateMechaSquadMechRec")

	l.Debug("creating mecha_squad_mech record >%#v<", rec)

	if err := m.validateMechaSquadMechRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_squad_mech record >%v<", err)
		return rec, err
	}

	r := m.MechaSquadMechRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaSquadMechRec(rec *mecha_record.MechaSquadMech) (*mecha_record.MechaSquadMech, error) {
	l := m.Logger("UpdateMechaSquadMechRec")

	currRec, err := m.GetMechaSquadMechRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_squad_mech record >%#v<", rec)

	if err := m.validateMechaSquadMechRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_squad_mech record >%v<", err)
		return rec, err
	}

	r := m.MechaSquadMechRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaSquadMechRec(recID string) error {
	l := m.Logger("DeleteMechaSquadMechRec")

	l.Debug("deleting mecha_squad_mech record ID >%s<", recID)

	_, err := m.GetMechaSquadMechRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaSquadMechRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaSquadMechRec(recID string) error {
	l := m.Logger("RemoveMechaSquadMechRec")

	l.Debug("removing mecha_squad_mech record ID >%s<", recID)

	r := m.MechaSquadMechRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
