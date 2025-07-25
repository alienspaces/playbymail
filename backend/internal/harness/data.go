package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

// Data -
type Data struct {
	AccountRecs                     []*record.Account
	GameRecs                        []*record.Game
	GameLocationRecs                []*record.AdventureGameLocation
	GameLocationLinkRecs            []*record.AdventureGameLocationLink
	GameCharacterRecs               []*record.AdventureGameCharacter
	GameCreatureRecs                []*record.AdventureGameCreature
	GameItemRecs                    []*record.AdventureGameItem
	GameLocationLinkRequirementRecs []*record.AdventureGameLocationLinkRequirement
	GameInstanceRecs                []*record.AdventureGameInstance
	GameLocationInstanceRecs        []*record.AdventureGameLocationInstance
	GameItemInstanceRecs            []*record.AdventureGameItemInstance
	GameCreatureInstanceRecs        []*record.AdventureGameCreatureInstance
	GameCharacterInstanceRecs       []*record.AdventureGameCharacterInstance
	GameSubscriptionRecs            []*record.GameSubscription
	GameAdministrationRecs          []*record.GameAdministration
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
	GameCreatureRefs                map[string]string // Map of game_creature refs to records
	GameItemRefs                    map[string]string // Map of game_item refs to records
	GameItemInstanceRefs            map[string]string // Map of refs to game_item_instance records
	GameInstanceRefs                map[string]string // Map of refs to game_instance records
	GameLocationInstanceRefs        map[string]string // Map of refs to game_location_instance records
	GameCreatureInstanceRefs        map[string]string // Map of refs to game_creature_instance records
	GameCharacterInstanceRefs       map[string]string // Map of refs to game_character_instance records
	GameSubscriptionRefs            map[string]string // Map of refs to game_subscription records
	GameAdministrationRefs          map[string]string // Map of refs to game_administration records
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
			GameCreatureRefs:                map[string]string{},
			GameItemRefs:                    map[string]string{},
			GameItemInstanceRefs:            map[string]string{},
			GameInstanceRefs:                map[string]string{},
			GameLocationInstanceRefs:        map[string]string{},
			GameCreatureInstanceRefs:        map[string]string{},
			GameCharacterInstanceRefs:       map[string]string{},
			GameSubscriptionRefs:            map[string]string{},
			GameAdministrationRefs:          map[string]string{},
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
func (d *Data) AddGameLocationRec(rec *record.AdventureGameLocation) {
	for idx := range d.GameLocationRecs {
		if d.GameLocationRecs[idx].ID == rec.ID {
			d.GameLocationRecs[idx] = rec
			return
		}
	}
	d.GameLocationRecs = append(d.GameLocationRecs, rec)
}

func (d *Data) GetGameLocationRecByID(locationID string) (*record.AdventureGameLocation, error) {
	for _, rec := range d.GameLocationRecs {
		if rec.ID == locationID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location with ID >%s<", locationID)
}

func (d *Data) GetGameLocationRecByRef(ref string) (*record.AdventureGameLocation, error) {
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
func (d *Data) AddGameLocationLinkRec(rec *record.AdventureGameLocationLink) {
	for idx := range d.GameLocationLinkRecs {
		if d.GameLocationLinkRecs[idx].ID == rec.ID {
			d.GameLocationLinkRecs[idx] = rec
			return
		}
	}
	d.GameLocationLinkRecs = append(d.GameLocationLinkRecs, rec)
}

func (d *Data) GetGameLocationLinkRecByID(linkID string) (*record.AdventureGameLocationLink, error) {
	for _, rec := range d.GameLocationLinkRecs {
		if rec.ID == linkID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location link with ID >%s<", linkID)
}

func (d *Data) GetGameLocationLinkRecByRef(ref string) (*record.AdventureGameLocationLink, error) {
	linkID, ok := d.Refs.GameLocationLinkRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting location link with ref >%s<", ref)
	}
	return d.GetGameLocationLinkRecByID(linkID)
}

// GameCreature
func (d *Data) AddGameCreatureRec(rec *record.AdventureGameCreature) {
	for idx := range d.GameCreatureRecs {
		if d.GameCreatureRecs[idx].ID == rec.ID {
			d.GameCreatureRecs[idx] = rec
			return
		}
	}
	d.GameCreatureRecs = append(d.GameCreatureRecs, rec)
}

func (d *Data) GetGameCreatureRecByID(creatureID string) (*record.AdventureGameCreature, error) {
	for _, rec := range d.GameCreatureRecs {
		if rec.ID == creatureID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting creature with ID >%s<", creatureID)
}

func (d *Data) GetGameCreatureRecByRef(ref string) (*record.AdventureGameCreature, error) {
	id, ok := d.Refs.GameCreatureRefs[ref]
	if !ok {
		return nil, fmt.Errorf("no creature with ref >%s<", ref)
	}
	for _, rec := range d.GameCreatureRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("no creature with id >%s< for ref >%s<", id, ref)
}

// GameCharacter
func (d *Data) AddGameCharacterRec(rec *record.AdventureGameCharacter) {
	for idx := range d.GameCharacterRecs {
		if d.GameCharacterRecs[idx].ID == rec.ID {
			d.GameCharacterRecs[idx] = rec
			return
		}
	}
	d.GameCharacterRecs = append(d.GameCharacterRecs, rec)
}

func (d *Data) GetGameCharacterRecByID(id string) (*record.AdventureGameCharacter, error) {
	for _, rec := range d.GameCharacterRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_character with ID >%s<", id)
}

func (d *Data) GetGameCharacterRecByRef(ref string) (*record.AdventureGameCharacter, error) {
	id, ok := d.Refs.GameCharacterRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_character with ref >%s<", ref)
	}
	return d.GetGameCharacterRecByID(id)
}

// GameItem
func (d *Data) AddGameItemRec(rec *record.AdventureGameItem) {
	for idx := range d.GameItemRecs {
		if d.GameItemRecs[idx].ID == rec.ID {
			d.GameItemRecs[idx] = rec
			return
		}
	}
	d.GameItemRecs = append(d.GameItemRecs, rec)
}

func (d *Data) GetGameItemRecByID(itemID string) (*record.AdventureGameItem, error) {
	for _, rec := range d.GameItemRecs {
		if rec.ID == itemID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item with ID >%s<", itemID)
}

func (d *Data) GetGameItemRecByRef(ref string) (*record.AdventureGameItem, error) {
	id, ok := d.Refs.GameItemRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_item with ref >%s<", ref)
	}
	return d.GetGameItemRecByID(id)
}

// GameLocationLinkRequirement
func (d *Data) AddGameLocationLinkRequirementRec(rec *record.AdventureGameLocationLinkRequirement) {
	for idx := range d.GameLocationLinkRequirementRecs {
		if d.GameLocationLinkRequirementRecs[idx].ID == rec.ID {
			d.GameLocationLinkRequirementRecs[idx] = rec
			return
		}
	}
	d.GameLocationLinkRequirementRecs = append(d.GameLocationLinkRequirementRecs, rec)
}

func (d *Data) GetGameLocationLinkRequirementRecByID(id string) (*record.AdventureGameLocationLinkRequirement, error) {
	for _, rec := range d.GameLocationLinkRequirementRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_link_requirement with ID >%s<", id)
}

func (d *Data) GetGameLocationLinkRequirementRecByRef(ref string) (*record.AdventureGameLocationLinkRequirement, error) {
	id, ok := d.Refs.GameLocationLinkRequirementRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_link_requirement with ref >%s<", ref)
	}
	return d.GetGameLocationLinkRequirementRecByID(id)
}

// GameInstance
func (d *Data) AddGameInstanceRec(rec *record.AdventureGameInstance) {
	for idx := range d.GameInstanceRecs {
		if d.GameInstanceRecs[idx].ID == rec.ID {
			d.GameInstanceRecs[idx] = rec
			return
		}
	}
	d.GameInstanceRecs = append(d.GameInstanceRecs, rec)
}

func (d *Data) GetGameInstanceRecByID(id string) (*record.AdventureGameInstance, error) {
	for _, rec := range d.GameInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_instance with ID >%s<", id)
}

func (d *Data) GetGameInstanceRecByRef(ref string) (*record.AdventureGameInstance, error) {
	id, ok := d.Refs.GameInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_instance with ref >%s<", ref)
	}
	return d.GetGameInstanceRecByID(id)
}

func (d *Data) GetGameInstanceRecByGameID(gameID string) []*record.AdventureGameInstance {
	var result []*record.AdventureGameInstance
	for _, rec := range d.GameInstanceRecs {
		if rec.GameID == gameID {
			result = append(result, rec)
		}
	}
	return result
}

// GameLocationInstance
func (d *Data) AddGameLocationInstanceRec(rec *record.AdventureGameLocationInstance) {
	for idx := range d.GameLocationInstanceRecs {
		if d.GameLocationInstanceRecs[idx].ID == rec.ID {
			d.GameLocationInstanceRecs[idx] = rec
			return
		}
	}
	d.GameLocationInstanceRecs = append(d.GameLocationInstanceRecs, rec)
}

func (d *Data) GetGameLocationInstanceRecByID(id string) (*record.AdventureGameLocationInstance, error) {
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_instance with ID >%s<", id)
}

func (d *Data) GetGameLocationInstanceRecByRef(ref string) (*record.AdventureGameLocationInstance, error) {
	id, ok := d.Refs.GameLocationInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_instance with ref >%s<", ref)
	}
	return d.GetGameLocationInstanceRecByID(id)
}

func (d *Data) GetGameLocationInstanceRecByLocationRef(ref string) (*record.AdventureGameLocationInstance, error) {
	id, ok := d.Refs.GameLocationRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_instance with ref >%s<", ref)
	}
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.AdventureGameLocationID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_instance with location ID >%s<", id)
}

func (d *Data) GetGameLocationInstanceRecByLocationID(locationID string) []*record.AdventureGameLocationInstance {
	var result []*record.AdventureGameLocationInstance
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.AdventureGameLocationID == locationID {
			result = append(result, rec)
		}
	}
	return result
}

// GameItemInstance
func (d *Data) AddGameItemInstanceRec(rec *record.AdventureGameItemInstance) {
	for idx := range d.GameItemInstanceRecs {
		if d.GameItemInstanceRecs[idx].ID == rec.ID {
			d.GameItemInstanceRecs[idx] = rec
			return
		}
	}
	d.GameItemInstanceRecs = append(d.GameItemInstanceRecs, rec)
}

func (d *Data) GetGameItemInstanceRecByID(id string) (*record.AdventureGameItemInstance, error) {
	for _, rec := range d.GameItemInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item_instance with ID >%s<", id)
}

func (d *Data) GetGameItemInstanceRecByRef(ref string) (*record.AdventureGameItemInstance, error) {
	id, ok := d.Refs.GameItemInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_item_instance with ref >%s<", ref)
	}
	return d.GetGameItemInstanceRecByID(id)
}

func (d *Data) GetGameItemInstanceRecByItemRef(ref string) (*record.AdventureGameItemInstance, error) {
	id, ok := d.Refs.GameItemRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_item_instance with ref >%s<", ref)
	}
	for _, rec := range d.GameItemInstanceRecs {
		if rec.AdventureGameItemID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item_instance with item ID >%s<", id)
}

func (d *Data) GetGameItemInstanceRecByItemID(itemID string) []*record.AdventureGameItemInstance {
	var result []*record.AdventureGameItemInstance
	for _, rec := range d.GameItemInstanceRecs {
		if rec.AdventureGameItemID == itemID {
			result = append(result, rec)
		}
	}
	return result
}

// GameCreatureInstance
func (d *Data) AddGameCreatureInstanceRec(rec *record.AdventureGameCreatureInstance) {
	for idx := range d.GameCreatureInstanceRecs {
		if d.GameCreatureInstanceRecs[idx].ID == rec.ID {
			d.GameCreatureInstanceRecs[idx] = rec
			return
		}
	}
	d.GameCreatureInstanceRecs = append(d.GameCreatureInstanceRecs, rec)
}

func (d *Data) GetGameCreatureInstanceRecByID(id string) (*record.AdventureGameCreatureInstance, error) {
	for _, rec := range d.GameCreatureInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_creature_instance with ID >%s<", id)
}

func (d *Data) GetGameCreatureInstanceRecByRef(ref string) (*record.AdventureGameCreatureInstance, error) {
	id, ok := d.Refs.GameCreatureInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_creature_instance with ref >%s<", ref)
	}
	return d.GetGameCreatureInstanceRecByID(id)
}

// GameCharacterInstance
func (d *Data) AddGameCharacterInstanceRec(rec *record.AdventureGameCharacterInstance) {
	for idx := range d.GameCharacterInstanceRecs {
		if d.GameCharacterInstanceRecs[idx].ID == rec.ID {
			d.GameCharacterInstanceRecs[idx] = rec
			return
		}
	}
	d.GameCharacterInstanceRecs = append(d.GameCharacterInstanceRecs, rec)
}

func (d *Data) GetGameCharacterInstanceRecByID(id string) (*record.AdventureGameCharacterInstance, error) {
	for _, rec := range d.GameCharacterInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_character_instance with ID >%s<", id)
}

func (d *Data) GetGameCharacterInstanceRecByRef(ref string) (*record.AdventureGameCharacterInstance, error) {
	id, ok := d.Refs.GameCharacterInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_character_instance with ref >%s<", ref)
	}
	return d.GetGameCharacterInstanceRecByID(id)
}

func (d *Data) AddGameSubscriptionRec(rec *record.GameSubscription) {
	for idx := range d.GameSubscriptionRecs {
		if d.GameSubscriptionRecs[idx].ID == rec.ID {
			d.GameSubscriptionRecs[idx] = rec
			return
		}
	}
	d.GameSubscriptionRecs = append(d.GameSubscriptionRecs, rec)
}

func (d *Data) GetGameSubscriptionRecByID(id string) (*record.GameSubscription, error) {
	for _, rec := range d.GameSubscriptionRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_subscription with ID >%s<", id)
}

func (d *Data) AddGameAdministrationRec(rec *record.GameAdministration) {
	for idx := range d.GameAdministrationRecs {
		if d.GameAdministrationRecs[idx].ID == rec.ID {
			d.GameAdministrationRecs[idx] = rec
			return
		}
	}
	d.GameAdministrationRecs = append(d.GameAdministrationRecs, rec)
}

func (d *Data) GetGameAdministrationRecByID(id string) (*record.GameAdministration, error) {
	for _, rec := range d.GameAdministrationRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_administration with ID >%s<", id)
}
