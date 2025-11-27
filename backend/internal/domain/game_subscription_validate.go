package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) validateGameSubscriptionRecForCreate(rec *game_record.GameSubscription) error {
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
		if !rec.AccountContactID.Valid || rec.AccountContactID.String == "" {
			return coreerror.NewInvalidDataError("account_contact_id is required for player subscriptions")
		}
	}

	if rec.SubscriptionType == "" {
		return coreerror.NewInvalidDataError("subscription_type is required")
	}

	if rec.Status == "" {
		rec.Status = game_record.GameSubscriptionStatusActive
	}

	if err := validateGameSubscriptionStatus(rec.Status); err != nil {
		return err
	}

	return nil
}

func (m *Domain) validateGameSubscriptionRecForUpdate(nextRec, currRec *game_record.GameSubscription) error {
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

	return nil
}

func (m *Domain) validateGameSubscriptionRecForDelete(rec *game_record.GameSubscription) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
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
