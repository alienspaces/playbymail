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

	// Enforce at-most-one starter lance per game at the application layer so
	// callers receive a clear error message instead of a DB unique-constraint violation.
	if rec.IsPlayerStarter {
		existing, err := m.GetManyMechaLanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaLanceGameID, Val: rec.GameID},
				{Col: mecha_record.FieldMechaLanceIsPlayerStarter, Val: true},
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

	if rec.IsPlayerStarter {
		// Starter template: no owner fields required or permitted.
	} else if rec.MechaComputerOpponentID.Valid {
		// Computer-opponent-owned lance: validate the opponent ID, skip account fields.
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceMechaComputerOpponentID, rec.MechaComputerOpponentID.String); err != nil {
			return err
		}
	} else {
		// Human-owned lance: account fields are required.
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceAccountID, rec.AccountID.String); err != nil {
			return err
		}

		if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceAccountUserID, rec.AccountUserID.String); err != nil {
			return err
		}
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaLanceName, rec.Name); err != nil {
		return err
	}

	return nil
}
