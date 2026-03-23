package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (m *Domain) GetManyMechWargameTurnSheetRecs(opts *coresql.Options) ([]*mech_wargame_record.MechWargameTurnSheet, error) {
	l := m.Logger("GetManyMechWargameTurnSheetRecs")

	l.Debug("getting many mech_wargame_turn_sheet records opts >%#v<", opts)

	r := m.MechWargameTurnSheetRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

func (m *Domain) GetMechWargameTurnSheetRec(recID string, lock *coresql.Lock) (*mech_wargame_record.MechWargameTurnSheet, error) {
	l := m.Logger("GetMechWargameTurnSheetRec")

	l.Debug("getting mech_wargame_turn_sheet record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.MechWargameTurnSheetRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(mech_wargame_record.TableMechWargameTurnSheet, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) CreateMechWargameTurnSheetRec(rec *mech_wargame_record.MechWargameTurnSheet) (*mech_wargame_record.MechWargameTurnSheet, error) {
	l := m.Logger("CreateMechWargameTurnSheetRec")

	l.Debug("creating mech_wargame_turn_sheet record >%#v<", rec)

	if err := m.validateMechWargameTurnSheetRecForCreate(rec); err != nil {
		l.Warn("failed to validate mech_wargame_turn_sheet record >%v<", err)
		return rec, err
	}

	r := m.MechWargameTurnSheetRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) RemoveMechWargameTurnSheetRec(recID string) error {
	l := m.Logger("RemoveMechWargameTurnSheetRec")

	l.Debug("removing mech_wargame_turn_sheet record ID >%s<", recID)

	r := m.MechWargameTurnSheetRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
