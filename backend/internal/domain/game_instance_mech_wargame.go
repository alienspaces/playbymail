package domain

import (
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

// MechWargameInstanceData holds all instance records created when a mech wargame instance starts.
type MechWargameInstanceData struct {
	SectorInstances []*mech_wargame_record.MechWargameSectorInstance
	LanceInstances  []*mech_wargame_record.MechWargameLanceInstance
	MechInstances   []*mech_wargame_record.MechWargameMechInstance
}

// PopulateMechWargameGameInstanceData creates all runtime records (sector instances, lance instances,
// mech instances) for a mech wargame game instance from its design definitions and player subscriptions.
func (m *Domain) PopulateMechWargameGameInstanceData(instanceID string) (*MechWargameInstanceData, error) {
	l := m.Logger("PopulateMechWargameGameInstanceData")

	l.Info("populating mech wargame instance data for instance >%s<", instanceID)

	instanceRec, err := m.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return nil, err
	}

	gameID := instanceRec.GameID
	out := &MechWargameInstanceData{}

	// 1. Create sector instances and build sectorID -> sectorInstanceID map
	sectorRecs, err := m.GetManyMechWargameSectorRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameSectorGameID, Val: gameID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech wargame sectors >%v<", err)
		return nil, err
	}

	sectorIDToInstanceID := make(map[string]string, len(sectorRecs))
	var startingSectorInstanceID string

	for _, sector := range sectorRecs {
		sectorInst, err := m.CreateMechWargameSectorInstanceRec(&mech_wargame_record.MechWargameSectorInstance{
			GameID:              gameID,
			GameInstanceID:      instanceID,
			MechWargameSectorID: sector.ID,
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

		lanceRecs, err := m.GetManyMechWargameLanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mech_wargame_record.FieldMechWargameLanceGameID, Val: gameID},
				{Col: mech_wargame_record.FieldMechWargameLanceAccountUserID, Val: sub.AccountUserID},
			},
			Limit: 1,
		})
		if err != nil || len(lanceRecs) == 0 {
			l.Warn("no mech wargame lance found for account_user_id >%s< game >%s< - skipping", sub.AccountUserID, gameID)
			continue
		}

		lanceID := lanceRecs[0].ID
		if createdLanceIDs[lanceID] {
			l.Info("lance instance already created for lance >%s< - skipping duplicate subscription", lanceID)
			continue
		}

		lanceInst, err := m.CreateMechWargameLanceInstanceRec(&mech_wargame_record.MechWargameLanceInstance{
			GameID:                     gameID,
			GameInstanceID:             instanceID,
			MechWargameLanceID:         lanceID,
			GameSubscriptionInstanceID: subInst.ID,
		})
		if err != nil {
			l.Warn("failed to create lance instance for lance >%s< >%v<", lanceID, err)
			return nil, err
		}
		createdLanceIDs[lanceID] = true
		out.LanceInstances = append(out.LanceInstances, lanceInst)

		// 3. Create mech instances for each mech in the lance
		lanceMechRecs, err := m.GetManyMechWargameLanceMechRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mech_wargame_record.FieldMechWargameLanceMechMechWargameLanceID, Val: lanceID},
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
			chassisRec, err := m.GetMechWargameChassisRec(lanceMech.MechWargameChassisID, nil)
			if err != nil {
				l.Warn("failed to get chassis >%s< for lance mech >%s< >%v<", lanceMech.MechWargameChassisID, lanceMech.ID, err)
				return nil, err
			}

			mechInst, err := m.CreateMechWargameMechInstanceRec(&mech_wargame_record.MechWargameMechInstance{
				GameID:                       gameID,
				GameInstanceID:               instanceID,
				MechWargameLanceInstanceID:   lanceInst.ID,
				MechWargameSectorInstanceID:  startingSectorInstanceID,
				MechWargameChassisID:         lanceMech.MechWargameChassisID,
				Callsign:                     lanceMech.Callsign,
				CurrentArmor:                 chassisRec.ArmorPoints,
				CurrentStructure:             chassisRec.StructurePoints,
				CurrentHeat:                  0,
				PilotSkill:                   4,
				Status:                       mech_wargame_record.MechInstanceStatusOperational,
			})
			if err != nil {
				l.Warn("failed to create mech instance for lance mech >%s< >%v<", lanceMech.ID, err)
				return nil, err
			}
			out.MechInstances = append(out.MechInstances, mechInst)
		}

		l.Info("created lance instance >%s< with %d mech instances for lance >%s<", lanceInst.ID, len(lanceMechRecs), lanceID)
	}

	l.Info("populated mech wargame instance data: sectors=%d lances=%d mechs=%d",
		len(out.SectorInstances), len(out.LanceInstances), len(out.MechInstances))

	return out, nil
}

// deleteMechWargameInstanceData removes mech wargame instance records (soft delete) for the given instanceID.
func (m *Domain) deleteMechWargameInstanceData(instanceID string) error {
	l := m.Logger("deleteMechWargameInstanceData")

	// Delete mech_wargame_turn_sheet records linked to lance instances for this game instance
	lanceInstances, err := m.GetManyMechWargameLanceInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameLanceInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get lance instances >%v<", err)
		return databaseError(err)
	}

	for _, lanceInst := range lanceInstances {
		turnSheets, err := m.GetManyMechWargameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mech_wargame_record.FieldMechWargameTurnSheetMechWargameLanceInstanceID, Val: lanceInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for lance instance >%s< >%v<", lanceInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.MechWargameTurnSheetRepository().DeleteOne(ts.ID); err != nil {
				l.Warn("failed to delete mech wargame turn sheet >%s< >%v<", ts.ID, err)
				return databaseError(err)
			}
		}
	}

	// Delete mech instances
	mechInstances, err := m.GetManyMechWargameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameMechInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return databaseError(err)
	}
	for _, mechInst := range mechInstances {
		if err := m.MechWargameMechInstanceRepository().DeleteOne(mechInst.ID); err != nil {
			l.Warn("failed to delete mech instance >%s< >%v<", mechInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete lance instances
	for _, lanceInst := range lanceInstances {
		if err := m.MechWargameLanceInstanceRepository().DeleteOne(lanceInst.ID); err != nil {
			l.Warn("failed to delete lance instance >%s< >%v<", lanceInst.ID, err)
			return databaseError(err)
		}
	}

	// Delete sector instances
	sectorInstances, err := m.GetManyMechWargameSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameSectorInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get sector instances >%v<", err)
		return databaseError(err)
	}
	for _, sectorInst := range sectorInstances {
		if err := m.MechWargameSectorInstanceRepository().DeleteOne(sectorInst.ID); err != nil {
			l.Warn("failed to delete sector instance >%s< >%v<", sectorInst.ID, err)
			return databaseError(err)
		}
	}

	return nil
}

// removeMechWargameInstanceData permanently removes mech wargame instance records for the given instanceID.
func (m *Domain) removeMechWargameInstanceData(instanceID string) error {
	l := m.Logger("removeMechWargameInstanceData")

	// Remove mech_wargame_turn_sheet records linked to lance instances
	lanceInstances, err := m.GetManyMechWargameLanceInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameLanceInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get lance instances >%v<", err)
		return databaseError(err)
	}

	for _, lanceInst := range lanceInstances {
		turnSheets, err := m.GetManyMechWargameTurnSheetRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: mech_wargame_record.FieldMechWargameTurnSheetMechWargameLanceInstanceID, Val: lanceInst.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get turn sheets for lance instance >%s< >%v<", lanceInst.ID, err)
			return databaseError(err)
		}
		for _, ts := range turnSheets {
			if err := m.RemoveMechWargameTurnSheetRec(ts.ID); err != nil {
				l.Warn("failed to remove mech wargame turn sheet >%s< >%v<", ts.ID, err)
				return err
			}
		}
	}

	// Remove mech instances
	mechInstances, err := m.GetManyMechWargameMechInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameMechInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get mech instances >%v<", err)
		return databaseError(err)
	}
	for _, mechInst := range mechInstances {
		if err := m.RemoveMechWargameMechInstanceRec(mechInst.ID); err != nil {
			l.Warn("failed to remove mech instance >%s< >%v<", mechInst.ID, err)
			return err
		}
	}

	// Remove lance instances
	for _, lanceInst := range lanceInstances {
		if err := m.RemoveMechWargameLanceInstanceRec(lanceInst.ID); err != nil {
			l.Warn("failed to remove lance instance >%s< >%v<", lanceInst.ID, err)
			return err
		}
	}

	// Remove sector instances
	sectorInstances, err := m.GetManyMechWargameSectorInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: mech_wargame_record.FieldMechWargameSectorInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		l.Warn("failed to get sector instances >%v<", err)
		return databaseError(err)
	}
	for _, sectorInst := range sectorInstances {
		if err := m.RemoveMechWargameSectorInstanceRec(sectorInst.ID); err != nil {
			l.Warn("failed to remove sector instance >%s< >%v<", sectorInst.ID, err)
			return err
		}
	}

	return nil
}
