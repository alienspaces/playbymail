package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameSectorInstanceRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameSectorInstance, error) {
	l := m.Logger("GetManyMechWargameSectorInstanceRecs")

	l.Debug("getting many mech_wargame_sector_instance records opts >%#v<", opts)

	r := m.MechWargameSectorInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameSectorInstanceRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameSectorInstance, error) {
	l := m.Logger("GetMechWargameSectorInstanceRec")

	l.Debug("getting mech_wargame_sector_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameSectorInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameSectorInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameSectorInstanceRec(rec *mech_wargame_record.MechWargameSectorInstance) (*mech_wargame_record.MechWargameSectorInstance, error) {
	l := m.Logger("CreateMechWargameSectorInstanceRec")

	l.Debug("creating mech_wargame_sector_instance record >%#v<", rec)

	if err := m.validateMechWargameSectorInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_sector_instance record >%v<", err)
		return rec, err
	}

	r := m.MechWargameSectorInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) RemoveMechWargameSectorInstanceRec(recID string) error {
	l := m.Logger("RemoveMechWargameSectorInstanceRec")

	l.Debug("removing mech_wargame_sector_instance record ID >%s<", recID)

	r := m.MechWargameSectorInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
