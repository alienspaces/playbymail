package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameSquadInstanceRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameSquadInstance, error) {
	l := m.Logger("GetManyMechaGameSquadInstanceRecs")

	l.Debug("getting many mecha_game_squad_instance records opts >%#v<", opts)

	r := m.MechaGameSquadInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameSquadInstanceRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameSquadInstance, error) {
	l := m.Logger("GetMechaGameSquadInstanceRec")

	l.Debug("getting mecha_game_squad_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameSquadInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameSquadInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameSquadInstanceRec(rec *mecha_game_record.MechaGameSquadInstance) (*mecha_game_record.MechaGameSquadInstance, error) {
	l := m.Logger("CreateMechaGameSquadInstanceRec")

	l.Debug("creating mecha_game_squad_instance record >%#v<", rec)

	if err := m.validateMechaGameSquadInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_squad_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSquadInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameSquadInstanceRec(rec *mecha_game_record.MechaGameSquadInstance) (*mecha_game_record.MechaGameSquadInstance, error) {
	l := m.Logger("UpdateMechaGameSquadInstanceRec")

	currRec, err := m.GetMechaGameSquadInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_squad_instance record >%#v<", rec)

	if err := m.validateMechaGameSquadInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_squad_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSquadInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) RemoveMechaGameSquadInstanceRec(recID string) error {
	l := m.Logger("RemoveMechaGameSquadInstanceRec")

	l.Debug("removing mecha_game_squad_instance record ID >%s<", recID)

	r := m.MechaGameSquadInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
