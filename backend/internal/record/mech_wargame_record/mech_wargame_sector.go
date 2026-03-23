package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameSector string = "mech_wargame_sector"
)

const (
	FieldMechWargameSectorID               string = "id"
	FieldMechWargameSectorGameID           string = "game_id"
	FieldMechWargameSectorName             string = "name"
	FieldMechWargameSectorDescription      string = "description"
	FieldMechWargameSectorTerrainType      string = "terrain_type"
	FieldMechWargameSectorElevation        string = "elevation"
	FieldMechWargameSectorIsStartingSector string = "is_starting_sector"
	FieldMechWargameSectorCreatedAt        string = "created_at"
	FieldMechWargameSectorUpdatedAt        string = "updated_at"
	FieldMechWargameSectorDeletedAt        string = "deleted_at"
)

const (
	SectorTerrainTypeOpen   string = "open"
	SectorTerrainTypeUrban  string = "urban"
	SectorTerrainTypeForest string = "forest"
	SectorTerrainTypeRough  string = "rough"
	SectorTerrainTypeWater  string = "water"
)

type MechWargameSector struct {
	record.Record
	GameID           string `db:"game_id"`
	Name             string `db:"name"`
	Description      string `db:"description"`
	TerrainType      string `db:"terrain_type"`
	Elevation        int    `db:"elevation"`
	IsStartingSector bool   `db:"is_starting_sector"`
}

func (r *MechWargameSector) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameSectorGameID] = r.GameID
	args[FieldMechWargameSectorName] = r.Name
	args[FieldMechWargameSectorDescription] = r.Description
	args[FieldMechWargameSectorTerrainType] = r.TerrainType
	args[FieldMechWargameSectorElevation] = r.Elevation
	args[FieldMechWargameSectorIsStartingSector] = r.IsStartingSector
	return args
}
