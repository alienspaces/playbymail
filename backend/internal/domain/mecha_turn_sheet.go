package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (m *Domain) GetManyMechaTurnSheetRecs(opts *coresql.Options) ([]*mecha_record.MechaTurnSheet, error) {
	l := m.Logger("GetManyMechaTurnSheetRecs")

	l.Debug("getting many mecha_turn_sheet records opts >%#v<", opts)

	r := m.MechaTurnSheetRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechaTurnSheetRec(recID string, lock *coresql.Lock) (*mecha_record.MechaTurnSheet, error) {
	l := m.Logger("GetMechaTurnSheetRec")

	l.Debug("getting mecha_turn_sheet record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechaTurnSheetRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mecha_record.TableMechaTurnSheet, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechaTurnSheetRec(rec *mecha_record.MechaTurnSheet) (*mecha_record.MechaTurnSheet, error) {
	l := m.Logger("CreateMechaTurnSheetRec")

	l.Debug("creating mecha_turn_sheet record >%#v<", rec)

	if err := m.validateMechaTurnSheetRecForCreate(rec); err != nil {
		l.Warn("failed to validate mecha_turn_sheet record >%v<", err)
		return rec, err
	}

	r := m.MechaTurnSheetRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) RemoveMechaTurnSheetRec(recID string) error {
	l := m.Logger("RemoveMechaTurnSheetRec")

	l.Debug("removing mecha_turn_sheet record ID >%s<", recID)

	r := m.MechaTurnSheetRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
