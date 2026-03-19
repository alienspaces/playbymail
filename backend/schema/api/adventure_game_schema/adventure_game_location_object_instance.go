package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationObjectInstanceResponseData -
type AdventureGameLocationObjectInstanceResponseData struct {
	ID                              string     `json:"id"`
	GameID                          string     `json:"game_id"`
	GameInstanceID                  string     `json:"game_instance_id"`
	AdventureGameLocationObjectID   string     `json:"adventure_game_location_object_id"`
	AdventureGameLocationInstanceID string     `json:"adventure_game_location_instance_id"`
	CurrentState                    string     `json:"current_state"`
	IsVisible                       bool       `json:"is_visible"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       *time.Time `json:"updated_at,omitempty"`
	DeletedAt                       *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationObjectInstanceResponse struct {
	Data       *AdventureGameLocationObjectInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError                     `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination                `json:"pagination,omitempty"`
}

type AdventureGameLocationObjectInstanceCollectionResponse struct {
	Data       []*AdventureGameLocationObjectInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError                       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination                  `json:"pagination,omitempty"`
}
