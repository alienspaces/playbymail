package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (m *Domain) GetManyMechaGameTurnSheetRecs(opts *coresql.Options) ([]*mecha_game_record.MechaGameTurnSheet, error) {
	l := m.Logger("GetManyMechaGameTurnSheetRecs")

	l.Debug("getting many mecha_game_turn_sheet records opts >%#v<", opts)

	r := m.MechaGameTurnSheetRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaGameTurnSheetRec(recID string, lock *coresql.Lock) (*mecha_game_record.MechaGameTurnSheet, error) {
	l := m.Logger("GetMechaGameTurnSheetRec")

	l.Debug("getting mecha_game_turn_sheet record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaGameTurnSheetRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_game_record.TableMechaGameTurnSheet, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaGameTurnSheetRec(rec *mecha_game_record.MechaGameTurnSheet) (*mecha_game_record.MechaGameTurnSheet, error) {
	l := m.Logger("CreateMechaGameTurnSheetRec")

	l.Debug("creating mecha_game_turn_sheet record >%#v<", rec)

	if err := m.validateMechaGameTurnSheetRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_game_turn_sheet record >%v<", err)
		return rec, err
	}

	r := m.MechaGameTurnSheetRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) RemoveMechaGameTurnSheetRec(recID string) error {
	l := m.Logger("RemoveMechaGameTurnSheetRec")

	l.Debug("removing mecha_game_turn_sheet record ID >%s<", recID)

	r := m.MechaGameTurnSheetRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
