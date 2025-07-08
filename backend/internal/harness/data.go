package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

// Data -
type Data struct {
	AccountRecs  []*record.Account
	GameRecs     []*record.Game
	LocationRecs []*record.Location
	// Data references
	Refs DataRefs
}

// DataRefs is a collection of maps of data references that were defined
// in test data configuration to the resulting record that was created.
// When adding new reference maps make sure to also initialise the map
// in the initialiseDataStores() function further below.
type DataRefs struct {
	AccountRefs  map[string]string // Map of account refs to account records
	GameRefs     map[string]string // Map of game refs to game records
	LocationRefs map[string]string // Map of location refs to location records
}

// initialiseDataStores - Data is required to maintain data references and
// may contain main test data and reference test data so may not be used
// as a source of teardown data.
func initialiseDataStores() Data {
	return Data{
		Refs: DataRefs{
			AccountRefs:  map[string]string{},
			GameRefs:     map[string]string{},
			LocationRefs: map[string]string{},
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
