package api

import "time"

// AdventureGameItemInstanceResponseData -
type AdventureGameItemInstanceResponseData struct {
	ID                      string     `json:"id"`
	GameID                  string     `json:"game_id"`
	GameItemID              string     `json:"game_item_id"`
	GameInstanceID          string     `json:"game_instance_id"`
	GameLocationInstanceID  string     `json:"game_location_instance_id,omitempty"`  // Leave this omitempty as we will be adding GameCharacterInstanceID and GameCreatureInstanceID as options the future
	GameCharacterInstanceID string     `json:"game_character_instance_id,omitempty"` // Leave this omitempty as we will be adding GameCreatureInstanceID as options the future
	GameCreatureInstanceID  string     `json:"game_creature_instance_id,omitempty"`  // Leave this omitempty as we will be adding GameCreatureInstanceID as options the future
	IsEquipped              bool       `json:"is_equipped"`
	IsUsed                  bool       `json:"is_used"`
	UsesRemaining           *int       `json:"uses_remaining,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameItemInstanceResponse struct {
	Data       *AdventureGameItemInstanceResponseData `json:"data"`
	Error      *ResponseError                         `json:"error,omitempty"`
	Pagination *ResponsePagination                    `json:"pagination,omitempty"`
}

type AdventureGameItemInstanceCollectionResponse struct {
	Data       []*AdventureGameItemInstanceResponseData `json:"data"`
	Error      *ResponseError                           `json:"error,omitempty"`
	Pagination *ResponsePagination                      `json:"pagination,omitempty"`
}

type AdventureGameItemInstanceRequest struct {
	Request
	GameID                  string `json:"game_id"`
	GameItemID              string `json:"game_item_id"`
	GameInstanceID          string `json:"game_instance_id"`
	GameLocationInstanceID  string `json:"game_location_instance_id,omitempty"`
	GameCharacterInstanceID string `json:"game_character_instance_id,omitempty"`
	GameCreatureInstanceID  string `json:"game_creature_instance_id,omitempty"`
	IsEquipped              bool   `json:"is_equipped"`
	IsUsed                  bool   `json:"is_used"`
	UsesRemaining           *int   `json:"uses_remaining,omitempty"`
}
