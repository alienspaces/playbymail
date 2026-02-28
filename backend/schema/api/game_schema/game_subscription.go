package game_schema

import "time"

type GameSubscriptionRequest struct {
	GameID               string  `json:"game_id"`
	AccountID            string  `json:"account_id"`
	AccountUserID        *string `json:"account_user_id,omitempty"`
	AccountUserContactID *string `json:"account_user_contact_id,omitempty"`
	SubscriptionType     string  `json:"subscription_type"`
	InstanceLimit        *int32  `json:"instance_limit,omitempty"`
}

type GameSubscriptionResponseData struct {
	ID               string     `json:"id"`
	GameID           string     `json:"game_id"`
	AccountID        string     `json:"account_id"`
	AccountUserID    *string    `json:"account_user_id,omitempty"`
	GameInstanceIDs  []string   `json:"game_instance_ids,omitempty"`
	SubscriptionType string     `json:"subscription_type"`
	InstanceLimit    *int32     `json:"instance_limit,omitempty"`
	Status           string     `json:"status"`
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
