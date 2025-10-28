package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// MoveCharacterToLocation moves a character to a specific location in the game
func (t *Testing) MoveCharacterToLocation(characterInstanceRef, locationInstanceRef string) error {
	l := t.Logger("MoveCharacterToLocation")

	characterInstance, err := t.Data.GetGameCharacterInstanceRecByRef(characterInstanceRef)
	if err != nil {
		l.Warn("failed to get character instance >%s< >%v<", characterInstanceRef, err)
		return err
	}

	locationInstance, err := t.Data.GetGameLocationInstanceRecByRef(locationInstanceRef)
	if err != nil {
		l.Warn("failed to get location instance >%s< >%v<", locationInstanceRef, err)
		return err
	}

	// Update character's location
	characterInstance.AdventureGameLocationInstanceID = locationInstance.ID

	_, err = t.Domain.(*domain.Domain).UpdateAdventureGameCharacterInstanceRec(characterInstance)
	if err != nil {
		l.Warn("failed to move character >%s< to location >%s< >%v<",
			characterInstanceRef, locationInstanceRef, err)
		return err
	}

	l.Info("moved character >%s< to location >%s<", characterInstanceRef, locationInstanceRef)

	return nil
}

// SetupCharacterAtLocation creates or updates a character instance at a specific location
func (t *Testing) SetupCharacterAtLocation(characterInstanceRef, locationInstanceRef string) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := t.Logger("SetupCharacterAtLocation")

	// Try to get existing character instance
	characterInstance, err := t.Data.GetGameCharacterInstanceRecByRef(characterInstanceRef)
	if err != nil {
		// Character doesn't exist yet, can't set it up
		return nil, fmt.Errorf("character instance >%s< does not exist", characterInstanceRef)
	}

	// Move to location
	err = t.MoveCharacterToLocation(characterInstanceRef, locationInstanceRef)
	if err != nil {
		l.Warn("failed to move character to location >%v<", err)
		return nil, err
	}

	// Get updated record
	characterInstance, err = t.Data.GetGameCharacterInstanceRecByRef(characterInstanceRef)
	if err != nil {
		return nil, err
	}

	return characterInstance, nil
}
