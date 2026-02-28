package account_subscription_schema

import "time"

type AccountSubscriptionRequest struct {
	AccountID          string     `json:"account_id"`
	SubscriptionType   string     `json:"subscription_type"`
	SubscriptionPeriod *string    `json:"subscription_period,omitempty"`
	Status             *string    `json:"status,omitempty"`
	AutoRenew          *bool      `json:"auto_renew,omitempty"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
}

type AccountSubscriptionResponseData struct {
	ID                 string     `json:"id"`
	AccountID          string     `json:"account_id"`
	SubscriptionType   string     `json:"subscription_type"`
	SubscriptionPeriod string     `json:"subscription_period"`
	Status             string     `json:"status"`
	AutoRenew          bool       `json:"auto_renew"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type AccountSubscriptionResponse struct {
	Data *AccountSubscriptionResponseData `json:"data"`
}

type AccountSubscriptionCollectionResponse struct {
	Data []*AccountSubscriptionResponseData `json:"data"`
}
