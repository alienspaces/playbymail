package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationObjectEffectArgs struct {
	nextRec *adventure_game_record.AdventureGameLocationObjectEffect
	currRec *adventure_game_record.AdventureGameLocationObjectEffect
}

func (m *Domain) populateAdventureGameLocationObjectEffectValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocationObjectEffect) (*validateAdventureGameLocationObjectEffectArgs, error) {
	args := &validateAdventureGameLocationObjectEffectArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationObjectEffectRecForCreate(rec *adventure_game_record.AdventureGameLocationObjectEffect) error {
	args, err := m.populateAdventureGameLocationObjectEffectValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectEffectRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationObjectEffectRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocationObjectEffect) error {
	args, err := m.populateAdventureGameLocationObjectEffectValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationObjectEffectRecForUpdate(args)
}

func validateAdventureGameLocationObjectEffectRecForCreate(args *validateAdventureGameLocationObjectEffectArgs) error {
	return validateAdventureGameLocationObjectEffectRec(args, false)
}

func validateAdventureGameLocationObjectEffectRecForUpdate(args *validateAdventureGameLocationObjectEffectArgs) error {
	return validateAdventureGameLocationObjectEffectRec(args, true)
}

func validateAdventureGameLocationObjectEffectRec(args *validateAdventureGameLocationObjectEffectArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectEffectID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectEffectGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationObjectEffectAdventureGameLocationObjectID, rec.AdventureGameLocationObjectID); err != nil {
		return err
	}

	if err := domain.ValidateEnumField(
		adventure_game_record.FieldAdventureGameLocationObjectEffectActionType,
		rec.ActionType,
		adventure_game_record.AdventureGameLocationObjectEffectActionTypes,
	); err != nil {
		return err
	}


	if err := domain.ValidateEnumField(
		adventure_game_record.FieldAdventureGameLocationObjectEffectEffectType,
		rec.EffectType,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypes,
	); err != nil {
		return err
	}

	// Conditional validation based on effect_type
	switch rec.EffectType {
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState:
		if !rec.ResultState.Valid || rec.ResultState.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultState,
				"",
				fmt.Sprintf("result_state is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState:
		if !rec.ResultState.Valid || rec.ResultState.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultState,
				"",
				fmt.Sprintf("result_state is required for effect_type %q", rec.EffectType),
			)
		}
		if !rec.ResultAdventureGameLocationObjectID.Valid || rec.ResultAdventureGameLocationObjectID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultAdventureGameLocationObjectID,
				"",
				fmt.Sprintf("result_adventure_game_location_object_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHideObject:
		if !rec.ResultAdventureGameLocationObjectID.Valid || rec.ResultAdventureGameLocationObjectID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultAdventureGameLocationObjectID,
				"",
				fmt.Sprintf("result_adventure_game_location_object_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHeal:
		if !rec.ResultValueMin.Valid {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultValueMin,
				"",
				fmt.Sprintf("result_value_min is required for effect_type %q", rec.EffectType),
			)
		}
		if !rec.ResultValueMax.Valid {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultValueMax,
				"",
				fmt.Sprintf("result_value_max is required for effect_type %q", rec.EffectType),
			)
		}
		if rec.ResultValueMax.Int32 < rec.ResultValueMin.Int32 {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultValueMax,
				fmt.Sprintf("%d", rec.ResultValueMax.Int32),
				"result_value_max must be >= result_value_min",
			)
		}
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeTeleport:
		if !rec.ResultAdventureGameLocationID.Valid || rec.ResultAdventureGameLocationID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultAdventureGameLocationID,
				"",
				fmt.Sprintf("result_adventure_game_location_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveItem:
		if !rec.ResultAdventureGameItemID.Valid || rec.ResultAdventureGameItemID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultAdventureGameItemID,
				"",
				fmt.Sprintf("result_adventure_game_item_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeOpenLink,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeCloseLink:
		if !rec.ResultAdventureGameLocationLinkID.Valid || rec.ResultAdventureGameLocationLinkID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultAdventureGameLocationLinkID,
				"",
				fmt.Sprintf("result_adventure_game_location_link_id is required for effect_type %q", rec.EffectType),
			)
		}
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeSummonCreature:
		if !rec.ResultAdventureGameCreatureID.Valid || rec.ResultAdventureGameCreatureID.String == "" {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationObjectEffectResultAdventureGameCreatureID,
				"",
				fmt.Sprintf("result_adventure_game_creature_id is required for effect_type %q", rec.EffectType),
			)
		}
	}

	return nil
}
