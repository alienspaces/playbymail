package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaSector string = "mecha_sector"
)

const (
	FieldMechaSectorID               string = "id"
	FieldMechaSectorGameID           string = "game_id"
	FieldMechaSectorName             string = "name"
	FieldMechaSectorDescription      string = "description"
	FieldMechaSectorTerrainType      string = "terrain_type"
	FieldMechaSectorElevation        string = "elevation"
	FieldMechaSectorIsStartingSector string = "is_starting_sector"
	FieldMechaSectorCreatedAt        string = "created_at"
	FieldMechaSectorUpdatedAt        string = "updated_at"
	FieldMechaSectorDeletedAt        string = "deleted_at"
)

const (
	SectorTerrainTypeOpen   string = "open"
	SectorTerrainTypeUrban  string = "urban"
	SectorTerrainTypeForest string = "forest"
	SectorTerrainTypeRough  string = "rough"
	SectorTerrainTypeWater  string = "water"
)

type MechaSector struct {
	record.Record
	GameID           string `db:"game_id"`
	Name             string `db:"name"`
	Description      string `db:"description"`
	TerrainType      string `db:"terrain_type"`
	Elevation        int    `db:"elevation"`
	IsStartingSector bool   `db:"is_starting_sector"`
}

func (r *MechaSector) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaSectorGameID] = r.GameID
	args[FieldMechaSectorName] = r.Name
	args[FieldMechaSectorDescription] = r.Description
	args[FieldMechaSectorTerrainType] = r.TerrainType
	args[FieldMechaSectorElevation] = r.Elevation
	args[FieldMechaSectorIsStartingSector] = r.IsStartingSector
	return args
}
