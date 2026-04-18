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

func MechaGameWeaponRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameWeapon) (*mecha_game_record.MechaGameWeapon, error) {
	l.Debug("mapping mecha_game_weapon request to record")

	var req mecha_game_schema.MechaGameWeaponRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.Damage = req.Damage
		if rec.Damage == 0 {
			rec.Damage = 5
		}
		rec.HeatCost = req.HeatCost
		if rec.HeatCost == 0 {
			rec.HeatCost = 3
		}
		rec.RangeBand = req.RangeBand
		if rec.RangeBand == "" {
			rec.RangeBand = mecha_game_record.WeaponRangeBandMedium
		}
		rec.MountSize = req.MountSize
		if rec.MountSize == "" {
			rec.MountSize = mecha_game_record.WeaponMountSizeMedium
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaGameWeaponRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameWeapon) (*mecha_game_schema.MechaGameWeaponResponseData, error) {
	l.Debug("mapping mecha_game_weapon record to response data")
	return &mecha_game_schema.MechaGameWeaponResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		Name:        rec.Name,
		Description: rec.Description,
		Damage:      rec.Damage,
		HeatCost:    rec.HeatCost,
		RangeBand:   rec.RangeBand,
		MountSize:   rec.MountSize,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaGameWeaponRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameWeapon) (*mecha_game_schema.MechaGameWeaponResponse, error) {
	l.Debug("mapping mecha_game_weapon record to response")
	data, err := MechaGameWeaponRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameWeaponResponse{
		Data: data,
	}, nil
}

func MechaGameWeaponRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameWeapon) (mecha_game_schema.MechaGameWeaponCollectionResponse, error) {
	l.Debug("mapping mecha_game_weapon records to collection response")
	data := []*mecha_game_schema.MechaGameWeaponResponseData{}
	for _, rec := range recs {
		d, err := MechaGameWeaponRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameWeaponCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameWeaponCollectionResponse{
		Data: data,
	}, nil
}
