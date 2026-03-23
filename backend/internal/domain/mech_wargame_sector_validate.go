package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameSectorArgs struct {
	currRec *mech_wargame_record.MechWargameSector
	nextRec *mech_wargame_record.MechWargameSector
}

func (m *Domain) validateMechWargameSectorRecForCreate(rec *mech_wargame_record.MechWargameSector) error {
	args := &validateMechWargameSectorArgs{nextRec: rec}
	return validateMechWargameSectorRec(args, false)
}

func (m *Domain) validateMechWargameSectorRecForUpdate(currRec, nextRec *mech_wargame_record.MechWargameSector) error {
	args := &validateMechWargameSectorArgs{currRec: currRec, nextRec: nextRec}
	return validateMechWargameSectorRec(args, true)
}

func validateMechWargameSectorRec(args *validateMechWargameSectorArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mech_wargame_record.FieldMechWargameSectorName, rec.Name); err != nil {
		return err
	}

	validTerrainTypes := map[string]bool{
		mech_wargame_record.SectorTerrainTypeOpen:   true,
		mech_wargame_record.SectorTerrainTypeUrban:  true,
		mech_wargame_record.SectorTerrainTypeForest: true,
		mech_wargame_record.SectorTerrainTypeRough:  true,
		mech_wargame_record.SectorTerrainTypeWater:  true,
	}
	if rec.TerrainType == "" {
		rec.TerrainType = mech_wargame_record.SectorTerrainTypeOpen
	}
	if !validTerrainTypes[rec.TerrainType] {
		return InvalidField(mech_wargame_record.FieldMechWargameSectorTerrainType, rec.TerrainType, "must be one of: open, urban, forest, rough, water")
	}

	return nil
}
