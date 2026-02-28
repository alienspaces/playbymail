package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationArgs struct {
	nextRec *adventure_game_record.AdventureGameLocation
	currRec *adventure_game_record.AdventureGameLocation
}

func (m *Domain) populateAdventureGameLocationValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocation) (*validateAdventureGameLocationArgs, error) {
	args := &validateAdventureGameLocationArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationRecForCreate(rec *adventure_game_record.AdventureGameLocation) error {
	args, err := m.populateAdventureGameLocationValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocation) error {
	args, err := m.populateAdventureGameLocationValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationRecForUpdate(args)
}

func validateAdventureGameLocationRecForCreate(args *validateAdventureGameLocationArgs) error {
	return validateAdventureGameLocationRec(args, false)
}

func validateAdventureGameLocationRecForUpdate(args *validateAdventureGameLocationArgs) error {
	return validateAdventureGameLocationRec(args, true)
}

func validateAdventureGameLocationRec(args *validateAdventureGameLocationArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationName, rec.Name); err != nil {
		return err
	}

	if len(rec.Name) > 255 {
		return InvalidField(adventure_game_record.FieldAdventureGameLocationName, rec.Name, "name exceeds maximum length of 255 characters")
	}

	return nil
}
