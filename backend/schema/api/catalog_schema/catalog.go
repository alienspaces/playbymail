package catalog_schema

import (
	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// CatalogSubscriptionData represents a manager subscription available in the public catalog.
// Each entry aggregates capacity and delivery info across the subscription's linked game instances.
type CatalogSubscriptionData struct {
	GameSubscriptionID    string `json:"game_subscription_id"`
	GameName              string `json:"game_name"`
	GameDescription       string `json:"game_description"`
	GameType              string `json:"game_type"`
	TurnDurationHours     int    `json:"turn_duration_hours"`
	TotalCapacity         int    `json:"total_capacity"`
	TotalPlayers          int    `json:"total_players"`
	DeliveryPhysicalPost  bool   `json:"delivery_physical_post"`
	DeliveryPhysicalLocal bool   `json:"delivery_physical_local"`
	DeliveryEmail         bool   `json:"delivery_email"`
}

// CatalogCollectionResponse is the response body for GET /api/v1/catalog/game-subscriptions.
type CatalogCollectionResponse struct {
	Data  []*CatalogSubscriptionData   `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}
