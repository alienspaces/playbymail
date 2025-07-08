package schema

import "time"

// LocationResponseData -
type LocationResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type LocationResponse struct {
	Response
	*LocationResponseData
}

type LocationCollectionResponse = []*LocationResponseData

type LocationRequest struct {
	Request
	GameID      string `json:"game_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LocationQueryParams struct {
	QueryParamsPagination
	LocationResponseData
}
