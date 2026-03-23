package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameChassisArgs struct {
	currRec *mech_wargame_record.MechWargameChassis
	nextRec *mech_wargame_record.MechWargameChassis
}

func (m *Domain) validateMechWargameChassisRecForCreate(rec *mech_wargame_record.MechWargameChassis) error {
	args := &validateMechWargameChassisArgs{nextRec: rec}
	return validateMechWargameChassisRec(args, false)
}

func (m *Domain) validateMechWargameChassisRecForUpdate(currRec, nextRec *mech_wargame_record.MechWargameChassis) error {
	args := &validateMechWargameChassisArgs{currRec: currRec, nextRec: nextRec}
	return validateMechWargameChassisRec(args, true)
}

func validateMechWargameChassisRec(args *validateMechWargameChassisArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameChassisID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameChassisGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mech_wargame_record.FieldMechWargameChassisName, rec.Name); err != nil {
		return err
	}

	validClasses := map[string]bool{
		mech_wargame_record.ChassisClassLight:   true,
		mech_wargame_record.ChassisClassMedium:  true,
		mech_wargame_record.ChassisClassHeavy:   true,
		mech_wargame_record.ChassisClassAssault: true,
	}
	if rec.ChassisClass == "" {
		rec.ChassisClass = mech_wargame_record.ChassisClassMedium
	}
	if !validClasses[rec.ChassisClass] {
		return InvalidField(mech_wargame_record.FieldMechWargameChassisChassisClass, rec.ChassisClass, "must be one of: light, medium, heavy, assault")
	}

	if rec.ArmorPoints <= 0 {
		return InvalidField(mech_wargame_record.FieldMechWargameChassisArmorPoints, "", "armor_points must be greater than 0")
	}

	if rec.StructurePoints <= 0 {
		return InvalidField(mech_wargame_record.FieldMechWargameChassisStructurePoints, "", "structure_points must be greater than 0")
	}

	if rec.HeatCapacity <= 0 {
		return InvalidField(mech_wargame_record.FieldMechWargameChassisHeatCapacity, "", "heat_capacity must be greater than 0")
	}

	if rec.Speed <= 0 {
		return InvalidField(mech_wargame_record.FieldMechWargameChassisSpeed, "", "speed must be greater than 0")
	}

	return nil
}
