package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameSectorInstanceRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameSectorInstance, error) {
	l := m.Logger("GetManyMechaGameSectorInstanceRecs")

	l.Debug("getting many mecha_game_sector_instance records opts >%#v<", opts)

	r := m.MechaGameSectorInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameSectorInstanceRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameSectorInstance, error) {
	l := m.Logger("GetMechaGameSectorInstanceRec")

	l.Debug("getting mecha_game_sector_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameSectorInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameSectorInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameSectorInstanceRec(rec *mecha_game_record.MechaGameSectorInstance) (*mecha_game_record.MechaGameSectorInstance, error) {
	l := m.Logger("CreateMechaGameSectorInstanceRec")

	l.Debug("creating mecha_game_sector_instance record >%#v<", rec)

	if err := m.validateMechaGameSectorInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_sector_instance record >%v<", err)
		return rec, err
	}

	r := m.MechaGameSectorInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) RemoveMechaGameSectorInstanceRec(recID string) error {
	l := m.Logger("RemoveMechaGameSectorInstanceRec")

	l.Debug("removing mecha_game_sector_instance record ID >%s<", recID)

	r := m.MechaGameSectorInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
