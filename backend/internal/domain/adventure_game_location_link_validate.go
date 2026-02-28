package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationLinkArgs struct {
	nextRec *adventure_game_record.AdventureGameLocationLink
	currRec *adventure_game_record.AdventureGameLocationLink
}

func (m *Domain) populateAdventureGameLocationLinkValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocationLink) (*validateAdventureGameLocationLinkArgs, error) {
	args := &validateAdventureGameLocationLinkArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationLinkRecForCreate(rec *adventure_game_record.AdventureGameLocationLink) error {
	args, err := m.populateAdventureGameLocationLinkValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationLinkRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationLinkRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocationLink) error {
	args, err := m.populateAdventureGameLocationLinkValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationLinkRecForUpdate(args)
}

func validateAdventureGameLocationLinkRecForCreate(args *validateAdventureGameLocationLinkArgs) error {
	return validateAdventureGameLocationLinkRec(args, false)
}

func validateAdventureGameLocationLinkRecForUpdate(args *validateAdventureGameLocationLinkArgs) error {
	return validateAdventureGameLocationLinkRec(args, true)
}

func validateAdventureGameLocationLinkRec(args *validateAdventureGameLocationLinkArgs, requireID bool) error {
	rec := args.nextRec

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkFromAdventureGameLocationID, rec.FromAdventureGameLocationID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkToAdventureGameLocationID, rec.ToAdventureGameLocationID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationLinkName, rec.Name); err != nil {
		return err
	}

	if len(rec.Name) > 255 {
		return InvalidField(adventure_game_record.FieldAdventureGameLocationLinkName, rec.Name, "name exceeds maximum length of 255 characters")
	}

	return nil
}
