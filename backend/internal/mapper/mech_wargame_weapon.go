package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
	"gitlab.com/alienspaces/playbymail/schema/api/mech_wargame_schema"
)

func MechWargameWeaponRequestToRecord(l logger.Logger, r *http.Request, rec *mech_wargame_record.MechWargameWeapon) (*mech_wargame_record.MechWargameWeapon, error) {
	l.Debug("mapping mech_wargame_weapon request to record")

	var req mech_wargame_schema.MechWargameWeaponRequest
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
			rec.RangeBand = mech_wargame_record.WeaponRangeBandMedium
		}
		rec.MountSize = req.MountSize
		if rec.MountSize == "" {
			rec.MountSize = mech_wargame_record.WeaponMountSizeMedium
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechWargameWeaponRecordToResponseData(l logger.Logger, rec *mech_wargame_record.MechWargameWeapon) (*mech_wargame_schema.MechWargameWeaponResponseData, error) {
	l.Debug("mapping mech_wargame_weapon record to response data")
	return &mech_wargame_schema.MechWargameWeaponResponseData{
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

func MechWargameWeaponRecordToResponse(l logger.Logger, rec *mech_wargame_record.MechWargameWeapon) (*mech_wargame_schema.MechWargameWeaponResponse, error) {
	l.Debug("mapping mech_wargame_weapon record to response")
	data, err := MechWargameWeaponRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mech_wargame_schema.MechWargameWeaponResponse{
		Data: data,
	}, nil
}

func MechWargameWeaponRecordsToCollectionResponse(l logger.Logger, recs []*mech_wargame_record.MechWargameWeapon) (mech_wargame_schema.MechWargameWeaponCollectionResponse, error) {
	l.Debug("mapping mech_wargame_weapon records to collection response")
	data := []*mech_wargame_schema.MechWargameWeaponResponseData{}
	for _, rec := range recs {
		d, err := MechWargameWeaponRecordToResponseData(l, rec)
		if err != nil {
			return mech_wargame_schema.MechWargameWeaponCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mech_wargame_schema.MechWargameWeaponCollectionResponse{
		Data: data,
	}, nil
}
