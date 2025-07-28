package api

import "time"

// AdventureGameLocationLinkResponseData -
type AdventureGameLocationLinkResponseData struct {
	ID                 string     `json:"id"`
	GameID             string     `json:"game_id"`
	FromGameLocationID string     `json:"from_game_location_id"`
	ToGameLocationID   string     `json:"to_game_location_id"`
	Description        string     `json:"description"`
	Name               string     `json:"name"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationLinkResponse struct {
	Data       *AdventureGameLocationLinkResponseData `json:"data"`
	Error      *ResponseError                         `json:"error,omitempty"`
	Pagination *ResponsePagination                    `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkCollectionResponse struct {
	Data       []*AdventureGameLocationLinkResponseData `json:"data"`
	Error      *ResponseError                           `json:"error,omitempty"`
	Pagination *ResponsePagination                      `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkRequest struct {
	Request
	FromGameLocationID string `json:"from_game_location_id"`
	ToGameLocationID   string `json:"to_game_location_id"`
	Description        string `json:"description"`
	Name               string `json:"name"`
}

type AdventureGameLocationLinkQueryParams struct {
	QueryParamsPagination
	AdventureGameLocationLinkResponseData
}
