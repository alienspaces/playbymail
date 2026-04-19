package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/mecha_game_schema"
)

func MechaGameEquipmentRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameEquipment) (*mecha_game_record.MechaGameEquipment, error) {
	l.Debug("mapping mecha_game_equipment request to record")

	var req mecha_game_schema.MechaGameEquipmentRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.MountSize = req.MountSize
		if rec.MountSize == "" {
			rec.MountSize = mecha_game_record.EquipmentMountSizeMedium
		}
		rec.EffectKind = req.EffectKind
		rec.Magnitude = req.Magnitude
		rec.HeatCost = req.HeatCost
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaGameEquipmentRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameEquipment) (*mecha_game_schema.MechaGameEquipmentResponseData, error) {
	l.Debug("mapping mecha_game_equipment record to response data")
	return &mecha_game_schema.MechaGameEquipmentResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		Name:        rec.Name,
		Description: rec.Description,
		MountSize:   rec.MountSize,
		EffectKind:  rec.EffectKind,
		Magnitude:   rec.Magnitude,
		HeatCost:    rec.HeatCost,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaGameEquipmentRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameEquipment) (*mecha_game_schema.MechaGameEquipmentResponse, error) {
	l.Debug("mapping mecha_game_equipment record to response")
	data, err := MechaGameEquipmentRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameEquipmentResponse{
		Data: data,
	}, nil
}

func MechaGameEquipmentRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameEquipment) (mecha_game_schema.MechaGameEquipmentCollectionResponse, error) {
	l.Debug("mapping mecha_game_equipment records to collection response")
	data := []*mecha_game_schema.MechaGameEquipmentResponseData{}
	for _, rec := range recs {
		d, err := MechaGameEquipmentRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameEquipmentCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameEquipmentCollectionResponse{
		Data: data,
	}, nil
}
