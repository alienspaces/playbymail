package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaSquadRecs(opts *coresql.Options) ([]*mecha_record.MechaSquad, error) {
	l := m.Logger("GetManyMechaSquadRecs")

	l.Debug("getting many mecha_squad records opts >%#v<", opts)

	r := m.MechaSquadRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaSquadRec(recID string, lock *coresql.Lock) (*mecha_record.MechaSquad, error) {
	l := m.Logger("GetMechaSquadRec")

	l.Debug("getting mecha_squad record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaSquadRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaSquad, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaSquadRec(rec *mecha_record.MechaSquad) (*mecha_record.MechaSquad, error) {
	l := m.Logger("CreateMechaSquadRec")

	l.Debug("creating mecha_squad record >%#v<", rec)

	if err := m.validateMechaSquadRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_squad record >%v<", err)
		return rec, err
	}

	r := m.MechaSquadRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaSquadRec(rec *mecha_record.MechaSquad) (*mecha_record.MechaSquad, error) {
	l := m.Logger("UpdateMechaSquadRec")

	currRec, err := m.GetMechaSquadRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_squad record >%#v<", rec)

	if err := m.validateMechaSquadRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_squad record >%v<", err)
		return rec, err
	}

	r := m.MechaSquadRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaSquadRec(recID string) error {
	l := m.Logger("DeleteMechaSquadRec")

	l.Debug("deleting mecha_squad record ID >%s<", recID)

	_, err := m.GetMechaSquadRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaSquadRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaSquadRec(recID string) error {
	l := m.Logger("RemoveMechaSquadRec")

	l.Debug("removing mecha_squad record ID >%s<", recID)

	r := m.MechaSquadRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
