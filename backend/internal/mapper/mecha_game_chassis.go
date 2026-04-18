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

func MechaGameChassisRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameChassis) (*mecha_game_record.MechaGameChassis, error) {
	l.Debug("mapping mecha_game_chassis request to record")

	var req mecha_game_schema.MechaGameChassisRequest
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
			rec.ChassisClass = mecha_game_record.ChassisClassMedium
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
		rec.SmallSlots = req.SmallSlots
		rec.MediumSlots = req.MediumSlots
		rec.LargeSlots = req.LargeSlots
		// Apply per-class defaults only when the client sent no slot values at
		// all. A client that intentionally sets (say) large_slots=0 on a light
		// mech still gets that respected.
		if rec.SmallSlots == 0 && rec.MediumSlots == 0 && rec.LargeSlots == 0 {
			rec.SmallSlots, rec.MediumSlots, rec.LargeSlots = mecha_game_record.DefaultSlotsForChassisClass(rec.ChassisClass)
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaGameChassisRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameChassis) (*mecha_game_schema.MechaGameChassisResponseData, error) {
	l.Debug("mapping mecha_game_chassis record to response data")
	return &mecha_game_schema.MechaGameChassisResponseData{
		ID:              rec.ID,
		GameID:          rec.GameID,
		Name:            rec.Name,
		Description:     rec.Description,
		ChassisClass:    rec.ChassisClass,
		ArmorPoints:     rec.ArmorPoints,
		StructurePoints: rec.StructurePoints,
		HeatCapacity:    rec.HeatCapacity,
		Speed:           rec.Speed,
		SmallSlots:      rec.SmallSlots,
		MediumSlots:     rec.MediumSlots,
		LargeSlots:      rec.LargeSlots,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:       nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaGameChassisRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameChassis) (*mecha_game_schema.MechaGameChassisResponse, error) {
	l.Debug("mapping mecha_game_chassis record to response")
	data, err := MechaGameChassisRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameChassisResponse{
		Data: data,
	}, nil
}

func MechaGameChassisRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameChassis) (mecha_game_schema.MechaGameChassisCollectionResponse, error) {
	l.Debug("mapping mecha_game_chassis records to collection response")
	data := []*mecha_game_schema.MechaGameChassisResponseData{}
	for _, rec := range recs {
		d, err := MechaGameChassisRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameChassisCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameChassisCollectionResponse{
		Data: data,
	}, nil
}
