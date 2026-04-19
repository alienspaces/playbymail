package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameSquadMechArgs struct {
	currRec *mecha_game_record.MechaGameSquadMech
	nextRec *mecha_game_record.MechaGameSquadMech
}

func (m *Domain) validateMechaGameSquadMechRecForCreate(rec *mecha_game_record.MechaGameSquadMech) error {
	args := &validateMechaGameSquadMechArgs{nextRec: rec}
	if err := validateMechaGameSquadMechRec(args, false); err != nil {
		return err
	}
	return m.validateMechaGameSquadMechLoadout(rec)
}

func (m *Domain) validateMechaGameSquadMechRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameSquadMech) error {
	args := &validateMechaGameSquadMechArgs{currRec: currRec, nextRec: nextRec}
	if err := validateMechaGameSquadMechRec(args, true); err != nil {
		return err
	}
	return m.validateMechaGameSquadMechLoadout(nextRec)
}

func validateMechaGameSquadMechRec(args *validateMechaGameSquadMechArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechMechaGameSquadID, rec.MechaGameSquadID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechMechaGameChassisID, rec.MechaGameChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameSquadMechCallsign, rec.Callsign); err != nil {
		return err
	}

	return nil
}

// validateMechaGameSquadMechLoadout ensures the mech's combined weapon +
// equipment config fits the chassis slot budget. An empty config on both
// sides is always allowed (a freshly created mech may be set up before any
// items are assigned).
func (m *Domain) validateMechaGameSquadMechLoadout(rec *mecha_game_record.MechaGameSquadMech) error {
	if rec == nil || (len(rec.WeaponConfig) == 0 && len(rec.EquipmentConfig) == 0) {
		return nil
	}

	chassisRec, err := m.GetMechaGameChassisRec(rec.MechaGameChassisID, nil)
	if err != nil {
		return err
	}

	weaponsByID := make(map[string]*mecha_game_record.MechaGameWeapon, len(rec.WeaponConfig))
	for _, entry := range rec.WeaponConfig {
		if _, ok := weaponsByID[entry.WeaponID]; ok {
			continue
		}
		w, err := m.GetMechaGameWeaponRec(entry.WeaponID, nil)
		if err != nil {
			return err
		}
		weaponsByID[entry.WeaponID] = w
	}

	equipmentByID := make(map[string]*mecha_game_record.MechaGameEquipment, len(rec.EquipmentConfig))
	for _, entry := range rec.EquipmentConfig {
		if _, ok := equipmentByID[entry.EquipmentID]; ok {
			continue
		}
		eq, err := m.GetMechaGameEquipmentRec(entry.EquipmentID, nil)
		if err != nil {
			return err
		}
		equipmentByID[entry.EquipmentID] = eq
	}

	return squadMechLoadoutFitResult(chassisRec, rec.WeaponConfig, weaponsByID, rec.EquipmentConfig, equipmentByID)
}

// squadMechLoadoutFitResult is the pure, record-driven part of the squad-mech
// loadout validator. It is kept separate from validateMechaGameSquadMechLoadout
// so it can be unit-tested without spinning up the domain's DB-backed record
// loaders. Returns nil when the loadout fits, or an InvalidField error
// pointing at weapon_config or equipment_config (whichever the user most
// likely needs to edit) when it does not.
func squadMechLoadoutFitResult(
	chassis *mecha_game_record.MechaGameChassis,
	weaponCfg []mecha_game_record.WeaponConfigEntry,
	weaponsByID map[string]*mecha_game_record.MechaGameWeapon,
	equipmentCfg []mecha_game_record.EquipmentConfigEntry,
	equipmentByID map[string]*mecha_game_record.MechaGameEquipment,
) error {
	if err := ValidateCombinedLoadoutFits(chassis, weaponCfg, weaponsByID, equipmentCfg, equipmentByID); err != nil {
		// Prefer pointing at the equipment field when there is equipment,
		// otherwise the weapon field. The returned message includes the
		// offending item's label so the user can find it either way.
		field := mecha_game_record.FieldMechaGameSquadMechWeaponConfig
		if len(equipmentCfg) > 0 {
			field = mecha_game_record.FieldMechaGameSquadMechEquipmentConfig
		}
		return InvalidField(field, "", err.Error())
	}
	return nil
}
