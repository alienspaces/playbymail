package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameCharacterRecs -
func (m *Domain) GetManyGameCharacterRecs(opts *coresql.Options) ([]*record.GameCharacter, error) {
	l := m.Logger("GetManyGameCharacterRecs")
	l.Debug("getting many game_character records opts >%#v<", opts)
	r := m.GameCharacterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameCharacterRec -
func (m *Domain) GetGameCharacterRec(recID string, lock *coresql.Lock) (*record.GameCharacter, error) {
	l := m.Logger("GetGameCharacterRec")
	l.Debug("getting game_character record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameCharacterRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGameCharacter, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameCharacterRec -
func (m *Domain) CreateGameCharacterRec(rec *record.GameCharacter) (*record.GameCharacter, error) {
	l := m.Logger("CreateGameCharacterRec")
	l.Debug("creating game_character record >%#v<", rec)
	if err := m.validateGameCharacterRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_character record >%v<", err)
		return rec, err
	}
	r := m.GameCharacterRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// DeleteGameCharacterRec -
func (m *Domain) DeleteGameCharacterRec(recID string) error {
	l := m.Logger("DeleteGameCharacterRec")
	l.Debug("deleting game_character record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return err
	}
	r := m.GameCharacterRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// validateGameCharacterRecForCreate -
func (m *Domain) validateGameCharacterRecForCreate(rec *record.GameCharacter) error {
	if err := domain.ValidateStringField(record.FieldGameCharacterName, rec.Name); err != nil {
		return err
	}
	if len(rec.Name) > 128 {
		return InvalidFieldValue("name")
	}
	return nil
}

func (m *Domain) RemoveGameCharacterRec(recID string) error {
	l := m.Logger("RemoveGameCharacterRec")
	l.Debug("removing game_character record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return err
	}
	r := m.GameCharacterRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
