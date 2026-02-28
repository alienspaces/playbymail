package game_schema

import "time"

type GameSubscriptionInstanceRequest struct {
	GameSubscriptionID string `json:"game_subscription_id"`
	GameInstanceID     string `json:"game_instance_id"`
}

type GameSubscriptionInstanceResponseData struct {
	ID                      string     `json:"id"`
	AccountID               string     `json:"account_id"`
	GameSubscriptionID      string     `json:"game_subscription_id"`
	GameInstanceID          string     `json:"game_instance_id"`
	TurnSheetTokenMasked    *string    `json:"turn_sheet_token_masked,omitempty"`
	TurnSheetTokenExpiresAt *time.Time `json:"turn_sheet_token_expires_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}

type GameSubscriptionInstanceResponse struct {
	Data *GameSubscriptionInstanceResponseData `json:"data"`
}

type GameSubscriptionInstanceCollectionResponse struct {
	Data []*GameSubscriptionInstanceResponseData `json:"data"`
}
