package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

func (t *Testing) processMechWargameConfig(gameConfig GameConfig, gameRec *game_record.Game) error {
	l := t.Logger("processMechWargameConfig")

	for _, cfg := range gameConfig.MechWargameChassisConfigs {
		if _, err := t.createMechWargameChassisRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mech wargame chassis record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechWargameWeaponConfigs {
		if _, err := t.createMechWargameWeaponRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mech wargame weapon record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechWargameSectorConfigs {
		if _, err := t.createMechWargameSectorRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mech wargame sector record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechWargameSectorLinkConfigs {
		if _, err := t.createMechWargameSectorLinkRec(cfg, gameRec); err != nil {
			l.Warn("failed creating mech wargame sector link record >%v<", err)
			return err
		}
	}

	for _, cfg := range gameConfig.MechWargameLanceConfigs {
		if cfg.AccountRef == "" {
			return fmt.Errorf("mech wargame lance config >%s< must have an AccountRef set", cfg.Reference)
		}
		accountUserRec, err := t.Data.GetAccountUserRecByRef(cfg.AccountRef)
		if err != nil {
			l.Warn("failed resolving account ref >%s<: %v", cfg.AccountRef, err)
			return err
		}
		if _, err := t.createMechWargameLanceRec(cfg, gameRec, accountUserRec); err != nil {
			l.Warn("failed creating mech wargame lance record >%v<", err)
			return err
		}
	}

	return nil
}

func (t *Testing) removeMechWargameRecords() error {
	l := t.Logger("removeMechWargameRecords")

	// Lance mechs must be removed before lances
	l.Debug("removing >%d< mech wargame lance mech records", len(t.teardownData.MechWargameLanceMechRecs))
	for _, rec := range t.teardownData.MechWargameLanceMechRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechWargameLanceMechRec(rec.ID); err != nil {
			l.Warn("failed removing mech wargame lance mech record >%v<", err)
			return err
		}
	}

	// Lances must be removed before chassis (FK dependency)
	l.Debug("removing >%d< mech wargame lance records", len(t.teardownData.MechWargameLanceRecs))
	for _, rec := range t.teardownData.MechWargameLanceRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechWargameLanceRec(rec.ID); err != nil {
			l.Warn("failed removing mech wargame lance record >%v<", err)
			return err
		}
	}

	// Sector links must be removed before sectors
	l.Debug("removing >%d< mech wargame sector link records", len(t.teardownData.MechWargameSectorLinkRecs))
	for _, rec := range t.teardownData.MechWargameSectorLinkRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechWargameSectorLinkRec(rec.ID); err != nil {
			l.Warn("failed removing mech wargame sector link record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mech wargame sector records", len(t.teardownData.MechWargameSectorRecs))
	for _, rec := range t.teardownData.MechWargameSectorRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechWargameSectorRec(rec.ID); err != nil {
			l.Warn("failed removing mech wargame sector record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mech wargame weapon records", len(t.teardownData.MechWargameWeaponRecs))
	for _, rec := range t.teardownData.MechWargameWeaponRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechWargameWeaponRec(rec.ID); err != nil {
			l.Warn("failed removing mech wargame weapon record >%v<", err)
			return err
		}
	}

	l.Debug("removing >%d< mech wargame chassis records", len(t.teardownData.MechWargameChassisRecs))
	for _, rec := range t.teardownData.MechWargameChassisRecs {
		if rec.ID == "" {
			continue
		}
		if err := t.Domain.(*domain.Domain).RemoveMechWargameChassisRec(rec.ID); err != nil {
			l.Warn("failed removing mech wargame chassis record >%v<", err)
			return err
		}
	}

	return nil
}

func (t *Testing) createMechWargameChassisRec(cfg MechWargameChassisConfig, gameRec *game_record.Game) (*mech_wargame_record.MechWargameChassis, error) {
	l := t.Logger("createMechWargameChassisRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mech wargame chassis config >%#v<", cfg)
	}

	var rec *mech_wargame_record.MechWargameChassis
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mech_wargame_record.MechWargameChassis{}
	}

	rec = applyMechWargameChassisDefaults(rec)
	rec.GameID = gameRec.ID

	l.Debug("creating mech wargame chassis record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechWargameChassisRec(rec)
	if err != nil {
		l.Warn("failed creating mech wargame chassis record >%v<", err)
		return nil, err
	}

	t.Data.AddMechWargameChassisRec(rec)
	t.teardownData.AddMechWargameChassisRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechWargameChassisRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func applyMechWargameChassisDefaults(rec *mech_wargame_record.MechWargameChassis) *mech_wargame_record.MechWargameChassis {
	if rec == nil {
		rec = &mech_wargame_record.MechWargameChassis{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.ChassisClass == "" {
		rec.ChassisClass = mech_wargame_record.ChassisClassMedium
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

func (t *Testing) createMechWargameWeaponRec(cfg MechWargameWeaponConfig, gameRec *game_record.Game) (*mech_wargame_record.MechWargameWeapon, error) {
	l := t.Logger("createMechWargameWeaponRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mech wargame weapon config >%#v<", cfg)
	}

	var rec *mech_wargame_record.MechWargameWeapon
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mech_wargame_record.MechWargameWeapon{}
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
		rec.RangeBand = mech_wargame_record.WeaponRangeBandMedium
	}
	if rec.MountSize == "" {
		rec.MountSize = mech_wargame_record.WeaponMountSizeMedium
	}
	rec.GameID = gameRec.ID

	l.Debug("creating mech wargame weapon record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechWargameWeaponRec(rec)
	if err != nil {
		l.Warn("failed creating mech wargame weapon record >%v<", err)
		return nil, err
	}

	t.Data.AddMechWargameWeaponRec(rec)
	t.teardownData.AddMechWargameWeaponRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechWargameWeaponRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechWargameSectorRec(cfg MechWargameSectorConfig, gameRec *game_record.Game) (*mech_wargame_record.MechWargameSector, error) {
	l := t.Logger("createMechWargameSectorRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mech wargame sector config >%#v<", cfg)
	}

	var rec *mech_wargame_record.MechWargameSector
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mech_wargame_record.MechWargameSector{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	if rec.TerrainType == "" {
		rec.TerrainType = mech_wargame_record.SectorTerrainTypeOpen
	}
	rec.GameID = gameRec.ID

	l.Debug("creating mech wargame sector record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechWargameSectorRec(rec)
	if err != nil {
		l.Warn("failed creating mech wargame sector record >%v<", err)
		return nil, err
	}

	t.Data.AddMechWargameSectorRec(rec)
	t.teardownData.AddMechWargameSectorRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechWargameSectorRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechWargameSectorLinkRec(cfg MechWargameSectorLinkConfig, gameRec *game_record.Game) (*mech_wargame_record.MechWargameSectorLink, error) {
	l := t.Logger("createMechWargameSectorLinkRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mech wargame sector link config >%#v<", cfg)
	}

	var rec *mech_wargame_record.MechWargameSectorLink
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mech_wargame_record.MechWargameSectorLink{}
	}

	rec.GameID = gameRec.ID

	if cfg.FromSectorRef != "" {
		fromID, ok := t.Data.Refs.MechWargameSectorRefs[cfg.FromSectorRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving from sector ref >%s<", cfg.FromSectorRef)
		}
		rec.FromMechWargameSectorID = fromID
	}

	if cfg.ToSectorRef != "" {
		toID, ok := t.Data.Refs.MechWargameSectorRefs[cfg.ToSectorRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving to sector ref >%s<", cfg.ToSectorRef)
		}
		rec.ToMechWargameSectorID = toID
	}

	l.Debug("creating mech wargame sector link record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechWargameSectorLinkRec(rec)
	if err != nil {
		l.Warn("failed creating mech wargame sector link record >%v<", err)
		return nil, err
	}

	t.Data.AddMechWargameSectorLinkRec(rec)
	t.teardownData.AddMechWargameSectorLinkRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechWargameSectorLinkRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) createMechWargameLanceRec(cfg MechWargameLanceConfig, gameRec *game_record.Game, accountUserRec *account_record.AccountUser) (*mech_wargame_record.MechWargameLance, error) {
	l := t.Logger("createMechWargameLanceRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for mech wargame lance config >%#v<", cfg)
	}

	if accountUserRec == nil {
		return nil, fmt.Errorf("account user record is nil for mech wargame lance config >%#v<", cfg)
	}

	var rec *mech_wargame_record.MechWargameLance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mech_wargame_record.MechWargameLance{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(8)
	}
	rec.GameID = gameRec.ID
	rec.AccountID = accountUserRec.AccountID
	rec.AccountUserID = accountUserRec.ID

	l.Debug("creating mech wargame lance record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechWargameLanceRec(rec)
	if err != nil {
		l.Warn("failed creating mech wargame lance record >%v<", err)
		return nil, err
	}

	t.Data.AddMechWargameLanceRec(rec)
	t.teardownData.AddMechWargameLanceRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechWargameLanceRefs[cfg.Reference] = rec.ID
	}

	// Create lance mechs
	for _, mechCfg := range cfg.LanceMechConfigs {
		_, err := t.createMechWargameLanceMechRec(mechCfg, gameRec, rec)
		if err != nil {
			return nil, err
		}
	}

	return rec, nil
}

func (t *Testing) createMechWargameLanceMechRec(cfg MechWargameLanceMechConfig, gameRec *game_record.Game, lanceRec *mech_wargame_record.MechWargameLance) (*mech_wargame_record.MechWargameLanceMech, error) {
	l := t.Logger("createMechWargameLanceMechRec")

	var rec *mech_wargame_record.MechWargameLanceMech
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &mech_wargame_record.MechWargameLanceMech{}
	}

	rec.GameID = gameRec.ID
	rec.MechWargameLanceID = lanceRec.ID

	if cfg.ChassisRef != "" {
		chassisID, ok := t.Data.Refs.MechWargameChassisRefs[cfg.ChassisRef]
		if !ok {
			return nil, fmt.Errorf("failed resolving chassis ref >%s<", cfg.ChassisRef)
		}
		rec.MechWargameChassisID = chassisID
	}

	if rec.Callsign == "" {
		rec.Callsign = UniqueName("Mech")
	}

	l.Debug("creating mech wargame lance mech record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateMechWargameLanceMechRec(rec)
	if err != nil {
		l.Warn("failed creating mech wargame lance mech record >%v<", err)
		return nil, err
	}

	t.Data.AddMechWargameLanceMechRec(rec)
	t.teardownData.AddMechWargameLanceMechRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.MechWargameLanceMechRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}
