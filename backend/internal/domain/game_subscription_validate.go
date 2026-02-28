package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type validateGameSubscriptionArgs struct {
	nextRec *game_record.GameSubscription
	currRec *game_record.GameSubscription
	gameRec *game_record.Game
}

func (m *Domain) populateGameSubscriptionValidateArgs(currRec, nextRec *game_record.GameSubscription) (*validateGameSubscriptionArgs, error) {
	args := &validateGameSubscriptionArgs{
		currRec: currRec,
		nextRec: nextRec,
	}

	// Get game record if game_id is provided
	if nextRec.GameID != "" {
		gameRec, err := m.GetGameRec(nextRec.GameID, nil)
		if err != nil {
			return nil, coreerror.NewInvalidDataError("game_id references invalid game")
		}
		args.gameRec = gameRec
	}

	return args, nil
}

func (m *Domain) validateGameSubscriptionRecForCreate(rec *game_record.GameSubscription) error {
	args, err := m.populateGameSubscriptionValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameSubscriptionRecForCreate(args)
}

func (m *Domain) validateGameSubscriptionRecForUpdate(currRec, nextRec *game_record.GameSubscription) error {
	args, err := m.populateGameSubscriptionValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateGameSubscriptionRecForUpdate(args)
}

func (m *Domain) validateGameSubscriptionRecForDelete(rec *game_record.GameSubscription) error {
	args, err := m.populateGameSubscriptionValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameSubscriptionRecForDelete(args)
}

func validateGameSubscriptionRecForCreate(args *validateGameSubscriptionArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if rec.GameID == "" {
		return coreerror.NewInvalidDataError("game_id is required")
	}

	if rec.AccountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	// Account contact is required for player subscriptions
	if rec.SubscriptionType == game_record.GameSubscriptionTypePlayer {
		if !rec.AccountUserContactID.Valid || rec.AccountUserContactID.String == "" {
			return coreerror.NewInvalidDataError("account_contact_id is required for player subscriptions")
		}
	}

	if rec.SubscriptionType == "" {
		return coreerror.NewInvalidDataError("subscription_type is required")
	}

	// Validate that game is published (only published games can be subscribed to)
	if args.gameRec == nil {
		return coreerror.NewInvalidDataError("game_id references invalid game")
	}

	if args.gameRec.Status != game_record.GameStatusPublished {
		return coreerror.NewInvalidDataError("only published games can be subscribed to")
	}

	if rec.Status == "" {
		rec.Status = game_record.GameSubscriptionStatusActive
	}

	if err := validateGameSubscriptionStatus(rec.Status); err != nil {
		return err
	}

	// Validate instance_limit if provided (must be positive)
	if rec.InstanceLimit.Valid {
		if rec.InstanceLimit.Int32 <= 0 {
			return coreerror.NewInvalidDataError("instance_limit must be positive if provided")
		}
	}

	return nil
}

func validateGameSubscriptionRecForUpdate(args *validateGameSubscriptionArgs) error {
	nextRec := args.nextRec
	currRec := args.currRec

	if nextRec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if nextRec.GameID == "" {
		return coreerror.NewInvalidDataError("game_id is required")
	}

	if nextRec.AccountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	if nextRec.SubscriptionType == "" {
		return coreerror.NewInvalidDataError("subscription_type is required")
	}

	if nextRec.Status == "" {
		nextRec.Status = currRec.Status
	}

	if err := validateGameSubscriptionStatus(nextRec.Status); err != nil {
		return err
	}

	// Validate instance_limit if provided (must be positive)
	if nextRec.InstanceLimit.Valid {
		if nextRec.InstanceLimit.Int32 <= 0 {
			return coreerror.NewInvalidDataError("instance_limit must be positive if provided")
		}
	}

	return nil
}

func validateGameSubscriptionRecForDelete(args *validateGameSubscriptionArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}

// validateDesignerSubscriptionForNewGame validates a designer subscription being
// auto-created as part of game creation. Unlike the standard create validation,
// this does not require the game to be published or look up the game via RLS,
// since the game record was just created in the same transaction.
func validateDesignerSubscriptionForNewGame(rec *game_record.GameSubscription) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if rec.GameID == "" {
		return coreerror.NewInvalidDataError("game_id is required")
	}

	if rec.AccountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	if rec.SubscriptionType != game_record.GameSubscriptionTypeDesigner {
		return coreerror.NewInvalidDataError("subscription_type must be designer for new game creation")
	}

	if rec.Status == "" {
		rec.Status = game_record.GameSubscriptionStatusActive
	}

	if err := validateGameSubscriptionStatus(rec.Status); err != nil {
		return err
	}

	return nil
}

// validateManagerSubscriptionForNewGame validates a manager subscription being
// auto-created as part of game creation. Like the designer variant, this does
// not require the game to be published or look up the game via RLS.
func validateManagerSubscriptionForNewGame(rec *game_record.GameSubscription) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if rec.GameID == "" {
		return coreerror.NewInvalidDataError("game_id is required")
	}

	if rec.AccountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	if rec.SubscriptionType != game_record.GameSubscriptionTypeManager {
		return coreerror.NewInvalidDataError("subscription_type must be manager for new game creation")
	}

	if rec.Status == "" {
		rec.Status = game_record.GameSubscriptionStatusActive
	}

	if err := validateGameSubscriptionStatus(rec.Status); err != nil {
		return err
	}

	return nil
}

func validateGameSubscriptionStatus(status string) error {
	switch status {
	case game_record.GameSubscriptionStatusPendingApproval,
		game_record.GameSubscriptionStatusActive,
		game_record.GameSubscriptionStatusRevoked:
		return nil
	default:
		return coreerror.NewInvalidDataError("invalid game subscription status >%s<", status)
	}
}
