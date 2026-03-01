package player_schema

import "gitlab.com/alienspaces/playbymail/schema/api/common_schema"

// JoinGameInfoInstanceData describes the game instance in the join game info response.
type JoinGameInfoInstanceData struct {
	ID                    string `json:"id"`
	RequiredPlayerCount   int    `json:"required_player_count"`
	PlayerCount           int    `json:"player_count"`
	DeliveryPhysicalPost  bool   `json:"delivery_physical_post"`
	DeliveryPhysicalLocal bool   `json:"delivery_physical_local"`
	DeliveryEmail         bool   `json:"delivery_email"`
}

// JoinGameInfoResponseData contains game and instance info for displaying the join form.
type JoinGameInfoResponseData struct {
	GameID            string                    `json:"game_id"`
	GameName          string                    `json:"game_name"`
	GameDescription   string                    `json:"game_description"`
	GameType          string                    `json:"game_type"`
	TurnDurationHours int                       `json:"turn_duration_hours"`
	Instance          *JoinGameInfoInstanceData `json:"instance"`
}

// JoinGameInfoResponse is the response for GET /api/v1/player/game-instances/:id/join-game.
type JoinGameInfoResponse struct {
	Data  *JoinGameInfoResponseData    `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}

// JoinGameVerifyEmailRequest is the request for POST /join-game/verify-email.
type JoinGameVerifyEmailRequest struct {
	Email string `json:"email"`
}

// JoinGameVerifyEmailResponseData is the data for the verify-email response.
type JoinGameVerifyEmailResponseData struct {
	HasAccount bool `json:"has_account"`
}

// JoinGameVerifyEmailResponse is the response for POST /join-game/verify-email.
type JoinGameVerifyEmailResponse struct {
	Data  *JoinGameVerifyEmailResponseData `json:"data"`
	Error *common_schema.ResponseError     `json:"error,omitempty"`
}

// JoinGameSubmitRequest is the request for POST /join-game/submit.
type JoinGameSubmitRequest struct {
	Email                  string `json:"email"`
	Name                   string `json:"name"`
	PostalAddressLine1     string `json:"postal_address_line1"`
	PostalAddressLine2     string `json:"postal_address_line2,omitempty"`
	StateProvince          string `json:"state_province"`
	Country                string `json:"country"`
	PostalCode             string `json:"postal_code"`
	DeliveryEmail          bool   `json:"delivery_email"`
	DeliveryPhysicalPost   bool   `json:"delivery_physical_post"`
	DeliveryPhysicalLocal  bool   `json:"delivery_physical_local"`
}

// JoinGameSubmitResponseData is the data for the join-game submit response.
type JoinGameSubmitResponseData struct {
	GameSubscriptionID string `json:"game_subscription_id"`
	GameInstanceID     string `json:"game_instance_id"`
	GameID             string `json:"game_id"`
}

// JoinGameSubmitResponse is the response for POST /join-game/submit.
type JoinGameSubmitResponse struct {
	Data  *JoinGameSubmitResponseData  `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}
