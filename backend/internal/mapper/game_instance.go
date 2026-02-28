package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameInstanceRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameInstance) (*game_record.GameInstance, error) {
	l.Debug("mapping game_instance request to record")

	var req game_schema.GameInstanceRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		l.Debug("mapping POST request: delivery_physical_post=%v, delivery_physical_local=%v, delivery_email=%v",
			req.DeliveryPhysicalPost, req.DeliveryPhysicalLocal, req.DeliveryEmail)
		rec.GameID = req.GameID
		rec.Status = req.Status
		rec.CurrentTurn = req.CurrentTurn
		rec.LastTurnProcessedAt = nulltime.FromTimePtr(req.LastTurnProcessedAt)
		rec.NextTurnDueAt = nulltime.FromTimePtr(req.NextTurnDueAt)
		rec.StartedAt = nulltime.FromTimePtr(req.StartedAt)
		rec.CompletedAt = nulltime.FromTimePtr(req.CompletedAt)
		rec.DeliveryPhysicalPost = req.DeliveryPhysicalPost
		rec.DeliveryPhysicalLocal = req.DeliveryPhysicalLocal
		rec.DeliveryEmail = req.DeliveryEmail
		l.Debug("after mapping: rec.DeliveryPhysicalPost=%v, rec.DeliveryPhysicalLocal=%v, rec.DeliveryEmail=%v",
			rec.DeliveryPhysicalPost, rec.DeliveryPhysicalLocal, rec.DeliveryEmail)
		if req.RequiredPlayerCount > 0 {
			rec.RequiredPlayerCount = req.RequiredPlayerCount
		}
		rec.IsClosedTesting = req.IsClosedTesting
	case server.HttpMethodPut, server.HttpMethodPatch:
		if req.Status != "" {
			rec.Status = req.Status
		}
		if req.CurrentTurn != 0 {
			rec.CurrentTurn = req.CurrentTurn
		}
		rec.LastTurnProcessedAt = nulltime.FromTimePtr(req.LastTurnProcessedAt)
		rec.NextTurnDueAt = nulltime.FromTimePtr(req.NextTurnDueAt)
		rec.StartedAt = nulltime.FromTimePtr(req.StartedAt)
		rec.CompletedAt = nulltime.FromTimePtr(req.CompletedAt)
		// Only update delivery flags if they're explicitly set (check if any are true)
		if req.DeliveryPhysicalPost || req.DeliveryPhysicalLocal || req.DeliveryEmail {
			rec.DeliveryPhysicalPost = req.DeliveryPhysicalPost
			rec.DeliveryPhysicalLocal = req.DeliveryPhysicalLocal
			rec.DeliveryEmail = req.DeliveryEmail
		}
		if req.RequiredPlayerCount > 0 {
			rec.RequiredPlayerCount = req.RequiredPlayerCount
		}
		rec.IsClosedTesting = req.IsClosedTesting
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameInstanceRecordToResponseData(l logger.Logger, rec *game_record.GameInstance, playerCount int) (*game_schema.GameInstanceResponseData, error) {
	l.Debug("mapping game_instance record to response data")
	return &game_schema.GameInstanceResponseData{
		ID:                                rec.ID,
		GameID:                            rec.GameID,
		Status:                            rec.Status,
		CurrentTurn:                       rec.CurrentTurn,
		LastTurnProcessedAt:               nulltime.ToTimePtr(rec.LastTurnProcessedAt),
		NextTurnDueAt:                     nulltime.ToTimePtr(rec.NextTurnDueAt),
		StartedAt:                         nulltime.ToTimePtr(rec.StartedAt),
		CompletedAt:                       nulltime.ToTimePtr(rec.CompletedAt),
		DeliveryPhysicalPost:              rec.DeliveryPhysicalPost,
		DeliveryPhysicalLocal:             rec.DeliveryPhysicalLocal,
		DeliveryEmail:                     rec.DeliveryEmail,
		RequiredPlayerCount:               rec.RequiredPlayerCount,
		PlayerCount:                       playerCount,
		IsClosedTesting:                   rec.IsClosedTesting,
		ClosedTestingJoinGameKey:          nullstring.ToStringPtr(rec.ClosedTestingJoinGameKey),
		ClosedTestingJoinGameKeyExpiresAt: nulltime.ToTimePtr(rec.ClosedTestingJoinGameKeyExpiresAt),
		CreatedAt:                         rec.CreatedAt,
		UpdatedAt:                         nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func GameInstanceRecordToResponse(l logger.Logger, rec *game_record.GameInstance, playerCount int) (*game_schema.GameInstanceResponse, error) {
	l.Debug("mapping game_instance record to response")
	data, err := GameInstanceRecordToResponseData(l, rec, playerCount)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameInstanceResponse{
		Data: data,
	}, nil
}

func GameInstanceRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameInstance, getPlayerCount func(string) (int, error)) (game_schema.GameInstanceCollectionResponse, error) {
	l.Debug("mapping game_instance records to collection response")
	data := []*game_schema.GameInstanceResponseData{}
	for _, rec := range recs {
		playerCount := 0
		if getPlayerCount != nil {
			count, err := getPlayerCount(rec.ID)
			if err != nil {
				l.Warn("failed to get player count for instance >%s< >%v<", rec.ID, err)
				// Continue with 0 if we can't get the count
			} else {
				playerCount = count
			}
		}
		d, err := GameInstanceRecordToResponseData(l, rec, playerCount)
		if err != nil {
			return game_schema.GameInstanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameInstanceCollectionResponse{
		Data: data,
	}, nil
}

func JoinGameLinkToResponse(l logger.Logger, joinGameURL, joinGameKey string) (*game_schema.JoinGameLinkResponse, error) {
	l.Debug("mapping join game link to response")
	return &game_schema.JoinGameLinkResponse{
		Data: &game_schema.JoinGameLinkResponseData{
			JoinGameURL: joinGameURL,
			JoinGameKey: joinGameKey,
		},
	}, nil
}

func InviteTesterRequestToEmail(l logger.Logger, r *http.Request) (string, error) {
	l.Debug("mapping invite tester request to email")

	var req game_schema.InviteTesterRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return "", err
	}

	if req.Email == "" {
		return "", fmt.Errorf("email is required")
	}

	return req.Email, nil
}

func InviteTesterToResponse(l logger.Logger, email string) (*game_schema.InviteTesterResponse, error) {
	l.Debug("mapping invite tester to response")
	return &game_schema.InviteTesterResponse{
		Data: &game_schema.InviteTesterResponseData{
			Message: "tester invitation queued",
			Email:   email,
		},
	}, nil
}
