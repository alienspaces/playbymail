package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationObjectArgs struct {
	nextRec *adventure_game_record.AdventureGameLocationObject
	currRec *adventure_game_record.AdventureGameLocationObject
}

func (m *Domain) populateAdventureGameLocationObjectValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocationObject) (*validateAdventureGameLocationObjectArgs, error) {
	args := &validateAdventureGameLocationObjectArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationObjectRecForCreate(rec *adventure_game_record.AdventureGameLocationObject) error {
	args, err := m.populateAdventureGameLocationObjectValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationObjectRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocationObject) error {
	args, err := m.populateAdventureGameLocationObjectValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectRecForUpdate(args)
}

func validateAdventureGameLocationObjectRecForCreate(args *validateAdventureGameLocationObjectArgs) error {
	return validateAdventureGameLocationObjectRec(args, false)
}

func validateAdventureGameLocationObjectRecForUpdate(args *validateAdventureGameLocationObjectArgs) error {
	return validateAdventureGameLocationObjectRec(args, true)
}

func validateAdventureGameLocationObjectRec(args *validateAdventureGameLocationObjectArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectAdventureGameLocationID, rec.AdventureGameLocationID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationObjectName, rec.Name); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationObjectDescription, rec.Description); err != nil {
		return err
	}

	if rec.InitialAdventureGameLocationObjectStateID.Valid {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectInitialAdventureGameLocationObjectStateID, rec.InitialAdventureGameLocationObjectStateID.String); err != nil {
			return err
		}
	}

	return nil
}
