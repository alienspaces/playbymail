package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaSquadArgs struct {
	currRec *mecha_record.MechaSquad
	nextRec *mecha_record.MechaSquad
}

func (m *Domain) validateMechaSquadRecForCreate(rec *mecha_record.MechaSquad) error {
	args := &validateMechaSquadArgs{nextRec: rec}
	if err := validateMechaSquadRec(args, false); err != nil {
		return err
	}

	if rec.SquadType == mecha_record.SquadTypeStarter {
		existing, err := m.GetManyMechaSquadRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaSquadGameID, Val: rec.GameID},
				{Col: mecha_record.FieldMechaSquadSquadType, Val: mecha_record.SquadTypeStarter},
			},
			Limit: 1,
		})
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return coreerror.NewInvalidDataError("this game already has a player starter squad")
		}
	}

	return nil
}

func (m *Domain) validateMechaSquadRecForUpdate(currRec, nextRec *mecha_record.MechaSquad) error {
	args := &validateMechaSquadArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaSquadRec(args, true)
}

func validateMechaSquadRec(args *validateMechaSquadArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadGameID, rec.GameID); err != nil {
		return err
	}

	if rec.SquadType != mecha_record.SquadTypeStarter && rec.SquadType != mecha_record.SquadTypeOpponent {
		return coreerror.NewInvalidDataError("squad_type must be 'starter' or 'opponent'")
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaSquadName, rec.Name); err != nil {
		return err
	}

	return nil
}
