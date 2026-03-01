package catalog_schema

import (
	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// CatalogGameInstanceData describes a game instance available for players to join.
type CatalogGameInstanceData struct {
	ID                    string `json:"id"`
	RequiredPlayerCount   int    `json:"required_player_count"`
	PlayerCount           int    `json:"player_count"`
	DeliveryPhysicalPost  bool   `json:"delivery_physical_post"`
	DeliveryPhysicalLocal bool   `json:"delivery_physical_local"`
	DeliveryEmail         bool   `json:"delivery_email"`
}

// CatalogGameResponseData is an entry in the public game catalog.
type CatalogGameResponseData struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	GameType           string                     `json:"game_type"`
	TurnDurationHours  int                        `json:"turn_duration_hours"`
	AvailableInstances []*CatalogGameInstanceData `json:"available_instances"`
}

// CatalogGameCollectionResponse is the response body for GET /api/v1/catalog/games.
type CatalogGameCollectionResponse struct {
	Data  []*CatalogGameResponseData    `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}
