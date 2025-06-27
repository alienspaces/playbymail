package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

// Data -
type Data struct {
	GameRecs []*record.Game
	// Data references
	Refs DataRefs
}

// DataRefs is a collection of maps of data references that were defined
// in test data configuration to the resulting record that was created.
// When adding new reference maps make sure to also initialise the map
// in the initialiseDataStores() function further below.
type DataRefs struct {
	GameRefs map[string]string // Map of game refs to game records
}

func initialiseDataStores() Data {
	return Data{
		Refs: DataRefs{
			GameRefs: map[string]string{},
		},
	}
}

// Example
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
