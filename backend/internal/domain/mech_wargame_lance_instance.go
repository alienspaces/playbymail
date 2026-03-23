package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameLanceInstanceRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameLanceInstance, error) {
	l := m.Logger("GetManyMechWargameLanceInstanceRecs")

	l.Debug("getting many mech_wargame_lance_instance records opts >%#v<", opts)

	r := m.MechWargameLanceInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameLanceInstanceRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameLanceInstance, error) {
	l := m.Logger("GetMechWargameLanceInstanceRec")

	l.Debug("getting mech_wargame_lance_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameLanceInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameLanceInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameLanceInstanceRec(rec *mech_wargame_record.MechWargameLanceInstance) (*mech_wargame_record.MechWargameLanceInstance, error) {
	l := m.Logger("CreateMechWargameLanceInstanceRec")

	l.Debug("creating mech_wargame_lance_instance record >%#v<", rec)

	if err := m.validateMechWargameLanceInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_lance_instance record >%v<", err)
		return rec, err
	}

	r := m.MechWargameLanceInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) RemoveMechWargameLanceInstanceRec(recID string) error {
	l := m.Logger("RemoveMechWargameLanceInstanceRec")

	l.Debug("removing mech_wargame_lance_instance record ID >%s<", recID)

	r := m.MechWargameLanceInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
