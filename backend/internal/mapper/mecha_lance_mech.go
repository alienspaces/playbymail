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

func MechaLanceMechRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_record.MechaLanceMech) (*mecha_record.MechaLanceMech, error) {
	l.Debug("mapping mecha_lance_mech request to record")

	var req mecha_schema.MechaLanceMechRequest
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

func MechaLanceMechRecordToResponseData(l logger.Logger, rec *mecha_record.MechaLanceMech) (*mecha_schema.MechaLanceMechResponseData, error) {
	l.Debug("mapping mecha_lance_mech record to response data")

	weaponConfig := make([]mecha_schema.WeaponConfigEntry, 0, len(rec.WeaponConfig))
	for _, wc := range rec.WeaponConfig {
		weaponConfig = append(weaponConfig, mecha_schema.WeaponConfigEntry{
			WeaponID:     wc.WeaponID,
			SlotLocation: wc.SlotLocation,
		})
	}

	return &mecha_schema.MechaLanceMechResponseData{
		ID:                   rec.ID,
		GameID:               rec.GameID,
		MechaLanceID:   rec.MechaLanceID,
		MechaChassisID: rec.MechaChassisID,
		Callsign:             rec.Callsign,
		WeaponConfig:         weaponConfig,
		CreatedAt:            rec.CreatedAt,
		UpdatedAt:            nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:            nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaLanceMechRecordToResponse(l logger.Logger, rec *mecha_record.MechaLanceMech) (*mecha_schema.MechaLanceMechResponse, error) {
	l.Debug("mapping mecha_lance_mech record to response")
	data, err := MechaLanceMechRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_schema.MechaLanceMechResponse{
		Data: data,
	}, nil
}

func MechaLanceMechRecordsToCollectionResponse(l logger.Logger, recs []*mecha_record.MechaLanceMech) (mecha_schema.MechaLanceMechCollectionResponse, error) {
	l.Debug("mapping mecha_lance_mech records to collection response")
	data := []*mecha_schema.MechaLanceMechResponseData{}
	for _, rec := range recs {
		d, err := MechaLanceMechRecordToResponseData(l, rec)
		if err != nil {
			return mecha_schema.MechaLanceMechCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_schema.MechaLanceMechCollectionResponse{
		Data: data,
	}, nil
}
