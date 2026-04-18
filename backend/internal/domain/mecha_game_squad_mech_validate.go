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

// validateMechaGameSquadMechLoadout ensures the mech's weapon config fits the
// chassis slot budget. An empty config is always allowed (a freshly created
// mech may be set up before any weapons are assigned).
func (m *Domain) validateMechaGameSquadMechLoadout(rec *mecha_game_record.MechaGameSquadMech) error {
	if rec == nil || len(rec.WeaponConfig) == 0 {
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

	items := MountablesFromWeaponConfig(rec.WeaponConfig, weaponsByID)
	if err := fitsLoadout(LoadoutCapacityFromChassis(chassisRec), items); err != nil {
		return InvalidField(mecha_game_record.FieldMechaGameSquadMechWeaponConfig, "", err.Error())
	}
	return nil
}
