package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameCharacterResponseData -
type AdventureGameCharacterResponseData struct {
	ID        string     `json:"id"`
	GameID    string     `json:"game_id"`
	AccountID string     `json:"account_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameCharacterResponse struct {
	Data       *AdventureGameCharacterResponseData `json:"data"`
	Error      *common_schema.ResponseError        `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination   `json:"pagination,omitempty"`
}

type AdventureGameCharacterCollectionResponse struct {
	Data       []*AdventureGameCharacterResponseData `json:"data"`
	Error      *common_schema.ResponseError          `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination     `json:"pagination,omitempty"`
}

type AdventureGameCharacterRequest struct {
	common_schema.Request
	AccountID     string `json:"account_id"`
	AccountUserID string `json:"account_user_id"`
	Name          string `json:"name"`
}

type AdventureGameCharacterQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameCharacterResponseData
}
