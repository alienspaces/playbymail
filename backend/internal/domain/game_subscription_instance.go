package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// GetManyGameSubscriptionInstanceRecs -
func (m *Domain) GetManyGameSubscriptionInstanceRecs(opts *sql.Options) ([]*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("GetManyGameSubscriptionInstanceRecs")

	l.Debug("getting many game_subscription_instance records opts >%#v<", opts)

	r := m.GameSubscriptionInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameSubscriptionInstanceRec -
func (m *Domain) GetGameSubscriptionInstanceRec(recID string, lock *sql.Lock) (*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("GetGameSubscriptionInstanceRec")

	l.Debug("getting game_subscription_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameSubscriptionInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameSubscriptionInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameSubscriptionInstanceRec -
func (m *Domain) CreateGameSubscriptionInstanceRec(rec *game_record.GameSubscriptionInstance) (*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("CreateGameSubscriptionInstanceRec")

	l.Debug("creating game_subscription_instance record >%#v<", rec)

	// Check for existing link first (idempotency)
	existingRecs, err := m.GetManyGameSubscriptionInstanceRecs(&sql.Options{
		Params: []sql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameSubscriptionID, Val: rec.GameSubscriptionID},
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: rec.GameInstanceID},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, databaseError(err)
	}

	if len(existingRecs) > 0 {
		// Link already exists, return it (idempotent behavior)
		l.Debug("link already exists, returning existing record")
		return existingRecs[0], nil
	}

	r := m.GameSubscriptionInstanceRepository()

	if err := m.validateGameSubscriptionInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_subscription_instance record >%v<", err)
		return nil, err
	}

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create game_subscription_instance record >%v<", err)
		return nil, databaseError(err)
	}

	return createdRec, nil
}

// UpdateGameSubscriptionInstanceRec -
func (m *Domain) UpdateGameSubscriptionInstanceRec(rec *game_record.GameSubscriptionInstance) (*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("UpdateGameSubscriptionInstanceRec")

	curr, err := m.GetGameSubscriptionInstanceRec(rec.ID, sql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game_subscription_instance record >%#v<", rec)

	if err := m.validateGameSubscriptionInstanceRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate game_subscription_instance record >%v<", err)
		return rec, err
	}

	r := m.GameSubscriptionInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteGameSubscriptionInstanceRec -
func (m *Domain) DeleteGameSubscriptionInstanceRec(recID string) error {
	l := m.Logger("DeleteGameSubscriptionInstanceRec")
	l.Debug("deleting game_subscription_instance record ID >%s<", recID)
	rec, err := m.GetGameSubscriptionInstanceRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameSubscriptionInstanceRepository()
	if err := m.validateGameSubscriptionInstanceRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameSubscriptionInstanceRec -
func (m *Domain) RemoveGameSubscriptionInstanceRec(recID string) error {
	l := m.Logger("RemoveGameSubscriptionInstanceRec")
	l.Debug("removing game_subscription_instance record ID >%s<", recID)
	rec, err := m.GetGameSubscriptionInstanceRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameSubscriptionInstanceRepository()
	if err := m.validateGameSubscriptionInstanceRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

func (m *Domain) validateGameSubscriptionInstanceRecForCreate(rec *game_record.GameSubscriptionInstance) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if rec.GameSubscriptionID == "" {
		return coreerror.NewInvalidDataError("game_subscription_id is required")
	}

	if rec.GameInstanceID == "" {
		return coreerror.NewInvalidDataError("game_instance_id is required")
	}

	// Validate subscription exists
	subscriptionRec, err := m.GetGameSubscriptionRec(rec.GameSubscriptionID, nil)
	if err != nil {
		return coreerror.NewInvalidDataError("game_subscription_id references invalid subscription")
	}

	// Derive account_id from subscription if not provided
	if rec.AccountID == "" {
		rec.AccountID = subscriptionRec.AccountID
	}

	// Validate account_id matches subscription's account_id
	if rec.AccountID != subscriptionRec.AccountID {
		return coreerror.NewInvalidDataError("account_id must match the subscription's account_id")
	}

	// Validate instance exists
	instanceRec, err := m.GetGameInstanceRec(rec.GameInstanceID, nil)
	if err != nil {
		return coreerror.NewInvalidDataError("game_instance_id references invalid instance")
	}

	// Validate instance belongs to same game as subscription
	if instanceRec.GameID != subscriptionRec.GameID {
		return coreerror.NewInvalidDataError("game instance must belong to the same game as the subscription")
	}

	// Validate instance limit
	if err := m.ValidateInstanceLimit(rec.GameSubscriptionID); err != nil {
		return err
	}

	return nil
}

func (m *Domain) validateGameSubscriptionInstanceRecForUpdate(currRec, nextRec *game_record.GameSubscriptionInstance) error {

	if nextRec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if nextRec.GameSubscriptionID != currRec.GameSubscriptionID {
		return coreerror.NewInvalidDataError("game_subscription_id cannot be updated")
	}

	if nextRec.GameInstanceID != currRec.GameInstanceID {
		return coreerror.NewInvalidDataError("game_instance_id cannot be updated")
	}

	if nextRec.AccountID != currRec.AccountID {
		return coreerror.NewInvalidDataError("account_id cannot be updated")
	}

	// Validate subscription exists
	subscriptionRec, err := m.GetGameSubscriptionRec(nextRec.GameSubscriptionID, nil)
	if err != nil {
		return coreerror.NewInvalidDataError("game_subscription_id references invalid subscription")
	}

	// Validate account_id matches subscription's account_id
	if nextRec.AccountID != subscriptionRec.AccountID {
		return coreerror.NewInvalidDataError("account_id must match the subscription's account_id")
	}

	// Validate instance exists
	instanceRec, err := m.GetGameInstanceRec(nextRec.GameInstanceID, nil)
	if err != nil {
		return coreerror.NewInvalidDataError("game_instance_id references invalid instance")
	}

	// Validate instance belongs to same game as subscription
	if instanceRec.GameID != subscriptionRec.GameID {
		return coreerror.NewInvalidDataError("game instance must belong to the same game as the subscription")
	}

	return nil
}

func (m *Domain) validateGameSubscriptionInstanceRecForDelete(rec *game_record.GameSubscriptionInstance) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}

// GenerateGameSubscriptionInstanceTurnSheetToken generates a UUID turn sheet key and sets expiration to 3 days from now.
// This invalidates any existing key by generating a new one.
func (m *Domain) GenerateGameSubscriptionInstanceTurnSheetToken(instanceID string) (string, error) {
	l := m.Logger("GenerateGameSubscriptionInstanceTurnSheetToken")

	rec, err := m.GetGameSubscriptionInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get game subscription instance record >%v<", err)
		return "", err
	}

	// Generate UUID for turn sheet key
	turnSheetKey := uuid.New().String()

	// Set turn sheet key and expiration (3 days from now)
	rec.TurnSheetToken = nullstring.FromString(turnSheetKey)
	rec.TurnSheetTokenExpiresAt = nulltime.FromTime(time.Now().Add(3 * 24 * time.Hour))

	// Update instance with new key
	_, err = m.UpdateGameSubscriptionInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating game subscription instance with turn sheet token >%v<", err)
		return "", err
	}

	l.Info("generated turn sheet token for game subscription instance >%s<", instanceID)

	return turnSheetKey, nil
}

// VerifyGameSubscriptionInstanceTurnSheetKey verifies that a turn sheet key exists, is not expired, and matches the instance.
func (m *Domain) VerifyGameSubscriptionInstanceTurnSheetKey(instanceID, turnSheetKey string) (*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("VerifyGameSubscriptionInstanceTurnSheetKey")
	l.Debug("verifying game subscription instance turn sheet key >%s<", turnSheetKey)

	if instanceID == "" {
		return nil, coreerror.NewInvalidDataError("instance_id is required")
	}

	if turnSheetKey == "" {
		return nil, coreerror.NewInvalidDataError("turn_sheet_key is required")
	}

	// Get instance by instance_id
	rec, err := m.GetGameSubscriptionInstanceRec(instanceID, sql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get game subscription instance record >%v<", err)
		return nil, err
	}

	// Check expiration
	if !nulltime.IsValid(rec.TurnSheetTokenExpiresAt) {
		return nil, fmt.Errorf("turn sheet token has no expiration set")
	}

	if time.Now().After(nulltime.ToTime(rec.TurnSheetTokenExpiresAt)) {
		return nil, fmt.Errorf("turn sheet token has expired")
	}

	// Verify key matches
	if !nullstring.IsValid(rec.TurnSheetToken) || nullstring.ToString(rec.TurnSheetToken) != turnSheetKey {
		return nil, fmt.Errorf("turn sheet token does not match instance")
	}

	l.Debug("turn sheet token verified successfully for instance >%s<", rec.ID)

	return rec, nil
}

// GetGameSubscriptionInstanceRecFromCodeData retrieves a game subscription instance from a turn sheet identifier.
func (m *Domain) GetGameSubscriptionInstanceRecFromCodeData(turnSheetCodeData *turnsheetutil.PlayGameTurnSheetCodeData) (*game_record.GameSubscriptionInstance, error) {
	l := m.Logger("GetGameSubscriptionInstanceRecFromCodeData")

	if turnSheetCodeData == nil {
		return nil, coreerror.NewInvalidDataError("turn sheet code data is required")
	}

	if turnSheetCodeData.CodeType == turnsheetutil.TurnSheetCodeTypeJoiningGame {
		return nil, coreerror.NewInvalidDataError("joining game codes do not have instances")
	}

	l.Info("game turn sheet ID >%s<", turnSheetCodeData.GameTurnSheetID)

	turnSheetRec, err := m.GetGameTurnSheetRec(turnSheetCodeData.GameTurnSheetID, nil)
	if err != nil {
		l.Warn("failed to get game turn sheet record >%v<", err)
		return nil, err
	}

	l.Info("fetching game subscription instance record for account ID >%s< and game instance ID >%s<", turnSheetRec.AccountID, turnSheetRec.GameInstanceID)

	recs, err := m.GetManyGameSubscriptionInstanceRecs(&sql.Options{
		Params: []sql.Param{
			{
				Col: game_record.FieldGameSubscriptionInstanceAccountID,
				Val: turnSheetRec.AccountID,
			},
			{
				Col: game_record.FieldGameSubscriptionInstanceGameInstanceID,
				Val: turnSheetRec.GameInstanceID,
			},
		},
	})
	if err != nil {
		l.Warn("failed to get game subscription instance records >%v<", err)
		return nil, err
	}

	if len(recs) == 0 {
		return nil, coreerror.NewInvalidDataError("game subscription instance not found for account ID >%s< and game instance ID >%s<", turnSheetRec.AccountID, nullstring.ToString(turnSheetRec.GameInstanceID))
	}

	if len(recs) > 1 {
		return nil, coreerror.NewInvalidDataError("multiple game subscription instances found for account ID >%s< and game instance ID >%s<", turnSheetRec.AccountID, nullstring.ToString(turnSheetRec.GameInstanceID))
	}

	return recs[0], nil
}
