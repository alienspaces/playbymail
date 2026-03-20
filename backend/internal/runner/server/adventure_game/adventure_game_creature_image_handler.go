package adventure_game

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	_ "golang.org/x/image/webp"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	UploadCreatureImage = "upload-creature-image"
	GetCreatureImage    = "get-creature-image"
	DeleteCreatureImage = "delete-creature-image"
)

func adventureGameCreatureImageHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameCreatureImageHandlerConfig")

	l.Debug("Adding adventure game creature image handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	cfg[UploadCreatureImage] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id/image",
		HandlerFunc: uploadCreatureImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Upload creature portrait image",
			Description: "Upload a portrait image for an adventure game creature. Accepts multipart form data with 'image' file. Images must be WebP, PNG, or JPEG format, max 1MB.",
		},
	}

	cfg[GetCreatureImage] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id/image",
		HandlerFunc: getCreatureImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Get creature portrait image",
			Description: "Get the portrait image for an adventure game creature.",
		},
	}

	cfg[DeleteCreatureImage] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id/image",
		HandlerFunc: deleteCreatureImageHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Delete creature portrait image",
			Description: "Delete the portrait image for an adventure game creature.",
		},
	}

	return cfg, nil
}

func uploadCreatureImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "uploadCreatureImageHandler")

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.RequiredPathParameter("game_id")
	}

	creatureID := pp.ByName("creature_id")
	if creatureID == "" {
		return coreerror.RequiredPathParameter("creature_id")
	}

	if _, err := authorizeAdventureGameDesigner(l, r, mm, gameID); err != nil {
		return err
	}

	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		return err
	}

	creatureRec, err := mm.GetAdventureGameCreatureRec(creatureID, nil)
	if err != nil {
		return err
	}
	if creatureRec.GameID != gameID {
		return coreerror.NewNotFoundError("creature", creatureID)
	}

	if err := r.ParseMultipartForm(2 << 20); err != nil {
		return coreerror.NewInvalidDataError("failed to parse multipart form: %v", err)
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return coreerror.NewInvalidDataError("image file is required")
	}
	defer file.Close()

	l.Info("received file >%s<", header.Filename)

	imageData, err := io.ReadAll(file)
	if err != nil {
		return coreerror.NewInvalidDataError("failed to read image data")
	}

	if len(imageData) > game_record.GameImageMaxSize {
		return coreerror.NewInvalidDataError("image file too large, max 1MB")
	}

	mimeType := http.DetectContentType(imageData)
	switch mimeType {
	case "image/webp":
		mimeType = game_record.GameImageMimeTypeWebP
	case "image/png":
		mimeType = game_record.GameImageMimeTypePNG
	case "image/jpeg":
		mimeType = game_record.GameImageMimeTypeJPEG
	default:
		return coreerror.NewInvalidDataError("invalid image format, must be WebP, PNG, or JPEG")
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return coreerror.NewInvalidDataError("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	rec := &game_record.GameImage{
		GameID:    gameID,
		RecordID:  nullstring.FromString(creatureID),
		Type:      game_record.GameImageTypeAsset,
		ImageData: imageData,
		MimeType:  mimeType,
		FileSize:  len(imageData),
		Width:     width,
		Height:    height,
	}

	rec, err = mm.UpsertGameImageRec(rec)
	if err != nil {
		return err
	}

	l.Info("successfully saved creature image record id >%s< for creature >%s<", rec.ID, creatureID)

	res, err := mapper.GameImageRecordToResponse(l, rec, "")
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func getCreatureImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getCreatureImageHandler")

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.RequiredPathParameter("game_id")
	}

	creatureID := pp.ByName("creature_id")
	if creatureID == "" {
		return coreerror.RequiredPathParameter("creature_id")
	}

	_, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		return err
	}

	creatureRec, err := mm.GetAdventureGameCreatureRec(creatureID, nil)
	if err != nil {
		return err
	}
	if creatureRec.GameID != gameID {
		return coreerror.NewNotFoundError("creature", creatureID)
	}

	recordID := nullstring.FromString(creatureID)
	img, err := mm.GetGameImageRecByGameAndType(gameID, recordID, game_record.GameImageTypeAsset, "")
	if err != nil {
		return err
	}

	res, err := mapper.CreatureImageToResponse(l, gameID, creatureID, img)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteCreatureImageHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteCreatureImageHandler")

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.RequiredPathParameter("game_id")
	}

	creatureID := pp.ByName("creature_id")
	if creatureID == "" {
		return coreerror.RequiredPathParameter("creature_id")
	}

	if _, err := authorizeAdventureGameDesigner(l, r, mm, gameID); err != nil {
		return err
	}

	_, err := mm.GetGameRec(gameID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	creatureRec, err := mm.GetAdventureGameCreatureRec(creatureID, nil)
	if err != nil {
		return err
	}
	if creatureRec.GameID != gameID {
		return coreerror.NewNotFoundError("creature", creatureID)
	}

	recordID := nullstring.FromString(creatureID)
	if err := mm.DeleteGameImageByGameAndType(gameID, recordID, game_record.GameImageTypeAsset, ""); err != nil {
		return err
	}

	l.Info("deleted portrait image for creature >%s<", creatureID)

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
