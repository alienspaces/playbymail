package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaChassisArgs struct {
	currRec *mecha_record.MechaChassis
	nextRec *mecha_record.MechaChassis
}

func (m *Domain) validateMechaChassisRecForCreate(rec *mecha_record.MechaChassis) error {
	args := &validateMechaChassisArgs{nextRec: rec}
	return validateMechaChassisRec(args, false)
}

func (m *Domain) validateMechaChassisRecForUpdate(currRec, nextRec *mecha_record.MechaChassis) error {
	args := &validateMechaChassisArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaChassisRec(args, true)
}

func validateMechaChassisRec(args *validateMechaChassisArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaChassisID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaChassisGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaChassisName, rec.Name); err != nil {
		return err
	}

	validClasses := map[string]bool{
		mecha_record.ChassisClassLight:   true,
		mecha_record.ChassisClassMedium:  true,
		mecha_record.ChassisClassHeavy:   true,
		mecha_record.ChassisClassAssault: true,
	}
	if rec.ChassisClass == "" {
		rec.ChassisClass = mecha_record.ChassisClassMedium
	}
	if !validClasses[rec.ChassisClass] {
		return InvalidField(mecha_record.FieldMechaChassisChassisClass, rec.ChassisClass, "must be one of: light, medium, heavy, assault")
	}

	if rec.ArmorPoints <= 0 {
		return InvalidField(mecha_record.FieldMechaChassisArmorPoints, "", "armor_points must be greater than 0")
	}

	if rec.StructurePoints <= 0 {
		return InvalidField(mecha_record.FieldMechaChassisStructurePoints, "", "structure_points must be greater than 0")
	}

	if rec.HeatCapacity <= 0 {
		return InvalidField(mecha_record.FieldMechaChassisHeatCapacity, "", "heat_capacity must be greater than 0")
	}

	if rec.Speed <= 0 {
		return InvalidField(mecha_record.FieldMechaChassisSpeed, "", "speed must be greater than 0")
	}

	return nil
}
