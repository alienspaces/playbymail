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

func MechaChassisRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaChassis) (*mecha_record.MechaChassis, error) {
	l.Debug("mapping mecha_chassis request to record")

	var req mecha_schema.MechaChassisRequest
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
			rec.ChassisClass = mecha_record.ChassisClassMedium
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

func MechaChassisRecordToResponseData(l logger.Logger, rec *mecha_record.MechaChassis) (*mecha_schema.MechaChassisResponseData, error) {
	l.Debug("mapping mecha_chassis record to response data")
	return &mecha_schema.MechaChassisResponseData{
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

func MechaChassisRecordToResponse(l logger.Logger, rec *mecha_record.MechaChassis) (*mecha_schema.MechaChassisResponse, error) {
	l.Debug("mapping mecha_chassis record to response")
	data, err := MechaChassisRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaChassisResponse{
		Data: data,
	}, nil
}

func MechaChassisRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaChassis) (mecha_schema.MechaChassisCollectionResponse, error) {
	l.Debug("mapping mecha_chassis records to collection response")
	data := []*mecha_schema.MechaChassisResponseData{}
	for _, rec := range recs {
		d, err := MechaChassisRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaChassisCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaChassisCollectionResponse{
		Data: data,
	}, nil
}
