package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameWeaponRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameWeapon, error) {
	l := m.Logger("GetManyMechWargameWeaponRecs")

	l.Debug("getting many mech_wargame_weapon records opts >%#v<", opts)

	r := m.MechWargameWeaponRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameWeaponRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameWeapon, error) {
	l := m.Logger("GetMechWargameWeaponRec")

	l.Debug("getting mech_wargame_weapon record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameWeaponRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameWeapon, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameWeaponRec(rec *mech_wargame_record.MechWargameWeapon) (*mech_wargame_record.MechWargameWeapon, error) {
	l := m.Logger("CreateMechWargameWeaponRec")

	l.Debug("creating mech_wargame_weapon record >%#v<", rec)

	if err := m.validateMechWargameWeaponRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_weapon record >%v<", err)
		return rec, err
	}

	r := m.MechWargameWeaponRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechWargameWeaponRec(rec *mech_wargame_record.MechWargameWeapon) (*mech_wargame_record.MechWargameWeapon, error) {
	l := m.Logger("UpdateMechWargameWeaponRec")

	currRec, err := m.GetMechWargameWeaponRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mech_wargame_weapon record >%#v<", rec)

	if err := m.validateMechWargameWeaponRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mech_wargame_weapon record >%v<", err)
		return rec, err
	}

	r := m.MechWargameWeaponRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechWargameWeaponRec(recID string) error {
	l := m.Logger("DeleteMechWargameWeaponRec")

	l.Debug("deleting mech_wargame_weapon record ID >%s<", recID)

	_, err := m.GetMechWargameWeaponRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechWargameWeaponRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechWargameWeaponRec(recID string) error {
	l := m.Logger("RemoveMechWargameWeaponRec")

	l.Debug("removing mech_wargame_weapon record ID >%s<", recID)

	r := m.MechWargameWeaponRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
