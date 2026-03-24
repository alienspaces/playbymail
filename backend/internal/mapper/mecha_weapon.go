package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
	"gitlab.com/alienspaces/playbymail/schema/api/mecha_schema"
)

func MechaWeaponRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaWeapon) (*mecha_record.MechaWeapon, error) {
	l.Debug("mapping mecha_weapon request to record")

	var req mecha_schema.MechaWeaponRequest
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
			rec.RangeBand = mecha_record.WeaponRangeBandMedium
		}
		rec.MountSize = req.MountSize
		if rec.MountSize == "" {
			rec.MountSize = mecha_record.WeaponMountSizeMedium
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaWeaponRecordToResponseData(l logger.Logger, rec *mecha_record.MechaWeapon) (*mecha_schema.MechaWeaponResponseData, error) {
	l.Debug("mapping mecha_weapon record to response data")
	return &mecha_schema.MechaWeaponResponseData{
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

func MechaWeaponRecordToResponse(l logger.Logger, rec *mecha_record.MechaWeapon) (*mecha_schema.MechaWeaponResponse, error) {
	l.Debug("mapping mecha_weapon record to response")
	data, err := MechaWeaponRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaWeaponResponse{
		Data: data,
	}, nil
}

func MechaWeaponRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaWeapon) (mecha_schema.MechaWeaponCollectionResponse, error) {
	l.Debug("mapping mecha_weapon records to collection response")
	data := []*mecha_schema.MechaWeaponResponseData{}
	for _, rec := range recs {
		d, err := MechaWeaponRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaWeaponCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaWeaponCollectionResponse{
		Data: data,
	}, nil
}
