package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameChassisArgs struct {
	currRec *mecha_game_record.MechaGameChassis
	nextRec *mecha_game_record.MechaGameChassis
}

func (m *Domain) validateMechaGameChassisRecForCreate(rec *mecha_game_record.MechaGameChassis) error {
	args := &validateMechaGameChassisArgs{nextRec: rec}
	return validateMechaGameChassisRec(args, false)
}

func (m *Domain) validateMechaGameChassisRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameChassis) error {
	args := &validateMechaGameChassisArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameChassisRec(args, true)
}

func validateMechaGameChassisRec(args *validateMechaGameChassisArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameChassisID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameChassisGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameChassisName, rec.Name); err != nil {
		return err
	}

	validClasses := map[string]bool{
		mecha_game_record.ChassisClassLight:   true,
		mecha_game_record.ChassisClassMedium:  true,
		mecha_game_record.ChassisClassHeavy:   true,
		mecha_game_record.ChassisClassAssault: true,
	}
	if rec.ChassisClass == "" {
		rec.ChassisClass = mecha_game_record.ChassisClassMedium
	}
	if !validClasses[rec.ChassisClass] {
		return InvalidField(mecha_game_record.FieldMechaGameChassisChassisClass, rec.ChassisClass, "must be one of: light, medium, heavy, assault")
	}

	if rec.ArmorPoints <= 0 || rec.ArmorPoints > maxChassisArmorPoints {
		return InvalidField(mecha_game_record.FieldMechaGameChassisArmorPoints, "", "armor_points must be between 1 and 1000")
	}

	if rec.StructurePoints <= 0 || rec.StructurePoints > maxChassisStructurePoints {
		return InvalidField(mecha_game_record.FieldMechaGameChassisStructurePoints, "", "structure_points must be between 1 and 1000")
	}

	if rec.HeatCapacity <= 0 || rec.HeatCapacity > maxChassisHeatCapacity {
		return InvalidField(mecha_game_record.FieldMechaGameChassisHeatCapacity, "", "heat_capacity must be between 1 and 200")
	}

	if rec.Speed <= 0 || rec.Speed > maxChassisSpeed {
		return InvalidField(mecha_game_record.FieldMechaGameChassisSpeed, "", "speed must be between 1 and 10")
	}

	if rec.SmallSlots < 0 || rec.SmallSlots > maxChassisSlotsPerSize {
		return InvalidField(mecha_game_record.FieldMechaGameChassisSmallSlots, "", "small_slots must be between 0 and 10")
	}
	if rec.MediumSlots < 0 || rec.MediumSlots > maxChassisSlotsPerSize {
		return InvalidField(mecha_game_record.FieldMechaGameChassisMediumSlots, "", "medium_slots must be between 0 and 10")
	}
	if rec.LargeSlots < 0 || rec.LargeSlots > maxChassisSlotsPerSize {
		return InvalidField(mecha_game_record.FieldMechaGameChassisLargeSlots, "", "large_slots must be between 0 and 10")
	}
	if rec.SmallSlots+rec.MediumSlots+rec.LargeSlots == 0 {
		return InvalidField(mecha_game_record.FieldMechaGameChassisSmallSlots, "", "chassis must have at least one slot")
	}

	return nil
}

// Upper bounds on chassis stats. These keep designer input within playable
// ranges — in particular Speed caps the BFS depth used by player reachability
// checks and AI movement, so uncapped values would make orders and AI decisions
// disproportionately expensive without adding gameplay value (long-range
// weapons only reach 2 hops).
const (
	maxChassisArmorPoints     = 1000
	maxChassisStructurePoints = 1000
	maxChassisHeatCapacity    = 200
	maxChassisSpeed           = 10
	// maxChassisSlotsPerSize caps each slot band (small/medium/large). Kept in
	// sync with the mecha_game_chassis_slot_bounds_check SQL constraint so
	// designer input is bounded both at the API and database layers.
	maxChassisSlotsPerSize = 10
)
