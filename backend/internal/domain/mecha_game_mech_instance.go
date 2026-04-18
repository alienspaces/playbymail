package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameMechInstanceRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameMechInstance, error) {
	l := m.Logger("GetManyMechaGameMechInstanceRecs")

	l.Debug("getting many mecha_game_mech_instance records opts >%#v<", opts)

	r := m.MechaGameMechInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameMechInstanceRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameMechInstance, error) {
	l := m.Logger("GetMechaGameMechInstanceRec")

	l.Debug("getting mecha_game_mech_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameMechInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameMechInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameMechInstanceRec(rec *mecha_game_record.MechaGameMechInstance) (*mecha_game_record.MechaGameMechInstance, error) {
	l := m.Logger("CreateMechaGameMechInstanceRec")

	l.Debug("creating mecha_game_mech_instance record >%#v<", rec)

	if err := m.validateMechaGameMechInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_mech_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaGameMechInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaGameMechInstanceRec(rec *mecha_game_record.MechaGameMechInstance) (*mecha_game_record.MechaGameMechInstance, error) {
	l := m.Logger("UpdateMechaGameMechInstanceRec")

	currRec, err := m.GetMechaGameMechInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_game_mech_instance record >%#v<", rec)

	if err := m.validateMechaGameMechInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_game_mech_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaGameMechInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) RemoveMechaGameMechInstanceRec(recID string) error {
	l := m.Logger("RemoveMechaGameMechInstanceRec")

	l.Debug("removing mecha_game_mech_instance record ID >%s<", recID)

	r := m.MechaGameMechInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
