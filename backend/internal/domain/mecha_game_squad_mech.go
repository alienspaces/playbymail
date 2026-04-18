package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameSquadMechRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameSquadMech, error) {
	l := m.Logger("GetManyMechaGameSquadMechRecs")

	l.Debug("getting many mecha_game_squad_mech records opts >%#v<", opts)

	r := m.MechaGameSquadMechRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameSquadMechRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameSquadMech, error) {
	l := m.Logger("GetMechaGameSquadMechRec")

	l.Debug("getting mecha_game_squad_mech record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameSquadMechRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameSquadMech, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameSquadMechRec(rec *mecha_game_record.MechaGameSquadMech) (*mecha_game_record.MechaGameSquadMech, error) {
	l := m.Logger("CreateMechaGameSquadMechRec")

	l.Debug("creating mecha_game_squad_mech record >%#v<", rec)

	if err := m.validateMechaGameSquadMechRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_squad_mech record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSquadMechRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameSquadMechRec(rec *mecha_game_record.MechaGameSquadMech) (*mecha_game_record.MechaGameSquadMech, error) {
	l := m.Logger("UpdateMechaGameSquadMechRec")

	currRec, err := m.GetMechaGameSquadMechRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_squad_mech record >%#v<", rec)

	if err := m.validateMechaGameSquadMechRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_squad_mech record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSquadMechRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) DeleteMechaGameSquadMechRec(recID string) error {
	l := m.Logger("DeleteMechaGameSquadMechRec")

	l.Debug("deleting mecha_game_squad_mech record ID >%s<", recID)

	_, err := m.GetMechaGameSquadMechRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.MechaGameSquadMechRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

func (m *Domain) RemoveMechaGameSquadMechRec(recID string) error {
	l := m.Logger("RemoveMechaGameSquadMechRec")

	l.Debug("removing mecha_game_squad_mech record ID >%s<", recID)

	r := m.MechaGameSquadMechRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
