package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameItemEffectResponseData -
type AdventureGameItemEffectResponseData struct {
	ID                                string     `json:"id"`
	GameID                            string     `json:"game_id"`
	AdventureGameItemID               string     `json:"adventure_game_item_id"`
	ActionType                        string     `json:"action_type"`
	RequiredAdventureGameItemID       *string    `json:"required_adventure_game_item_id,omitempty"`
	RequiredAdventureGameLocationID   *string    `json:"required_adventure_game_location_id,omitempty"`
	ResultDescription                 string     `json:"result_description"`
	EffectType                        string     `json:"effect_type"`
	ResultAdventureGameItemID         *string    `json:"result_adventure_game_item_id,omitempty"`
	ResultAdventureGameLocationLinkID *string    `json:"result_adventure_game_location_link_id,omitempty"`
	ResultAdventureGameCreatureID     *string    `json:"result_adventure_game_creature_id,omitempty"`
	ResultAdventureGameLocationID     *string    `json:"result_adventure_game_location_id,omitempty"`
	ResultValueMin                    *int32     `json:"result_value_min,omitempty"`
	ResultValueMax                    *int32     `json:"result_value_max,omitempty"`
	IsRepeatable                      bool       `json:"is_repeatable"`
	CreatedAt                         time.Time  `json:"created_at"`
	UpdatedAt                         *time.Time `json:"updated_at,omitempty"`
	DeletedAt                         *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameItemEffectResponse struct {
	Data       *AdventureGameItemEffectResponseData `json:"data"`
	Error      *common_schema.ResponseError         `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination    `json:"pagination,omitempty"`
}

type AdventureGameItemEffectCollectionResponse struct {
	Data       []*AdventureGameItemEffectResponseData `json:"data"`
	Error      *common_schema.ResponseError           `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination      `json:"pagination,omitempty"`
}

type AdventureGameItemEffectRequest struct {
	common_schema.Request
	AdventureGameItemID               string  `json:"adventure_game_item_id"`
	ActionType                        string  `json:"action_type"`
	RequiredAdventureGameItemID       *string `json:"required_adventure_game_item_id,omitempty"`
	RequiredAdventureGameLocationID   *string `json:"required_adventure_game_location_id,omitempty"`
	ResultDescription                 string  `json:"result_description"`
	EffectType                        string  `json:"effect_type"`
	ResultAdventureGameItemID         *string `json:"result_adventure_game_item_id,omitempty"`
	ResultAdventureGameLocationLinkID *string `json:"result_adventure_game_location_link_id,omitempty"`
	ResultAdventureGameCreatureID     *string `json:"result_adventure_game_creature_id,omitempty"`
	ResultAdventureGameLocationID     *string `json:"result_adventure_game_location_id,omitempty"`
	ResultValueMin                    *int32  `json:"result_value_min,omitempty"`
	ResultValueMax                    *int32  `json:"result_value_max,omitempty"`
	IsRepeatable                      bool    `json:"is_repeatable,omitempty"`
}
