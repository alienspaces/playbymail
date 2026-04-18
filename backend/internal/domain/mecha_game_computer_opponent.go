package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameComputerOpponentRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameComputerOpponent, error) {
	l := m.Logger("GetManyMechaGameComputerOpponentRecs")

	l.Debug("getting many mecha_game_computer_opponent records opts >%#v<", opts)

	r := m.MechaGameComputerOpponentRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameComputerOpponentRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameComputerOpponent, error) {
	l := m.Logger("GetMechaGameComputerOpponentRec")

	l.Debug("getting mecha_game_computer_opponent record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameComputerOpponentRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameComputerOpponent, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameComputerOpponentRec(rec *mecha_game_record.MechaGameComputerOpponent) (*mecha_game_record.MechaGameComputerOpponent, error) {
	l := m.Logger("CreateMechaGameComputerOpponentRec")

	l.Debug("creating mecha_game_computer_opponent record >%#v<", rec)

	if err := m.validateMechaGameComputerOpponentRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_computer_opponent record >%v<", err)
		return rec, err
	}

	r := m.MechaGameComputerOpponentRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameComputerOpponentRec(rec *mecha_game_record.MechaGameComputerOpponent) (*mecha_game_record.MechaGameComputerOpponent, error) {
	l := m.Logger("UpdateMechaGameComputerOpponentRec")

	currRec, err := m.GetMechaGameComputerOpponentRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_computer_opponent record >%#v<", rec)

	if err := m.validateMechaGameComputerOpponentRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_computer_opponent record >%v<", err)
		return rec, err
	}

	r := m.MechaGameComputerOpponentRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameComputerOpponentRec(recID string) error {
	l := m.Logger("DeleteMechaGameComputerOpponentRec")

	l.Debug("deleting mecha_game_computer_opponent record ID >%s<", recID)

	_, err := m.GetMechaGameComputerOpponentRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameComputerOpponentRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameComputerOpponentRec(recID string) error {
	l := m.Logger("RemoveMechaGameComputerOpponentRec")

	l.Debug("removing mecha_game_computer_opponent record ID >%s<", recID)

	r := m.MechaGameComputerOpponentRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
