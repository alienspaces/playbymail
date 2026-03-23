package domain

import (
	"database/sql"
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// CreateDefaultMechaLanceForPlayer clones the game's player starter lance for a new
// player. It is idempotent: if the player already has a lance the existing record is
// returned. The starter lance must be configured by the game designer; there is no
// auto-pick fallback.
func (m *Domain) CreateDefaultMechaLanceForPlayer(gameID, accountID, accountUserID, commanderName, playerName string) (*mecha_record.MechaLance, error) {
	l := m.Logger("CreateDefaultMechaLanceForPlayer")

	// Idempotent: return an existing lance if the player already has one.
	existingRecs, err := m.GetManyMechaLanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceGameID, Val: gameID},
			{Col: mecha_record.FieldMechaLanceAccountUserID, Val: accountUserID},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(existingRecs) > 0 {
		l.Info("player already has a lance >%s< for game >%s< - returning existing", existingRecs[0].ID, gameID)
		return existingRecs[0], nil
	}

	// Locate the designer-configured starter lance for this game.
	starterRecs, err := m.GetManyMechaLanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceGameID, Val: gameID},
			{Col: mecha_record.FieldMechaLanceIsPlayerStarter, Val: true},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(starterRecs) == 0 {
		return nil, coreerror.NewInvalidDataError("game >%s< has no player starter lance configured — the game designer must create one before players can join", gameID)
	}
	starterLance := starterRecs[0]

	// Build the player's lance name from the provided names, falling back to the starter's name.
	lanceBaseName := commanderName
	if lanceBaseName == "" {
		lanceBaseName = playerName
	}
	if lanceBaseName == "" {
		lanceBaseName = "Player"
	}

	playerLance, err := m.CreateMechaLanceRec(&mecha_record.MechaLance{
		GameID:        gameID,
		AccountID:     sql.NullString{Valid: true, String: accountID},
		AccountUserID: sql.NullString{Valid: true, String: accountUserID},
		Name:          lanceBaseName + "'s Lance",
		Description:   starterLance.Description,
	})
	if err != nil {
		l.Warn("failed to create player mecha lance >%v<", err)
		return nil, err
	}

	// Clone the starter's mech roster into the player's lance.
	starterMechs, err := m.GetManyMechaLanceMechRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceMechMechaLanceID, Val: starterLance.ID},
		},
	})
	if err != nil {
		l.Warn("failed to get starter lance mechs >%v<", err)
		return nil, err
	}

	for i, mech := range starterMechs {
		callsign := fmt.Sprintf("%s-%d", lanceBaseName, i+1)
		_, err := m.CreateMechaLanceMechRec(&mecha_record.MechaLanceMech{
			GameID:          gameID,
			MechaLanceID:    playerLance.ID,
			MechaChassisID:  mech.MechaChassisID,
			Callsign:        callsign,
			WeaponConfig:    mech.WeaponConfig,
			WeaponConfigJSON: mech.WeaponConfigJSON,
		})
		if err != nil {
			l.Warn("failed to create player lance mech >%v<", err)
			return nil, err
		}
	}

	l.Info("cloned starter lance >%s< into player lance >%s< with >%d< mechs for player >%s<", starterLance.ID, playerLance.ID, len(starterMechs), accountUserID)

	return playerLance, nil
}

// MechaInstanceData holds all instance records created when a mecha instance starts.
type MechaInstanceData struct {
	SectorInstances []*mecha_record.MechaSectorInstance
	LanceInstances  []*mecha_record.MechaLanceInstance
	MechInstances   []*mecha_record.MechaMechInstance
}

// PopulateMechaGameInstanceData creates all runtime records (sector instances, lance instances,
// mech instances) for a mecha game instance from its design definitions and player subscriptions.
func (m *Domain) PopulateMechaGameInstanceData(instanceID string) (*MechaInstanceData, error) {
	l := m.Logger("PopulateMechaGameInstanceData")

	l.Info("populating mecha instance data for instance >%s<", instanceID)

	instanceRec, err := m.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return nil, err
	}

	gameID := instanceRec.GameID
	out := &MechaInstanceData{}

	// 1. Create sector instances and build sectorID -> sectorInstanceID map
	sectorRecs, err := m.GetManyMechaSectorRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSectorGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get mecha sectors >%v<", err)
		return nil, err
	}

	sectorIDToInstanceID := make(map[string]string, len(sectorRecs))
	var startingSectorInstanceID string

	for _, sector := range sectorRecs {
		sectorInst, err := m.CreateMechaSectorInstanceRec(&mecha_record.MechaSectorInstance{
			GameID:              gameID,
			GameInstanceID:      instanceID,
			MechaSectorID: sector.ID,
		})
		if err != nil {
			l.Warn("failed to create sector instance for sector >%s< >%v<", sector.ID, err)
			return nil, err
		}
		sectorIDToInstanceID[sector.ID] = sectorInst.ID
		out.SectorInstances = append(out.SectorInstances, sectorInst)
		if sector.IsStartingSector && startingSectorInstanceID == "" {
			startingSectorInstanceID = sectorInst.ID
		}
	}

	if startingSectorInstanceID == "" && len(sectorRecs) > 0 {
		return nil, fmt.Errorf("no starting sector found for game >%s<: at least one sector must have is_starting_sector = true", gameID)
	}

	// 2. Create lance instances for all subscribed players
	subscriptionInstances, err := m.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get subscription instances >%v<", err)
		return nil, err
	}

	createdLanceIDs := make(map[string]bool)
	for _, subInst := range subscriptionInstances {
		subRecs, err := m.GetManyGameSubscriptionRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: game_record.FieldGameSubscriptionID, Val: subInst.GameSubscriptionID},
			},
			Limit: 1,
		})
		if err != nil || len(subRecs) == 0 {
			l.Warn("failed to get game subscription >%s< >%v<", subInst.GameSubscriptionID, err)
			continue
		}
		sub := subRecs[0]

		if sub.SubscriptionType != game_record.GameSubscriptionTypePlayer {
			continue
		}

		lanceRecs, err := m.GetManyMechaLanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaLanceGameID, Val: gameID},
				{Col: mecha_record.FieldMechaLanceAccountUserID, Val: sub.AccountUserID},
			},
			Limit: 1,
		})
		if err != nil || len(lanceRecs) == 0 {
			l.Warn("no mecha lance found for account_user_id >%s< game >%s< - skipping", sub.AccountUserID, gameID)
			continue
		}

		lanceID := lanceRecs[0].ID
		if createdLanceIDs[lanceID] {
			l.Info("lance instance already created for lance >%s< - skipping duplicate subscription", lanceID)
			continue
		}

		lanceInst, err := m.CreateMechaLanceInstanceRec(&mecha_record.MechaLanceInstance{
			GameID:         gameID,
			GameInstanceID: instanceID,
			MechaLanceID:   lanceID,
			GameSubscriptionInstanceID: sql.NullString{
				String: subInst.ID,
				Valid:  true,
			},
		})
		if err != nil {
			l.Warn("failed to create lance instance for lance >%s< >%v<", lanceID, err)
			return nil, err
		}
		createdLanceIDs[lanceID] = true
		out.LanceInstances = append(out.LanceInstances, lanceInst)

		// 3. Create mech instances for each mech in the lance
		lanceMechRecs, err := m.GetManyMechaLanceMechRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaLanceMechMechaLanceID, Val: lanceID},
			},
		})
		if err != nil {
			l.Warn("failed to get lance mechs for lance >%s< >%v<", lanceID, err)
			return nil, err
		}

		if startingSectorInstanceID == "" {
			return nil, fmt.Errorf("no starting sector instance found for game >%s<: cannot create mech instances", gameID)
		}

		for _, lanceMech := range lanceMechRecs {
			chassisRec, err := m.GetMechaChassisRec(lanceMech.MechaChassisID, nil)
			if err != nil {
				l.Warn("failed to get chassis >%s< for lance mech >%s< >%v<", lanceMech.MechaChassisID, lanceMech.ID, err)
				return nil, err
			}

			mechInst, err := m.CreateMechaMechInstanceRec(&mecha_record.MechaMechInstance{
				GameID:                gameID,
				GameInstanceID:        instanceID,
				MechaLanceInstanceID:  lanceInst.ID,
				MechaSectorInstanceID: startingSectorInstanceID,
				MechaChassisID:        lanceMech.MechaChassisID,
				Callsign:              lanceMech.Callsign,
				CurrentArmor:          chassisRec.ArmorPoints,
				CurrentStructure:      chassisRec.StructurePoints,
				CurrentHeat:           0,
				PilotSkill:            4,
				Status:                mecha_record.MechInstanceStatusOperational,
				WeaponConfig:          lanceMech.WeaponConfig,
				WeaponConfigJSON:      lanceMech.WeaponConfigJSON,
			})
			if err != nil {
				l.Warn("failed to create mech instance for lance mech >%s< >%v<", lanceMech.ID, err)
				return nil, err
			}
			out.MechInstances = append(out.MechInstances, mechInst)
		}

		l.Info("created lance instance >%s< with %d mech instances for lance >%s<", lanceInst.ID, len(lanceMechRecs), lanceID)
	}

	// 4. Create lance instances for computer-opponent-owned lances
	computerLanceRecs, err := m.GetManyMechaLanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get mecha lances for computer opponent check >%v<", err)
		return nil, err
	}

	for _, lanceRec := range computerLanceRecs {
		if !lanceRec.MechaComputerOpponentID.Valid {
			continue
		}
		if createdLanceIDs[lanceRec.ID] {
			continue
		}

		lanceInst, err := m.CreateMechaLanceInstanceRec(&mecha_record.MechaLanceInstance{
			GameID:                     gameID,
			GameInstanceID:             instanceID,
			MechaLanceID:               lanceRec.ID,
			GameSubscriptionInstanceID: sql.NullString{Valid: false},
		})
		if err != nil {
			l.Warn("failed to create lance instance for computer opponent lance >%s< >%v<", lanceRec.ID, err)
			return nil, err
		}
		createdLanceIDs[lanceRec.ID] = true
		out.LanceInstances = append(out.LanceInstances, lanceInst)

		lanceMechRecs, err := m.GetManyMechaLanceMechRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaLanceMechMechaLanceID, Val: lanceRec.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get lance mechs for computer opponent lance >%s< >%v<", lanceRec.ID, err)
			return nil, err
		}

		if startingSectorInstanceID == "" {
			return nil, fmt.Errorf("no starting sector instance found for game >%s<: cannot create mech instances for computer opponent", gameID)
		}

		for _, lanceMech := range lanceMechRecs {
			chassisRec, err := m.GetMechaChassisRec(lanceMech.MechaChassisID, nil)
			if err != nil {
				l.Warn("failed to get chassis >%s< for computer opponent lance mech >%s< >%v<", lanceMech.MechaChassisID, lanceMech.ID, err)
				return nil, err
			}

			mechInst, err := m.CreateMechaMechInstanceRec(&mecha_record.MechaMechInstance{
				GameID:                gameID,
				GameInstanceID:        instanceID,
				MechaLanceInstanceID:  lanceInst.ID,
				MechaSectorInstanceID: startingSectorInstanceID,
				MechaChassisID:        lanceMech.MechaChassisID,
				Callsign:              lanceMech.Callsign,
				CurrentArmor:          chassisRec.ArmorPoints,
				CurrentStructure:      chassisRec.StructurePoints,
				CurrentHeat:           0,
				PilotSkill:            4,
				Status:                mecha_record.MechInstanceStatusOperational,
				WeaponConfig:          lanceMech.WeaponConfig,
				WeaponConfigJSON:      lanceMech.WeaponConfigJSON,
			})
			if err != nil {
				l.Warn("failed to create mech instance for computer opponent lance mech >%s< >%v<", lanceMech.ID, err)
				return nil, err
			}
			out.MechInstances = append(out.MechInstances, mechInst)
		}

		l.Info("created computer opponent lance instance >%s< with %d mech instances for lance >%s<", lanceInst.ID, len(lanceMechRecs), lanceRec.ID)
	}

	l.Info("populated mecha instance data: sectors=%d lances=%d mechs=%d",
		len(out.SectorInstances), len(out.LanceInstances), len(out.MechInstances))

	return out, nil
}

// deleteMechaInstanceData removes mecha instance records (soft delete) for the given instanceID.
func (m *Domain) deleteMechaInstanceData(instanceID string) error {
	l := m.Logger("deleteMechaInstanceData")

	// Delete mecha_turn_sheet records linked to lance instances for this game instance
	lanceInstances, err := m.GetManyMechaLanceInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get lance instances >%v<", err)
		return databaseError(err)
	}

	for _, lanceInst := range lanceInstances {
		turnSheets, err := m.GetManyMechaTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaTurnSheetMechaLanceInstanceID, Val: lanceInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for lance instance >%s< >%v<", lanceInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.MechaTurnSheetRepository().DeleteOne(ts.ID); err != nil {
				l.Warn("failed to delete mecha turn sheet >%s< >%v<", ts.ID, err)
				return databaseError(err)
			}
		}
	}

	// Delete mech instances
	mechInstances, err := m.GetManyMechaMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaMechInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return databaseError(err)
	}
	for _, mechInst := range mechInstances {
		if err := m.MechaMechInstanceRepository().DeleteOne(mechInst.ID); err != nil {
			l.Warn("failed to delete mech instance >%s< >%v<", mechInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete lance instances
	for _, lanceInst := range lanceInstances {
		if err := m.MechaLanceInstanceRepository().DeleteOne(lanceInst.ID); err != nil {
			l.Warn("failed to delete lance instance >%s< >%v<", lanceInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete sector instances
	sectorInstances, err := m.GetManyMechaSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSectorInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get sector instances >%v<", err)
		return databaseError(err)
	}
	for _, sectorInst := range sectorInstances {
		if err := m.MechaSectorInstanceRepository().DeleteOne(sectorInst.ID); err != nil {
			l.Warn("failed to delete sector instance >%s< >%v<", sectorInst.ID, err)
			return databaseError(err)
		}
	}

	return nil
}

// removeMechaInstanceData permanently removes mecha instance records for the given instanceID.
func (m *Domain) removeMechaInstanceData(instanceID string) error {
	l := m.Logger("removeMechaInstanceData")

	// Remove mecha_turn_sheet records linked to lance instances
	lanceInstances, err := m.GetManyMechaLanceInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaLanceInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get lance instances >%v<", err)
		return databaseError(err)
	}

	for _, lanceInst := range lanceInstances {
		turnSheets, err := m.GetManyMechaTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaTurnSheetMechaLanceInstanceID, Val: lanceInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for lance instance >%s< >%v<", lanceInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.RemoveMechaTurnSheetRec(ts.ID); err != nil {
				l.Warn("failed to remove mecha turn sheet >%s< >%v<", ts.ID, err)
				return err
			}
		}
	}

	// Remove mech instances
	mechInstances, err := m.GetManyMechaMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaMechInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return databaseError(err)
	}
	for _, mechInst := range mechInstances {
		if err := m.RemoveMechaMechInstanceRec(mechInst.ID); err != nil {
			l.Warn("failed to remove mech instance >%s< >%v<", mechInst.ID, err)
			return err
		}
	}

	// Remove lance instances
	for _, lanceInst := range lanceInstances {
		if err := m.RemoveMechaLanceInstanceRec(lanceInst.ID); err != nil {
			l.Warn("failed to remove lance instance >%s< >%v<", lanceInst.ID, err)
			return err
		}
	}

	// Remove sector instances
	sectorInstances, err := m.GetManyMechaSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSectorInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get sector instances >%v<", err)
		return databaseError(err)
	}
	for _, sectorInst := range sectorInstances {
		if err := m.RemoveMechaSectorInstanceRec(sectorInst.ID); err != nil {
			l.Warn("failed to remove sector instance >%s< >%v<", sectorInst.ID, err)
			return err
		}
	}

	return nil
}
