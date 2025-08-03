package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameItemInstanceResponseData -
type AdventureGameItemInstanceResponseData struct {
	ID                               string     `json:"id"`
	GameID                           string     `json:"game_id"`
	GameInstanceID                   string     `json:"game_instance_id"`
	AdventureGameItemID              string     `json:"adventure_game_item_id"`
	AdventureGameCharacterInstanceID string     `json:"adventure_game_character_instance_id"`
	AdventureGameCreatureInstanceID  string     `json:"adventure_game_creature_instance_id"`
	AdventureGameLocationInstanceID  string     `json:"adventure_game_location_instance_id"`
	IsEquipped                       bool       `json:"is_equipped"`
	IsUsed                           bool       `json:"is_used"`
	UsesRemaining                    int        `json:"uses_remaining"`
	CreatedAt                        time.Time  `json:"created_at"`
	UpdatedAt                        *time.Time `json:"updated_at,omitempty"`
	DeletedAt                        *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameItemInstanceResponse struct {
	Data       *AdventureGameItemInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError           `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination      `json:"pagination,omitempty"`
}

type AdventureGameItemInstanceCollectionResponse struct {
	Data       []*AdventureGameItemInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError             `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination        `json:"pagination,omitempty"`
}
