package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationObjectResponseData -
type AdventureGameLocationObjectResponseData struct {
	ID                                        string     `json:"id"`
	GameID                                    string     `json:"game_id"`
	AdventureGameLocationID                   string     `json:"adventure_game_location_id"`
	Name                                      string     `json:"name"`
	Description                               string     `json:"description"`
	InitialAdventureGameLocationObjectStateID *string    `json:"initial_adventure_game_location_object_state_id,omitempty"`
	IsHidden                                  bool       `json:"is_hidden"`
	CreatedAt                                 time.Time  `json:"created_at"`
	UpdatedAt                                 *time.Time `json:"updated_at,omitempty"`
	DeletedAt                                 *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationObjectResponse struct {
	Data       *AdventureGameLocationObjectResponseData `json:"data"`
	Error      *common_schema.ResponseError             `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination        `json:"pagination,omitempty"`
}

type AdventureGameLocationObjectCollectionResponse struct {
	Data       []*AdventureGameLocationObjectResponseData `json:"data"`
	Error      *common_schema.ResponseError               `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination          `json:"pagination,omitempty"`
}

type AdventureGameLocationObjectRequest struct {
	common_schema.Request
	AdventureGameLocationID                   string `json:"adventure_game_location_id"`
	Name                                      string `json:"name"`
	Description                               string `json:"description"`
	InitialAdventureGameLocationObjectStateID string `json:"initial_adventure_game_location_object_state_id,omitempty"`
	IsHidden                                  bool   `json:"is_hidden,omitempty"`
}
