package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameSectorRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameSector, error) {
	l := m.Logger("GetManyMechaGameSectorRecs")

	l.Debug("getting many mecha_game_sector records opts >%#v<", opts)

	r := m.MechaGameSectorRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameSectorRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameSector, error) {
	l := m.Logger("GetMechaGameSectorRec")

	l.Debug("getting mecha_game_sector record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameSectorRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameSector, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameSectorRec(rec *mecha_game_record.MechaGameSector) (*mecha_game_record.MechaGameSector, error) {
	l := m.Logger("CreateMechaGameSectorRec")

	l.Debug("creating mecha_game_sector record >%#v<", rec)

	if err := m.validateMechaGameSectorRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_sector record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSectorRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameSectorRec(rec *mecha_game_record.MechaGameSector) (*mecha_game_record.MechaGameSector, error) {
	l := m.Logger("UpdateMechaGameSectorRec")

	currRec, err := m.GetMechaGameSectorRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_sector record >%#v<", rec)

	if err := m.validateMechaGameSectorRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_sector record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSectorRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameSectorRec(recID string) error {
	l := m.Logger("DeleteMechaGameSectorRec")

	l.Debug("deleting mecha_game_sector record ID >%s<", recID)

	_, err := m.GetMechaGameSectorRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameSectorRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameSectorRec(recID string) error {
	l := m.Logger("RemoveMechaGameSectorRec")

	l.Debug("removing mecha_game_sector record ID >%s<", recID)

	r := m.MechaGameSectorRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
