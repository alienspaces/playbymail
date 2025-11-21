package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetManyGameSubscriptionRecs -
func (m *Domain) GetManyGameSubscriptionRecs(opts *sql.Options) ([]*game_record.GameSubscription, error) {
	l := m.Logger("GetManyGameSubscriptionRecs")
	l.Debug("getting many game_subscription records opts >%#v<", opts)
	r := m.GameSubscriptionRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameSubscriptionRec -
func (m *Domain) GetGameSubscriptionRec(recID string, lock *sql.Lock) (*game_record.GameSubscription, error) {
	l := m.Logger("GetGameSubscriptionRec")
	l.Debug("getting game_subscription record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameSubscriptionRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscription, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameSubscriptionRec -
func (m *Domain) CreateGameSubscriptionRec(rec *game_record.GameSubscription) (*game_record.GameSubscription, error) {
	l := m.Logger("CreateGameSubscriptionRec")
	l.Debug("creating game_subscription record >%#v<", rec)
	r := m.GameSubscriptionRepository()
	if err := m.validateGameSubscriptionRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_subscription record >%v<", err)
		return rec, err
	}
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameSubscriptionRec -
func (m *Domain) UpdateGameSubscriptionRec(next *game_record.GameSubscription) (*game_record.GameSubscription, error) {
	l := m.Logger("UpdateGameSubscriptionRec")
	curr, err := m.GetGameSubscriptionRec(next.ID, sql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_subscription record >%#v<", next)
	if err := m.validateGameSubscriptionRecForUpdate(next, curr); err != nil {
		l.Warn("failed to validate game_subscription record >%v<", err)
		return next, err
	}
	r := m.GameSubscriptionRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteGameSubscriptionRec -
func (m *Domain) DeleteGameSubscriptionRec(recID string) error {
	l := m.Logger("DeleteGameSubscriptionRec")
	l.Debug("deleting game_subscription record ID >%s<", recID)
	rec, err := m.GetGameSubscriptionRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameSubscriptionRepository()
	if err := m.validateGameSubscriptionRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameSubscriptionRec -
func (m *Domain) RemoveGameSubscriptionRec(recID string) error {
	l := m.Logger("RemoveGameSubscriptionRec")
	l.Debug("removing game_subscription record ID >%s<", recID)
	rec, err := m.GetGameSubscriptionRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameSubscriptionRepository()
	if err := m.validateGameSubscriptionRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// Validation stubs
func (m *Domain) validateGameSubscriptionRecForCreate(rec *game_record.GameSubscription) error {
	l := m.Logger("validateGameSubscriptionRecForCreate")
	l.Debug("validating game_subscription record >%#v<", rec)

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
func (m *Domain) validateGameSubscriptionRecForUpdate(next, curr *game_record.GameSubscription) error {
	l := m.Logger("validateGameSubscriptionRecForUpdate")
	l.Debug("validating current game_subscription record >%#v< against next >%#v<", curr, next)

	if next == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if next.GameID == "" {
		return coreerror.NewInvalidDataError("game_id is required")
	}

	if next.AccountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	if next.SubscriptionType == "" {
		return coreerror.NewInvalidDataError("subscription_type is required")
	}

	if next.Status == "" {
		next.Status = curr.Status
	}

	if err := validateGameSubscriptionStatus(next.Status); err != nil {
		return err
	}

	return nil
}
func (m *Domain) validateGameSubscriptionRecForDelete(rec *game_record.GameSubscription) error {
	l := m.Logger("validateGameSubscriptionRecForDelete")
	l.Debug("validating game_subscription record >%#v<", rec)

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

func (m *Domain) UpsertPendingGameSubscription(gameID, accountID, accountContactID, subscriptionType string) (*game_record.GameSubscription, error) {
	l := m.Logger("UpsertPendingGameSubscription")

	l.Debug("upserting pending subscription game_id >%s< account_id >%s< account_contact_id >%s< type >%s<", gameID, accountID, accountContactID, subscriptionType)

	switch {
	case gameID == "":
		return nil, coreerror.NewInvalidDataError("game_id is required")
	case accountID == "":
		return nil, coreerror.NewInvalidDataError("account_id is required")
	case accountContactID == "":
		return nil, coreerror.NewInvalidDataError("account_contact_id is required")
	case subscriptionType == "":
		return nil, coreerror.NewInvalidDataError("subscription_type is required")
	}

	repo := m.GameSubscriptionRepository()

	recs, err := repo.GetMany(&sql.Options{
		Params: []sql.Param{
			{Col: game_record.FieldGameSubscriptionGameID, Val: gameID},
			{Col: game_record.FieldGameSubscriptionAccountID, Val: accountID},
			{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: subscriptionType},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, databaseError(err)
	}

	if len(recs) > 0 {
		rec := recs[0]
		// Update account_contact_id if needed
		if !rec.AccountContactID.Valid || rec.AccountContactID.String != accountContactID {
			rec.AccountContactID = nullstring.FromString(accountContactID)
		}
		if rec.Status != game_record.GameSubscriptionStatusPendingApproval {
			rec.Status = game_record.GameSubscriptionStatusPendingApproval
		}
		updated, err := m.UpdateGameSubscriptionRec(rec)
		if err != nil {
			return nil, err
		}
		return updated, nil
	}

	rec := &game_record.GameSubscription{
		GameID:           gameID,
		AccountID:        accountID,
		AccountContactID: nullstring.FromString(accountContactID),
		SubscriptionType: subscriptionType,
		Status:           game_record.GameSubscriptionStatusPendingApproval,
	}

	return m.CreateGameSubscriptionRec(rec)
}

// ApproveGameSubscription approves a pending game subscription by verifying the email
// matches the subscription's account and updating the status to active.
func (m *Domain) ApproveGameSubscription(subscriptionID, email string) (*game_record.GameSubscription, error) {
	l := m.Logger("ApproveGameSubscription")

	l.Debug("approving game subscription ID >%s< for email >%s<", subscriptionID, email)

	if subscriptionID == "" {
		return nil, coreerror.NewInvalidDataError("subscription_id is required")
	}

	if email == "" {
		return nil, coreerror.NewInvalidDataError("email is required")
	}

	// Get the subscription record
	rec, err := m.GetGameSubscriptionRec(subscriptionID, sql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	// Verify the subscription is in pending_approval status
	if rec.Status != game_record.GameSubscriptionStatusPendingApproval {
		return nil, coreerror.NewInvalidDataError("subscription is not pending approval, current status: %s", rec.Status)
	}

	// Get the account record to verify email matches
	accountRec, err := m.GetAccountRec(rec.AccountID, nil)
	if err != nil {
		l.Warn("failed to get account record >%v<", err)
		return nil, err
	}

	// Verify email matches
	if accountRec.Email != email {
		l.Warn("email mismatch: subscription account email >%s< does not match provided email >%s<", accountRec.Email, email)
		return nil, coreerror.NewInvalidDataError("email does not match subscription")
	}

	// Update status to active
	rec.Status = game_record.GameSubscriptionStatusActive

	updated, err := m.UpdateGameSubscriptionRec(rec)
	if err != nil {
		l.Warn("failed to update subscription status >%v<", err)
		return nil, err
	}

	l.Info("approved game subscription ID >%s< for account ID >%s<", subscriptionID, rec.AccountID)

	return updated, nil
}
