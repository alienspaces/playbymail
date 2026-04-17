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

func MechaSquadMechRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaSquadMech) (*mecha_record.MechaSquadMech, error) {
	l.Debug("mapping mecha_squad_mech request to record")

	var req mecha_schema.MechaSquadMechRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.MechaChassisID = req.MechaChassisID
		rec.Callsign = req.Callsign
		rec.WeaponConfig = make([]mecha_record.WeaponConfigEntry, 0, len(req.WeaponConfig))
		for _, wc := range req.WeaponConfig {
			rec.WeaponConfig = append(rec.WeaponConfig, mecha_record.WeaponConfigEntry{
				WeaponID:     wc.WeaponID,
				SlotLocation: wc.SlotLocation,
			})
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaSquadMechRecordToResponseData(l logger.Logger, rec *mecha_record.MechaSquadMech) (*mecha_schema.MechaSquadMechResponseData, error) {
	l.Debug("mapping mecha_squad_mech record to response data")

	weaponConfig := make([]mecha_schema.WeaponConfigEntry, 0, len(rec.WeaponConfig))
	for _, wc := range rec.WeaponConfig {
		weaponConfig = append(weaponConfig, mecha_schema.WeaponConfigEntry{
			WeaponID:     wc.WeaponID,
			SlotLocation: wc.SlotLocation,
		})
	}

	return &mecha_schema.MechaSquadMechResponseData{
		ID:             rec.ID,
		GameID:         rec.GameID,
		MechaSquadID:   rec.MechaSquadID,
		MechaChassisID: rec.MechaChassisID,
		Callsign:       rec.Callsign,
		WeaponConfig:   weaponConfig,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:      nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaSquadMechRecordToResponse(l logger.Logger, rec *mecha_record.MechaSquadMech) (*mecha_schema.MechaSquadMechResponse, error) {
	l.Debug("mapping mecha_squad_mech record to response")
	data, err := MechaSquadMechRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaSquadMechResponse{
		Data: data,
	}, nil
}

func MechaSquadMechRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaSquadMech) (mecha_schema.MechaSquadMechCollectionResponse, error) {
	l.Debug("mapping mecha_squad_mech records to collection response")
	data := []*mecha_schema.MechaSquadMechResponseData{}
	for _, rec := range recs {
		d, err := MechaSquadMechRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaSquadMechCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaSquadMechCollectionResponse{
		Data: data,
	}, nil
}
