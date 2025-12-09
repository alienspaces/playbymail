package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) validateGameInstanceRecForCreate(rec *game_record.GameInstance) error {
	if err := validateGameInstanceRec(rec, false); err != nil {
		return err
	}

	// Validate that game_subscription_id references a Manager subscription
	subscriptionRec, err := m.GetGameSubscriptionRec(rec.GameSubscriptionID, nil)
	if err != nil {
		return coreerror.NewInvalidDataError("game_subscription_id references invalid game subscription")
	}

	if subscriptionRec.SubscriptionType != game_record.GameSubscriptionTypeManager {
		return coreerror.NewInvalidDataError("game_subscription_id must reference a Manager subscription")
	}

	// Enforce 10 instance limit per manager subscription
	existingInstances, err := m.GetManyGameInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameInstanceGameSubscriptionID,
				Val: rec.GameSubscriptionID,
			},
		},
	})
	if err != nil {
		return coreerror.NewInvalidDataError("failed to check existing game instances")
	}

	// Count non-deleted instances
	activeInstanceCount := 0
	for _, instance := range existingInstances {
		if !instance.DeletedAt.Valid {
			activeInstanceCount++
		}
	}

	const maxInstancesPerManager = 10
	if activeInstanceCount >= maxInstancesPerManager {
		return coreerror.NewInvalidDataError("manager subscription has reached the maximum of %d game instances", maxInstancesPerManager)
	}

	return nil
}

func (m *Domain) validateGameInstanceRecForUpdate(rec *game_record.GameInstance) error {
	return validateGameInstanceRec(rec, true)
}

func validateGameInstanceRec(rec *game_record.GameInstance, requireID bool) error {
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

	if err := domain.ValidateUUIDField(game_record.FieldGameInstanceGameSubscriptionID, rec.GameSubscriptionID); err != nil {
		return err
	}

	if rec.Status == "" {
		rec.Status = game_record.GameInstanceStatusCreated
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
