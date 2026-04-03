package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaLanceArgs struct {
	currRec *mecha_record.MechaLance
	nextRec *mecha_record.MechaLance
}

func (m *Domain) validateMechaLanceRecForCreate(rec *mecha_record.MechaLance) error {
	args := &validateMechaLanceArgs{nextRec: rec}
	if err := validateMechaLanceRec(args, false); err != nil {
		return err
	}

	if rec.LanceType == mecha_record.LanceTypeStarter {
		existing, err := m.GetManyMechaLanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaLanceGameID, Val: rec.GameID},
				{Col: mecha_record.FieldMechaLanceLanceType, Val: mecha_record.LanceTypeStarter},
			},
			Limit: 1,
		})
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return coreerror.NewInvalidDataError("this game already has a player starter lance")
		}
	}

	return nil
}

func (m *Domain) validateMechaLanceRecForUpdate(currRec, nextRec *mecha_record.MechaLance) error {
	args := &validateMechaLanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaLanceRec(args, true)
}

func validateMechaLanceRec(args *validateMechaLanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceGameID, rec.GameID); err != nil {
		return err
	}

	if rec.LanceType != mecha_record.LanceTypeStarter && rec.LanceType != mecha_record.LanceTypeOpponent {
		return coreerror.NewInvalidDataError("lance_type must be 'starter' or 'opponent'")
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaLanceName, rec.Name); err != nil {
		return err
	}

	return nil
}
