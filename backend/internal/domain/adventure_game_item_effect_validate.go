package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameItemEffectArgs struct {
	nextRec *adventure_game_record.AdventureGameItemEffect
	currRec *adventure_game_record.AdventureGameItemEffect
}

func (m *Domain) populateAdventureGameItemEffectValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameItemEffect) (*validateAdventureGameItemEffectArgs, error) {
	args := &validateAdventureGameItemEffectArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameItemEffectRecForCreate(rec *adventure_game_record.AdventureGameItemEffect) error {
	args, err := m.populateAdventureGameItemEffectValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemEffectRecForCreate(args)
}

func (m *Domain) validateAdventureGameItemEffectRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameItemEffect) error {
	args, err := m.populateAdventureGameItemEffectValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemEffectRecForUpdate(args)
}

func validateAdventureGameItemEffectRecForCreate(args *validateAdventureGameItemEffectArgs) error {
	return validateAdventureGameItemEffectRec(args, false)
}

func validateAdventureGameItemEffectRecForUpdate(args *validateAdventureGameItemEffectArgs) error {
	return validateAdventureGameItemEffectRec(args, true)
}

func validateAdventureGameItemEffectRec(args *validateAdventureGameItemEffectArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemEffectID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemEffectGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemEffectAdventureGameItemID, rec.AdventureGameItemID); err != nil {
		return err
	}

	if err := domain.ValidateEnumField(
		adventure_game_record.FieldAdventureGameItemEffectActionType,
		rec.ActionType,
		adventure_game_record.AdventureGameItemEffectActionTypes,
	); err != nil {
		return err
	}

	if err := domain.ValidateEnumField(
		adventure_game_record.FieldAdventureGameItemEffectEffectType,
		rec.EffectType,
		adventure_game_record.AdventureGameItemEffectEffectTypes,
	); err != nil {
		return err
	}

	// Conditional validation based on effect_type
	switch rec.EffectType {
	case adventure_game_record.AdventureGameItemEffectEffectTypeDamageTarget,
		adventure_game_record.AdventureGameItemEffectEffectTypeDamageWielder,
		adventure_game_record.AdventureGameItemEffectEffectTypeHealTarget,
		adventure_game_record.AdventureGameItemEffectEffectTypeHealWielder,
		adventure_game_record.AdventureGameItemEffectEffectTypeWeaponDamage,
		adventure_game_record.AdventureGameItemEffectEffectTypeArmorDefense:
		if !rec.ResultValueMin.Valid {
			return InvalidField(
				adventure_game_record.FieldAdventureGameItemEffectResultValueMin,
				"",
				fmt.Sprintf("result_value_min is required for effect_type %q", rec.EffectType),
			)
		}
		if !rec.ResultValueMax.Valid {
			return InvalidField(
				adventure_game_record.FieldAdventureGameItemEffectResultValueMax,
				"",
				fmt.Sprintf("result_value_max is required for effect_type %q", rec.EffectType),
			)
		}
		if rec.ResultValueMax.Int32 < rec.ResultValueMin.Int32 {
			return InvalidField(
				adventure_game_record.FieldAdventureGameItemEffectResultValueMax,
				fmt.Sprintf("%d", rec.ResultValueMax.Int32),
				"result_value_max must be >= result_value_min",
			)
		}
	case adventure_game_record.AdventureGameItemEffectEffectTypeTeleport:
		if !rec.ResultAdventureGameLocationID.Valid || rec.ResultAdventureGameLocationID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameItemEffectResultAdventureGameLocationID,
				"",
				fmt.Sprintf("result_adventure_game_location_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameItemEffectEffectTypeOpenLink,
		adventure_game_record.AdventureGameItemEffectEffectTypeCloseLink:
		if !rec.ResultAdventureGameLocationLinkID.Valid || rec.ResultAdventureGameLocationLinkID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameItemEffectResultAdventureGameLocationLinkID,
				"",
				fmt.Sprintf("result_adventure_game_location_link_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameItemEffectEffectTypeGiveItem,
		adventure_game_record.AdventureGameItemEffectEffectTypeRemoveItem:
		if !rec.ResultAdventureGameItemID.Valid || rec.ResultAdventureGameItemID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameItemEffectResultAdventureGameItemID,
				"",
				fmt.Sprintf("result_adventure_game_item_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameItemEffectEffectTypeSummonCreature:
		if !rec.ResultAdventureGameCreatureID.Valid || rec.ResultAdventureGameCreatureID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameItemEffectResultAdventureGameCreatureID,
				"",
				fmt.Sprintf("result_adventure_game_creature_id is required for effect_type %q", rec.EffectType),
			)
		}
	}

	return nil
}
