package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameCreatureInstanceArgs struct {
	nextRec *adventure_game_record.AdventureGameCreatureInstance
	currRec *adventure_game_record.AdventureGameCreatureInstance
}

func (m *Domain) populateAdventureGameCreatureInstanceValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameCreatureInstance) (*validateAdventureGameCreatureInstanceArgs, error) {
	args := &validateAdventureGameCreatureInstanceArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameCreatureInstanceRecForCreate(rec *adventure_game_record.AdventureGameCreatureInstance) error {
	args, err := m.populateAdventureGameCreatureInstanceValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameCreatureInstanceRecForCreate(args)
}

func (m *Domain) validateAdventureGameCreatureInstanceRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameCreatureInstance) error {
	args, err := m.populateAdventureGameCreatureInstanceValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameCreatureInstanceRecForUpdate(args)
}

func validateAdventureGameCreatureInstanceRecForCreate(args *validateAdventureGameCreatureInstanceArgs) error {
	return validateAdventureGameCreatureInstanceRec(args, false)
}

func validateAdventureGameCreatureInstanceRecForUpdate(args *validateAdventureGameCreatureInstanceArgs) error {
	return validateAdventureGameCreatureInstanceRec(args, true)
}

func validateAdventureGameCreatureInstanceRec(args *validateAdventureGameCreatureInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, rec.AdventureGameCreatureID); err != nil {
		return err
	}

	if rec.AdventureGameLocationInstanceID != "" {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, rec.AdventureGameLocationInstanceID); err != nil {
			return err
		}
	}

	if rec.Health < 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameCreatureInstanceHealth,
			fmt.Sprintf("%d", rec.Health),
			"health must be zero or greater",
		)
	}

	return nil
}
