package domain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

// mechInitialLoadoutState captures the runtime loadout data derived at clone
// time from a squad mech and its referenced weapons + equipment. It is used
// exclusively by the clone site in PopulateMechaGameInstanceData so the
// player and computer-opponent paths stay in lock-step.
type mechInitialLoadoutState struct {
	WeaponConfig     []mecha_game_record.WeaponConfigEntry
	WeaponConfigJSON []byte
	EquipmentConfig  []mecha_game_record.EquipmentConfigEntry
	EquipmentJSON    []byte
	ArmorBonus       int
	AmmoRemaining    int
}

// resolveMechInitialLoadout hydrates a squad mech's persisted JSON loadouts,
// looks up the referenced weapon / equipment records, and returns the initial
// armor bonus (from armor_upgrade magnitudes) and ammo pool (Σ weapon.ammo_capacity
// + Σ ammo_bin magnitudes) to seed the new mech instance with.
func (m *Domain) resolveMechInitialLoadout(squadMech *mecha_game_record.MechaGameSquadMech) (*mechInitialLoadoutState, error) {
	state := &mechInitialLoadoutState{
		WeaponConfigJSON: squadMech.WeaponConfigJSON,
		EquipmentJSON:    squadMech.EquipmentConfigJSON,
	}

	if len(squadMech.WeaponConfigJSON) > 0 {
		if err := json.Unmarshal(squadMech.WeaponConfigJSON, &state.WeaponConfig); err != nil {
			return nil, fmt.Errorf("resolveMechInitialLoadout: unmarshal weapon_config for squad mech >%s<: %w", squadMech.ID, err)
		}
	}
	if len(squadMech.EquipmentConfigJSON) > 0 {
		if err := json.Unmarshal(squadMech.EquipmentConfigJSON, &state.EquipmentConfig); err != nil {
			return nil, fmt.Errorf("resolveMechInitialLoadout: unmarshal equipment_config for squad mech >%s<: %w", squadMech.ID, err)
		}
	}

	for _, entry := range state.WeaponConfig {
		if entry.WeaponID == "" {
			continue
		}
		w, err := m.GetMechaGameWeaponRec(entry.WeaponID, nil)
		if err != nil {
			return nil, err
		}
		state.AmmoRemaining += w.AmmoCapacity
	}

	for _, entry := range state.EquipmentConfig {
		if entry.EquipmentID == "" {
			continue
		}
		eq, err := m.GetMechaGameEquipmentRec(entry.EquipmentID, nil)
		if err != nil {
			return nil, err
		}
		switch eq.EffectKind {
		case mecha_game_record.EquipmentEffectKindArmorUpgrade:
			state.ArmorBonus += eq.Magnitude
		case mecha_game_record.EquipmentEffectKindAmmoBin:
			state.AmmoRemaining += eq.Magnitude
		}
	}

	return state, nil
}

// MechaGameInstanceData holds all instance records created when a mecha instance starts.
type MechaGameInstanceData struct {
	SectorInstances []*mecha_game_record.MechaGameSectorInstance
	SquadInstances  []*mecha_game_record.MechaGameSquadInstance
	MechInstances   []*mecha_game_record.MechaGameMechInstance
}

// PopulateMechaGameInstanceData creates all runtime records (sector instances, squad instances,
// mech instances) for a mecha game instance from its design definitions and player subscriptions.
//
// Player squads: each subscribed player gets a squad instance cloned from the starter template.
// Opponent squads: each computer opponent is randomly assigned an opponent squad template.
func (m *Domain) PopulateMechaGameInstanceData(instanceID string) (*MechaGameInstanceData, error) {
	l := m.Logger("PopulateMechaGameInstanceData")

	l.Info("populating mecha instance data for instance >%s<", instanceID)

	instanceRec, err := m.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return nil, err
	}

	gameID := instanceRec.GameID
	out := &MechaGameInstanceData{}

	// 1. Create sector instances and build sectorID -> sectorInstanceID map
	sectorRecs, err := m.GetManyMechaGameSectorRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSectorGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get mecha sectors >%v<", err)
		return nil, err
	}

	var startingSectorInstanceID string

	for _, sector := range sectorRecs {
		sectorInst, err := m.CreateMechaGameSectorInstanceRec(&mecha_game_record.MechaGameSectorInstance{
			GameID:         gameID,
			GameInstanceID: instanceID,
			MechaGameSectorID:  sector.ID,
		})
		if err != nil {
			l.Warn("failed to create sector instance for sector >%s< >%v<", sector.ID, err)
			return nil, err
		}
		out.SectorInstances = append(out.SectorInstances, sectorInst)
		if sector.IsStartingSector && startingSectorInstanceID == "" {
			startingSectorInstanceID = sectorInst.ID
		}
	}

	if startingSectorInstanceID == "" && len(sectorRecs) > 0 {
		return nil, fmt.Errorf("no starting sector found for game >%s<: at least one sector must have is_starting_sector = true", gameID)
	}

	// 2. Locate the starter squad template
	starterRecs, err := m.GetManyMechaGameSquadRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSquadGameID, Val: gameID},
			{Col: mecha_game_record.FieldMechaGameSquadSquadType, Val: mecha_game_record.SquadTypeStarter},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get starter squad for game >%s< >%v<", gameID, err)
		return nil, err
	}
	if len(starterRecs) == 0 {
		return nil, coreerror.NewInvalidDataError("game >%s< has no player starter squad configured", gameID)
	}
	starterSquad := starterRecs[0]

	starterMechs, err := m.GetManyMechaGameSquadMechRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSquadMechMechaGameSquadID, Val: starterSquad.ID},
		},
	})
	if err != nil {
		l.Warn("failed to get mechs for starter squad >%s< >%v<", starterSquad.ID, err)
		return nil, err
	}

	// 3. Create a squad instance for each subscribed player, cloning mechs from the starter template
	subscriptionInstances, err := m.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get subscription instances >%v<", err)
		return nil, err
	}

	playerNumber := 0
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

		playerNumber++

		squadInst, err := m.CreateMechaGameSquadInstanceRec(&mecha_game_record.MechaGameSquadInstance{
			GameID:                     gameID,
			GameInstanceID:             instanceID,
			MechaGameSquadID:               starterSquad.ID,
			GameSubscriptionInstanceID: sql.NullString{String: subInst.ID, Valid: true},
		})
		if err != nil {
			l.Warn("failed to create player squad instance for subscription >%s< >%v<", subInst.ID, err)
			return nil, err
		}
		out.SquadInstances = append(out.SquadInstances, squadInst)

		if startingSectorInstanceID == "" {
			return nil, fmt.Errorf("no starting sector instance found for game >%s<: cannot create mech instances", gameID)
		}

		for i, squadMech := range starterMechs {
			chassisRec, err := m.GetMechaGameChassisRec(squadMech.MechaGameChassisID, nil)
			if err != nil {
				l.Warn("failed to get chassis >%s< for starter mech >%s< >%v<", squadMech.MechaGameChassisID, squadMech.ID, err)
				return nil, err
			}

			loadout, err := m.resolveMechInitialLoadout(squadMech)
			if err != nil {
				l.Warn("failed to resolve initial loadout for starter mech >%s< >%v<", squadMech.ID, err)
				return nil, err
			}

			callsign := fmt.Sprintf("P%d-%d", playerNumber, i+1)
			mechInst, err := m.CreateMechaGameMechInstanceRec(&mecha_game_record.MechaGameMechInstance{
				GameID:                    gameID,
				GameInstanceID:            instanceID,
				MechaGameSquadInstanceID:  squadInst.ID,
				MechaGameSectorInstanceID: startingSectorInstanceID,
				MechaGameChassisID:        squadMech.MechaGameChassisID,
				Callsign:                  callsign,
				CurrentArmor:              chassisRec.ArmorPoints + loadout.ArmorBonus,
				CurrentStructure:          chassisRec.StructurePoints,
				CurrentHeat:               0,
				PilotSkill:                0,
				Status:                    mecha_game_record.MechInstanceStatusOperational,
				WeaponConfig:              loadout.WeaponConfig,
				WeaponConfigJSON:          loadout.WeaponConfigJSON,
				EquipmentConfig:           loadout.EquipmentConfig,
				EquipmentConfigJSON:       loadout.EquipmentJSON,
				AmmoRemaining:             loadout.AmmoRemaining,
			})
			if err != nil {
				l.Warn("failed to create mech instance for player >%d< >%v<", playerNumber, err)
				return nil, err
			}
			out.MechInstances = append(out.MechInstances, mechInst)
		}

		l.Info("created player squad instance >%s< with >%d< mechs for subscription >%s<", squadInst.ID, len(starterMechs), subInst.ID)
	}

	// 4. Fetch opponent squad templates and computer opponents; randomly assign one template per opponent
	opponentTemplates, err := m.GetManyMechaGameSquadRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSquadGameID, Val: gameID},
			{Col: mecha_game_record.FieldMechaGameSquadSquadType, Val: mecha_game_record.SquadTypeOpponent},
		},
	})
	if err != nil {
		l.Warn("failed to get opponent squad templates >%v<", err)
		return nil, err
	}

	if len(opponentTemplates) > 0 {
		computerOpponents, err := m.GetManyMechaGameComputerOpponentRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameComputerOpponentGameID, Val: gameID},
			},
		})
		if err != nil {
			l.Warn("failed to get computer opponents >%v<", err)
			return nil, err
		}

		for i, opponent := range computerOpponents {
			template := opponentTemplates[i%len(opponentTemplates)]

			templateMechs, err := m.GetManyMechaGameSquadMechRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: mecha_game_record.FieldMechaGameSquadMechMechaGameSquadID, Val: template.ID},
				},
			})
			if err != nil {
				l.Warn("failed to get mechs for opponent template >%s< >%v<", template.ID, err)
				return nil, err
			}

			squadInst, err := m.CreateMechaGameSquadInstanceRec(&mecha_game_record.MechaGameSquadInstance{
				GameID:                     gameID,
				GameInstanceID:             instanceID,
				MechaGameSquadID:               template.ID,
				GameSubscriptionInstanceID: sql.NullString{Valid: false},
				MechaGameComputerOpponentID:    sql.NullString{String: opponent.ID, Valid: true},
			})
			if err != nil {
				l.Warn("failed to create opponent squad instance for opponent >%s< >%v<", opponent.ID, err)
				return nil, err
			}
			out.SquadInstances = append(out.SquadInstances, squadInst)

			if startingSectorInstanceID == "" {
				return nil, fmt.Errorf("no starting sector instance found for game >%s<: cannot create mech instances for computer opponent", gameID)
			}

			// Shuffle callsign suffix to give each clone distinct names
			offset := rand.Intn(100)
			for j, squadMech := range templateMechs {
				chassisRec, err := m.GetMechaGameChassisRec(squadMech.MechaGameChassisID, nil)
				if err != nil {
					l.Warn("failed to get chassis >%s< for opponent squad mech >%s< >%v<", squadMech.MechaGameChassisID, squadMech.ID, err)
					return nil, err
				}

				loadout, err := m.resolveMechInitialLoadout(squadMech)
				if err != nil {
					l.Warn("failed to resolve initial loadout for opponent mech >%s< >%v<", squadMech.ID, err)
					return nil, err
				}

				callsign := fmt.Sprintf("AI%d-%d", offset+j+1, i+1)
				mechInst, err := m.CreateMechaGameMechInstanceRec(&mecha_game_record.MechaGameMechInstance{
					GameID:                    gameID,
					GameInstanceID:            instanceID,
					MechaGameSquadInstanceID:  squadInst.ID,
					MechaGameSectorInstanceID: startingSectorInstanceID,
					MechaGameChassisID:        squadMech.MechaGameChassisID,
					Callsign:                  callsign,
					CurrentArmor:              chassisRec.ArmorPoints + loadout.ArmorBonus,
					CurrentStructure:          chassisRec.StructurePoints,
					CurrentHeat:               0,
					PilotSkill:                0,
					Status:                    mecha_game_record.MechInstanceStatusOperational,
					WeaponConfig:              loadout.WeaponConfig,
					WeaponConfigJSON:          loadout.WeaponConfigJSON,
					EquipmentConfig:           loadout.EquipmentConfig,
					EquipmentConfigJSON:       loadout.EquipmentJSON,
					AmmoRemaining:             loadout.AmmoRemaining,
				})
				if err != nil {
					l.Warn("failed to create mech instance for computer opponent >%s< >%v<", opponent.ID, err)
					return nil, err
				}
				out.MechInstances = append(out.MechInstances, mechInst)
			}

			l.Info("created opponent squad instance >%s< with >%d< mechs for opponent >%s< using template >%s<", squadInst.ID, len(templateMechs), opponent.ID, template.ID)
		}
	}

	l.Info("populated mecha instance data: sectors=%d squads=%d mechs=%d",
		len(out.SectorInstances), len(out.SquadInstances), len(out.MechInstances))

	return out, nil
}

// deleteMechaGameInstanceData removes mecha instance records (soft delete) for the given instanceID.
func (m *Domain) deleteMechaGameInstanceData(instanceID string) error {
	l := m.Logger("deleteMechaGameInstanceData")

	// Delete mecha_game_turn_sheet records linked to squad instances for this game instance
	squadInstances, err := m.GetManyMechaGameSquadInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSquadInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get squad instances >%v<", err)
		return databaseError(err)
	}

	for _, squadInst := range squadInstances {
		turnSheets, err := m.GetManyMechaGameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameTurnSheetMechaGameSquadInstanceID, Val: squadInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for squad instance >%s< >%v<", squadInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.MechaGameTurnSheetRepository().DeleteOne(ts.ID); err != nil {
				l.Warn("failed to delete mecha turn sheet >%s< >%v<", ts.ID, err)
				return databaseError(err)
			}
		}
	}

	// Delete mech instances
	mechInstances, err := m.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return databaseError(err)
	}
	for _, mechInst := range mechInstances {
		if err := m.MechaGameMechInstanceRepository().DeleteOne(mechInst.ID); err != nil {
			l.Warn("failed to delete mech instance >%s< >%v<", mechInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete squad instances
	for _, squadInst := range squadInstances {
		if err := m.MechaGameSquadInstanceRepository().DeleteOne(squadInst.ID); err != nil {
			l.Warn("failed to delete squad instance >%s< >%v<", squadInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete sector instances
	sectorInstances, err := m.GetManyMechaGameSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get sector instances >%v<", err)
		return databaseError(err)
	}
	for _, sectorInst := range sectorInstances {
		if err := m.MechaGameSectorInstanceRepository().DeleteOne(sectorInst.ID); err != nil {
			l.Warn("failed to delete sector instance >%s< >%v<", sectorInst.ID, err)
			return databaseError(err)
		}
	}

	return nil
}

// removeMechaGameInstanceData permanently removes mecha instance records for the given instanceID.
func (m *Domain) removeMechaGameInstanceData(instanceID string) error {
	l := m.Logger("removeMechaGameInstanceData")

	// Remove mecha_game_turn_sheet records linked to squad instances
	squadInstances, err := m.GetManyMechaGameSquadInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSquadInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get squad instances >%v<", err)
		return databaseError(err)
	}

	for _, squadInst := range squadInstances {
		turnSheets, err := m.GetManyMechaGameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_game_record.FieldMechaGameTurnSheetMechaGameSquadInstanceID, Val: squadInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for squad instance >%s< >%v<", squadInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.RemoveMechaGameTurnSheetRec(ts.ID); err != nil {
				l.Warn("failed to remove mecha turn sheet >%s< >%v<", ts.ID, err)
				return err
			}
		}
	}

	// Remove mech instances
	mechInstances, err := m.GetManyMechaGameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return databaseError(err)
	}
	for _, mechInst := range mechInstances {
		if err := m.RemoveMechaGameMechInstanceRec(mechInst.ID); err != nil {
			l.Warn("failed to remove mech instance >%s< >%v<", mechInst.ID, err)
			return err
		}
	}

	// Remove squad instances
	for _, squadInst := range squadInstances {
		if err := m.RemoveMechaGameSquadInstanceRec(squadInst.ID); err != nil {
			l.Warn("failed to remove squad instance >%s< >%v<", squadInst.ID, err)
			return err
		}
	}

	// Remove sector instances
	sectorInstances, err := m.GetManyMechaGameSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get sector instances >%v<", err)
		return databaseError(err)
	}
	for _, sectorInst := range sectorInstances {
		if err := m.RemoveMechaGameSectorInstanceRec(sectorInst.ID); err != nil {
			l.Warn("failed to remove sector instance >%s< >%v<", sectorInst.ID, err)
			return err
		}
	}

	return nil
}
