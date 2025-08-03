package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationLinkRequirementResponseData -
type AdventureGameLocationLinkRequirementResponseData struct {
	ID                 string     `json:"id"`
	GameID             string     `json:"game_id"`
	GameLocationLinkID string     `json:"game_location_link_id"`
	RequirementType    string     `json:"requirement_type"`
	RequirementValue   string     `json:"requirement_value"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationLinkRequirementResponse struct {
	Data       *AdventureGameLocationLinkRequirementResponseData `json:"data"`
	Error      *common_schema.ResponseError                      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination                 `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkRequirementCollectionResponse struct {
	Data       []*AdventureGameLocationLinkRequirementResponseData `json:"data"`
	Error      *common_schema.ResponseError                        `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination                   `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkRequirementRequest struct {
	common_schema.Request
	GameLocationLinkID string `json:"game_location_link_id"`
	RequirementType    string `json:"requirement_type"`
	RequirementValue   string `json:"requirement_value"`
}
