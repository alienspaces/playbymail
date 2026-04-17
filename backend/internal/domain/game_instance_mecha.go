package domain

import (
	"database/sql"
	"fmt"
	"math/rand"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// MechaInstanceData holds all instance records created when a mecha instance starts.
type MechaInstanceData struct {
	SectorInstances []*mecha_record.MechaSectorInstance
	SquadInstances  []*mecha_record.MechaSquadInstance
	MechInstances   []*mecha_record.MechaMechInstance
}

// PopulateMechaGameInstanceData creates all runtime records (sector instances, squad instances,
// mech instances) for a mecha game instance from its design definitions and player subscriptions.
//
// Player squads: each subscribed player gets a squad instance cloned from the starter template.
// Opponent squads: each computer opponent is randomly assigned an opponent squad template.
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

	var startingSectorInstanceID string

	for _, sector := range sectorRecs {
		sectorInst, err := m.CreateMechaSectorInstanceRec(&mecha_record.MechaSectorInstance{
			GameID:         gameID,
			GameInstanceID: instanceID,
			MechaSectorID:  sector.ID,
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
	starterRecs, err := m.GetManyMechaSquadRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSquadGameID, Val: gameID},
			{Col: mecha_record.FieldMechaSquadSquadType, Val: mecha_record.SquadTypeStarter},
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

	starterMechs, err := m.GetManyMechaSquadMechRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSquadMechMechaSquadID, Val: starterSquad.ID},
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

		squadInst, err := m.CreateMechaSquadInstanceRec(&mecha_record.MechaSquadInstance{
			GameID:                     gameID,
			GameInstanceID:             instanceID,
			MechaSquadID:               starterSquad.ID,
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
			chassisRec, err := m.GetMechaChassisRec(squadMech.MechaChassisID, nil)
			if err != nil {
				l.Warn("failed to get chassis >%s< for starter mech >%s< >%v<", squadMech.MechaChassisID, squadMech.ID, err)
				return nil, err
			}

			callsign := fmt.Sprintf("P%d-%d", playerNumber, i+1)
			mechInst, err := m.CreateMechaMechInstanceRec(&mecha_record.MechaMechInstance{
				GameID:                gameID,
				GameInstanceID:        instanceID,
				MechaSquadInstanceID:  squadInst.ID,
				MechaSectorInstanceID: startingSectorInstanceID,
				MechaChassisID:        squadMech.MechaChassisID,
				Callsign:              callsign,
				CurrentArmor:          chassisRec.ArmorPoints,
				CurrentStructure:      chassisRec.StructurePoints,
				CurrentHeat:           0,
				PilotSkill:            0,
				Status:                mecha_record.MechInstanceStatusOperational,
				WeaponConfig:          squadMech.WeaponConfig,
				WeaponConfigJSON:      squadMech.WeaponConfigJSON,
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
	opponentTemplates, err := m.GetManyMechaSquadRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSquadGameID, Val: gameID},
			{Col: mecha_record.FieldMechaSquadSquadType, Val: mecha_record.SquadTypeOpponent},
		},
	})
	if err != nil {
		l.Warn("failed to get opponent squad templates >%v<", err)
		return nil, err
	}

	if len(opponentTemplates) > 0 {
		computerOpponents, err := m.GetManyMechaComputerOpponentRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaComputerOpponentGameID, Val: gameID},
			},
		})
		if err != nil {
			l.Warn("failed to get computer opponents >%v<", err)
			return nil, err
		}

		for i, opponent := range computerOpponents {
			template := opponentTemplates[i%len(opponentTemplates)]

			templateMechs, err := m.GetManyMechaSquadMechRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: mecha_record.FieldMechaSquadMechMechaSquadID, Val: template.ID},
				},
			})
			if err != nil {
				l.Warn("failed to get mechs for opponent template >%s< >%v<", template.ID, err)
				return nil, err
			}

			squadInst, err := m.CreateMechaSquadInstanceRec(&mecha_record.MechaSquadInstance{
				GameID:                     gameID,
				GameInstanceID:             instanceID,
				MechaSquadID:               template.ID,
				GameSubscriptionInstanceID: sql.NullString{Valid: false},
				MechaComputerOpponentID:    sql.NullString{String: opponent.ID, Valid: true},
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
				chassisRec, err := m.GetMechaChassisRec(squadMech.MechaChassisID, nil)
				if err != nil {
					l.Warn("failed to get chassis >%s< for opponent squad mech >%s< >%v<", squadMech.MechaChassisID, squadMech.ID, err)
					return nil, err
				}

				callsign := fmt.Sprintf("AI%d-%d", offset+j+1, i+1)
				mechInst, err := m.CreateMechaMechInstanceRec(&mecha_record.MechaMechInstance{
					GameID:                gameID,
					GameInstanceID:        instanceID,
					MechaSquadInstanceID:  squadInst.ID,
					MechaSectorInstanceID: startingSectorInstanceID,
					MechaChassisID:        squadMech.MechaChassisID,
					Callsign:              callsign,
					CurrentArmor:          chassisRec.ArmorPoints,
					CurrentStructure:      chassisRec.StructurePoints,
					CurrentHeat:           0,
					PilotSkill:            0,
					Status:                mecha_record.MechInstanceStatusOperational,
					WeaponConfig:          squadMech.WeaponConfig,
					WeaponConfigJSON:      squadMech.WeaponConfigJSON,
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

// deleteMechaInstanceData removes mecha instance records (soft delete) for the given instanceID.
func (m *Domain) deleteMechaInstanceData(instanceID string) error {
	l := m.Logger("deleteMechaInstanceData")

	// Delete mecha_turn_sheet records linked to squad instances for this game instance
	squadInstances, err := m.GetManyMechaSquadInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSquadInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get squad instances >%v<", err)
		return databaseError(err)
	}

	for _, squadInst := range squadInstances {
		turnSheets, err := m.GetManyMechaTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaTurnSheetMechaSquadInstanceID, Val: squadInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for squad instance >%s< >%v<", squadInst.ID, err)
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

	// Delete squad instances
	for _, squadInst := range squadInstances {
		if err := m.MechaSquadInstanceRepository().DeleteOne(squadInst.ID); err != nil {
			l.Warn("failed to delete squad instance >%s< >%v<", squadInst.ID, err)
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

	// Remove mecha_turn_sheet records linked to squad instances
	squadInstances, err := m.GetManyMechaSquadInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mecha_record.FieldMechaSquadInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get squad instances >%v<", err)
		return databaseError(err)
	}

	for _, squadInst := range squadInstances {
		turnSheets, err := m.GetManyMechaTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mecha_record.FieldMechaTurnSheetMechaSquadInstanceID, Val: squadInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for squad instance >%s< >%v<", squadInst.ID, err)
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

	// Remove squad instances
	for _, squadInst := range squadInstances {
		if err := m.RemoveMechaSquadInstanceRec(squadInst.ID); err != nil {
			l.Warn("failed to remove squad instance >%s< >%v<", squadInst.ID, err)
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
