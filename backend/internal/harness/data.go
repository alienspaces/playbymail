package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// Data -
type Data struct {
	AccountRecs               []*account_record.Account
	AccountContactRecs        []*account_record.AccountContact
	GameRecs                  []*game_record.Game
	GameSubscriptionRecs      []*game_record.GameSubscription
	GameAdministrationRecs    []*game_record.GameAdministration
	GameInstanceRecs          []*game_record.GameInstance
	GameInstanceParameterRecs []*game_record.GameInstanceParameter
	GameTurnSheetRecs         []*game_record.GameTurnSheet
	// Adventure game specific resources
	AdventureGameLocationRecs                []*adventure_game_record.AdventureGameLocation
	AdventureGameLocationLinkRecs            []*adventure_game_record.AdventureGameLocationLink
	AdventureGameCharacterRecs               []*adventure_game_record.AdventureGameCharacter
	AdventureGameCreatureRecs                []*adventure_game_record.AdventureGameCreature
	AdventureGameItemRecs                    []*adventure_game_record.AdventureGameItem
	AdventureGameLocationLinkRequirementRecs []*adventure_game_record.AdventureGameLocationLinkRequirement
	AdventureGameLocationInstanceRecs        []*adventure_game_record.AdventureGameLocationInstance
	AdventureGameItemInstanceRecs            []*adventure_game_record.AdventureGameItemInstance
	AdventureGameCreatureInstanceRecs        []*adventure_game_record.AdventureGameCreatureInstance
	AdventureGameCharacterInstanceRecs       []*adventure_game_record.AdventureGameCharacterInstance
	AdventureGameTurnSheetRecs               []*adventure_game_record.AdventureGameTurnSheet
	// Data references
	Refs DataRefs
}

// DataRefs is a collection of maps of data references that were defined
// in test data configuration to the resulting record that was created.
// When adding new reference maps make sure to also initialise the map
// in the initialiseDataStores() function further below.
type DataRefs struct {
	AccountRefs               map[string]string // Map of account refs to account records
	GameRefs                  map[string]string // Map of game refs to game records
	GameSubscriptionRefs      map[string]string // Map of refs to game_subscription records
	GameAdministrationRefs    map[string]string // Map of refs to game_administration records
	GameInstanceRefs          map[string]string // Map of refs to game_instance records
	GameInstanceParameterRefs map[string]string // Map of refs to game_instance_parameter records
	GameTurnSheetRefs         map[string]string // Map of refs to game_turn_sheet records
	// Adventure game specific resources
	AdventureGameLocationRefs                map[string]string // Map of adventure game location refs to adventure game location records
	AdventureGameLocationLinkRefs            map[string]string // Map of adventure game location link refs to adventure game location link records
	AdventureGameLocationLinkRequirementRefs map[string]string // Map of adventure game location link requirement refs to adventure game location link requirement records
	AdventureGameLocationInstanceRefs        map[string]string // Map of adventure game location instance refs to adventure game location instance records
	AdventureGameCharacterRefs               map[string]string // Map of adventure game character refs to adventure game character records
	AdventureGameCreatureRefs                map[string]string // Map of adventure game creature refs to adventure game creature records
	AdventureGameItemRefs                    map[string]string // Map of adventure game item refs to adventure game item records
	AdventureGameItemInstanceRefs            map[string]string // Map of adventure game item instance refs to adventure game item instance records
	AdventureGameCreatureInstanceRefs        map[string]string // Map of adventure game creature instance refs to adventure game creature instance records
	AdventureGameCharacterInstanceRefs       map[string]string // Map of adventure game character instance refs to adventure game character instance records
}

// initialiseDataStores - Data is required to maintain data references and
// may contain main test data and reference test data so may not be used
// as a source of teardown data.
func initialiseDataStores() Data {
	return Data{
		Refs: DataRefs{
			AccountRefs:               map[string]string{},
			GameRefs:                  map[string]string{},
			GameSubscriptionRefs:      map[string]string{},
			GameAdministrationRefs:    map[string]string{},
			GameInstanceRefs:          map[string]string{},
			GameInstanceParameterRefs: map[string]string{},
			GameTurnSheetRefs:         map[string]string{},
			// Adventure game specific resources
			AdventureGameLocationRefs:                map[string]string{},
			AdventureGameLocationLinkRefs:            map[string]string{},
			AdventureGameLocationLinkRequirementRefs: map[string]string{},
			AdventureGameLocationInstanceRefs:        map[string]string{},
			AdventureGameCharacterRefs:               map[string]string{},
			AdventureGameCreatureRefs:                map[string]string{},
			AdventureGameItemRefs:                    map[string]string{},
			AdventureGameItemInstanceRefs:            map[string]string{},
			AdventureGameCreatureInstanceRefs:        map[string]string{},
			AdventureGameCharacterInstanceRefs:       map[string]string{},
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

func (d *Data) AddAccountContactRec(rec *account_record.AccountContact) {
	for idx := range d.AccountContactRecs {
		if d.AccountContactRecs[idx].ID == rec.ID {
			d.AccountContactRecs[idx] = rec
			return
		}
	}
	d.AccountContactRecs = append(d.AccountContactRecs, rec)
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

func (d *Data) GetAccountContactRecByAccountID(accountID string) (*account_record.AccountContact, error) {
	for _, rec := range d.AccountContactRecs {
		if rec.AccountID == accountID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting account contact with account ID >%s<", accountID)
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

// GameTurnSheet
func (d *Data) AddGameTurnSheetRec(rec *game_record.GameTurnSheet) {
	for idx := range d.GameTurnSheetRecs {
		if d.GameTurnSheetRecs[idx].ID == rec.ID {
			d.GameTurnSheetRecs[idx] = rec
			return
		}
	}
	d.GameTurnSheetRecs = append(d.GameTurnSheetRecs, rec)
}

func (d *Data) GetGameTurnSheetRecByID(id string) (*game_record.GameTurnSheet, error) {
	for _, rec := range d.GameTurnSheetRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_turn_sheet with ID >%s<", id)
}

func (d *Data) GetGameTurnSheetRecByRef(ref string) (*game_record.GameTurnSheet, error) {
	id, exists := d.Refs.GameTurnSheetRefs[ref]
	if !exists {
		return nil, fmt.Errorf("failed getting game_turn_sheet with ref >%s<", ref)
	}
	return d.GetGameTurnSheetRecByID(id)
}

// ------------------------------------------------------------
// Adventure game specific resources
// ------------------------------------------------------------

// AdventureGameLocation
func (d *Data) AddAdventureGameLocationRec(rec *adventure_game_record.AdventureGameLocation) {
	for idx := range d.AdventureGameLocationRecs {
		if d.AdventureGameLocationRecs[idx].ID == rec.ID {
			d.AdventureGameLocationRecs[idx] = rec
			return
		}
	}
	d.AdventureGameLocationRecs = append(d.AdventureGameLocationRecs, rec)
}

func (d *Data) GetAdventureGameLocationRecByID(locationID string) (*adventure_game_record.AdventureGameLocation, error) {
	for _, rec := range d.AdventureGameLocationRecs {
		if rec.ID == locationID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location with ID >%s<", locationID)
}

func (d *Data) GetAdventureGameLocationRecByRef(ref string) (*adventure_game_record.AdventureGameLocation, error) {
	id, ok := d.Refs.AdventureGameLocationRefs[ref]
	if !ok {
		return nil, fmt.Errorf("no adventure game location with ref >%s<", ref)
	}
	for _, rec := range d.AdventureGameLocationRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("no adventure game location with id >%s< for ref >%s<", id, ref)
}

// AdventureGameLocationLink
func (d *Data) AddAdventureGameLocationLinkRec(rec *adventure_game_record.AdventureGameLocationLink) {
	for idx := range d.AdventureGameLocationLinkRecs {
		if d.AdventureGameLocationLinkRecs[idx].ID == rec.ID {
			d.AdventureGameLocationLinkRecs[idx] = rec
			return
		}
	}
	d.AdventureGameLocationLinkRecs = append(d.AdventureGameLocationLinkRecs, rec)
}

func (d *Data) GetAdventureGameLocationLinkRecByID(linkID string) (*adventure_game_record.AdventureGameLocationLink, error) {
	for _, rec := range d.AdventureGameLocationLinkRecs {
		if rec.ID == linkID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game location link with ID >%s<", linkID)
}

func (d *Data) GetAdventureGameLocationLinkRecByRef(ref string) (*adventure_game_record.AdventureGameLocationLink, error) {
	linkID, ok := d.Refs.AdventureGameLocationLinkRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game location link with ref >%s<", ref)
	}
	return d.GetAdventureGameLocationLinkRecByID(linkID)
}

// AdventureGameCreature
func (d *Data) AddAdventureGameCreatureRec(rec *adventure_game_record.AdventureGameCreature) {
	for idx := range d.AdventureGameCreatureRecs {
		if d.AdventureGameCreatureRecs[idx].ID == rec.ID {
			d.AdventureGameCreatureRecs[idx] = rec
			return
		}
	}
	d.AdventureGameCreatureRecs = append(d.AdventureGameCreatureRecs, rec)
}

func (d *Data) GetAdventureGameCreatureRecByID(creatureID string) (*adventure_game_record.AdventureGameCreature, error) {
	for _, rec := range d.AdventureGameCreatureRecs {
		if rec.ID == creatureID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game creature with ID >%s<", creatureID)
}

func (d *Data) GetAdventureGameCreatureRecByRef(ref string) (*adventure_game_record.AdventureGameCreature, error) {
	id, ok := d.Refs.AdventureGameCreatureRefs[ref]
	if !ok {
		return nil, fmt.Errorf("no adventure game creature with ref >%s<", ref)
	}
	for _, rec := range d.AdventureGameCreatureRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("no adventure game creature with id >%s< for ref >%s<", id, ref)
}

// AdventureGameCharacter
func (d *Data) AddAdventureGameCharacterRec(rec *adventure_game_record.AdventureGameCharacter) {
	for idx := range d.AdventureGameCharacterRecs {
		if d.AdventureGameCharacterRecs[idx].ID == rec.ID {
			d.AdventureGameCharacterRecs[idx] = rec
			return
		}
	}
	d.AdventureGameCharacterRecs = append(d.AdventureGameCharacterRecs, rec)
}

func (d *Data) GetAdventureGameCharacterRecByID(id string) (*adventure_game_record.AdventureGameCharacter, error) {
	for _, rec := range d.AdventureGameCharacterRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game character with ID >%s<", id)
}

func (d *Data) GetAdventureGameCharacterRecByRef(ref string) (*adventure_game_record.AdventureGameCharacter, error) {
	id, ok := d.Refs.AdventureGameCharacterRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game character with ref >%s<", ref)
	}
	return d.GetAdventureGameCharacterRecByID(id)
}

// AdventureGameItem
func (d *Data) AddAdventureGameItemRec(rec *adventure_game_record.AdventureGameItem) {
	for idx := range d.AdventureGameItemRecs {
		if d.AdventureGameItemRecs[idx].ID == rec.ID {
			d.AdventureGameItemRecs[idx] = rec
			return
		}
	}
	d.AdventureGameItemRecs = append(d.AdventureGameItemRecs, rec)
}

func (d *Data) GetAdventureGameItemRecByID(itemID string) (*adventure_game_record.AdventureGameItem, error) {
	for _, rec := range d.AdventureGameItemRecs {
		if rec.ID == itemID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game item with ID >%s<", itemID)
}

func (d *Data) GetAdventureGameItemRecByRef(ref string) (*adventure_game_record.AdventureGameItem, error) {
	id, ok := d.Refs.AdventureGameItemRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game item with ref >%s<", ref)
	}
	return d.GetAdventureGameItemRecByID(id)
}

// AdventureGameLocationLinkRequirement
func (d *Data) AddAdventureGameLocationLinkRequirementRec(rec *adventure_game_record.AdventureGameLocationLinkRequirement) {
	for idx := range d.AdventureGameLocationLinkRequirementRecs {
		if d.AdventureGameLocationLinkRequirementRecs[idx].ID == rec.ID {
			d.AdventureGameLocationLinkRequirementRecs[idx] = rec
			return
		}
	}
	d.AdventureGameLocationLinkRequirementRecs = append(d.AdventureGameLocationLinkRequirementRecs, rec)
}

func (d *Data) GetAdventureGameLocationLinkRequirementRecByID(id string) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	for _, rec := range d.AdventureGameLocationLinkRequirementRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game location link requirement with ID >%s<", id)
}

func (d *Data) GetAdventureGameLocationLinkRequirementRecByRef(ref string) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	id, ok := d.Refs.AdventureGameLocationLinkRequirementRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game location link requirement with ref >%s<", ref)
	}
	return d.GetAdventureGameLocationLinkRequirementRecByID(id)
}

// AdventureGameLocationInstance
func (d *Data) AddAdventureGameLocationInstanceRec(rec *adventure_game_record.AdventureGameLocationInstance) {
	for idx := range d.AdventureGameLocationInstanceRecs {
		if d.AdventureGameLocationInstanceRecs[idx].ID == rec.ID {
			d.AdventureGameLocationInstanceRecs[idx] = rec
			return
		}
	}
	d.AdventureGameLocationInstanceRecs = append(d.AdventureGameLocationInstanceRecs, rec)
}

func (d *Data) GetAdventureGameLocationInstanceRecByID(id string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	for _, rec := range d.AdventureGameLocationInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game location instance with ID >%s<", id)
}

func (d *Data) GetAdventureGameLocationInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	id, ok := d.Refs.AdventureGameLocationInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_location_instance with ref >%s<", ref)
	}
	return d.GetAdventureGameLocationInstanceRecByID(id)
}

func (d *Data) GetAdventureGameLocationInstanceRecByLocationRef(ref string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	id, ok := d.Refs.AdventureGameLocationRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game location instance with ref >%s<", ref)
	}
	for _, rec := range d.AdventureGameLocationInstanceRecs {
		if rec.AdventureGameLocationID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game location instance with location ID >%s<", id)
}

func (d *Data) GetAdventureGameLocationInstanceRecByLocationID(locationID string) []*adventure_game_record.AdventureGameLocationInstance {
	var result []*adventure_game_record.AdventureGameLocationInstance
	for _, rec := range d.AdventureGameLocationInstanceRecs {
		if rec.AdventureGameLocationID == locationID {
			result = append(result, rec)
		}
	}
	return result
}

// AdventureGameItemInstance
func (d *Data) AddAdventureGameItemInstanceRec(rec *adventure_game_record.AdventureGameItemInstance) {
	for idx := range d.AdventureGameItemInstanceRecs {
		if d.AdventureGameItemInstanceRecs[idx].ID == rec.ID {
			d.AdventureGameItemInstanceRecs[idx] = rec
			return
		}
	}
	d.AdventureGameItemInstanceRecs = append(d.AdventureGameItemInstanceRecs, rec)
}

func (d *Data) GetAdventureGameItemInstanceRecByID(id string) (*adventure_game_record.AdventureGameItemInstance, error) {
	for _, rec := range d.AdventureGameItemInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game item instance with ID >%s<", id)
}

func (d *Data) GetAdventureGameItemInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameItemInstance, error) {
	id, ok := d.Refs.AdventureGameItemInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game item instance with ref >%s<", ref)
	}
	return d.GetAdventureGameItemInstanceRecByID(id)
}

func (d *Data) GetAdventureGameItemInstanceRecByItemRef(ref string) (*adventure_game_record.AdventureGameItemInstance, error) {
	id, ok := d.Refs.AdventureGameItemRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game item instance with ref >%s<", ref)
	}
	for _, rec := range d.AdventureGameItemInstanceRecs {
		if rec.AdventureGameItemID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game item instance with item ID >%s<", id)
}

func (d *Data) GetAdventureGameItemInstanceRecByItemID(itemID string) []*adventure_game_record.AdventureGameItemInstance {
	var result []*adventure_game_record.AdventureGameItemInstance
	for _, rec := range d.AdventureGameItemInstanceRecs {
		if rec.AdventureGameItemID == itemID {
			result = append(result, rec)
		}
	}
	return result
}

// AdventureGameCreatureInstance
func (d *Data) AddAdventureGameCreatureInstanceRec(rec *adventure_game_record.AdventureGameCreatureInstance) {
	for idx := range d.AdventureGameCreatureInstanceRecs {
		if d.AdventureGameCreatureInstanceRecs[idx].ID == rec.ID {
			d.AdventureGameCreatureInstanceRecs[idx] = rec
			return
		}
	}
	d.AdventureGameCreatureInstanceRecs = append(d.AdventureGameCreatureInstanceRecs, rec)
}

func (d *Data) GetAdventureGameCreatureInstanceRecByID(id string) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	for _, rec := range d.AdventureGameCreatureInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game creature instance with ID >%s<", id)
}

func (d *Data) GetAdventureGameCreatureInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	id, ok := d.Refs.AdventureGameCreatureInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game creature instance with ref >%s<", ref)
	}
	return d.GetAdventureGameCreatureInstanceRecByID(id)
}

// AdventureGameCharacterInstance
func (d *Data) AddAdventureGameCharacterInstanceRec(rec *adventure_game_record.AdventureGameCharacterInstance) {
	for idx := range d.AdventureGameCharacterInstanceRecs {
		if d.AdventureGameCharacterInstanceRecs[idx].ID == rec.ID {
			d.AdventureGameCharacterInstanceRecs[idx] = rec
			return
		}
	}
	d.AdventureGameCharacterInstanceRecs = append(d.AdventureGameCharacterInstanceRecs, rec)
}

func (d *Data) GetAdventureGameCharacterInstanceRecByID(id string) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	for _, rec := range d.AdventureGameCharacterInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting adventure game character instance with ID >%s<", id)
}

func (d *Data) GetAdventureGameCharacterInstanceRecByRef(ref string) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	id, ok := d.Refs.AdventureGameCharacterInstanceRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting adventure game character instance with ref >%s<", ref)
	}
	return d.GetAdventureGameCharacterInstanceRecByID(id)
}

// AdventureGameTurnSheet
func (d *Data) AddAdventureGameTurnSheetRec(rec *adventure_game_record.AdventureGameTurnSheet) {
	for idx := range d.AdventureGameTurnSheetRecs {
		if d.AdventureGameTurnSheetRecs[idx].ID == rec.ID {
			d.AdventureGameTurnSheetRecs[idx] = rec
			return
		}
	}
	d.AdventureGameTurnSheetRecs = append(d.AdventureGameTurnSheetRecs, rec)
}
