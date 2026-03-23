package harness

import (
	"database/sql"
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

func (t *Testing) processMechaConfig(gameConfig GameConfig, gameRec *game_record.Game) error {
	l := t.Logger("processMechaConfig")

	for _, cfg := range gameConfig.MechaChassisConfigs {
		if _, err := t.createMechaChassisRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha chassis record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaWeaponConfigs {
		if _, err := t.createMechaWeaponRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha weapon record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaSectorConfigs {
		if _, err := t.createMechaSectorRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha sector record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaSectorLinkConfigs {
		if _, err := t.createMechaSectorLinkRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha sector link record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaComputerOpponentConfigs {
		if _, err := t.createMechaComputerOpponentRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha computer opponent record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaLanceConfigs {
		ownerCount := 0
		if cfg.AccountRef != "" {
			ownerCount++
		}
		if cfg.ComputerOpponentRef != "" {
			ownerCount++
		}
		if cfg.IsPlayerStarter {
			ownerCount++
		}
		if ownerCount == 0 {
			return fmt.Errorf("mecha lance config >%s< must have AccountRef, ComputerOpponentRef, or IsPlayerStarter set", cfg.Reference)
		}
		if ownerCount > 1 {
			return fmt.Errorf("mecha lance config >%s< must have exactly one of AccountRef, ComputerOpponentRef, or IsPlayerStarter set", cfg.Reference)
		}

		if cfg.ComputerOpponentRef != "" {
			opponentRec, err := t.Data.GetMechaComputerOpponentRecByRef(cfg.ComputerOpponentRef)
			if err != nil {
				l.Warn("failed resolving computer opponent ref >%s<: %v", cfg.ComputerOpponentRef, err)
				return err
			}
			if _, err := t.createMechaComputerOpponentLanceRec(cfg, gameRec, opponentRec); err != nil {
				l.Warn("failed creating mecha computer opponent lance record >%v<", err)
				return err
			}
		} else if cfg.IsPlayerStarter {
			if _, err := t.createMechaPlayerStarterLanceRec(cfg, gameRec); err != nil {
				l.Warn("failed creating mecha player starter lance record >%v<", err)
				return err
			}
		} else {
			accountUserRec, err := t.Data.GetAccountUserRecByRef(cfg.AccountRef)
			if err != nil {
				l.Warn("failed resolving account ref >%s<: %v", cfg.AccountRef, err)
				return err
			}
			if _, err := t.createMechaHumanLanceRec(cfg, gameRec, accountUserRec.ID, accountUserRec.AccountID); err != nil {
				l.Warn("failed creating mecha lance record >%v<", err)
				return err
			}
		}
	}

	return nil
}

func (t *Testing) removeMechaRecords() error {
	l := t.Logger("removeMechaRecords")

	// Lance mechs must be removed before lances
	l.Debug("removing >%d< mecha lance mech records", len(t.teardownData.MechaLanceMechRecs))
	for _, rec := range t.teardownData.MechaLanceMechRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaLanceMechRec(rec.ID); err != nil {
			l.Warn("failed removing mecha lance mech record >%v<", err)
			return err
		}
	}

	// Lances must be removed before computer opponents (FK dependency)
	l.Debug("removing >%d< mecha lance records", len(t.teardownData.MechaLanceRecs))
	for _, rec := range t.teardownData.MechaLanceRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaLanceRec(rec.ID); err != nil {
			l.Warn("failed removing mecha lance record >%v<", err)
			return err
		}
	}

	// Computer opponents removed after lances
	l.Debug("removing >%d< mecha computer opponent records", len(t.teardownData.MechaComputerOpponentRecs))
	for _, rec := range t.teardownData.MechaComputerOpponentRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaComputerOpponentRec(rec.ID); err != nil {
			l.Warn("failed removing mecha computer opponent record >%v<", err)
			return err
		}
	}

	// Sector links must be removed before sectors
	l.Debug("removing >%d< mecha sector link records", len(t.teardownData.MechaSectorLinkRecs))
	for _, rec := range t.teardownData.MechaSectorLinkRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaSectorLinkRec(rec.ID); err != nil {
			l.Warn("failed removing mecha sector link record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mecha sector records", len(t.teardownData.MechaSectorRecs))
	for _, rec := range t.teardownData.MechaSectorRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaSectorRec(rec.ID); err != nil {
			l.Warn("failed removing mecha sector record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mecha weapon records", len(t.teardownData.MechaWeaponRecs))
	for _, rec := range t.teardownData.MechaWeaponRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaWeaponRec(rec.ID); err != nil {
			l.Warn("failed removing mecha weapon record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mecha chassis records", len(t.teardownData.MechaChassisRecs))
	for _, rec := range t.teardownData.MechaChassisRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaChassisRec(rec.ID); err != nil {
			l.Warn("failed removing mecha chassis record >%v<", err)
			return err
		}
	}

	return nil
}

func (t *Testing) createMechaChassisRec(cfg MechaChassisConfig, gameRec *game_record.Game) (*mecha_record.MechaChassis, error) {
	l := t.Logger("createMechaChassisRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha chassis config >%#v<", cfg)
	}

	var rec *mecha_record.MechaChassis
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaChassis{}
	}

	rec = applyMechaChassisDefaults(rec)
	rec.GameID = gameRec.ID

	l.Debug("creating mecha chassis record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaChassisRec(rec)
	if err != nil {
		l.Warn("failed creating mecha chassis record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaChassisRec(rec)
	t.teardownData.AddMechaChassisRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaChassisRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func applyMechaChassisDefaults(rec *mecha_record.MechaChassis) *mecha_record.MechaChassis {
	if rec == nil {
		rec = &mecha_record.MechaChassis{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.ChassisClass == "" {
		rec.ChassisClass = mecha_record.ChassisClassMedium
	}
	if rec.ArmorPoints == 0 {
		rec.ArmorPoints = 100
	}
	if rec.StructurePoints == 0 {
		rec.StructurePoints = 50
	}
	if rec.HeatCapacity == 0 {
		rec.HeatCapacity = 30
	}
	if rec.Speed == 0 {
		rec.Speed = 3
	}
	return rec
}

func (t *Testing) createMechaWeaponRec(cfg MechaWeaponConfig, gameRec *game_record.Game) (*mecha_record.MechaWeapon, error) {
	l := t.Logger("createMechaWeaponRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha weapon config >%#v<", cfg)
	}

	var rec *mecha_record.MechaWeapon
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaWeapon{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.Damage == 0 {
		rec.Damage = 5
	}
	if rec.HeatCost == 0 {
		rec.HeatCost = 3
	}
	if rec.RangeBand == "" {
		rec.RangeBand = mecha_record.WeaponRangeBandMedium
	}
	if rec.MountSize == "" {
		rec.MountSize = mecha_record.WeaponMountSizeMedium
	}
	rec.GameID = gameRec.ID

	l.Debug("creating mecha weapon record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaWeaponRec(rec)
	if err != nil {
		l.Warn("failed creating mecha weapon record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaWeaponRec(rec)
	t.teardownData.AddMechaWeaponRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaWeaponRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaSectorRec(cfg MechaSectorConfig, gameRec *game_record.Game) (*mecha_record.MechaSector, error) {
	l := t.Logger("createMechaSectorRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha sector config >%#v<", cfg)
	}

	var rec *mecha_record.MechaSector
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaSector{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.TerrainType == "" {
		rec.TerrainType = mecha_record.SectorTerrainTypeOpen
	}
	rec.GameID = gameRec.ID

	l.Debug("creating mecha sector record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaSectorRec(rec)
	if err != nil {
		l.Warn("failed creating mecha sector record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaSectorRec(rec)
	t.teardownData.AddMechaSectorRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaSectorRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaSectorLinkRec(cfg MechaSectorLinkConfig, gameRec *game_record.Game) (*mecha_record.MechaSectorLink, error) {
	l := t.Logger("createMechaSectorLinkRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha sector link config >%#v<", cfg)
	}

	var rec *mecha_record.MechaSectorLink
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaSectorLink{}
	}

	rec.GameID = gameRec.ID

	if cfg.FromSectorRef != "" {
		fromID, ok := t.Data.Refs.MechaSectorRefs[cfg.FromSectorRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving from sector ref >%s<", cfg.FromSectorRef)
		}
		rec.FromMechaSectorID = fromID
	}

	if cfg.ToSectorRef != "" {
		toID, ok := t.Data.Refs.MechaSectorRefs[cfg.ToSectorRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving to sector ref >%s<", cfg.ToSectorRef)
		}
		rec.ToMechaSectorID = toID
	}

	l.Debug("creating mecha sector link record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaSectorLinkRec(rec)
	if err != nil {
		l.Warn("failed creating mecha sector link record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaSectorLinkRec(rec)
	t.teardownData.AddMechaSectorLinkRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaSectorLinkRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaComputerOpponentRec(cfg MechaComputerOpponentConfig, gameRec *game_record.Game) (*mecha_record.MechaComputerOpponent, error) {
	l := t.Logger("createMechaComputerOpponentRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha computer opponent config >%#v<", cfg)
	}

	var rec *mecha_record.MechaComputerOpponent
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaComputerOpponent{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.Aggression == 0 {
		rec.Aggression = 5
	}
	if rec.IQ == 0 {
		rec.IQ = 5
	}
	rec.GameID = gameRec.ID

	l.Debug("creating mecha computer opponent record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaComputerOpponentRec(rec)
	if err != nil {
		l.Warn("failed creating mecha computer opponent record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaComputerOpponentRec(rec)
	t.teardownData.AddMechaComputerOpponentRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaComputerOpponentRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaHumanLanceRec(cfg MechaLanceConfig, gameRec *game_record.Game, accountUserID, accountID string) (*mecha_record.MechaLance, error) {
	l := t.Logger("createMechaHumanLanceRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha lance config >%#v<", cfg)
	}

	var rec *mecha_record.MechaLance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaLance{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	rec.GameID = gameRec.ID
	rec.AccountID = sql.NullString{String: accountID, Valid: true}
	rec.AccountUserID = sql.NullString{String: accountUserID, Valid: true}

	l.Debug("creating mecha human lance record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaLanceRec(rec)
	if err != nil {
		l.Warn("failed creating mecha lance record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaLanceRec(rec)
	t.teardownData.AddMechaLanceRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaLanceRefs[cfg.Reference] = rec.ID
	}

	for _, mechCfg := range cfg.LanceMechConfigs {
		if _, err := t.createMechaLanceMechRec(mechCfg, gameRec, rec); err != nil {
			return nil, err
		}
	}

	return rec, nil
}

func (t *Testing) createMechaPlayerStarterLanceRec(cfg MechaLanceConfig, gameRec *game_record.Game) (*mecha_record.MechaLance, error) {
	l := t.Logger("createMechaPlayerStarterLanceRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha player starter lance config >%#v<", cfg)
	}

	var rec *mecha_record.MechaLance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaLance{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName("Starter Lance")
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	rec.GameID = gameRec.ID
	rec.IsPlayerStarter = true

	l.Debug("creating mecha player starter lance record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaLanceRec(rec)
	if err != nil {
		l.Warn("failed creating mecha player starter lance record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaLanceRec(rec)
	t.teardownData.AddMechaLanceRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaLanceRefs[cfg.Reference] = rec.ID
	}

	for _, mechCfg := range cfg.LanceMechConfigs {
		if _, err := t.createMechaLanceMechRec(mechCfg, gameRec, rec); err != nil {
			return nil, err
		}
	}

	return rec, nil
}

func (t *Testing) createMechaComputerOpponentLanceRec(cfg MechaLanceConfig, gameRec *game_record.Game, opponentRec *mecha_record.MechaComputerOpponent) (*mecha_record.MechaLance, error) {
	l := t.Logger("createMechaComputerOpponentLanceRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha computer opponent lance config >%#v<", cfg)
	}

	if opponentRec == nil {
		return nil, fmt.Errorf("computer opponent record is nil for mecha lance config >%#v<", cfg)
	}

	var rec *mecha_record.MechaLance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaLance{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	rec.GameID = gameRec.ID
	rec.MechaComputerOpponentID = sql.NullString{String: opponentRec.ID, Valid: true}

	l.Debug("creating mecha computer opponent lance record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaLanceRec(rec)
	if err != nil {
		l.Warn("failed creating mecha computer opponent lance record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaLanceRec(rec)
	t.teardownData.AddMechaLanceRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaLanceRefs[cfg.Reference] = rec.ID
	}

	for _, mechCfg := range cfg.LanceMechConfigs {
		if _, err := t.createMechaLanceMechRec(mechCfg, gameRec, rec); err != nil {
			return nil, err
		}
	}

	return rec, nil
}

func (t *Testing) createMechaLanceMechRec(cfg MechaLanceMechConfig, gameRec *game_record.Game, lanceRec *mecha_record.MechaLance) (*mecha_record.MechaLanceMech, error) {
	l := t.Logger("createMechaLanceMechRec")

	var rec *mecha_record.MechaLanceMech
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_record.MechaLanceMech{}
	}

	rec.GameID = gameRec.ID
	rec.MechaLanceID = lanceRec.ID

	if cfg.ChassisRef != "" {
		chassisID, ok := t.Data.Refs.MechaChassisRefs[cfg.ChassisRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving chassis ref >%s<", cfg.ChassisRef)
		}
		rec.MechaChassisID = chassisID
	}

	if len(cfg.WeaponConfigRefs) > 0 {
		entries := make([]mecha_record.WeaponConfigEntry, 0, len(cfg.WeaponConfigRefs))
		for _, wRef := range cfg.WeaponConfigRefs {
			weaponID, ok := t.Data.Refs.MechaWeaponRefs[wRef.WeaponRef]
			if !ok {
				return nil, fmt.Errorf("failed resolving weapon ref >%s<", wRef.WeaponRef)
			}
			entries = append(entries, mecha_record.WeaponConfigEntry{
				WeaponID:     weaponID,
				SlotLocation: wRef.SlotLocation,
			})
		}
		rec.WeaponConfig = entries
	}

	if rec.Callsign == "" {
		rec.Callsign = UniqueName("Mech")
	}

	l.Debug("creating mecha lance mech record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaLanceMechRec(rec)
	if err != nil {
		l.Warn("failed creating mecha lance mech record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaLanceMechRec(rec)
	t.teardownData.AddMechaLanceMechRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaLanceMechRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}
