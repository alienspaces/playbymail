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

func MechaGameSquadMechRequestToRecord(l logger.Logger, r *http.Request, rec *mecha_game_record.MechaGameSquadMech) (*mecha_game_record.MechaGameSquadMech, error) {
	l.Debug("mapping mecha_game_squad_mech request to record")

	var req mecha_game_schema.MechaGameSquadMechRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost, server.HttpMethodPut, server.HttpMethodPatch:
		rec.MechaGameChassisID = req.MechaGameChassisID
		rec.Callsign = req.Callsign
		rec.WeaponConfig = make([]mecha_game_record.WeaponConfigEntry, 0, len(req.WeaponConfig))
		for _, wc := range req.WeaponConfig {
			rec.WeaponConfig = append(rec.WeaponConfig, mecha_game_record.WeaponConfigEntry{
				WeaponID:     wc.WeaponID,
				SlotLocation: wc.SlotLocation,
			})
		}
		rec.EquipmentConfig = make([]mecha_game_record.EquipmentConfigEntry, 0, len(req.EquipmentConfig))
		for _, ec := range req.EquipmentConfig {
			rec.EquipmentConfig = append(rec.EquipmentConfig, mecha_game_record.EquipmentConfigEntry{
				EquipmentID:  ec.EquipmentID,
				SlotLocation: ec.SlotLocation,
			})
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func MechaGameSquadMechRecordToResponseData(l logger.Logger, rec *mecha_game_record.MechaGameSquadMech) (*mecha_game_schema.MechaGameSquadMechResponseData, error) {
	l.Debug("mapping mecha_game_squad_mech record to response data")

	weaponConfig := make([]mecha_game_schema.WeaponConfigEntry, 0, len(rec.WeaponConfig))
	for _, wc := range rec.WeaponConfig {
		weaponConfig = append(weaponConfig, mecha_game_schema.WeaponConfigEntry{
			WeaponID:     wc.WeaponID,
			SlotLocation: wc.SlotLocation,
		})
	}
	equipmentConfig := make([]mecha_game_schema.EquipmentConfigEntry, 0, len(rec.EquipmentConfig))
	for _, ec := range rec.EquipmentConfig {
		equipmentConfig = append(equipmentConfig, mecha_game_schema.EquipmentConfigEntry{
			EquipmentID:  ec.EquipmentID,
			SlotLocation: ec.SlotLocation,
		})
	}

	return &mecha_game_schema.MechaGameSquadMechResponseData{
		ID:             rec.ID,
		GameID:         rec.GameID,
		MechaGameSquadID:   rec.MechaGameSquadID,
		MechaGameChassisID: rec.MechaGameChassisID,
		Callsign:        rec.Callsign,
		WeaponConfig:    weaponConfig,
		EquipmentConfig: equipmentConfig,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:      nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func MechaGameSquadMechRecordToResponse(l logger.Logger, rec *mecha_game_record.MechaGameSquadMech) (*mecha_game_schema.MechaGameSquadMechResponse, error) {
	l.Debug("mapping mecha_game_squad_mech record to response")
	data, err := MechaGameSquadMechRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &mecha_game_schema.MechaGameSquadMechResponse{
		Data: data,
	}, nil
}

func MechaGameSquadMechRecordsToCollectionResponse(l logger.Logger, recs []*mecha_game_record.MechaGameSquadMech) (mecha_game_schema.MechaGameSquadMechCollectionResponse, error) {
	l.Debug("mapping mecha_game_squad_mech records to collection response")
	data := []*mecha_game_schema.MechaGameSquadMechResponseData{}
	for _, rec := range recs {
		d, err := MechaGameSquadMechRecordToResponseData(l, rec)
		if err != nil {
			return mecha_game_schema.MechaGameSquadMechCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return mecha_game_schema.MechaGameSquadMechCollectionResponse{
		Data: data,
	}, nil
}
