package player_schema

import "gitlab.com/alienspaces/playbymail/schema/api/common_schema"

// JoinGameInfoResponseData contains game and subscription info for displaying the join form.
type JoinGameInfoResponseData struct {
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

// JoinGameInfoResponse is the response for GET /api/v1/game-subscriptions/:id/join.
type JoinGameInfoResponse struct {
	Data  *JoinGameInfoResponseData    `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}

// JoinGameSubmitRequest is the request for POST /join.
type JoinGameSubmitRequest struct {
	Email                 string `json:"email"`
	Name                  string `json:"name"`
	CharacterName         string `json:"character_name,omitempty"`
	PostalAddressLine1    string `json:"postal_address_line1"`
	PostalAddressLine2    string `json:"postal_address_line2,omitempty"`
	StateProvince         string `json:"state_province"`
	Country               string `json:"country"`
	PostalCode            string `json:"postal_code"`
	DeliveryEmail         bool   `json:"delivery_email"`
	DeliveryPhysicalPost  bool   `json:"delivery_physical_post"`
	DeliveryPhysicalLocal bool   `json:"delivery_physical_local"`
}

// JoinGameSubmitResponseData is the data for the join submit response.
type JoinGameSubmitResponseData struct {
	GameSubscriptionID string `json:"game_subscription_id"`
	GameInstanceID     string `json:"game_instance_id"`
	GameID             string `json:"game_id"`
}

// JoinGameSubmitResponse is the response for POST /join.
type JoinGameSubmitResponse struct {
	Data  *JoinGameSubmitResponseData  `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}
