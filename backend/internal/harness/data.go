package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// Data -
type Data struct {
	AccountRecs                     []*account_record.Account
	GameRecs                        []*game_record.Game
	GameLocationRecs                []*adventure_game_record.AdventureGameLocation
	GameLocationLinkRecs            []*adventure_game_record.AdventureGameLocationLink
	GameCharacterRecs               []*adventure_game_record.AdventureGameCharacter
	GameCreatureRecs                []*adventure_game_record.AdventureGameCreature
	GameItemRecs                    []*adventure_game_record.AdventureGameItem
	GameLocationLinkRequirementRecs []*adventure_game_record.AdventureGameLocationLinkRequirement
	GameInstanceRecs                []*game_record.GameInstance
	GameLocationInstanceRecs        []*adventure_game_record.AdventureGameLocationInstance
	GameItemInstanceRecs            []*adventure_game_record.AdventureGameItemInstance
	GameCreatureInstanceRecs        []*adventure_game_record.AdventureGameCreatureInstance
	GameCharacterInstanceRecs       []*adventure_game_record.AdventureGameCharacterInstance
	GameSubscriptionRecs            []*game_record.GameSubscription
	GameAdministrationRecs          []*game_record.GameAdministration
	GameInstanceParameterRecs       []*game_record.GameInstanceParameter
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
	GameInstanceParameterRefs       map[string]string // Map of refs to game_instance_parameter records
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
			GameInstanceParameterRefs:       map[string]string{},
		},
	}
}

// initialiseTeardownDataStores - Teardown data is not required to maintain
// data references but is used for cleaning up data after tests.
func initialiseTeardownDataStores() Data {
	return Data{}
}

// Account
func (d *Data) AddAccountRec(rec *account_record.Account) {
	for idx := range d.AccountRecs {
		if d.AccountRecs[idx].ID == rec.ID {
			d.AccountRecs[idx] = rec
			return
		}
	}
	d.AccountRecs = append(d.AccountRecs, rec)
}

func (d *Data) GetAccountRecByID(accountID string) (*account_record.Account, error) {
	for _, rec := range d.AccountRecs {
		if rec.ID == accountID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting account with ID >%s<", accountID)
}

func (d *Data) GetAccountRecByRef(ref string) (*account_record.Account, error) {
	accountID, ok := d.Refs.AccountRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting account with ref >%s<", ref)
	}
	return d.GetAccountRecByID(accountID)
}

// Game
func (d *Data) AddGameRec(rec *game_record.Game) {
	for idx := range d.GameRecs {
		if d.GameRecs[idx].ID == rec.ID {
			d.GameRecs[idx] = rec
			return
		}
	}
	d.GameRecs = append(d.GameRecs, rec)
}

func (d *Data) GetGameRecByID(gameID string) (*game_record.Game, error) {
	for _, rec := range d.GameRecs {
		if rec.ID == gameID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game with ID >%s<", gameID)
}

func (d *Data) GetGameRecByRef(ref string) (*game_record.Game, error) {
	gameID, ok := d.Refs.GameRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game with ref >%s<", ref)
	}
	return d.GetGameRecByID(gameID)
}

// GameLocation
func (d *Data) AddGameLocationRec(rec *adventure_game_record.AdventureGameLocation) {
	for idx := range d.GameLocationRecs {
		if d.GameLocationRecs[idx].ID == rec.ID {
			d.GameLocationRecs[idx] = rec
			return
		}
	}
	d.GameLocationRecs = append(d.GameLocationRecs, rec)
}

func (d *Data) GetGameLocationRecByID(locationID string) (*adventure_game_record.AdventureGameLocation, error) {
	for _, rec := range d.GameLocationRecs {
		if rec.ID == locationID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location with ID >%s<", locationID)
}

func (d *Data) GetGameLocationRecByRef(ref string) (*adventure_game_record.AdventureGameLocation, error) {
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
func (d *Data) AddGameLocationLinkRec(rec *adventure_game_record.AdventureGameLocationLink) {
	for idx := range d.GameLocationLinkRecs {
		if d.GameLocationLinkRecs[idx].ID == rec.ID {
			d.GameLocationLinkRecs[idx] = rec
			return
		}
	}
	d.GameLocationLinkRecs = append(d.GameLocationLinkRecs, rec)
}

func (d *Data) GetGameLocationLinkRecByID(linkID string) (*adventure_game_record.AdventureGameLocationLink, error) {
	for _, rec := range d.GameLocationLinkRecs {
		if rec.ID == linkID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location link with ID >%s<", linkID)
}

func (d *Data) GetGameLocationLinkRecByRef(ref string) (*adventure_game_record.AdventureGameLocationLink, error) {
	linkID, ok := d.Refs.GameLocationLinkRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting location link with ref >%s<", ref)
	}
	return d.GetGameLocationLinkRecByID(linkID)
}

// GameCreature
func (d *Data) AddGameCreatureRec(rec *adventure_game_record.AdventureGameCreature) {
	for idx := range d.GameCreatureRecs {
		if d.GameCreatureRecs[idx].ID == rec.ID {
			d.GameCreatureRecs[idx] = rec
			return
		}
	}
	d.GameCreatureRecs = append(d.GameCreatureRecs, rec)
}

func (d *Data) GetGameCreatureRecByID(creatureID string) (*adventure_game_record.AdventureGameCreature, error) {
	for _, rec := range d.GameCreatureRecs {
		if rec.ID == creatureID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting creature with ID >%s<", creatureID)
}

func (d *Data) GetGameCreatureRecByRef(ref string) (*adventure_game_record.AdventureGameCreature, error) {
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
func (d *Data) AddGameCharacterRec(rec *adventure_game_record.AdventureGameCharacter) {
	for idx := range d.GameCharacterRecs {
		if d.GameCharacterRecs[idx].ID == rec.ID {
			d.GameCharacterRecs[idx] = rec
			return
		}
	}
	d.GameCharacterRecs = append(d.GameCharacterRecs, rec)
}

func (d *Data) GetGameCharacterRecByID(id string) (*adventure_game_record.AdventureGameCharacter, error) {
	for _, rec := range d.GameCharacterRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_character with ID >%s<", id)
}

func (d *Data) GetGameCharacterRecByRef(ref string) (*adventure_game_record.AdventureGameCharacter, error) {
	id, ok := d.Refs.GameCharacterRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_character with ref >%s<", ref)
	}
	return d.GetGameCharacterRecByID(id)
}

// GameItem
func (d *Data) AddGameItemRec(rec *adventure_game_record.AdventureGameItem) {
	for idx := range d.GameItemRecs {
		if d.GameItemRecs[idx].ID == rec.ID {
			d.GameItemRecs[idx] = rec
			return
		}
	}
	d.GameItemRecs = append(d.GameItemRecs, rec)
}

func (d *Data) GetGameItemRecByID(itemID string) (*adventure_game_record.AdventureGameItem, error) {
	for _, rec := range d.GameItemRecs {
		if rec.ID == itemID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item with ID >%s<", itemID)
}

func (d *Data) GetGameItemRecByRef(ref string) (*adventure_game_record.AdventureGameItem, error) {
	id, ok := d.Refs.GameItemRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_item with ref >%s<", ref)
	}
	return d.GetGameItemRecByID(id)
}

// GameLocationLinkRequirement
func (d *Data) AddGameLocationLinkRequirementRec(rec *adventure_game_record.AdventureGameLocationLinkRequirement) {
	for idx := range d.GameLocationLinkRequirementRecs {
		if d.GameLocationLinkRequirementRecs[idx].ID == rec.ID {
			d.GameLocationLinkRequirementRecs[idx] = rec
			return
		}
	}
	d.GameLocationLinkRequirementRecs = append(d.GameLocationLinkRequirementRecs, rec)
}

func (d *Data) GetGameLocationLinkRequirementRecByID(id string) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	for _, rec := range d.GameLocationLinkRequirementRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_link_requirement with ID >%s<", id)
}

func (d *Data) GetGameLocationLinkRequirementRecByRef(ref string) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	id, ok := d.Refs.GameLocationLinkRequirementRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_link_requirement with ref >%s<", ref)
	}
	return d.GetGameLocationLinkRequirementRecByID(id)
}

// GameInstance
func (d *Data) AddGameInstanceRec(rec *game_record.GameInstance) {
	for idx := range d.GameInstanceRecs {
		if d.GameInstanceRecs[idx].ID == rec.ID {
			d.GameInstanceRecs[idx] = rec
			return
		}
	}
	d.GameInstanceRecs = append(d.GameInstanceRecs, rec)
}

func (d *Data) GetGameInstanceRecByID(id string) (*game_record.GameInstance, error) {
	for _, rec := range d.GameInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_instance with ID >%s<", id)
}

func (d *Data) GetGameInstanceRecByRef(ref string) (*game_record.GameInstance, error) {
	id, ok := d.Refs.GameInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_instance with ref >%s<", ref)
	}
	return d.GetGameInstanceRecByID(id)
}

func (d *Data) GetGameInstanceRecByGameID(gameID string) []*game_record.GameInstance {
	var result []*game_record.GameInstance
	for _, rec := range d.GameInstanceRecs {
		if rec.GameID == gameID {
			result = append(result, rec)
		}
	}
	return result
}

// GameLocationInstance
func (d *Data) AddGameLocationInstanceRec(rec *adventure_game_record.AdventureGameLocationInstance) {
	for idx := range d.GameLocationInstanceRecs {
		if d.GameLocationInstanceRecs[idx].ID == rec.ID {
			d.GameLocationInstanceRecs[idx] = rec
			return
		}
	}
	d.GameLocationInstanceRecs = append(d.GameLocationInstanceRecs, rec)
}

func (d *Data) GetGameLocationInstanceRecByID(id string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_location_instance with ID >%s<", id)
}

func (d *Data) GetGameLocationInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	id, ok := d.Refs.GameLocationInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_instance with ref >%s<", ref)
	}
	return d.GetGameLocationInstanceRecByID(id)
}

func (d *Data) GetGameLocationInstanceRecByLocationRef(ref string) (*adventure_game_record.AdventureGameLocationInstance, error) {
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

func (d *Data) GetGameLocationInstanceRecByLocationID(locationID string) []*adventure_game_record.AdventureGameLocationInstance {
	var result []*adventure_game_record.AdventureGameLocationInstance
	for _, rec := range d.GameLocationInstanceRecs {
		if rec.AdventureGameLocationID == locationID {
			result = append(result, rec)
		}
	}
	return result
}

// GameItemInstance
func (d *Data) AddGameItemInstanceRec(rec *adventure_game_record.AdventureGameItemInstance) {
	for idx := range d.GameItemInstanceRecs {
		if d.GameItemInstanceRecs[idx].ID == rec.ID {
			d.GameItemInstanceRecs[idx] = rec
			return
		}
	}
	d.GameItemInstanceRecs = append(d.GameItemInstanceRecs, rec)
}

func (d *Data) GetGameItemInstanceRecByID(id string) (*adventure_game_record.AdventureGameItemInstance, error) {
	for _, rec := range d.GameItemInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_item_instance with ID >%s<", id)
}

func (d *Data) GetGameItemInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameItemInstance, error) {
	id, ok := d.Refs.GameItemInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_item_instance with ref >%s<", ref)
	}
	return d.GetGameItemInstanceRecByID(id)
}

func (d *Data) GetGameItemInstanceRecByItemRef(ref string) (*adventure_game_record.AdventureGameItemInstance, error) {
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

func (d *Data) GetGameItemInstanceRecByItemID(itemID string) []*adventure_game_record.AdventureGameItemInstance {
	var result []*adventure_game_record.AdventureGameItemInstance
	for _, rec := range d.GameItemInstanceRecs {
		if rec.AdventureGameItemID == itemID {
			result = append(result, rec)
		}
	}
	return result
}

// GameCreatureInstance
func (d *Data) AddGameCreatureInstanceRec(rec *adventure_game_record.AdventureGameCreatureInstance) {
	for idx := range d.GameCreatureInstanceRecs {
		if d.GameCreatureInstanceRecs[idx].ID == rec.ID {
			d.GameCreatureInstanceRecs[idx] = rec
			return
		}
	}
	d.GameCreatureInstanceRecs = append(d.GameCreatureInstanceRecs, rec)
}

func (d *Data) GetGameCreatureInstanceRecByID(id string) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	for _, rec := range d.GameCreatureInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_creature_instance with ID >%s<", id)
}

func (d *Data) GetGameCreatureInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	id, ok := d.Refs.GameCreatureInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_creature_instance with ref >%s<", ref)
	}
	return d.GetGameCreatureInstanceRecByID(id)
}

// GameCharacterInstance
func (d *Data) AddGameCharacterInstanceRec(rec *adventure_game_record.AdventureGameCharacterInstance) {
	for idx := range d.GameCharacterInstanceRecs {
		if d.GameCharacterInstanceRecs[idx].ID == rec.ID {
			d.GameCharacterInstanceRecs[idx] = rec
			return
		}
	}
	d.GameCharacterInstanceRecs = append(d.GameCharacterInstanceRecs, rec)
}

func (d *Data) GetGameCharacterInstanceRecByID(id string) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	for _, rec := range d.GameCharacterInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_character_instance with ID >%s<", id)
}

func (d *Data) GetGameCharacterInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	id, ok := d.Refs.GameCharacterInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_character_instance with ref >%s<", ref)
	}
	return d.GetGameCharacterInstanceRecByID(id)
}

func (d *Data) AddGameSubscriptionRec(rec *game_record.GameSubscription) {
	for idx := range d.GameSubscriptionRecs {
		if d.GameSubscriptionRecs[idx].ID == rec.ID {
			d.GameSubscriptionRecs[idx] = rec
			return
		}
	}
	d.GameSubscriptionRecs = append(d.GameSubscriptionRecs, rec)
}

func (d *Data) GetGameSubscriptionRecByID(id string) (*game_record.GameSubscription, error) {
	for _, rec := range d.GameSubscriptionRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_subscription with ID >%s<", id)
}

func (d *Data) AddGameAdministrationRec(rec *game_record.GameAdministration) {
	for idx := range d.GameAdministrationRecs {
		if d.GameAdministrationRecs[idx].ID == rec.ID {
			d.GameAdministrationRecs[idx] = rec
			return
		}
	}
	d.GameAdministrationRecs = append(d.GameAdministrationRecs, rec)
}

func (d *Data) GetGameAdministrationRecByID(id string) (*game_record.GameAdministration, error) {
	for _, rec := range d.GameAdministrationRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("game administration record not found for id >%s<", id)
}

// GameInstanceParameter
func (d *Data) AddGameInstanceParameterRec(rec *game_record.GameInstanceParameter) {
	for idx := range d.GameInstanceParameterRecs {
		if d.GameInstanceParameterRecs[idx].ID == rec.ID {
			d.GameInstanceParameterRecs[idx] = rec
			return
		}
	}
	d.GameInstanceParameterRecs = append(d.GameInstanceParameterRecs, rec)
}

func (d *Data) GetGameInstanceParameterRecByID(id string) (*game_record.GameInstanceParameter, error) {
	for _, rec := range d.GameInstanceParameterRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("game instance parameter record not found for id >%s<", id)
}

func (d *Data) GetGameInstanceParameterRecByRef(ref string) (*game_record.GameInstanceParameter, error) {
	id, exists := d.Refs.GameInstanceParameterRefs[ref]
	if !exists {
		return nil, fmt.Errorf("game instance parameter reference not found for ref >%s<", ref)
	}
	return d.GetGameInstanceParameterRecByID(id)
}
