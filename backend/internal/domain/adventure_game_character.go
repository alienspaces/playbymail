package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	adventure_game_record "gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

// GetManyAdventureGameCharacterRecs -
func (m *Domain) GetManyAdventureGameCharacterRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameCharacter, error) {
	l := m.Logger("GetManyAdventureGameCharacterRecs")
	l.Debug("getting many adventure_game_character records opts >%#v<", opts)
	r := m.AdventureGameCharacterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameCharacterRec -
func (m *Domain) GetAdventureGameCharacterRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameCharacter, error) {
	l := m.Logger("GetAdventureGameCharacterRec")
	l.Debug("getting adventure_game_character record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameCharacterRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameCharacter, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameCharacterRec -
func (m *Domain) CreateAdventureGameCharacterRec(rec *adventure_game_record.AdventureGameCharacter) (*adventure_game_record.AdventureGameCharacter, error) {
	l := m.Logger("CreateAdventureGameCharacterRec")
	l.Debug("creating adventure_game_character record >%#v<", rec)
	if err := m.validateAdventureGameCharacterRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_character record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameCharacterRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameCharacterRec -
func (m *Domain) UpdateAdventureGameCharacterRec(rec *adventure_game_record.AdventureGameCharacter) (*adventure_game_record.AdventureGameCharacter, error) {
	l := m.Logger("UpdateAdventureGameCharacterRec")
	l.Debug("updating adventure_game_character record >%#v<", rec)
	if err := m.validateAdventureGameCharacterRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_character record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameCharacterRepository()
	var err error
	rec, err = r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// DeleteAdventureGameCharacterRec -
func (m *Domain) DeleteAdventureGameCharacterRec(recID string) error {
	l := m.Logger("DeleteAdventureGameCharacterRec")
	l.Debug("deleting adventure_game_character record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return err
	}
	r := m.AdventureGameCharacterRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// validateAdventureGameCharacterRecForCreate -
func (m *Domain) validateAdventureGameCharacterRecForCreate(rec *adventure_game_record.AdventureGameCharacter) error {
	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameCharacterName, rec.Name); err != nil {
		return err
	}
	if len(rec.Name) > 128 {
		return InvalidFieldValue("name")
	}
	return nil
}

func (m *Domain) RemoveAdventureGameCharacterRec(recID string) error {
	l := m.Logger("RemoveAdventureGameCharacterRec")
	l.Debug("removing adventure_game_character record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return err
	}
	r := m.AdventureGameCharacterRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
