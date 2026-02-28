package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// Data -
type Data struct {
	AccountRecs                  []*account_record.Account
	AccountUserRecs              []*account_record.AccountUser
	AccountUserContactRecs       []*account_record.AccountUserContact
	AccountSubscriptionRecs      []*account_record.AccountSubscription
	GameRecs                     []*game_record.Game
	GameImageRecs                []*game_record.GameImage
	GameSubscriptionRecs         []*game_record.GameSubscription
	GameSubscriptionInstanceRecs []*game_record.GameSubscriptionInstance
	GameInstanceRecs             []*game_record.GameInstance
	GameInstanceParameterRecs    []*game_record.GameInstanceParameter
	GameTurnSheetRecs            []*game_record.GameTurnSheet
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
	// Session tokens by account ID
	AccountSessionTokens map[string]string
	// Data references
	Refs DataRefs
}

// DataRefs is a collection of maps of data references that were defined
// in test data configuration to the resulting record that was created.
// When adding new reference maps make sure to also initialise the map
// in the initialiseDataStores() function further below.
type DataRefs struct {
	AccountUserRefs           map[string]string // Map of account user refs to account user records
	AccountSubscriptionRefs   map[string]string // Map of account subscription refs to account subscription records
	GameRefs                  map[string]string // Map of game refs to game records
	GameImageRefs             map[string]string // Map of refs to game_image records
	GameSubscriptionRefs      map[string]string // Map of refs to game_subscription records
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
		AccountSessionTokens: map[string]string{},
		Refs: DataRefs{
			AccountUserRefs:           map[string]string{},
			AccountSubscriptionRefs:   map[string]string{},
			GameRefs:                  map[string]string{},
			GameImageRefs:             map[string]string{},
			GameSubscriptionRefs:      map[string]string{},
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

func (d *Data) GetAccountRecByID(accountID string) (*account_record.Account, error) {
	for _, rec := range d.AccountRecs {
		if rec.ID == accountID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting account with ID >%s<", accountID)
}

// AccountUser
func (d *Data) AddAccountUserRec(rec *account_record.AccountUser) {
	for idx := range d.AccountUserRecs {
		if d.AccountUserRecs[idx].ID == rec.ID {
			d.AccountUserRecs[idx] = rec
			return
		}
	}
	d.AccountUserRecs = append(d.AccountUserRecs, rec)
}

func (d *Data) GetAccountUserRecByID(accountUserID string) (*account_record.AccountUser, error) {
	for _, rec := range d.AccountUserRecs {
		if rec.ID == accountUserID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting account user with ID >%s<", accountUserID)
}

func (d *Data) GetAccountUserRecByRef(ref string) (*account_record.AccountUser, error) {
	accountUserID, ok := d.Refs.AccountUserRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting account user with ref >%s<", ref)
	}
	return d.GetAccountUserRecByID(accountUserID)
}

func (d *Data) AddAccountUserContactRec(rec *account_record.AccountUserContact) {
	for idx := range d.AccountUserContactRecs {
		if d.AccountUserContactRecs[idx].ID == rec.ID {
			d.AccountUserContactRecs[idx] = rec
			return
		}
	}
	d.AccountUserContactRecs = append(d.AccountUserContactRecs, rec)
}


// AddAccountSessionToken stores a session token for an account by account ID
func (d *Data) AddAccountSessionToken(accountID, sessionToken string) {
	if d.AccountSessionTokens == nil {
		d.AccountSessionTokens = make(map[string]string)
	}
	d.AccountSessionTokens[accountID] = sessionToken
}

// GetAccountSessionToken returns the session token for an account by account ID
func (d *Data) GetAccountSessionToken(accountID string) (string, error) {
	if d.AccountSessionTokens == nil {
		return "", fmt.Errorf("failed getting session token for account ID >%s< (session tokens map not initialized)", accountID)
	}
	token, ok := d.AccountSessionTokens[accountID]
	if !ok {
		return "", fmt.Errorf("failed getting session token for account ID >%s<", accountID)
	}
	return token, nil
}

// GetAccountSessionTokenByAccountRef returns the session token for an account user by reference
func (d *Data) GetAccountSessionTokenByAccountRef(ref string) (string, error) {
	accountUserID, ok := d.Refs.AccountUserRefs[ref]
	if !ok {
		return "", fmt.Errorf("failed getting account user with ref >%s<", ref)
	}
	return d.GetAccountSessionToken(accountUserID)
}

func (d *Data) GetAccountUserContactRecByAccountUserID(accountUserID string) (*account_record.AccountUserContact, error) {
	for _, rec := range d.AccountUserContactRecs {
		if rec.AccountUserID == accountUserID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting account contact with account user ID >%s<", accountUserID)
}

func (d *Data) AddAccountSubscriptionRec(rec *account_record.AccountSubscription) {
	for idx := range d.AccountSubscriptionRecs {
		if d.AccountSubscriptionRecs[idx].ID == rec.ID {
			d.AccountSubscriptionRecs[idx] = rec
			return
		}
	}
	d.AccountSubscriptionRecs = append(d.AccountSubscriptionRecs, rec)
}

func (d *Data) GetAccountSubscriptionRecByID(id string) (*account_record.AccountSubscription, error) {
	for _, rec := range d.AccountSubscriptionRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting account_subscription with ID >%s<", id)
}

func (d *Data) GetAccountSubscriptionRecByRef(ref string) (*account_record.AccountSubscription, error) {
	subscriptionID, ok := d.Refs.AccountSubscriptionRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting account_subscription with ref >%s<", ref)
	}
	return d.GetAccountSubscriptionRecByID(subscriptionID)
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

// GameImage
func (d *Data) AddGameImageRec(rec *game_record.GameImage) {
	for idx := range d.GameImageRecs {
		if d.GameImageRecs[idx].ID == rec.ID {
			d.GameImageRecs[idx] = rec
			return
		}
	}
	d.GameImageRecs = append(d.GameImageRecs, rec)
}

func (d *Data) GetGameImageRecByID(imageID string) (*game_record.GameImage, error) {
	for _, rec := range d.GameImageRecs {
		if rec.ID == imageID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game image with ID >%s<", imageID)
}

func (d *Data) GetGameImageRecByRef(ref string) (*game_record.GameImage, error) {
	imageID, ok := d.Refs.GameImageRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game image with ref >%s<", ref)
	}
	return d.GetGameImageRecByID(imageID)
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

func (d *Data) GetGameSubscriptionRecByRef(ref string) (*game_record.GameSubscription, error) {
	subscriptionID, ok := d.Refs.GameSubscriptionRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting game_subscription with ref >%s<", ref)
	}
	return d.GetGameSubscriptionRecByID(subscriptionID)
}

// GameSubscriptionInstance
func (d *Data) AddGameSubscriptionInstanceRec(rec *game_record.GameSubscriptionInstance) {
	for idx := range d.GameSubscriptionInstanceRecs {
		if d.GameSubscriptionInstanceRecs[idx].ID == rec.ID {
			d.GameSubscriptionInstanceRecs[idx] = rec
			return
		}
	}
	d.GameSubscriptionInstanceRecs = append(d.GameSubscriptionInstanceRecs, rec)
}

func (d *Data) GetGameSubscriptionInstanceRecByID(id string) (*game_record.GameSubscriptionInstance, error) {
	for _, rec := range d.GameSubscriptionInstanceRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting game_subscription_instance with ID >%s<", id)
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

// GetAdventureGameLocationInstanceRecByGameInstanceAndLocationRef gets a location instance
// by game instance reference and location reference. This is useful when location instances
// are auto-generated and don't have explicit references.
func (d *Data) GetAdventureGameLocationInstanceRecByGameInstanceAndLocationRef(gameInstanceRef, locationRef string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	// Get game instance
	gameInstanceID, ok := d.Refs.GameInstanceRefs[gameInstanceRef]
	if !ok {
		return nil, fmt.Errorf("failed getting game instance with ref >%s<", gameInstanceRef)
	}

	// Get location
	locationID, ok := d.Refs.AdventureGameLocationRefs[locationRef]
	if !ok {
		return nil, fmt.Errorf("failed getting game location with ref >%s<", locationRef)
	}

	// Find location instance that matches both
	for _, rec := range d.AdventureGameLocationInstanceRecs {
		if rec.GameInstanceID == gameInstanceID && rec.AdventureGameLocationID == locationID {
			return rec, nil
		}
	}

	return nil, fmt.Errorf("failed getting game_location_instance with game_instance ref >%s< and location ref >%s<", gameInstanceRef, locationRef)
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
