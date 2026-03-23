package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaComputerOpponentRecs(opts *coresql.Options) ([]*mecha_record.MechaComputerOpponent, error) {
	l := m.Logger("GetManyMechaComputerOpponentRecs")

	l.Debug("getting many mecha_computer_opponent records opts >%#v<", opts)

	r := m.MechaComputerOpponentRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaComputerOpponentRec(recID string, lock *coresql.Lock) (*mecha_record.MechaComputerOpponent, error) {
	l := m.Logger("GetMechaComputerOpponentRec")

	l.Debug("getting mecha_computer_opponent record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaComputerOpponentRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaComputerOpponent, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaComputerOpponentRec(rec *mecha_record.MechaComputerOpponent) (*mecha_record.MechaComputerOpponent, error) {
	l := m.Logger("CreateMechaComputerOpponentRec")

	l.Debug("creating mecha_computer_opponent record >%#v<", rec)

	if err := m.validateMechaComputerOpponentRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_computer_opponent record >%v<", err)
		return rec, err
	}

	r := m.MechaComputerOpponentRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaComputerOpponentRec(rec *mecha_record.MechaComputerOpponent) (*mecha_record.MechaComputerOpponent, error) {
	l := m.Logger("UpdateMechaComputerOpponentRec")

	currRec, err := m.GetMechaComputerOpponentRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_computer_opponent record >%#v<", rec)

	if err := m.validateMechaComputerOpponentRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_computer_opponent record >%v<", err)
		return rec, err
	}

	r := m.MechaComputerOpponentRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaComputerOpponentRec(recID string) error {
	l := m.Logger("DeleteMechaComputerOpponentRec")

	l.Debug("deleting mecha_computer_opponent record ID >%s<", recID)

	_, err := m.GetMechaComputerOpponentRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaComputerOpponentRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaComputerOpponentRec(recID string) error {
	l := m.Logger("RemoveMechaComputerOpponentRec")

	l.Debug("removing mecha_computer_opponent record ID >%s<", recID)

	r := m.MechaComputerOpponentRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
