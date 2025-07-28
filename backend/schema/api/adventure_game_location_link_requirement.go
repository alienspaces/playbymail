package api

import "time"

// AdventureGameLocationLinkRequirementResponseData -
type AdventureGameLocationLinkRequirementResponseData struct {
	ID                 string     `json:"id"`
	GameID             string     `json:"game_id"`
	GameLocationLinkID string     `json:"game_location_link_id"`
	GameItemID         string     `json:"game_item_id"`
	Quantity           int        `json:"quantity"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationLinkRequirementResponse struct {
	Data       *AdventureGameLocationLinkRequirementResponseData `json:"data"`
	Error      *ResponseError                                    `json:"error,omitempty"`
	Pagination *ResponsePagination                               `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkRequirementCollectionResponse struct {
	Data       []*AdventureGameLocationLinkRequirementResponseData `json:"data"`
	Error      *ResponseError                                      `json:"error,omitempty"`
	Pagination *ResponsePagination                                 `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkRequirementRequest struct {
	Request
	GameID             string `json:"game_id"`
	GameLocationLinkID string `json:"game_location_link_id"`
	GameItemID         string `json:"game_item_id"`
	Quantity           int    `json:"quantity"`
}
