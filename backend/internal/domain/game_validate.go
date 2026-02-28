package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const (
	ValidationSeverityError   = "error"
	ValidationSeverityWarning = "warning"
)

type GameValidationIssue struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// ValidateGameReadyForInstance checks whether a game is ready to create an
// instance. It collects all issues rather than stopping at the first, so that
// game designers can see everything that needs fixing at once.
func (m *Domain) ValidateGameReadyForInstance(gameID string) ([]GameValidationIssue, error) {
	l := m.Logger("ValidateGameReadyForInstance")

	gameRec, err := m.GetGameRec(gameID, nil)
	if err != nil {
		return nil, err
	}

	var issues []GameValidationIssue

	if gameRec.GameType == game_record.GameTypeAdventure {
		issues, err = m.validateAdventureGameReadyForInstance(gameID)
		if err != nil {
			l.Warn("failed validating adventure game >%s< >%v<", gameID, err)
			return nil, err
		}
	}

	return issues, nil
}

func (m *Domain) validateAdventureGameReadyForInstance(gameID string) ([]GameValidationIssue, error) {
	var issues []GameValidationIssue

	locationRecs, err := m.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(locationRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "locations",
			Message:  "Adventure game must have at least one location",
			Severity: ValidationSeverityError,
		})
		return issues, nil
	}

	startingLocationRecs, err := m.GetManyAdventureGameLocationRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
			{Col: adventure_game_record.FieldAdventureGameLocationIsStartingLocation, Val: true},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(startingLocationRecs) == 0 {
		issues = append(issues, GameValidationIssue{
			Field:    "starting_location",
			Message:  "Adventure game must have at least one starting location before creating an instance",
			Severity: ValidationSeverityError,
		})
	}

	return issues, nil
}

type validateGameArgs struct {
	nextRec *game_record.Game
	currRec *game_record.Game
}

func (m *Domain) populateGameValidateArgs(currRec, nextRec *game_record.Game) (*validateGameArgs, error) {
	args := &validateGameArgs{
		currRec: currRec,
		nextRec: nextRec,
	}

	// Note: account_id validation and subscription limit checks are now handled
	// at the handler level where we have access to authenticated account information.
	// Games no longer have account_id field, so we can't validate it here.

	return args, nil
}

func (m *Domain) validateGameRecForCreate(rec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameRecForCreate(args)
}

func (m *Domain) validateGameRecForUpdate(currRec, nextRec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateGameRecForUpdate(args)
}

func (m *Domain) validateGameRecForDelete(rec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameRecForDelete(args)
}

func validateGameRecForCreate(args *validateGameArgs) error {
	rec := args.nextRec

	// Set default status if not provided
	if rec.Status == "" {
		rec.Status = game_record.GameStatusDraft
	}

	// Note: Subscription limit validation has been moved to handler level
	// where we have access to authenticated account information.
	// Games no longer have account_id field.

	return validateGameRec(args)
}

func validateGameRecForUpdate(args *validateGameArgs) error {
	currRec := args.currRec
	nextRec := args.nextRec

	// Prevent modifications to published games (except status changes)
	if currRec.Status == game_record.GameStatusPublished {
		// Allow status change from published to published (no-op)
		if nextRec.Status == game_record.GameStatusPublished {
			// Check if any other fields changed
			if currRec.Name != nextRec.Name ||
				currRec.Description != nextRec.Description ||
				currRec.GameType != nextRec.GameType ||
				currRec.TurnDurationHours != nextRec.TurnDurationHours {
				return InvalidField(game_record.FieldGameStatus, currRec.Status, "published games cannot be modified")
			}
			// No changes, allow the update
			return nil
		}
		// Cannot change status from published to draft
		return InvalidField(game_record.FieldGameStatus, currRec.Status, "published games cannot be modified")
	}

	// Validate status transitions
	// Only allow: draft -> published, draft -> draft
	if currRec.Status == game_record.GameStatusDraft {
		if nextRec.Status != game_record.GameStatusDraft && nextRec.Status != game_record.GameStatusPublished {
			return InvalidField(game_record.FieldGameStatus, nextRec.Status, "invalid status transition from draft")
		}
	}

	return validateGameRec(args)
}

func validateGameRec(args *validateGameArgs) error {
	rec := args.nextRec

	// Note: account_id validation removed - games no longer have account_id field

	if err := domain.ValidateStringField(game_record.FieldGameName, rec.Name); err != nil {
		return err
	}

	if rec.GameType != game_record.GameTypeAdventure {
		return InvalidField(game_record.FieldGameType, rec.GameType, "game type is not valid")
	}

	if rec.TurnDurationHours <= 0 {
		return InvalidField(game_record.FieldGameTurnDurationHours, fmt.Sprintf("%d", rec.TurnDurationHours), "turn duration hours must be greater than 0")
	}

	if err := domain.ValidateStringField(game_record.FieldGameDescription, rec.Description); err != nil {
		return err
	}

	if rec.Status != game_record.GameStatusDraft && rec.Status != game_record.GameStatusPublished {
		return InvalidField(game_record.FieldGameStatus, rec.Status, "status is not valid")
	}

	return nil
}

func validateGameRecForDelete(args *validateGameArgs) error {
	rec := args.nextRec

	if err := domain.ValidateUUIDField(game_record.FieldGameID, rec.ID); err != nil {
		return err
	}

	return nil
}
