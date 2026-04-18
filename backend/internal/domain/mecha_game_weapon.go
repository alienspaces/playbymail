package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameWeaponRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameWeapon, error) {
	l := m.Logger("GetManyMechaGameWeaponRecs")

	l.Debug("getting many mecha_game_weapon records opts >%#v<", opts)

	r := m.MechaGameWeaponRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameWeaponRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameWeapon, error) {
	l := m.Logger("GetMechaGameWeaponRec")

	l.Debug("getting mecha_game_weapon record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameWeaponRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameWeapon, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameWeaponRec(rec *mecha_game_record.MechaGameWeapon) (*mecha_game_record.MechaGameWeapon, error) {
	l := m.Logger("CreateMechaGameWeaponRec")

	l.Debug("creating mecha_game_weapon record >%#v<", rec)

	if err := m.validateMechaGameWeaponRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_weapon record >%v<", err)
		return rec, err
	}

	r := m.MechaGameWeaponRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameWeaponRec(rec *mecha_game_record.MechaGameWeapon) (*mecha_game_record.MechaGameWeapon, error) {
	l := m.Logger("UpdateMechaGameWeaponRec")

	currRec, err := m.GetMechaGameWeaponRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_weapon record >%#v<", rec)

	if err := m.validateMechaGameWeaponRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_weapon record >%v<", err)
		return rec, err
	}

	r := m.MechaGameWeaponRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameWeaponRec(recID string) error {
	l := m.Logger("DeleteMechaGameWeaponRec")

	l.Debug("deleting mecha_game_weapon record ID >%s<", recID)

	_, err := m.GetMechaGameWeaponRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameWeaponRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameWeaponRec(recID string) error {
	l := m.Logger("RemoveMechaGameWeaponRec")

	l.Debug("removing mecha_game_weapon record ID >%s<", recID)

	r := m.MechaGameWeaponRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
