package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaLanceInstanceRecs(opts *coresql.Options) ([]*mecha_record.MechaLanceInstance, error) {
	l := m.Logger("GetManyMechaLanceInstanceRecs")

	l.Debug("getting many mecha_lance_instance records opts >%#v<", opts)

	r := m.MechaLanceInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaLanceInstanceRec(recID string, lock *coresql.Lock) (*mecha_record.MechaLanceInstance, error) {
	l := m.Logger("GetMechaLanceInstanceRec")

	l.Debug("getting mecha_lance_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaLanceInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaLanceInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaLanceInstanceRec(rec *mecha_record.MechaLanceInstance) (*mecha_record.MechaLanceInstance, error) {
	l := m.Logger("CreateMechaLanceInstanceRec")

	l.Debug("creating mecha_lance_instance record >%#v<", rec)

	if err := m.validateMechaLanceInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_lance_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaLanceInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) UpdateMechaLanceInstanceRec(rec *mecha_record.MechaLanceInstance) (*mecha_record.MechaLanceInstance, error) {
	l := m.Logger("UpdateMechaLanceInstanceRec")

	currRec, err := m.GetMechaLanceInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating mecha_lance_instance record >%#v<", rec)

	if err := m.validateMechaLanceInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate mecha_lance_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaLanceInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

func (m *Domain) RemoveMechaLanceInstanceRec(recID string) error {
	l := m.Logger("RemoveMechaLanceInstanceRec")

	l.Debug("removing mecha_lance_instance record ID >%s<", recID)

	r := m.MechaLanceInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
