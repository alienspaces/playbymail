package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaWeaponRecs(opts *coresql.Options) ([]*mecha_record.MechaWeapon, error) {
	l := m.Logger("GetManyMechaWeaponRecs")

	l.Debug("getting many mecha_weapon records opts >%#v<", opts)

	r := m.MechaWeaponRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaWeaponRec(recID string, lock *coresql.Lock) (*mecha_record.MechaWeapon, error) {
	l := m.Logger("GetMechaWeaponRec")

	l.Debug("getting mecha_weapon record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaWeaponRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaWeapon, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaWeaponRec(rec *mecha_record.MechaWeapon) (*mecha_record.MechaWeapon, error) {
	l := m.Logger("CreateMechaWeaponRec")

	l.Debug("creating mecha_weapon record >%#v<", rec)

	if err := m.validateMechaWeaponRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_weapon record >%v<", err)
		return rec, err
	}

	r := m.MechaWeaponRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaWeaponRec(rec *mecha_record.MechaWeapon) (*mecha_record.MechaWeapon, error) {
	l := m.Logger("UpdateMechaWeaponRec")

	currRec, err := m.GetMechaWeaponRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_weapon record >%#v<", rec)

	if err := m.validateMechaWeaponRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_weapon record >%v<", err)
		return rec, err
	}

	r := m.MechaWeaponRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaWeaponRec(recID string) error {
	l := m.Logger("DeleteMechaWeaponRec")

	l.Debug("deleting mecha_weapon record ID >%s<", recID)

	_, err := m.GetMechaWeaponRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaWeaponRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaWeaponRec(recID string) error {
	l := m.Logger("RemoveMechaWeaponRec")

	l.Debug("removing mecha_weapon record ID >%s<", recID)

	r := m.MechaWeaponRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
