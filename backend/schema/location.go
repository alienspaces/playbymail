package schema

import "time"

// GameLocationResponseData -
type GameLocationResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type GameLocationResponse struct {
	Response
	*GameLocationResponseData
}

type GameLocationCollectionResponse = []*GameLocationResponseData

type GameLocationRequest struct {
	Request
	GameID      string `json:"game_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GameLocationQueryParams struct {
	QueryParamsPagination
	GameLocationResponseData
}
