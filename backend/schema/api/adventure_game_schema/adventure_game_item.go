package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameItemResponseData -
type AdventureGameItemResponseData struct {
	ID             string     `json:"id"`
	GameID         string     `json:"game_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	CanBeEquipped  bool       `json:"can_be_equipped"`
	ItemCategory   string     `json:"item_category,omitempty"`
	EquipmentSlot  string     `json:"equipment_slot,omitempty"`
	IsStartingItem bool       `json:"is_starting_item"`
	CanBeUsed      bool       `json:"can_be_used"`
	Damage         int        `json:"damage"`
	Defense        int        `json:"defense"`
	HealAmount     int        `json:"heal_amount"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameItemResponse struct {
	Data       *AdventureGameItemResponseData    `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AdventureGameItemCollectionResponse struct {
	Data       []*AdventureGameItemResponseData  `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AdventureGameItemRequest struct {
	common_schema.Request
	Name           string `json:"name"`
	Description    string `json:"description"`
	CanBeEquipped  bool   `json:"can_be_equipped"`
	ItemCategory   string `json:"item_category,omitempty"`
	EquipmentSlot  string `json:"equipment_slot,omitempty"`
	IsStartingItem bool   `json:"is_starting_item"`
	CanBeUsed      bool   `json:"can_be_used"`
	Damage         int    `json:"damage"`
	Defense        int    `json:"defense"`
	HealAmount     int    `json:"heal_amount"`
}

type AdventureGameItemQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameItemResponseData
}
