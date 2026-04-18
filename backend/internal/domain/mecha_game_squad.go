package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameSquadRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameSquad, error) {
	l := m.Logger("GetManyMechaGameSquadRecs")

	l.Debug("getting many mecha_game_squad records opts >%#v<", opts)

	r := m.MechaGameSquadRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameSquadRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameSquad, error) {
	l := m.Logger("GetMechaGameSquadRec")

	l.Debug("getting mecha_game_squad record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameSquadRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameSquad, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameSquadRec(rec *mecha_game_record.MechaGameSquad) (*mecha_game_record.MechaGameSquad, error) {
	l := m.Logger("CreateMechaGameSquadRec")

	l.Debug("creating mecha_game_squad record >%#v<", rec)

	if err := m.validateMechaGameSquadRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_squad record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSquadRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameSquadRec(rec *mecha_game_record.MechaGameSquad) (*mecha_game_record.MechaGameSquad, error) {
	l := m.Logger("UpdateMechaGameSquadRec")

	currRec, err := m.GetMechaGameSquadRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_squad record >%#v<", rec)

	if err := m.validateMechaGameSquadRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_squad record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSquadRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameSquadRec(recID string) error {
	l := m.Logger("DeleteMechaGameSquadRec")

	l.Debug("deleting mecha_game_squad record ID >%s<", recID)

	_, err := m.GetMechaGameSquadRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameSquadRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameSquadRec(recID string) error {
	l := m.Logger("RemoveMechaGameSquadRec")

	l.Debug("removing mecha_game_squad record ID >%s<", recID)

	r := m.MechaGameSquadRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
