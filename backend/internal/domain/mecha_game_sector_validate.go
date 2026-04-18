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

	if rec.Elevation < minSectorElevation || rec.Elevation > maxSectorElevation {
		return InvalidField(mecha_game_record.FieldMechaGameSectorElevation, "", "elevation must be between -10 and 10")
	}

	if rec.CoverModifier < minSectorCoverModifier || rec.CoverModifier > maxSectorCoverModifier {
		return InvalidField(mecha_game_record.FieldMechaGameSectorCoverModifier, "", "cover_modifier must be between -50 and 50")
	}

	return nil
}

// Bounds on sector stats. Cover modifier is added directly to a 50% base hit
// chance (see combat_resolution.go), so values outside ±50 would push the
// clamp at 0/95 without adding play value. Elevation is only used as a
// tactical preference input by the AI (see computer_opponent_decision.go), so
// a compact ±10 range is more than enough for meaningful differentiation.
const (
	minSectorElevation     = -10
	maxSectorElevation     = 10
	minSectorCoverModifier = -50
	maxSectorCoverModifier = 50
)
