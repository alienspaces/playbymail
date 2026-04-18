package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameSectorLinkRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameSectorLink, error) {
	l := m.Logger("GetManyMechaGameSectorLinkRecs")

	l.Debug("getting many mecha_game_sector_link records opts >%#v<", opts)

	r := m.MechaGameSectorLinkRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameSectorLinkRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameSectorLink, error) {
	l := m.Logger("GetMechaGameSectorLinkRec")

	l.Debug("getting mecha_game_sector_link record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameSectorLinkRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameSectorLink, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameSectorLinkRec(rec *mecha_game_record.MechaGameSectorLink) (*mecha_game_record.MechaGameSectorLink, error) {
	l := m.Logger("CreateMechaGameSectorLinkRec")

	l.Debug("creating mecha_game_sector_link record >%#v<", rec)

	if err := m.validateMechaGameSectorLinkRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_sector_link record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSectorLinkRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameSectorLinkRec(rec *mecha_game_record.MechaGameSectorLink) (*mecha_game_record.MechaGameSectorLink, error) {
	l := m.Logger("UpdateMechaGameSectorLinkRec")

	currRec, err := m.GetMechaGameSectorLinkRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_sector_link record >%#v<", rec)

	if err := m.validateMechaGameSectorLinkRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_sector_link record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSectorLinkRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameSectorLinkRec(recID string) error {
	l := m.Logger("DeleteMechaGameSectorLinkRec")

	l.Debug("deleting mecha_game_sector_link record ID >%s<", recID)

	_, err := m.GetMechaGameSectorLinkRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameSectorLinkRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameSectorLinkRec(recID string) error {
	l := m.Logger("RemoveMechaGameSectorLinkRec")

	l.Debug("removing mecha_game_sector_link record ID >%s<", recID)

	r := m.MechaGameSectorLinkRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
