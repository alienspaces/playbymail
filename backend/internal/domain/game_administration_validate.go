package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) validateGameAdministrationRecForCreate(rec *game_record.GameAdministration) error {
	return validateGameAdministrationRec(rec, false)
}

func (m *Domain) validateGameAdministrationRecForUpdate(nextRec *game_record.GameAdministration) error {
	return validateGameAdministrationRec(nextRec, true)
}

func (m *Domain) validateGameAdministrationRecForDelete(rec *game_record.GameAdministration) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}
	return nil
}

func validateGameAdministrationRec(rec *game_record.GameAdministration, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(game_record.FieldGameAdministrationID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameAdministrationGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameAdministrationAccountID, rec.AccountID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameAdministrationGrantedByAccountID, rec.GrantedByAccountID); err != nil {
		return err
	}

	return nil
}
