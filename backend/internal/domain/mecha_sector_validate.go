package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaSectorArgs struct {
	currRec *mecha_record.MechaSector
	nextRec *mecha_record.MechaSector
}

func (m *Domain) validateMechaSectorRecForCreate(rec *mecha_record.MechaSector) error {
	args := &validateMechaSectorArgs{nextRec: rec}
	return validateMechaSectorRec(args, false)
}

func (m *Domain) validateMechaSectorRecForUpdate(currRec, nextRec *mecha_record.MechaSector) error {
	args := &validateMechaSectorArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaSectorRec(args, true)
}

func validateMechaSectorRec(args *validateMechaSectorArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaSectorName, rec.Name); err != nil {
		return err
	}

	validTerrainTypes := map[string]bool{
		mecha_record.SectorTerrainTypeOpen:   true,
		mecha_record.SectorTerrainTypeUrban:  true,
		mecha_record.SectorTerrainTypeForest: true,
		mecha_record.SectorTerrainTypeRough:  true,
		mecha_record.SectorTerrainTypeWater:  true,
	}
	if rec.TerrainType == "" {
		rec.TerrainType = mecha_record.SectorTerrainTypeOpen
	}
	if !validTerrainTypes[rec.TerrainType] {
		return InvalidField(mecha_record.FieldMechaSectorTerrainType, rec.TerrainType, "must be one of: open, urban, forest, rough, water")
	}

	return nil
}
