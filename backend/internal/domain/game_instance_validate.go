package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type validateGameInstanceArgs struct {
	currRec *game_record.GameInstance
	nextRec *game_record.GameInstance
	gameRec *game_record.Game
}

func (m *Domain) populateGameInstanceValidateArgs(currRec, nextRec *game_record.GameInstance) (*validateGameInstanceArgs, error) {
	args := &validateGameInstanceArgs{
		currRec: currRec,
		nextRec: nextRec,
	}

	// Get game record
	if nextRec.GameID != "" {
		gameRec, err := m.GetGameRec(nextRec.GameID, nil)
		if err != nil {
			return nil, coreerror.NewInvalidDataError("game_id references invalid game")
		}
		args.gameRec = gameRec
	}

	return args, nil
}

func (m *Domain) validateGameInstanceRecForCreate(rec *game_record.GameInstance) error {
	args, err := m.populateGameInstanceValidateArgs(nil, rec)
	if err != nil {
		return err
	}

	// Draft games require closed testing — instances of unpublished games
	// cannot be joined by the public and require a specific invitation.
	if args.gameRec != nil && args.gameRec.Status == game_record.GameStatusDraft {
		rec.IsClosedTesting = true
		if !rec.DeliveryEmail {
			rec.DeliveryEmail = true
		}
	}

	// Basic validation first
	if err := validateGameInstanceRec(args, false); err != nil {
		return err
	}

	// Validate create-specific rules
	if err := validateGameInstanceRecForCreate(args); err != nil {
		return err
	}

	// Validate game is ready for instance creation
	issues, err := m.ValidateGameReadyForInstance(rec.GameID)
	if err != nil {
		return coreerror.NewInvalidDataError("failed to validate game ID >%s< >%v<", rec.GameID, err)
	}

	for _, issue := range issues {
		if issue.Severity == ValidationSeverityError {
			return InvalidField(issue.Field, rec.GameID, issue.Message)
		}
	}

	return nil
}

func (m *Domain) validateGameInstanceRecForUpdate(currRec *game_record.GameInstance, rec *game_record.GameInstance) error {
	args, err := m.populateGameInstanceValidateArgs(currRec, rec)
	if err != nil {
		return err
	}

	// Basic validation first
	if err := validateGameInstanceRec(args, true); err != nil {
		return err
	}

	return validateGameInstanceRecForUpdate(args)
}

func validateGameInstanceRecForCreate(args *validateGameInstanceArgs) error {
	rec := args.nextRec

	if rec.CurrentTurn > 0 {
		return InvalidField(
			game_record.FieldGameInstanceCurrentTurn,
			fmt.Sprintf("%d", rec.CurrentTurn),
			"current_turn must be zero for a new game instance",
		)
	}

	return nil
}

func validateGameInstanceRecForUpdate(args *validateGameInstanceArgs) error {
	currRec := args.currRec
	rec := args.nextRec

	if rec.CurrentTurn < currRec.CurrentTurn {
		return InvalidField(
			game_record.FieldGameInstanceCurrentTurn,
			fmt.Sprintf("%d", rec.CurrentTurn),
			"current_turn cannot be less than the current turn",
		)
	}

	return nil
}

func validateGameInstanceRec(args *validateGameInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(game_record.FieldGameInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := validateGameInstanceStatus(rec.Status); err != nil {
		return err
	}

	if rec.CurrentTurn < 0 {
		return InvalidField(
			game_record.FieldGameInstanceCurrentTurn,
			fmt.Sprintf("%d", rec.CurrentTurn),
			"current_turn must be zero or greater",
		)
	}

	// Validate at least one delivery method is enabled
	if !rec.DeliveryPhysicalPost && !rec.DeliveryPhysicalLocal && !rec.DeliveryEmail {
		return InvalidField(
			game_record.FieldGameInstanceDeliveryPhysicalPost,
			"false",
			"at least one delivery method must be enabled (delivery_physical_post, delivery_physical_local, or delivery_email)",
		)
	}

	// Validate closed testing requires email delivery
	if rec.IsClosedTesting && !rec.DeliveryEmail {
		return InvalidField(
			game_record.FieldGameInstanceIsClosedTesting,
			"true",
			"closed testing requires email delivery to be enabled",
		)
	}

	// Validate required_player_count (0 means no check, >= 1 means check is enforced)
	if rec.RequiredPlayerCount < 0 {
		return InvalidField(
			game_record.FieldGameInstanceRequiredPlayerCount,
			fmt.Sprintf("%d", rec.RequiredPlayerCount),
			"required_player_count must be 0 or greater",
		)
	}

	return nil
}

func validateGameInstanceStatus(status string) error {
	switch status {
	case game_record.GameInstanceStatusCreated,
		game_record.GameInstanceStatusStarted,
		game_record.GameInstanceStatusPaused,
		game_record.GameInstanceStatusCompleted,
		game_record.GameInstanceStatusCancelled:
		return nil
	default:
		return coreerror.NewInvalidDataError("invalid game instance status >%s<", status)
	}
}
