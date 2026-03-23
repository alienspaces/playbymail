package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameLanceRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameLance, error) {
	l := m.Logger("GetManyMechWargameLanceRecs")

	l.Debug("getting many mech_wargame_lance records opts >%#v<", opts)

	r := m.MechWargameLanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameLanceRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameLance, error) {
	l := m.Logger("GetMechWargameLanceRec")

	l.Debug("getting mech_wargame_lance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameLanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameLance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameLanceRec(rec *mech_wargame_record.MechWargameLance) (*mech_wargame_record.MechWargameLance, error) {
	l := m.Logger("CreateMechWargameLanceRec")

	l.Debug("creating mech_wargame_lance record >%#v<", rec)

	if err := m.validateMechWargameLanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_lance record >%v<", err)
		return rec, err
	}

	r := m.MechWargameLanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechWargameLanceRec(rec *mech_wargame_record.MechWargameLance) (*mech_wargame_record.MechWargameLance, error) {
	l := m.Logger("UpdateMechWargameLanceRec")

	currRec, err := m.GetMechWargameLanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mech_wargame_lance record >%#v<", rec)

	if err := m.validateMechWargameLanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mech_wargame_lance record >%v<", err)
		return rec, err
	}

	r := m.MechWargameLanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechWargameLanceRec(recID string) error {
	l := m.Logger("DeleteMechWargameLanceRec")

	l.Debug("deleting mech_wargame_lance record ID >%s<", recID)

	_, err := m.GetMechWargameLanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechWargameLanceRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechWargameLanceRec(recID string) error {
	l := m.Logger("RemoveMechWargameLanceRec")

	l.Debug("removing mech_wargame_lance record ID >%s<", recID)

	r := m.MechWargameLanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
