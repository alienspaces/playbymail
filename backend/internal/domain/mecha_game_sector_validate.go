package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameSectorArgs struct {
	currRec *mecha_game_record.MechaGameSector
	nextRec *mecha_game_record.MechaGameSector
}

func (m *Domain) validateMechaGameSectorRecForCreate(rec *mecha_game_record.MechaGameSector) error {
	args := &validateMechaGameSectorArgs{nextRec: rec}
	return validateMechaGameSectorRec(args, false)
}

func (m *Domain) validateMechaGameSectorRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameSector) error {
	args := &validateMechaGameSectorArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameSectorRec(args, true)
}

func validateMechaGameSectorRec(args *validateMechaGameSectorArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameSectorName, rec.Name); err != nil {
		return err
	}

	validTerrainTypes := map[string]bool{
		mecha_game_record.SectorTerrainTypeOpen:   true,
		mecha_game_record.SectorTerrainTypeUrban:  true,
		mecha_game_record.SectorTerrainTypeForest: true,
		mecha_game_record.SectorTerrainTypeRough:  true,
		mecha_game_record.SectorTerrainTypeWater:  true,
	}
	if rec.TerrainType == "" {
		rec.TerrainType = mecha_game_record.SectorTerrainTypeOpen
	}
	if !validTerrainTypes[rec.TerrainType] {
		return InvalidField(mecha_game_record.FieldMechaGameSectorTerrainType, rec.TerrainType, "must be one of: open, urban, forest, rough, water")
	}

	return nil
}
