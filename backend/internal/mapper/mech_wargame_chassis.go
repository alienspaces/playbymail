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

func MechWargameChassisRequestToRecord(l logger.Logger, r *http.Request, rec *mech_wargame_record.MechWargameChassis) (*mech_wargame_record.MechWargameChassis, error) {
	l.Debug("mapping mech_wargame_chassis request to record")

	var req mech_wargame_schema.MechWargameChassisRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.Description = req.Description
		rec.ChassisClass = req.ChassisClass
		if rec.ChassisClass == "" {
			rec.ChassisClass = mech_wargame_record.ChassisClassMedium
		}
		rec.ArmorPoints = req.ArmorPoints
		if rec.ArmorPoints == 0 {
			rec.ArmorPoints = 100
		}
		rec.StructurePoints = req.StructurePoints
		if rec.StructurePoints == 0 {
			rec.StructurePoints = 50
		}
		rec.HeatCapacity = req.HeatCapacity
		if rec.HeatCapacity == 0 {
			rec.HeatCapacity = 30
		}
		rec.Speed = req.Speed
		if rec.Speed == 0 {
			rec.Speed = 3
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechWargameChassisRecordToResponseData(l logger.Logger, rec *mech_wargame_record.MechWargameChassis) (*mech_wargame_schema.MechWargameChassisResponseData, error) {
	l.Debug("mapping mech_wargame_chassis record to response data")
	return &mech_wargame_schema.MechWargameChassisResponseData{
		ID:              rec.ID,
		GameID:          rec.GameID,
		Name:            rec.Name,
		Description:     rec.Description,
		ChassisClass:    rec.ChassisClass,
		ArmorPoints:     rec.ArmorPoints,
		StructurePoints: rec.StructurePoints,
		HeatCapacity:    rec.HeatCapacity,
		Speed:           rec.Speed,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:       nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechWargameChassisRecordToResponse(l logger.Logger, rec *mech_wargame_record.MechWargameChassis) (*mech_wargame_schema.MechWargameChassisResponse, error) {
	l.Debug("mapping mech_wargame_chassis record to response")
	data, err := MechWargameChassisRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mech_wargame_schema.MechWargameChassisResponse{
		Data: data,
	}, nil
}

func MechWargameChassisRecordsToCollectionResponse(l logger.Logger, recs []*mech_wargame_record.MechWargameChassis) (mech_wargame_schema.MechWargameChassisCollectionResponse, error) {
	l.Debug("mapping mech_wargame_chassis records to collection response")
	data := []*mech_wargame_schema.MechWargameChassisResponseData{}
	for _, rec := range recs {
		d, err := MechWargameChassisRecordToResponseData(l, rec)
		if err != nil {
			return mech_wargame_schema.MechWargameChassisCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mech_wargame_schema.MechWargameChassisCollectionResponse{
		Data: data,
	}, nil
}
