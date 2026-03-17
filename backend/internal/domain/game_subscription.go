package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetManyGameSubscriptionRecs -
func (m *Domain) GetManyGameSubscriptionRecs(opts *sql.Options) ([]*game_record.GameSubscription, error) {
	l := m.Logger("GetManyGameSubscriptionRecs")
	l.Debug("getting many game_subscription records opts >%#v<", opts)

	r := m.GameSubscriptionRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		l.Warn("failed to get many game_subscription records >%v<", err)
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetManyGameSubscriptionViewRecs returns game subscription view records with aggregated instance IDs
// This is the preferred method for API queries that need instance IDs
func (m *Domain) GetManyGameSubscriptionViewRecs(opts *sql.Options) ([]*game_record.GameSubscriptionView, error) {
	l := m.Logger("GetManyGameSubscriptionViewRecs")
	l.Debug("getting many game_subscription_view records opts >%#v<", opts)

	r := m.GameSubscriptionViewRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		l.Warn("failed to get many game_subscription view records >%v<", err)
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameSubscriptionViewRec returns a single game subscription view record with aggregated instance IDs
// This is the preferred method for API queries that need instance IDs
func (m *Domain) GetGameSubscriptionViewRec(recID string, lock *sql.Lock) (*game_record.GameSubscriptionView, error) {
	l := m.Logger("GetGameSubscriptionViewRec")
	l.Debug("getting game_subscription_view record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameSubscriptionViewRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscriptionView, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
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

// GetGameSubscriptionRecByAccountAndGame finds a subscription for a specific account, game, and subscription type.
func (m *Domain) GetGameSubscriptionRecByAccountAndGame(accountID, gameID, subscriptionType string) (*game_record.GameSubscription, error) {
	l := m.Logger("GetGameSubscriptionRecByAccountAndGame")
	l.Debug("getting game_subscription for account >%s< game >%s< type >%s<", accountID, gameID, subscriptionType)

	if err := domain.ValidateUUIDField("account_id", accountID); err != nil {
		return nil, err
	}
	if err := domain.ValidateUUIDField("game_id", gameID); err != nil {
		return nil, err
	}

	r := m.GameSubscriptionRepository()
	recs, err := r.GetMany(&sql.Options{
		Params: []sql.Param{
			{Col: game_record.FieldGameSubscriptionAccountID, Val: accountID},
			{Col: game_record.FieldGameSubscriptionGameID, Val: gameID},
			{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: subscriptionType},
			{Col: game_record.FieldGameSubscriptionStatus, Val: game_record.GameSubscriptionStatusActive},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get game_subscription record >%v<", err)
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscription,
			"account_id="+accountID+", game_id="+gameID+", type="+subscriptionType)
	}

	return recs[0], nil
}

// CreateDesignerSubscriptionForNewGame creates the initial designer subscription
// as part of game creation. This uses validation appropriate for the creation
// context and does not require the game to be published.
func (m *Domain) CreateDesignerSubscriptionForNewGame(gameRec *game_record.Game, accountID, accountUserID string) (*game_record.GameSubscription, error) {
	l := m.Logger("CreateDesignerSubscriptionForNewGame")

	rec := &game_record.GameSubscription{
		GameID:           gameRec.ID,
		AccountID:        accountID,
		AccountUserID:    accountUserID,
		SubscriptionType: game_record.GameSubscriptionTypeDesigner,
		Status:           game_record.GameSubscriptionStatusActive,
	}

	l.Debug("creating designer subscription for new game >%s< account >%s<", gameRec.ID, accountID)

	if err := validateDesignerSubscriptionForNewGame(rec); err != nil {
		l.Warn("failed to validate designer subscription >%v<", err)
		return nil, err
	}

	r := m.GameSubscriptionRepository()

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create designer subscription >%v<", err)
		return nil, databaseError(err)
	}

	return createdRec, nil
}

// CreateManagerSubscriptionForNewGame creates the initial manager subscription
// as part of game creation. This allows the game designer to create and manage
// instances of their game before it is published.
func (m *Domain) CreateManagerSubscriptionForNewGame(gameRec *game_record.Game, accountID, accountUserID string) (*game_record.GameSubscription, error) {
	l := m.Logger("CreateManagerSubscriptionForNewGame")

	rec := &game_record.GameSubscription{
		GameID:           gameRec.ID,
		AccountID:        accountID,
		AccountUserID:    accountUserID,
		SubscriptionType: game_record.GameSubscriptionTypeManager,
		Status:           game_record.GameSubscriptionStatusActive,
	}

	l.Debug("creating manager subscription for new game >%s< account >%s<", gameRec.ID, accountID)

	if err := validateManagerSubscriptionForNewGame(rec); err != nil {
		l.Warn("failed to validate manager subscription >%v<", err)
		return nil, err
	}

	r := m.GameSubscriptionRepository()

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create manager subscription >%v<", err)
		return nil, databaseError(err)
	}

	return createdRec, nil
}

// CreateGameSubscriptionRec -
func (m *Domain) CreateGameSubscriptionRec(rec *game_record.GameSubscription) (*game_record.GameSubscription, error) {
	l := m.Logger("CreateGameSubscriptionRec")

	l.Debug("creating game_subscription record >%#v<", rec)

	r := m.GameSubscriptionRepository()

	// Player subscriptions require an explicit delivery method.
	if rec.SubscriptionType == game_record.GameSubscriptionTypePlayer && !rec.DeliveryMethod.Valid {
		return nil, coreerror.NewInvalidDataError("delivery_method is required for player subscriptions")
	}

	if err := m.validateGameSubscriptionRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_subscription record >%v<", err)
		return nil, err
	}

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create game_subscription record >%v<", err)
		return nil, databaseError(err)
	}

	return createdRec, nil
}

// UpdateGameSubscriptionRec -
func (m *Domain) UpdateGameSubscriptionRec(rec *game_record.GameSubscription) (*game_record.GameSubscription, error) {
	l := m.Logger("UpdateGameSubscriptionRec")

	curr, err := m.GetGameSubscriptionRec(rec.ID, sql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game_subscription record >%#v<", rec)

	if rec.Status == "" {
		rec.Status = curr.Status
	}

	if err := m.validateGameSubscriptionRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate game_subscription record >%v<", err)
		return rec, err
	}

	r := m.GameSubscriptionRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
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
	r := m.GameSubscriptionRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

func (m *Domain) UpsertPendingGameSubscription(gameID, accountID, accountUserID, accountUserContactID, subscriptionType string) (*game_record.GameSubscription, error) {
	l := m.Logger("UpsertPendingGameSubscription")

	l.Debug("upserting pending subscription game_id >%s< account_id >%s< account_user_id >%s< type >%s<", gameID, accountID, accountUserID, subscriptionType)

	switch {
	case gameID == "":
		return nil, coreerror.NewInvalidDataError("game_id is required")
	case accountID == "":
		return nil, coreerror.NewInvalidDataError("account_id is required")
	case accountUserID == "":
		return nil, coreerror.NewInvalidDataError("account_user_id is required")
	case accountUserContactID == "":
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
		if !rec.AccountUserContactID.Valid || rec.AccountUserContactID.String != accountUserContactID {
			rec.AccountUserContactID = nullstring.FromString(accountUserContactID)
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
		GameID:               gameID,
		AccountID:            accountID,
		AccountUserID:        accountUserID,
		AccountUserContactID: nullstring.FromString(accountUserContactID),
		SubscriptionType:     subscriptionType,
		Status:               game_record.GameSubscriptionStatusPendingApproval,
	}

	return m.CreateGameSubscriptionRec(rec)
}

// CreateOrGetPendingGameSubscriptionForJoinProcess creates a new pending game subscription for the join process,
// or returns an existing one if the same game, account user, and subscription type already has a subscription
// in pending_approval status. Used during join game turn sheet processing.
// Caller must pass a GameSubscription record with GameID, AccountID, AccountUserID, AccountUserContactID (required), and SubscriptionType set.
// Caps on how many times a player can join will be enforced later via account-level subscriptions.
func (m *Domain) CreateOrGetPendingGameSubscriptionForJoinProcess(rec *game_record.GameSubscription) (*game_record.GameSubscription, error) {
	l := m.Logger("CreateOrGetPendingGameSubscriptionForJoinProcess")

	if rec == nil {
		return nil, coreerror.NewInvalidDataError("game_subscription record is required")
	}

	l.Debug("create or get pending subscription game_id >%s< account_id >%s< account_user_id >%s< type >%s<",
		rec.GameID, rec.AccountID, rec.AccountUserID, rec.SubscriptionType)

	switch {
	case rec.GameID == "":
		return nil, coreerror.NewInvalidDataError("game_id is required")
	case rec.AccountID == "":
		return nil, coreerror.NewInvalidDataError("account_id is required")
	case rec.AccountUserID == "":
		return nil, coreerror.NewInvalidDataError("account_user_id is required")
	case !nullstring.IsValid(rec.AccountUserContactID):
		return nil, coreerror.NewInvalidDataError("account_user_contact_id is required")
	case rec.SubscriptionType == "":
		return nil, coreerror.NewInvalidDataError("subscription_type is required")
	}

	existingRecs, err := m.GetManyGameSubscriptionRecs(&sql.Options{
		Params: []sql.Param{
			{Col: game_record.FieldGameSubscriptionGameID, Val: rec.GameID},
			{Col: game_record.FieldGameSubscriptionAccountUserID, Val: rec.AccountUserID},
			{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: rec.SubscriptionType},
			{Col: game_record.FieldGameSubscriptionStatus, Val: game_record.GameSubscriptionStatusPendingApproval},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, databaseError(err)
	}

	if len(existingRecs) > 0 {
		l.Warn("reusing existing pending game subscription id >%s< for game >%s< account_user >%s<",
			existingRecs[0].ID, rec.GameID, rec.AccountUserID)
		return existingRecs[0], nil
	}

	// Default to pending approval status
	rec.Status = game_record.GameSubscriptionStatusPendingApproval

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

	// Get the account user record(s) linked to this tenant
	// Note: We use the sql alias for core/sql as defined in imports
	accountUserRecs, err := m.GetManyAccountUserRecs(&sql.Options{
		Params: []sql.Param{
			{Col: account_record.FieldAccountUserAccountID, Val: rec.AccountID},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account records >%v<", err)
		return nil, err
	}
	if len(accountUserRecs) == 0 {
		l.Warn("no account user found for tenant >%s<", rec.AccountID)
		return nil, coreerror.NewNotFoundError("account_user", "for tenant "+rec.AccountID)
	}
	accountUserRec := accountUserRecs[0]

	// Verify email matches
	if accountUserRec.Email != email {
		l.Warn("email mismatch: subscription account email >%s< does not match provided email >%s<", accountUserRec.Email, email)
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

// GetGameSubscriptionInstanceRecsBySubscription retrieves all instance links for a subscription
func (m *Domain) GetGameSubscriptionInstanceRecsBySubscription(subscriptionID string) ([]*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("GetGameSubscriptionInstanceRecsBySubscription")
	l.Debug("getting game_subscription_instance records for subscription >%s<", subscriptionID)

	if err := domain.ValidateUUIDField("subscription_id", subscriptionID); err != nil {
		return nil, err
	}

	r := m.GameSubscriptionInstanceRepository()
	recs, err := r.GetMany(&sql.Options{
		Params: []sql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameSubscriptionID, Val: subscriptionID},
		},
	})
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameSubscriptionInstanceRecsByInstance retrieves all subscription links for an instance
func (m *Domain) GetGameSubscriptionInstanceRecsByInstance(instanceID string) ([]*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("GetGameSubscriptionInstanceRecsByInstance")
	l.Debug("getting game_subscription_instance records for instance >%s<", instanceID)

	if err := domain.ValidateUUIDField("instance_id", instanceID); err != nil {
		return nil, err
	}

	r := m.GameSubscriptionInstanceRepository()
	recs, err := r.GetMany(&sql.Options{
		Params: []sql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID},
		},
	})
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// ValidateInstanceLimit checks if adding another instance would exceed the subscription's limit
func (m *Domain) ValidateInstanceLimit(subscriptionID string) error {
	l := m.Logger("ValidateInstanceLimit")
	l.Debug("validating instance limit for subscription >%s<", subscriptionID)

	subscriptionRec, err := m.GetGameSubscriptionRec(subscriptionID, nil)
	if err != nil {
		return err
	}

	// If instance_limit is NULL, unlimited instances allowed
	if !subscriptionRec.InstanceLimit.Valid {
		return nil
	}

	// Get current instance count
	instanceRecs, err := m.GetGameSubscriptionInstanceRecsBySubscription(subscriptionID)
	if err != nil {
		return err
	}

	currentCount := len(instanceRecs)
	limit := int(subscriptionRec.InstanceLimit.Int32)

	if currentCount >= limit {
		return coreerror.NewInvalidDataError("instance limit reached: %d/%d instances", currentCount, limit)
	}

	return nil
}
