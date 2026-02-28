package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationInstanceArgs struct {
	nextRec *adventure_game_record.AdventureGameLocationInstance
	currRec *adventure_game_record.AdventureGameLocationInstance
}

func (m *Domain) populateAdventureGameLocationInstanceValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocationInstance) (*validateAdventureGameLocationInstanceArgs, error) {
	args := &validateAdventureGameLocationInstanceArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationInstanceRecForCreate(rec *adventure_game_record.AdventureGameLocationInstance) error {
	args, err := m.populateAdventureGameLocationInstanceValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationInstanceRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationInstanceRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocationInstance) error {
	args, err := m.populateAdventureGameLocationInstanceValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationInstanceRecForUpdate(args)
}

func validateAdventureGameLocationInstanceRecForCreate(args *validateAdventureGameLocationInstanceArgs) error {
	return validateAdventureGameLocationInstanceRec(args, false)
}

func validateAdventureGameLocationInstanceRecForUpdate(args *validateAdventureGameLocationInstanceArgs) error {
	return validateAdventureGameLocationInstanceRec(args, true)
}

func validateAdventureGameLocationInstanceRec(args *validateAdventureGameLocationInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceAdventureGameLocationID, rec.AdventureGameLocationID); err != nil {
		return err
	}

	return nil
}

