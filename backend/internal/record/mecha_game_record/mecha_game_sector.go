package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameSector string = "mecha_game_sector"
)

const (
	FieldMechaGameSectorID               string = "id"
	FieldMechaGameSectorGameID           string = "game_id"
	FieldMechaGameSectorName             string = "name"
	FieldMechaGameSectorDescription      string = "description"
	FieldMechaGameSectorTerrainType      string = "terrain_type"
	FieldMechaGameSectorElevation        string = "elevation"
	FieldMechaGameSectorCoverModifier    string = "cover_modifier"
	FieldMechaGameSectorIsStartingSector string = "is_starting_sector"
	FieldMechaGameSectorCreatedAt        string = "created_at"
	FieldMechaGameSectorUpdatedAt        string = "updated_at"
	FieldMechaGameSectorDeletedAt        string = "deleted_at"
)

const (
	SectorTerrainTypeOpen   string = "open"
	SectorTerrainTypeUrban  string = "urban"
	SectorTerrainTypeForest string = "forest"
	SectorTerrainTypeRough  string = "rough"
	SectorTerrainTypeWater  string = "water"
)

type MechaGameSector struct {
	record.Record
	GameID           string `db:"game_id"`
	Name             string `db:"name"`
	Description      string `db:"description"`
	TerrainType      string `db:"terrain_type"`
	Elevation        int    `db:"elevation"`
	CoverModifier    int    `db:"cover_modifier"`
	IsStartingSector bool   `db:"is_starting_sector"`
}

func (r *MechaGameSector) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameSectorGameID] = r.GameID
	args[FieldMechaGameSectorName] = r.Name
	args[FieldMechaGameSectorDescription] = r.Description
	args[FieldMechaGameSectorTerrainType] = r.TerrainType
	args[FieldMechaGameSectorElevation] = r.Elevation
	args[FieldMechaGameSectorCoverModifier] = r.CoverModifier
	args[FieldMechaGameSectorIsStartingSector] = r.IsStartingSector
	return args
}
