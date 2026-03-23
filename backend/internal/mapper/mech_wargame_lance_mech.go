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

func MechWargameLanceMechRequestToRecord(l logger.Logger, r *http.Request, rec *mech_wargame_record.MechWargameLanceMech) (*mech_wargame_record.MechWargameLanceMech, error) {
	l.Debug("mapping mech_wargame_lance_mech request to record")

	var req mech_wargame_schema.MechWargameLanceMechRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.MechWargameChassisID = req.MechWargameChassisID
		rec.Callsign = req.Callsign
		rec.WeaponConfig = make([]mech_wargame_record.WeaponConfigEntry, 0, len(req.WeaponConfig))
		for _, wc := range req.WeaponConfig {
			rec.WeaponConfig = append(rec.WeaponConfig, mech_wargame_record.WeaponConfigEntry{
				WeaponID:     wc.WeaponID,
				SlotLocation: wc.SlotLocation,
			})
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechWargameLanceMechRecordToResponseData(l logger.Logger, rec *mech_wargame_record.MechWargameLanceMech) (*mech_wargame_schema.MechWargameLanceMechResponseData, error) {
	l.Debug("mapping mech_wargame_lance_mech record to response data")

	weaponConfig := make([]mech_wargame_schema.WeaponConfigEntry, 0, len(rec.WeaponConfig))
	for _, wc := range rec.WeaponConfig {
		weaponConfig = append(weaponConfig, mech_wargame_schema.WeaponConfigEntry{
			WeaponID:     wc.WeaponID,
			SlotLocation: wc.SlotLocation,
		})
	}

	return &mech_wargame_schema.MechWargameLanceMechResponseData{
		ID:                   rec.ID,
		GameID:               rec.GameID,
		MechWargameLanceID:   rec.MechWargameLanceID,
		MechWargameChassisID: rec.MechWargameChassisID,
		Callsign:             rec.Callsign,
		WeaponConfig:         weaponConfig,
		CreatedAt:            rec.CreatedAt,
		UpdatedAt:            nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:            nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechWargameLanceMechRecordToResponse(l logger.Logger, rec *mech_wargame_record.MechWargameLanceMech) (*mech_wargame_schema.MechWargameLanceMechResponse, error) {
	l.Debug("mapping mech_wargame_lance_mech record to response")
	data, err := MechWargameLanceMechRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mech_wargame_schema.MechWargameLanceMechResponse{
		Data: data,
	}, nil
}

func MechWargameLanceMechRecordsToCollectionResponse(l logger.Logger, recs []*mech_wargame_record.MechWargameLanceMech) (mech_wargame_schema.MechWargameLanceMechCollectionResponse, error) {
	l.Debug("mapping mech_wargame_lance_mech records to collection response")
	data := []*mech_wargame_schema.MechWargameLanceMechResponseData{}
	for _, rec := range recs {
		d, err := MechWargameLanceMechRecordToResponseData(l, rec)
		if err != nil {
			return mech_wargame_schema.MechWargameLanceMechCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mech_wargame_schema.MechWargameLanceMechCollectionResponse{
		Data: data,
	}, nil
}
