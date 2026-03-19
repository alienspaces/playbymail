package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationObjectInstanceArgs struct {
	nextRec *adventure_game_record.AdventureGameLocationObjectInstance
	currRec *adventure_game_record.AdventureGameLocationObjectInstance
}

func (m *Domain) populateAdventureGameLocationObjectInstanceValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocationObjectInstance) (*validateAdventureGameLocationObjectInstanceArgs, error) {
	args := &validateAdventureGameLocationObjectInstanceArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationObjectInstanceRecForCreate(rec *adventure_game_record.AdventureGameLocationObjectInstance) error {
	args, err := m.populateAdventureGameLocationObjectInstanceValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectInstanceRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationObjectInstanceRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocationObjectInstance) error {
	args, err := m.populateAdventureGameLocationObjectInstanceValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectInstanceRecForUpdate(args)
}

func validateAdventureGameLocationObjectInstanceRecForCreate(args *validateAdventureGameLocationObjectInstanceArgs) error {
	return validateAdventureGameLocationObjectInstanceRec(args, false)
}

func validateAdventureGameLocationObjectInstanceRecForUpdate(args *validateAdventureGameLocationObjectInstanceArgs) error {
	return validateAdventureGameLocationObjectInstanceRec(args, true)
}

func validateAdventureGameLocationObjectInstanceRec(args *validateAdventureGameLocationObjectInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationObjectID, rec.AdventureGameLocationObjectID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationInstanceID, rec.AdventureGameLocationInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectInstanceCurrentAdventureGameLocationObjectStateID, rec.CurrentAdventureGameLocationObjectStateID); err != nil {
		return err
	}

	return nil
}
