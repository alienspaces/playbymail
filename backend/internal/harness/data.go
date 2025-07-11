package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

// Data -
type Data struct {
	AccountRecs                     []*record.Account
	GameRecs                        []*record.Game
	GameLocationRecs                []*record.GameLocation
	GameLocationLinkRecs            []*record.GameLocationLink
	GameCharacterRecs               []*record.GameCharacter
	GameItemRecs                    []*record.GameItem
	GameLocationLinkRequirementRecs []*record.GameLocationLinkRequirement
	GameInstanceRecs                []*record.GameInstance
	GameLocationInstanceRecs        []*record.GameLocationInstance
	GameItemInstanceRecs            []*record.GameItemInstance
	// Data references
	Refs DataRefs
}

// DataRefs is a collection of maps of data references that were defined
// in test data configuration to the resulting record that was created.
// When adding new reference maps make sure to also initialise the map
// in the initialiseDataStores() function further below.
type DataRefs struct {
	AccountRefs                     map[string]string // Map of account refs to account records
	GameRefs                        map[string]string // Map of game refs to game records
	GameLocationRefs                map[string]string // Map of game location refs to game location records
	GameLocationLinkRefs            map[string]string // Map of location link refs to location link records
	GameLocationLinkRequirementRefs map[string]string // Map of refs to game_location_link_requirement records
	GameCharacterRefs               map[string]string // Map of game_character refs to records
	GameItemRefs                    map[string]string // Map of game_item refs to records
	GameItemInstanceRefs            map[string]string // Map of refs to game_item_instance records
	GameInstanceRefs                map[string]string // Map of refs to game_instance records
	GameLocationInstanceRefs        map[string]string // Map of refs to game_location_instance records
}

// initialiseDataStores - Data is required to maintain data references and
// may contain main test data and reference test data so may not be used
// as a source of teardown data.
func initialiseDataStores() Data {
	return Data{
		Refs: DataRefs{
			AccountRefs:                     map[string]string{},
			GameRefs:                        map[string]string{},
			GameLocationRefs:                map[string]string{},
			GameLocationLinkRefs:            map[string]string{},
			GameLocationLinkRequirementRefs: map[string]string{},
			GameCharacterRefs:               map[string]string{},
			GameItemRefs:                    map[string]string{},
			GameItemInstanceRefs:            map[string]string{},
			GameInstanceRefs:                map[string]string{},
			GameLocationInstanceRefs:        map[string]string{},
		},
	}
}

// initialiseTeardownDataStores - Teardown data is not required to maintain
// data references but is used for cleaning up data after tests.
func initialiseTeardownDataStores() Data {
	return Data{}
}

// Account
func (d *Data) AddAccountRec(rec *record.Account) {
	for idx := range d.AccountRecs {
		if d.AccountRecs[idx].ID == rec.ID {
			d.AccountRecs[idx] = rec
			return
		}
	}
	d.AccountRecs = append(d.AccountRecs, rec)
}

func (d *Data) GetAccountRecByID(accountID string) (*record.Account, error) {
	for _, rec := range d.AccountRecs {
		if rec.ID == accountID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting account with ID >%s<", accountID)
}

func (d *Data) GetAccountRecByRef(ref string) (*record.Account, error) {
	accountID, ok := d.Refs.AccountRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting account with ref >%s<", ref)
	}
	return d.GetAccountRecByID(accountID)
}

// Game
func (d *Data) AddGameRec(rec *record.Game) {
	for idx := range d.GameRecs {
		if d.GameRecs[idx].ID == rec.ID {
			d.GameRecs[idx] = rec
			return
		}
	}
	d.GameRecs = append(d.GameRecs, rec)
}

func (d *Data) GetGameRecByID(gameID string) (*record.Game, error) {
	for _, rec := range d.GameRecs {
		if rec.ID == gameID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game with ID >%s<", gameID)
}

func (d *Data) GetGameRecByRef(ref string) (*record.Game, error) {
	gameID, ok := d.Refs.GameRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game with ref >%s<", ref)
	}
	return d.GetGameRecByID(gameID)
}

// GameLocation
func (d *Data) AddGameLocationRec(rec *record.GameLocation) {
	for idx := range d.GameLocationRecs {
		if d.GameLocationRecs[idx].ID == rec.ID {
			d.GameLocationRecs[idx] = rec
			return
		}
	}
	d.GameLocationRecs = append(d.GameLocationRecs, rec)
}

func (d *Data) GetGameLocationRecByID(locationID string) (*record.GameLocation, error) {
	for _, rec := range d.GameLocationRecs {
		if rec.ID == locationID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location with ID >%s<", locationID)
}

func (d *Data) GetGameLocationRecByRef(ref string) (*record.GameLocation, error) {
	id, ok := d.Refs.GameLocationRefs[ref]
	if !ok {
		return nil, fmt.Errorf("no location with ref >%s<", ref)
	}
	for _, rec := range d.GameLocationRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("no location with id >%s< for ref >%s<", id, ref)
}

// GameLocationLink
func (d *Data) AddGameLocationLinkRec(rec *record.GameLocationLink) {
	for idx := range d.GameLocationLinkRecs {
		if d.GameLocationLinkRecs[idx].ID == rec.ID {
			d.GameLocationLinkRecs[idx] = rec
			return
		}
	}
	d.GameLocationLinkRecs = append(d.GameLocationLinkRecs, rec)
}

func (d *Data) GetGameLocationLinkRecByID(linkID string) (*record.GameLocationLink, error) {
	for _, rec := range d.GameLocationLinkRecs {
		if rec.ID == linkID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location link with ID >%s<", linkID)
}

func (d *Data) GetGameLocationLinkRecByRef(ref string) (*record.GameLocationLink, error) {
	linkID, ok := d.Refs.GameLocationLinkRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting location link with ref >%s<", ref)
	}
	return d.GetGameLocationLinkRecByID(linkID)
}

// GameCharacter
func (d *Data) AddGameCharacterRec(rec *record.GameCharacter) {
	for idx := range d.GameCharacterRecs {
		if d.GameCharacterRecs[idx].ID == rec.ID {
			d.GameCharacterRecs[idx] = rec
			return
		}
	}
	d.GameCharacterRecs = append(d.GameCharacterRecs, rec)
}

func (d *Data) GetGameCharacterRecByID(id string) (*record.GameCharacter, error) {
	for _, rec := range d.GameCharacterRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_character with ID >%s<", id)
}

func (d *Data) GetGameCharacterRecByRef(ref string) (*record.GameCharacter, error) {
	id, ok := d.Refs.GameCharacterRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_character with ref >%s<", ref)
	}
	return d.GetGameCharacterRecByID(id)
}

// GameItem
func (d *Data) AddGameItemRec(rec *record.GameItem) {
	for idx := range d.GameItemRecs {
		if d.GameItemRecs[idx].ID == rec.ID {
			d.GameItemRecs[idx] = rec
			return
		}
	}
	d.GameItemRecs = append(d.GameItemRecs, rec)
}

func (d *Data) GetGameItemRecByID(itemID string) (*record.GameItem, error) {
	for _, rec := range d.GameItemRecs {
		if rec.ID == itemID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item with ID >%s<", itemID)
}

func (d *Data) GetGameItemRecByRef(ref string) (*record.GameItem, error) {
	id, ok := d.Refs.GameItemRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_item with ref >%s<", ref)
	}
	return d.GetGameItemRecByID(id)
}

// GameLocationLinkRequirement
func (d *Data) AddGameLocationLinkRequirementRec(rec *record.GameLocationLinkRequirement) {
	for idx := range d.GameLocationLinkRequirementRecs {
		if d.GameLocationLinkRequirementRecs[idx].ID == rec.ID {
			d.GameLocationLinkRequirementRecs[idx] = rec
			return
		}
	}
	d.GameLocationLinkRequirementRecs = append(d.GameLocationLinkRequirementRecs, rec)
}

func (d *Data) GetGameLocationLinkRequirementRecByID(id string) (*record.GameLocationLinkRequirement, error) {
	for _, rec := range d.GameLocationLinkRequirementRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_link_requirement with ID >%s<", id)
}

func (d *Data) GetGameLocationLinkRequirementRecByRef(ref string) (*record.GameLocationLinkRequirement, error) {
	id, ok := d.Refs.GameLocationLinkRequirementRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_link_requirement with ref >%s<", ref)
	}
	return d.GetGameLocationLinkRequirementRecByID(id)
}

// GameInstance
func (d *Data) AddGameInstanceRec(rec *record.GameInstance) {
	for idx := range d.GameInstanceRecs {
		if d.GameInstanceRecs[idx].ID == rec.ID {
			d.GameInstanceRecs[idx] = rec
			return
		}
	}
	d.GameInstanceRecs = append(d.GameInstanceRecs, rec)
}

func (d *Data) GetGameInstanceRecByID(id string) (*record.GameInstance, error) {
	for _, rec := range d.GameInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_instance with ID >%s<", id)
}

func (d *Data) GetGameInstanceRecByRef(ref string) (*record.GameInstance, error) {
	id, ok := d.Refs.GameInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_instance with ref >%s<", ref)
	}
	return d.GetGameInstanceRecByID(id)
}

func (d *Data) GetGameInstanceRecByGameID(gameID string) []*record.GameInstance {
	var result []*record.GameInstance
	for _, rec := range d.GameInstanceRecs {
		if rec.GameID == gameID {
			result = append(result, rec)
		}
	}
	return result
}

// GameLocationInstance
func (d *Data) AddGameLocationInstanceRec(rec *record.GameLocationInstance) {
	for idx := range d.GameLocationInstanceRecs {
		if d.GameLocationInstanceRecs[idx].ID == rec.ID {
			d.GameLocationInstanceRecs[idx] = rec
			return
		}
	}
	d.GameLocationInstanceRecs = append(d.GameLocationInstanceRecs, rec)
}

func (d *Data) GetGameLocationInstanceRecByID(id string) (*record.GameLocationInstance, error) {
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_instance with ID >%s<", id)
}

func (d *Data) GetGameLocationInstanceRecByRef(ref string) (*record.GameLocationInstance, error) {
	id, ok := d.Refs.GameLocationInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_instance with ref >%s<", ref)
	}
	return d.GetGameLocationInstanceRecByID(id)
}

func (d *Data) GetGameLocationInstanceRecByLocationRef(ref string) (*record.GameLocationInstance, error) {
	id, ok := d.Refs.GameLocationRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_instance with ref >%s<", ref)
	}
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.GameLocationID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_instance with location ID >%s<", id)
}

func (d *Data) GetGameLocationInstanceRecByLocationID(locationID string) []*record.GameLocationInstance {
	var result []*record.GameLocationInstance
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.GameLocationID == locationID {
			result = append(result, rec)
		}
	}
	return result
}

// GameItemInstance
func (d *Data) AddGameItemInstanceRec(rec *record.GameItemInstance) {
	for idx := range d.GameItemInstanceRecs {
		if d.GameItemInstanceRecs[idx].ID == rec.ID {
			d.GameItemInstanceRecs[idx] = rec
			return
		}
	}
	d.GameItemInstanceRecs = append(d.GameItemInstanceRecs, rec)
}

func (d *Data) GetGameItemInstanceRecByID(id string) (*record.GameItemInstance, error) {
	for _, rec := range d.GameItemInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item_instance with ID >%s<", id)
}

func (d *Data) GetGameItemInstanceRecByItemRef(ref string) (*record.GameItemInstance, error) {
	id, ok := d.Refs.GameItemRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_item_instance with ref >%s<", ref)
	}
	for _, rec := range d.GameItemInstanceRecs {
		if rec.GameItemID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item_instance with item ID >%s<", id)
}

func (d *Data) GetGameItemInstanceRecByItemID(itemID string) []*record.GameItemInstance {
	var result []*record.GameItemInstance
	for _, rec := range d.GameItemInstanceRecs {
		if rec.GameItemID == itemID {
			result = append(result, rec)
		}
	}
	return result
}
