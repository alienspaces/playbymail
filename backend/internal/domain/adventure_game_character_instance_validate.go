package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameCharacterInstanceArgs struct {
	nextRec *adventure_game_record.AdventureGameCharacterInstance
	currRec *adventure_game_record.AdventureGameCharacterInstance
}

func (m *Domain) populateAdventureGameCharacterInstanceValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameCharacterInstance) (*validateAdventureGameCharacterInstanceArgs, error) {
	args := &validateAdventureGameCharacterInstanceArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameCharacterInstanceRecForCreate(rec *adventure_game_record.AdventureGameCharacterInstance) error {
	args, err := m.populateAdventureGameCharacterInstanceValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameCharacterInstanceRecForCreate(args)
}

func (m *Domain) validateAdventureGameCharacterInstanceRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameCharacterInstance) error {
	args, err := m.populateAdventureGameCharacterInstanceValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameCharacterInstanceRecForUpdate(args)
}

func validateAdventureGameCharacterInstanceRecForCreate(args *validateAdventureGameCharacterInstanceArgs) error {
	return validateAdventureGameCharacterInstanceRec(args, false)
}

func validateAdventureGameCharacterInstanceRecForUpdate(args *validateAdventureGameCharacterInstanceArgs) error {
	return validateAdventureGameCharacterInstanceRec(args, true)
}

func validateAdventureGameCharacterInstanceRec(args *validateAdventureGameCharacterInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceAdventureGameCharacterID, rec.AdventureGameCharacterID); err != nil {
		return err
	}

	if rec.AdventureGameLocationInstanceID != "" {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceAdventureGameLocationInstanceID, rec.AdventureGameLocationInstanceID); err != nil {
			return err
		}
	}

	if rec.Health < 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameCharacterInstanceHealth,
			fmt.Sprintf("%d", rec.Health),
			"health must be zero or greater",
		)
	}

	return nil
}
