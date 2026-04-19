package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameEquipmentArgs struct {
	currRec *mecha_game_record.MechaGameEquipment
	nextRec *mecha_game_record.MechaGameEquipment
}

func (m *Domain) validateMechaGameEquipmentRecForCreate(rec *mecha_game_record.MechaGameEquipment) error {
	args := &validateMechaGameEquipmentArgs{nextRec: rec}
	return validateMechaGameEquipmentRec(args, false)
}

func (m *Domain) validateMechaGameEquipmentRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameEquipment) error {
	args := &validateMechaGameEquipmentArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameEquipmentRec(args, true)
}

// validateMechaGameEquipmentRec enforces the shared slot vocabulary, the
// closed effect-kind enum, and the per-kind magnitude caps. heat_cost is
// allowed on every kind (the engine applies it only when the kind's
// "applied-this-turn" predicate fires, so letting designers configure a
// nonzero heat_cost on any kind is intentional).
func validateMechaGameEquipmentRec(args *validateMechaGameEquipmentArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameEquipmentID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameEquipmentGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameEquipmentName, rec.Name); err != nil {
		return err
	}

	if rec.MountSize == "" {
		rec.MountSize = mecha_game_record.EquipmentMountSizeMedium
	}
	if !mecha_game_record.ValidEquipmentMountSize(rec.MountSize) {
		return InvalidField(mecha_game_record.FieldMechaGameEquipmentMountSize, rec.MountSize, "must be one of: small, medium, large")
	}

	if !mecha_game_record.ValidEquipmentEffectKind(rec.EffectKind) {
		return InvalidField(mecha_game_record.FieldMechaGameEquipmentEffectKind, rec.EffectKind,
			"must be one of: heat_sink, targeting_computer, armor_upgrade, jump_jets, ecm, ammo_bin")
	}

	maxMag := mecha_game_record.MagnitudeMaxForEffectKind(rec.EffectKind)
	if rec.Magnitude < 1 || rec.Magnitude > maxMag {
		return InvalidField(mecha_game_record.FieldMechaGameEquipmentMagnitude, "",
			fmt.Sprintf("magnitude must be between 1 and %d for effect_kind %q", maxMag, rec.EffectKind))
	}

	if rec.HeatCost < 0 || rec.HeatCost > maxEquipmentHeatCost {
		return InvalidField(mecha_game_record.FieldMechaGameEquipmentHeatCost, "",
			fmt.Sprintf("heat_cost must be between 0 and %d", maxEquipmentHeatCost))
	}

	return nil
}

const (
	maxEquipmentHeatCost = 20
)
