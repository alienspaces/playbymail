package schema

import "time"

// GameLocationLinkRequirementResponseData -
type GameLocationLinkRequirementResponseData struct {
	ID                 string     `json:"id"`
	GameLocationLinkID string     `json:"game_location_link_id"`
	GameItemID         string     `json:"game_item_id"`
	Quantity           int        `json:"quantity"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type GameLocationLinkRequirementResponse struct {
	Data       *GameLocationLinkRequirementResponseData `json:"data"`
	Error      *ResponseError                           `json:"error,omitempty"`
	Pagination *ResponsePagination                      `json:"pagination,omitempty"`
}

type GameLocationLinkRequirementCollectionResponse struct {
	Data       []*GameLocationLinkRequirementResponseData `json:"data"`
	Error      *ResponseError                             `json:"error,omitempty"`
	Pagination *ResponsePagination                        `json:"pagination,omitempty"`
}

type GameLocationLinkRequirementRequest struct {
	Request
	GameLocationLinkID string `json:"game_location_link_id"`
	GameItemID         string `json:"game_item_id"`
	Quantity           int    `json:"quantity"`
}
