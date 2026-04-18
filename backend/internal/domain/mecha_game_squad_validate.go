package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameSquadArgs struct {
	currRec *mecha_game_record.MechaGameSquad
	nextRec *mecha_game_record.MechaGameSquad
}

func (m *Domain) validateMechaGameSquadRecForCreate(rec *mecha_game_record.MechaGameSquad) error {
	args := &validateMechaGameSquadArgs{nextRec: rec}
	if err := validateMechaGameSquadRec(args, false); err != nil {
		return err
	}

	if rec.SquadType == mecha_game_record.SquadTypeStarter {
		existing, err := m.GetManyMechaGameSquadRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameSquadGameID, Val: rec.GameID},
				{Col: mecha_game_record.FieldMechaGameSquadSquadType, Val: mecha_game_record.SquadTypeStarter},
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

func (m *Domain) validateMechaGameSquadRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameSquad) error {
	args := &validateMechaGameSquadArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameSquadRec(args, true)
}

func validateMechaGameSquadRec(args *validateMechaGameSquadArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadGameID, rec.GameID); err != nil {
		return err
	}

	if rec.SquadType != mecha_game_record.SquadTypeStarter && rec.SquadType != mecha_game_record.SquadTypeOpponent {
		return coreerror.NewInvalidDataError("squad_type must be 'starter' or 'opponent'")
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameSquadName, rec.Name); err != nil {
		return err
	}

	return nil
}
