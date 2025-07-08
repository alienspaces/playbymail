package schema

import "time"

// LocationLinkResponseData -
type LocationLinkResponseData struct {
	ID             string     `json:"id"`
	FromLocationID string     `json:"from_location_id"`
	ToLocationID   string     `json:"to_location_id"`
	Description    string     `json:"description"`
	Name           string     `json:"name"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type LocationLinkResponse struct {
	Response
	*LocationLinkResponseData
}

type LocationLinkCollectionResponse = []*LocationLinkResponseData

type LocationLinkRequest struct {
	Request
	FromLocationID string `json:"from_location_id"`
	ToLocationID   string `json:"to_location_id"`
	Description    string `json:"description"`
	Name           string `json:"name"`
}

type LocationLinkQueryParams struct {
	QueryParamsPagination
	LocationLinkResponseData
}
