package domain

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameCharacterInstanceRecs -
func (m *Domain) GetManyAdventureGameCharacterInstanceRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("GetManyAdventureGameCharacterInstanceRecs")

	l.Debug("getting many adventure_game_character_instance records opts >%#v<", opts)

	r := m.AdventureGameCharacterInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAdventureGameCharacterInstanceRec -
func (m *Domain) GetAdventureGameCharacterInstanceRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("GetAdventureGameCharacterInstanceRec")

	l.Debug("getting adventure_game_character_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AdventureGameCharacterInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameCharacterInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAdventureGameCharacterInstanceRec -
func (m *Domain) CreateAdventureGameCharacterInstanceRec(rec *adventure_game_record.AdventureGameCharacterInstance) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("CreateAdventureGameCharacterInstanceRec")

	l.Debug("creating adventure_game_character_instance record >%#v<", rec)

	if err := m.validateAdventureGameCharacterInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_character_instance record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameCharacterInstanceRepository()

	var err error
	createdRec, err := r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return createdRec, nil
}

// UpdateAdventureGameCharacterInstanceRec -
func (m *Domain) UpdateAdventureGameCharacterInstanceRec(rec *adventure_game_record.AdventureGameCharacterInstance) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("UpdateAdventureGameCharacterInstanceRec")

	_, err := m.GetAdventureGameCharacterInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_character_instance record >%#v<", rec)

	if err := m.validateAdventureGameCharacterInstanceRecForUpdate(rec); err != nil {
		l.Warn("failed to validate adventure_game_character_instance record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameCharacterInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameCharacterInstanceRec -
func (m *Domain) DeleteAdventureGameCharacterInstanceRec(recID string) error {
	l := m.Logger("DeleteAdventureGameCharacterInstanceRec")

	l.Debug("deleting adventure_game_character_instance record ID >%s<", recID)

	_, err := m.GetAdventureGameCharacterInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameCharacterInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAdventureGameCharacterInstanceRec -
func (m *Domain) RemoveAdventureGameCharacterInstanceRec(recID string) error {
	l := m.Logger("RemoveAdventureGameCharacterInstanceRec")

	l.Debug("removing adventure_game_character_instance record ID >%s<", recID)

	_, err := m.GetAdventureGameCharacterInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameCharacterInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// AssignStartingItemsToCharacterInstance assigns all starting items for a game to a character instance
func (m *Domain) AssignStartingItemsToCharacterInstance(characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance) error {
	l := m.Logger("AssignStartingItemsToCharacterInstance")

	l.Info("assigning starting items to character instance >%s< for game >%s<", characterInstanceRec.ID, characterInstanceRec.GameID)

	// Get all starting items for this game
	startingItems, err := m.GetManyAdventureGameItemRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemGameID, Val: characterInstanceRec.GameID},
			{Col: adventure_game_record.FieldAdventureGameItemIsStartingItem, Val: true},
		},
	})
	if err != nil {
		l.Warn("failed to get starting items for game >%s< >%v<", characterInstanceRec.GameID, err)
		return err
	}

	if len(startingItems) == 0 {
		l.Info("no starting items found for game >%s<", characterInstanceRec.GameID)
		return nil
	}

	l.Info("found >%d< starting items for game >%s<", len(startingItems), characterInstanceRec.GameID)

	// Create item instances for each starting item
	for _, item := range startingItems {
		itemInstance := &adventure_game_record.AdventureGameItemInstance{
			GameID:                           characterInstanceRec.GameID,
			GameInstanceID:                   characterInstanceRec.GameInstanceID,
			AdventureGameItemID:              item.ID,
			AdventureGameCharacterInstanceID:  nullstring.FromString(characterInstanceRec.ID),
			AdventureGameLocationInstanceID:  sql.NullString{Valid: false},
			AdventureGameCreatureInstanceID:   sql.NullString{Valid: false},
			IsEquipped:                       false,
			EquipmentSlot:                    sql.NullString{Valid: false},
			IsUsed:                           false,
			UsesRemaining:                    0,
		}

		_, err := m.CreateAdventureGameItemInstanceRec(itemInstance)
		if err != nil {
			l.Warn("failed to create item instance for starting item >%s< >%v<", item.ID, err)
			return err
		}

		l.Info("assigned starting item >%s< (%s) to character instance >%s<", item.ID, item.Name, characterInstanceRec.ID)
	}

	l.Info("successfully assigned >%d< starting items to character instance >%s<", len(startingItems), characterInstanceRec.ID)

	return nil
}
