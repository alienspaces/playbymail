package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaSectorLinkRecs(opts *coresql.Options) ([]*mecha_record.MechaSectorLink, error) {
	l := m.Logger("GetManyMechaSectorLinkRecs")

	l.Debug("getting many mecha_sector_link records opts >%#v<", opts)

	r := m.MechaSectorLinkRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaSectorLinkRec(recID string, lock *coresql.Lock) (*mecha_record.MechaSectorLink, error) {
	l := m.Logger("GetMechaSectorLinkRec")

	l.Debug("getting mecha_sector_link record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaSectorLinkRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaSectorLink, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaSectorLinkRec(rec *mecha_record.MechaSectorLink) (*mecha_record.MechaSectorLink, error) {
	l := m.Logger("CreateMechaSectorLinkRec")

	l.Debug("creating mecha_sector_link record >%#v<", rec)

	if err := m.validateMechaSectorLinkRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_sector_link record >%v<", err)
		return rec, err
	}

	r := m.MechaSectorLinkRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaSectorLinkRec(rec *mecha_record.MechaSectorLink) (*mecha_record.MechaSectorLink, error) {
	l := m.Logger("UpdateMechaSectorLinkRec")

	currRec, err := m.GetMechaSectorLinkRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_sector_link record >%#v<", rec)

	if err := m.validateMechaSectorLinkRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_sector_link record >%v<", err)
		return rec, err
	}

	r := m.MechaSectorLinkRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaSectorLinkRec(recID string) error {
	l := m.Logger("DeleteMechaSectorLinkRec")

	l.Debug("deleting mecha_sector_link record ID >%s<", recID)

	_, err := m.GetMechaSectorLinkRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaSectorLinkRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaSectorLinkRec(recID string) error {
	l := m.Logger("RemoveMechaSectorLinkRec")

	l.Debug("removing mecha_sector_link record ID >%s<", recID)

	r := m.MechaSectorLinkRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
