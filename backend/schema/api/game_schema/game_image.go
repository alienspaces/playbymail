package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// GameImageResponseData -
type GameImageResponseData struct {
	ID            string     `json:"id"`
	GameID        string     `json:"game_id"`
	RecordID      string     `json:"record_id,omitempty"`
	Type          string     `json:"type"`
	TurnSheetType string     `json:"turn_sheet_type,omitempty"`
	MimeType      string     `json:"mime_type"`
	FileSize      int        `json:"file_size"`
	Width         int        `json:"width"`
	Height        int        `json:"height"`
	Warning       string     `json:"warning,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type GameImageResponse struct {
	Data       *GameImageResponseData            `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameImageCollectionResponse struct {
	Data       []*GameImageResponseData          `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

// GameTurnSheetImageResponse returns the turn sheet background image for a game
type GameTurnSheetImageResponse struct {
	Data       *GameTurnSheetImageData           `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameTurnSheetImageData struct {
	GameID     string                 `json:"game_id"`
	Background *GameImageResponseData `json:"background,omitempty"`
}

// LocationTurnSheetImageResponse returns the turn sheet background image for a location
type LocationTurnSheetImageResponse struct {
	Data       *LocationTurnSheetImageData       `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type LocationTurnSheetImageData struct {
	GameID     string                 `json:"game_id"`
	LocationID string                 `json:"location_id"`
	Background *GameImageResponseData `json:"background,omitempty"`
}
