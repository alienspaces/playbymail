package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaSquadInstanceRecs(opts *coresql.Options) ([]*mecha_record.MechaSquadInstance, error) {
	l := m.Logger("GetManyMechaSquadInstanceRecs")

	l.Debug("getting many mecha_squad_instance records opts >%#v<", opts)

	r := m.MechaSquadInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaSquadInstanceRec(recID string, lock *coresql.Lock) (*mecha_record.MechaSquadInstance, error) {
	l := m.Logger("GetMechaSquadInstanceRec")

	l.Debug("getting mecha_squad_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaSquadInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaSquadInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaSquadInstanceRec(rec *mecha_record.MechaSquadInstance) (*mecha_record.MechaSquadInstance, error) {
	l := m.Logger("CreateMechaSquadInstanceRec")

	l.Debug("creating mecha_squad_instance record >%#v<", rec)

	if err := m.validateMechaSquadInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_squad_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaSquadInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaSquadInstanceRec(rec *mecha_record.MechaSquadInstance) (*mecha_record.MechaSquadInstance, error) {
	l := m.Logger("UpdateMechaSquadInstanceRec")

	currRec, err := m.GetMechaSquadInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_squad_instance record >%#v<", rec)

	if err := m.validateMechaSquadInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_squad_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaSquadInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) RemoveMechaSquadInstanceRec(recID string) error {
	l := m.Logger("RemoveMechaSquadInstanceRec")

	l.Debug("removing mecha_squad_instance record ID >%s<", recID)

	r := m.MechaSquadInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
