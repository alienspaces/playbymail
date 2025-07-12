package schema

import "time"

// GameLocationLinkResponseData -
type GameLocationLinkResponseData struct {
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

type GameLocationLinkResponse struct {
	Data       *GameLocationLinkResponseData `json:"data"`
	Error      *ResponseError                `json:"error,omitempty"`
	Pagination *ResponsePagination           `json:"pagination,omitempty"`
}

type GameLocationLinkCollectionResponse struct {
	Data       []*GameLocationLinkResponseData `json:"data"`
	Error      *ResponseError                  `json:"error,omitempty"`
	Pagination *ResponsePagination             `json:"pagination,omitempty"`
}

type GameLocationLinkRequest struct {
	Request
	GameID             string `json:"game_id"`
	FromGameLocationID string `json:"from_game_location_id"`
	ToGameLocationID   string `json:"to_game_location_id"`
	Description        string `json:"description"`
	Name               string `json:"name"`
}

type GameLocationLinkQueryParams struct {
	QueryParamsPagination
	GameLocationLinkResponseData
}
