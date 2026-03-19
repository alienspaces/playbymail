package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameLocationObjectInstanceRec(
	objectRec *adventure_game_record.AdventureGameLocationObject,
	locationInstanceRec *adventure_game_record.AdventureGameLocationInstance,
	gameInstanceRec *game_record.GameInstance,
) (*adventure_game_record.AdventureGameLocationObjectInstance, error) {
	l := t.Logger("createAdventureGameLocationObjectInstanceRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil")
	}
	if objectRec == nil {
		return nil, fmt.Errorf("object record is nil")
	}
	if locationInstanceRec == nil {
		return nil, fmt.Errorf("location instance record is nil")
	}

	rec := &adventure_game_record.AdventureGameLocationObjectInstance{
		GameID:                          gameInstanceRec.GameID,
		GameInstanceID:                  gameInstanceRec.ID,
		AdventureGameLocationObjectID:   objectRec.ID,
		AdventureGameLocationInstanceID: locationInstanceRec.ID,
		CurrentState:                    objectRec.InitialState,
		IsVisible:                       !objectRec.IsHidden,
	}

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationObjectInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating adventure_game_location_object_instance record >%v<", err)
		return nil, err
	}

	t.Data.AddAdventureGameLocationObjectInstanceRec(createdRec)
	t.teardownData.AddAdventureGameLocationObjectInstanceRec(createdRec)

	l.Debug("created adventure_game_location_object_instance record ID >%s< for object >%s< at location instance >%s<",
		createdRec.ID, objectRec.ID, locationInstanceRec.ID)

	return createdRec, nil
}
