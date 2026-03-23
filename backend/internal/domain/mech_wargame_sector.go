package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameSectorRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameSector, error) {
	l := m.Logger("GetManyMechWargameSectorRecs")

	l.Debug("getting many mech_wargame_sector records opts >%#v<", opts)

	r := m.MechWargameSectorRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameSectorRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameSector, error) {
	l := m.Logger("GetMechWargameSectorRec")

	l.Debug("getting mech_wargame_sector record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameSectorRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameSector, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameSectorRec(rec *mech_wargame_record.MechWargameSector) (*mech_wargame_record.MechWargameSector, error) {
	l := m.Logger("CreateMechWargameSectorRec")

	l.Debug("creating mech_wargame_sector record >%#v<", rec)

	if err := m.validateMechWargameSectorRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_sector record >%v<", err)
		return rec, err
	}

	r := m.MechWargameSectorRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechWargameSectorRec(rec *mech_wargame_record.MechWargameSector) (*mech_wargame_record.MechWargameSector, error) {
	l := m.Logger("UpdateMechWargameSectorRec")

	currRec, err := m.GetMechWargameSectorRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mech_wargame_sector record >%#v<", rec)

	if err := m.validateMechWargameSectorRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mech_wargame_sector record >%v<", err)
		return rec, err
	}

	r := m.MechWargameSectorRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechWargameSectorRec(recID string) error {
	l := m.Logger("DeleteMechWargameSectorRec")

	l.Debug("deleting mech_wargame_sector record ID >%s<", recID)

	_, err := m.GetMechWargameSectorRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechWargameSectorRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechWargameSectorRec(recID string) error {
	l := m.Logger("RemoveMechWargameSectorRec")

	l.Debug("removing mech_wargame_sector record ID >%s<", recID)

	r := m.MechWargameSectorRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
