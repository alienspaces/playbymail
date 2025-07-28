package api

import "time"

type GameSubscriptionRequest struct {
	GameID           string `json:"game_id"`
	AccountID        string `json:"account_id"`
	SubscriptionType string `json:"subscription_type"`
}

type GameSubscriptionResponseData struct {
	ID               string     `json:"id"`
	GameID           string     `json:"game_id"`
	AccountID        string     `json:"account_id"`
	SubscriptionType string     `json:"subscription_type"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type GameSubscriptionResponse struct {
	Data *GameSubscriptionResponseData `json:"data"`
}

type GameSubscriptionCollectionResponse struct {
	Data []*GameSubscriptionResponseData `json:"data"`
}
