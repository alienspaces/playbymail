package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameImage table name
const TableGameImage string = "game_image"

// GameImage field names
const (
	FieldGameImageID        string = "id"
	FieldGameImageGameID    string = "game_id"
	FieldGameImageRecordID  string = "record_id"
	FieldGameImageType      string = "type"
	FieldGameImageImageData string = "image_data"
	FieldGameImageMimeType  string = "mime_type"
	FieldGameImageFileSize  string = "file_size"
	FieldGameImageWidth     string = "width"
	FieldGameImageHeight    string = "height"
)

// GameImage type constants
const (
	GameImageTypeTurnSheetBackground string = "turn_sheet_background"
	GameImageTypeAsset               string = "asset"
)

// GameImageTypes is the set of valid image types
var GameImageTypes = set.New(
	GameImageTypeTurnSheetBackground,
	GameImageTypeAsset,
)

// GameImage MIME type constants
const (
	GameImageMimeTypeWebP string = "image/webp"
	GameImageMimeTypePNG  string = "image/png"
	GameImageMimeTypeJPEG string = "image/jpeg"
)

// GameImageMimeTypes is the set of valid MIME types
var GameImageMimeTypes = set.New(
	GameImageMimeTypeWebP,
	GameImageMimeTypePNG,
	GameImageMimeTypeJPEG,
)

// GameImage dimension constraints
// Note: Constraints are intentionally relaxed to allow flexible image sizes.
// The turn sheet preview will show how images actually appear.
// A client-side cropping tool may be added in the future for better UX.
const (
	GameImageMinWidth  int = 400     // Minimum to prevent tiny unusable images
	GameImageMaxWidth  int = 4000    // Maximum to prevent huge uploads
	GameImageMinHeight int = 200     // Allows banner/header style images
	GameImageMaxHeight int = 6000    // Maximum to prevent huge uploads
	GameImageMaxSize   int = 1048576 // 1MB in bytes
)

// GameImage recommended dimensions (A4 @ 300 DPI)
const (
	GameImageRecommendedWidth  int = 2480
	GameImageRecommendedHeight int = 3508
)

// GameImage struct represents a game image record
type GameImage struct {
	record.Record
	GameID    string         `db:"game_id"`
	RecordID  sql.NullString `db:"record_id"`
	Type      string         `db:"type"`
	ImageData []byte         `db:"image_data"`
	MimeType  string         `db:"mime_type"`
	FileSize  int            `db:"file_size"`
	Width     int            `db:"width"`
	Height    int            `db:"height"`
}

// ToNamedArgs converts the GameImage record to named arguments for database operations
func (r *GameImage) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameImageGameID] = r.GameID
	args[FieldGameImageRecordID] = r.RecordID
	args[FieldGameImageType] = r.Type
	args[FieldGameImageImageData] = r.ImageData
	args[FieldGameImageMimeType] = r.MimeType
	args[FieldGameImageFileSize] = r.FileSize
	args[FieldGameImageWidth] = r.Width
	args[FieldGameImageHeight] = r.Height
	return args
}
