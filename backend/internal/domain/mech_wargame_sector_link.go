package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameSectorLinkRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameSectorLink, error) {
	l := m.Logger("GetManyMechWargameSectorLinkRecs")

	l.Debug("getting many mech_wargame_sector_link records opts >%#v<", opts)

	r := m.MechWargameSectorLinkRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameSectorLinkRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameSectorLink, error) {
	l := m.Logger("GetMechWargameSectorLinkRec")

	l.Debug("getting mech_wargame_sector_link record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameSectorLinkRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameSectorLink, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameSectorLinkRec(rec *mech_wargame_record.MechWargameSectorLink) (*mech_wargame_record.MechWargameSectorLink, error) {
	l := m.Logger("CreateMechWargameSectorLinkRec")

	l.Debug("creating mech_wargame_sector_link record >%#v<", rec)

	if err := m.validateMechWargameSectorLinkRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_sector_link record >%v<", err)
		return rec, err
	}

	r := m.MechWargameSectorLinkRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechWargameSectorLinkRec(rec *mech_wargame_record.MechWargameSectorLink) (*mech_wargame_record.MechWargameSectorLink, error) {
	l := m.Logger("UpdateMechWargameSectorLinkRec")

	currRec, err := m.GetMechWargameSectorLinkRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mech_wargame_sector_link record >%#v<", rec)

	if err := m.validateMechWargameSectorLinkRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mech_wargame_sector_link record >%v<", err)
		return rec, err
	}

	r := m.MechWargameSectorLinkRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechWargameSectorLinkRec(recID string) error {
	l := m.Logger("DeleteMechWargameSectorLinkRec")

	l.Debug("deleting mech_wargame_sector_link record ID >%s<", recID)

	_, err := m.GetMechWargameSectorLinkRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechWargameSectorLinkRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechWargameSectorLinkRec(recID string) error {
	l := m.Logger("RemoveMechWargameSectorLinkRec")

	l.Debug("removing mech_wargame_sector_link record ID >%s<", recID)

	r := m.MechWargameSectorLinkRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
