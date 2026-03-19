package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationObjectStateArgs struct {
	nextRec *adventure_game_record.AdventureGameLocationObjectState
	currRec *adventure_game_record.AdventureGameLocationObjectState
}

func (m *Domain) populateAdventureGameLocationObjectStateValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocationObjectState) (*validateAdventureGameLocationObjectStateArgs, error) {
	args := &validateAdventureGameLocationObjectStateArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationObjectStateRecForCreate(rec *adventure_game_record.AdventureGameLocationObjectState) error {
	args, err := m.populateAdventureGameLocationObjectStateValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectStateRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationObjectStateRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocationObjectState) error {
	args, err := m.populateAdventureGameLocationObjectStateValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectStateRecForUpdate(args)
}

func validateAdventureGameLocationObjectStateRecForCreate(args *validateAdventureGameLocationObjectStateArgs) error {
	return validateAdventureGameLocationObjectStateRec(args, false)
}

func validateAdventureGameLocationObjectStateRecForUpdate(args *validateAdventureGameLocationObjectStateArgs) error {
	return validateAdventureGameLocationObjectStateRec(args, true)
}

func validateAdventureGameLocationObjectStateRec(args *validateAdventureGameLocationObjectStateArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectStateID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectStateGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectStateAdventureGameLocationObjectID, rec.AdventureGameLocationObjectID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationObjectStateName, rec.Name); err != nil {
		return err
	}

	return nil
}
