package game_schema

import "time"

type GameAdministrationRequest struct {
	GameID             string `json:"game_id"`
	AccountID          string `json:"account_id"`
	GrantedByAccountID string `json:"granted_by_account_id"`
}

type GameAdministrationResponseData struct {
	ID                 string     `json:"id"`
	GameID             string     `json:"game_id"`
	AccountID          string     `json:"account_id"`
	GrantedByAccountID string     `json:"granted_by_account_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type GameAdministrationResponse struct {
	Data *GameAdministrationResponseData `json:"data"`
}

type GameAdministrationCollectionResponse struct {
	Data []*GameAdministrationResponseData `json:"data"`
}
