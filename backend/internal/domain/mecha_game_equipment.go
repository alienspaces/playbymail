package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameEquipmentRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameEquipment, error) {
	l := m.Logger("GetManyMechaGameEquipmentRecs")

	l.Debug("getting many mecha_game_equipment records opts >%#v<", opts)

	r := m.MechaGameEquipmentRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameEquipmentRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameEquipment, error) {
	l := m.Logger("GetMechaGameEquipmentRec")

	l.Debug("getting mecha_game_equipment record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameEquipmentRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameEquipment, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameEquipmentRec(rec *mecha_game_record.MechaGameEquipment) (*mecha_game_record.MechaGameEquipment, error) {
	l := m.Logger("CreateMechaGameEquipmentRec")

	l.Debug("creating mecha_game_equipment record >%#v<", rec)

	if err := m.validateMechaGameEquipmentRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_equipment record >%v<", err)
		return rec, err
	}

	r := m.MechaGameEquipmentRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameEquipmentRec(rec *mecha_game_record.MechaGameEquipment) (*mecha_game_record.MechaGameEquipment, error) {
	l := m.Logger("UpdateMechaGameEquipmentRec")

	currRec, err := m.GetMechaGameEquipmentRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_equipment record >%#v<", rec)

	if err := m.validateMechaGameEquipmentRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_equipment record >%v<", err)
		return rec, err
	}

	r := m.MechaGameEquipmentRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameEquipmentRec(recID string) error {
	l := m.Logger("DeleteMechaGameEquipmentRec")

	l.Debug("deleting mecha_game_equipment record ID >%s<", recID)

	_, err := m.GetMechaGameEquipmentRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameEquipmentRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameEquipmentRec(recID string) error {
	l := m.Logger("RemoveMechaGameEquipmentRec")

	l.Debug("removing mecha_game_equipment record ID >%s<", recID)

	r := m.MechaGameEquipmentRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
