package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

func (t *Testing) processMechaGameConfig(gameConfig GameConfig, gameRec *game_record.Game) error {
	l := t.Logger("processMechaGameConfig")

	for _, cfg := range gameConfig.MechaGameChassisConfigs {
		if _, err := t.createMechaGameChassisRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha chassis record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaGameWeaponConfigs {
		if _, err := t.createMechaGameWeaponRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha weapon record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaGameSectorConfigs {
		if _, err := t.createMechaGameSectorRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha sector record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaGameSectorLinkConfigs {
		if _, err := t.createMechaGameSectorLinkRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha sector link record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaGameComputerOpponentConfigs {
		if _, err := t.createMechaGameComputerOpponentRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha computer opponent record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechaGameSquadConfigs {
		if cfg.SquadType != mecha_game_record.SquadTypeStarter && cfg.SquadType != mecha_game_record.SquadTypeOpponent {
			return fmt.Errorf("mecha squad config >%s< must have SquadType set to 'starter' or 'opponent'", cfg.Reference)
		}
		if _, err := t.createMechaGameSquadRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mecha squad record >%v<", err)
			return err
		}
	}

	return nil
}

func (t *Testing) removeMechaGameRecords() error {
	l := t.Logger("removeMechaGameRecords")

	// Squad mechs must be removed before squads
	l.Debug("removing >%d< mecha squad mech records", len(t.teardownData.MechaGameSquadMechRecs))
	for _, rec := range t.teardownData.MechaGameSquadMechRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaGameSquadMechRec(rec.ID); err != nil {
			l.Warn("failed removing mecha squad mech record >%v<", err)
			return err
		}
	}

	// Squads must be removed before computer opponents (FK dependency)
	l.Debug("removing >%d< mecha squad records", len(t.teardownData.MechaGameSquadRecs))
	for _, rec := range t.teardownData.MechaGameSquadRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaGameSquadRec(rec.ID); err != nil {
			l.Warn("failed removing mecha squad record >%v<", err)
			return err
		}
	}

	// Computer opponents removed after squads
	l.Debug("removing >%d< mecha computer opponent records", len(t.teardownData.MechaGameComputerOpponentRecs))
	for _, rec := range t.teardownData.MechaGameComputerOpponentRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaGameComputerOpponentRec(rec.ID); err != nil {
			l.Warn("failed removing mecha computer opponent record >%v<", err)
			return err
		}
	}

	// Sector links must be removed before sectors
	l.Debug("removing >%d< mecha sector link records", len(t.teardownData.MechaGameSectorLinkRecs))
	for _, rec := range t.teardownData.MechaGameSectorLinkRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaGameSectorLinkRec(rec.ID); err != nil {
			l.Warn("failed removing mecha sector link record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mecha sector records", len(t.teardownData.MechaGameSectorRecs))
	for _, rec := range t.teardownData.MechaGameSectorRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaGameSectorRec(rec.ID); err != nil {
			l.Warn("failed removing mecha sector record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mecha weapon records", len(t.teardownData.MechaGameWeaponRecs))
	for _, rec := range t.teardownData.MechaGameWeaponRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaGameWeaponRec(rec.ID); err != nil {
			l.Warn("failed removing mecha weapon record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mecha chassis records", len(t.teardownData.MechaGameChassisRecs))
	for _, rec := range t.teardownData.MechaGameChassisRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechaGameChassisRec(rec.ID); err != nil {
			l.Warn("failed removing mecha chassis record >%v<", err)
			return err
		}
	}

	return nil
}

func (t *Testing) createMechaGameChassisRec(cfg MechaGameChassisConfig, gameRec *game_record.Game) (*mecha_game_record.MechaGameChassis, error) {
	l := t.Logger("createMechaGameChassisRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha chassis config >%#v<", cfg)
	}

	var rec *mecha_game_record.MechaGameChassis
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_game_record.MechaGameChassis{}
	}

	rec = applyMechaGameChassisDefaults(rec)
	rec.GameID = gameRec.ID

	l.Debug("creating mecha chassis record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaGameChassisRec(rec)
	if err != nil {
		l.Warn("failed creating mecha chassis record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaGameChassisRec(rec)
	t.teardownData.AddMechaGameChassisRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaGameChassisRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func applyMechaGameChassisDefaults(rec *mecha_game_record.MechaGameChassis) *mecha_game_record.MechaGameChassis {
	if rec == nil {
		rec = &mecha_game_record.MechaGameChassis{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.ChassisClass == "" {
		rec.ChassisClass = mecha_game_record.ChassisClassMedium
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

func (t *Testing) createMechaGameWeaponRec(cfg MechaGameWeaponConfig, gameRec *game_record.Game) (*mecha_game_record.MechaGameWeapon, error) {
	l := t.Logger("createMechaGameWeaponRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha weapon config >%#v<", cfg)
	}

	var rec *mecha_game_record.MechaGameWeapon
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_game_record.MechaGameWeapon{}
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
		rec.RangeBand = mecha_game_record.WeaponRangeBandMedium
	}
	if rec.MountSize == "" {
		rec.MountSize = mecha_game_record.WeaponMountSizeMedium
	}
	rec.GameID = gameRec.ID

	l.Debug("creating mecha weapon record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaGameWeaponRec(rec)
	if err != nil {
		l.Warn("failed creating mecha weapon record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaGameWeaponRec(rec)
	t.teardownData.AddMechaGameWeaponRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaGameWeaponRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaGameSectorRec(cfg MechaGameSectorConfig, gameRec *game_record.Game) (*mecha_game_record.MechaGameSector, error) {
	l := t.Logger("createMechaGameSectorRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha sector config >%#v<", cfg)
	}

	var rec *mecha_game_record.MechaGameSector
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_game_record.MechaGameSector{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.TerrainType == "" {
		rec.TerrainType = mecha_game_record.SectorTerrainTypeOpen
	}
	rec.GameID = gameRec.ID

	l.Debug("creating mecha sector record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaGameSectorRec(rec)
	if err != nil {
		l.Warn("failed creating mecha sector record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaGameSectorRec(rec)
	t.teardownData.AddMechaGameSectorRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaGameSectorRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaGameSectorLinkRec(cfg MechaGameSectorLinkConfig, gameRec *game_record.Game) (*mecha_game_record.MechaGameSectorLink, error) {
	l := t.Logger("createMechaGameSectorLinkRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha sector link config >%#v<", cfg)
	}

	var rec *mecha_game_record.MechaGameSectorLink
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_game_record.MechaGameSectorLink{}
	}

	rec.GameID = gameRec.ID

	if cfg.FromSectorRef != "" {
		fromID, ok := t.Data.Refs.MechaGameSectorRefs[cfg.FromSectorRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving from sector ref >%s<", cfg.FromSectorRef)
		}
		rec.FromMechaGameSectorID = fromID
	}

	if cfg.ToSectorRef != "" {
		toID, ok := t.Data.Refs.MechaGameSectorRefs[cfg.ToSectorRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving to sector ref >%s<", cfg.ToSectorRef)
		}
		rec.ToMechaGameSectorID = toID
	}

	l.Debug("creating mecha sector link record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaGameSectorLinkRec(rec)
	if err != nil {
		l.Warn("failed creating mecha sector link record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaGameSectorLinkRec(rec)
	t.teardownData.AddMechaGameSectorLinkRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaGameSectorLinkRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaGameComputerOpponentRec(cfg MechaGameComputerOpponentConfig, gameRec *game_record.Game) (*mecha_game_record.MechaGameComputerOpponent, error) {
	l := t.Logger("createMechaGameComputerOpponentRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha computer opponent config >%#v<", cfg)
	}

	var rec *mecha_game_record.MechaGameComputerOpponent
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_game_record.MechaGameComputerOpponent{}
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

	rec, err := t.Domain.(*domain.Domain).CreateMechaGameComputerOpponentRec(rec)
	if err != nil {
		l.Warn("failed creating mecha computer opponent record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaGameComputerOpponentRec(rec)
	t.teardownData.AddMechaGameComputerOpponentRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaGameComputerOpponentRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechaGameSquadRec(cfg MechaGameSquadConfig, gameRec *game_record.Game) (*mecha_game_record.MechaGameSquad, error) {
	l := t.Logger("createMechaGameSquadRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mecha squad config >%#v<", cfg)
	}

	var rec *mecha_game_record.MechaGameSquad
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_game_record.MechaGameSquad{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	rec.GameID = gameRec.ID
	rec.SquadType = cfg.SquadType

	l.Debug("creating mecha squad record type >%s< >%#v<", cfg.SquadType, rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaGameSquadRec(rec)
	if err != nil {
		l.Warn("failed creating mecha squad record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaGameSquadRec(rec)
	t.teardownData.AddMechaGameSquadRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaGameSquadRefs[cfg.Reference] = rec.ID
	}

	for _, mechCfg := range cfg.SquadMechConfigs {
		if _, err := t.createMechaGameSquadMechRec(mechCfg, gameRec, rec); err != nil {
			return nil, err
		}
	}

	return rec, nil
}

func (t *Testing) createMechaGameSquadMechRec(cfg MechaGameSquadMechConfig, gameRec *game_record.Game, squadRec *mecha_game_record.MechaGameSquad) (*mecha_game_record.MechaGameSquadMech, error) {
	l := t.Logger("createMechaGameSquadMechRec")

	var rec *mecha_game_record.MechaGameSquadMech
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mecha_game_record.MechaGameSquadMech{}
	}

	rec.GameID = gameRec.ID
	rec.MechaGameSquadID = squadRec.ID

	if cfg.ChassisRef != "" {
		chassisID, ok := t.Data.Refs.MechaGameChassisRefs[cfg.ChassisRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving chassis ref >%s<", cfg.ChassisRef)
		}
		rec.MechaGameChassisID = chassisID
	}

	if len(cfg.WeaponConfigRefs) > 0 {
		entries := make([]mecha_game_record.WeaponConfigEntry, 0, len(cfg.WeaponConfigRefs))
		for _, wRef := range cfg.WeaponConfigRefs {
			weaponID, ok := t.Data.Refs.MechaGameWeaponRefs[wRef.WeaponRef]
			if !ok {
				return nil, fmt.Errorf("failed resolving weapon ref >%s<", wRef.WeaponRef)
			}
			entries = append(entries, mecha_game_record.WeaponConfigEntry{
				WeaponID:     weaponID,
				SlotLocation: wRef.SlotLocation,
			})
		}
		rec.WeaponConfig = entries
	}

	if rec.Callsign == "" {
		rec.Callsign = UniqueName("Mech")
	}

	l.Debug("creating mecha squad mech record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechaGameSquadMechRec(rec)
	if err != nil {
		l.Warn("failed creating mecha squad mech record >%v<", err)
		return nil, err
	}

	t.Data.AddMechaGameSquadMechRec(rec)
	t.teardownData.AddMechaGameSquadMechRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechaGameSquadMechRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}
