package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationObjectEffectResponseData -
type AdventureGameLocationObjectEffectResponseData struct {
	ID                                              string     `json:"id"`
	GameID                                          string     `json:"game_id"`
	AdventureGameLocationObjectID                   string     `json:"adventure_game_location_object_id"`
	ActionType                                      string     `json:"action_type"`
	RequiredAdventureGameLocationObjectStateID      *string    `json:"required_adventure_game_location_object_state_id,omitempty"`
	RequiredAdventureGameItemID                     *string    `json:"required_adventure_game_item_id,omitempty"`
	ResultDescription                               string     `json:"result_description"`
	EffectType                                      string     `json:"effect_type"`
	ResultAdventureGameLocationObjectStateID        *string    `json:"result_adventure_game_location_object_state_id,omitempty"`
	ResultAdventureGameItemID                       *string    `json:"result_adventure_game_item_id,omitempty"`
	ResultAdventureGameLocationLinkID               *string    `json:"result_adventure_game_location_link_id,omitempty"`
	ResultAdventureGameCreatureID                   *string    `json:"result_adventure_game_creature_id,omitempty"`
	ResultAdventureGameLocationObjectID             *string    `json:"result_adventure_game_location_object_id,omitempty"`
	ResultAdventureGameLocationID                   *string    `json:"result_adventure_game_location_id,omitempty"`
	ResultValueMin                                  *int32     `json:"result_value_min,omitempty"`
	ResultValueMax                                  *int32     `json:"result_value_max,omitempty"`
	IsRepeatable                                    bool       `json:"is_repeatable"`
	CreatedAt                                       time.Time  `json:"created_at"`
	UpdatedAt                                       *time.Time `json:"updated_at,omitempty"`
	DeletedAt                                       *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationObjectEffectResponse struct {
	Data       *AdventureGameLocationObjectEffectResponseData `json:"data"`
	Error      *common_schema.ResponseError                   `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination              `json:"pagination,omitempty"`
}

type AdventureGameLocationObjectEffectCollectionResponse struct {
	Data       []*AdventureGameLocationObjectEffectResponseData `json:"data"`
	Error      *common_schema.ResponseError                     `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination                `json:"pagination,omitempty"`
}

type AdventureGameLocationObjectEffectRequest struct {
	common_schema.Request
	AdventureGameLocationObjectID                   string  `json:"adventure_game_location_object_id"`
	ActionType                                      string  `json:"action_type"`
	RequiredAdventureGameLocationObjectStateID      *string `json:"required_adventure_game_location_object_state_id,omitempty"`
	RequiredAdventureGameItemID                     *string `json:"required_adventure_game_item_id,omitempty"`
	ResultDescription                               string  `json:"result_description"`
	EffectType                                      string  `json:"effect_type"`
	ResultAdventureGameLocationObjectStateID        *string `json:"result_adventure_game_location_object_state_id,omitempty"`
	ResultAdventureGameItemID                       *string `json:"result_adventure_game_item_id,omitempty"`
	ResultAdventureGameLocationLinkID               *string `json:"result_adventure_game_location_link_id,omitempty"`
	ResultAdventureGameCreatureID                   *string `json:"result_adventure_game_creature_id,omitempty"`
	ResultAdventureGameLocationObjectID             *string `json:"result_adventure_game_location_object_id,omitempty"`
	ResultAdventureGameLocationID                   *string `json:"result_adventure_game_location_id,omitempty"`
	ResultValueMin                                  *int32  `json:"result_value_min,omitempty"`
	ResultValueMax                                  *int32  `json:"result_value_max,omitempty"`
	IsRepeatable                                    bool    `json:"is_repeatable,omitempty"`
}
