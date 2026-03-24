package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaSectorInstanceRecs(opts *coresql.Options) ([]*mecha_record.MechaSectorInstance, error) {
	l := m.Logger("GetManyMechaSectorInstanceRecs")

	l.Debug("getting many mecha_sector_instance records opts >%#v<", opts)

	r := m.MechaSectorInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaSectorInstanceRec(recID string, lock *coresql.Lock) (*mecha_record.MechaSectorInstance, error) {
	l := m.Logger("GetMechaSectorInstanceRec")

	l.Debug("getting mecha_sector_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaSectorInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaSectorInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaSectorInstanceRec(rec *mecha_record.MechaSectorInstance) (*mecha_record.MechaSectorInstance, error) {
	l := m.Logger("CreateMechaSectorInstanceRec")

	l.Debug("creating mecha_sector_instance record >%#v<", rec)

	if err := m.validateMechaSectorInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_sector_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaSectorInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) RemoveMechaSectorInstanceRec(recID string) error {
	l := m.Logger("RemoveMechaSectorInstanceRec")

	l.Debug("removing mecha_sector_instance record ID >%s<", recID)

	r := m.MechaSectorInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
