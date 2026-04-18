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

	if rec.ArmorPoints <= 0 {
		return InvalidField(mecha_game_record.FieldMechaGameChassisArmorPoints, "", "armor_points must be greater than 0")
	}

	if rec.StructurePoints <= 0 {
		return InvalidField(mecha_game_record.FieldMechaGameChassisStructurePoints, "", "structure_points must be greater than 0")
	}

	if rec.HeatCapacity <= 0 {
		return InvalidField(mecha_game_record.FieldMechaGameChassisHeatCapacity, "", "heat_capacity must be greater than 0")
	}

	if rec.Speed <= 0 {
		return InvalidField(mecha_game_record.FieldMechaGameChassisSpeed, "", "speed must be greater than 0")
	}

	return nil
}
