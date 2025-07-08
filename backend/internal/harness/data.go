package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

// Data -
type Data struct {
	AccountRecs       []*record.Account
	GameRecs          []*record.Game
	LocationRecs      []*record.Location
	LocationLinkRecs  []*record.LocationLink
	GameCharacterRecs []*record.GameCharacter // Add this line
	// Data references
	Refs DataRefs
}

// DataRefs is a collection of maps of data references that were defined
// in test data configuration to the resulting record that was created.
// When adding new reference maps make sure to also initialise the map
// in the initialiseDataStores() function further below.
type DataRefs struct {
	AccountRefs       map[string]string // Map of account refs to account records
	GameRefs          map[string]string // Map of game refs to game records
	LocationRefs      map[string]string // Map of location refs to location records
	LocationLinkRefs  map[string]string // Map of location link refs to location link records
	GameCharacterRefs map[string]string // Map of game_character refs to records
}

// initialiseDataStores - Data is required to maintain data references and
// may contain main test data and reference test data so may not be used
// as a source of teardown data.
func initialiseDataStores() Data {
	return Data{
		Refs: DataRefs{
			AccountRefs:       map[string]string{},
			GameRefs:          map[string]string{},
			LocationRefs:      map[string]string{},
			LocationLinkRefs:  map[string]string{},
			GameCharacterRefs: map[string]string{}, // Add this line
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

// Location
func (d *Data) AddLocationRec(rec *record.Location) {
	for idx := range d.LocationRecs {
		if d.LocationRecs[idx].ID == rec.ID {
			d.LocationRecs[idx] = rec
			return
		}
	}
	d.LocationRecs = append(d.LocationRecs, rec)
}

func (d *Data) GetLocationRecByID(locationID string) (*record.Location, error) {
	for _, rec := range d.LocationRecs {
		if rec.ID == locationID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location with ID >%s<", locationID)
}

func (d *Data) GetLocationRecByRef(ref string) (*record.Location, error) {
	id, ok := d.Refs.LocationRefs[ref]
	if !ok {
		return nil, fmt.Errorf("no location with ref >%s<", ref)
	}
	for _, rec := range d.LocationRecs {
		if rec.ID == id {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("no location with id >%s< for ref >%s<", id, ref)
}

// LocationLink
func (d *Data) AddLocationLinkRec(rec *record.LocationLink) {
	for idx := range d.LocationLinkRecs {
		if d.LocationLinkRecs[idx].ID == rec.ID {
			d.LocationLinkRecs[idx] = rec
			return
		}
	}
	d.LocationLinkRecs = append(d.LocationLinkRecs, rec)
}

func (d *Data) GetLocationLinkRecByID(linkID string) (*record.LocationLink, error) {
	for _, rec := range d.LocationLinkRecs {
		if rec.ID == linkID {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("failed getting location link with ID >%s<", linkID)
}

func (d *Data) GetLocationLinkRecByRef(ref string) (*record.LocationLink, error) {
	linkID, ok := d.Refs.LocationLinkRefs[ref]
	if !ok {
		return nil, fmt.Errorf("failed getting location link with ref >%s<", ref)
	}
	return d.GetLocationLinkRecByID(linkID)
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
