package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameCreatureResponseData -
type AdventureGameCreatureResponseData struct {
	ID                string     `json:"id"`
	GameID            string     `json:"game_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	AttackDamage      int        `json:"attack_damage"`
	Defense           int        `json:"defense"`
	MaxHealth         int        `json:"max_health"`
	Disposition       string     `json:"disposition"`
	AttackMethod      string     `json:"attack_method"`
	AttackDescription string     `json:"attack_description"`
	BodyDecayTurns    int        `json:"body_decay_turns"`
	RespawnTurns      int        `json:"respawn_turns"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameCreatureResponse struct {
	Data       *AdventureGameCreatureResponseData `json:"data"`
	Error      *common_schema.ResponseError       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination  `json:"pagination,omitempty"`
}

type AdventureGameCreatureCollectionResponse struct {
	Data       []*AdventureGameCreatureResponseData `json:"data"`
	Error      *common_schema.ResponseError         `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination    `json:"pagination,omitempty"`
}

type AdventureGameCreatureRequest struct {
	common_schema.Request
	Name              string `json:"name"`
	Description       string `json:"description"`
	AttackDamage      int    `json:"attack_damage,omitempty"`
	Defense           int    `json:"defense,omitempty"`
	MaxHealth         int    `json:"max_health,omitempty"`
	Disposition       string `json:"disposition,omitempty"`
	AttackMethod      string `json:"attack_method,omitempty"`
	AttackDescription string `json:"attack_description,omitempty"`
	BodyDecayTurns    int    `json:"body_decay_turns,omitempty"`
	RespawnTurns      int    `json:"respawn_turns,omitempty"`
}

type AdventureGameCreatureQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameCreatureResponseData
}
