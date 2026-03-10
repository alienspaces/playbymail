package domain

import (
	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

type validateAccountSubscriptionArgs struct {
	nextRec               *account_record.AccountSubscription
	currRec               *account_record.AccountSubscription
	existingSubscriptions []*account_record.AccountSubscription
}

func (m *Domain) populateAccountSubscriptionValidateArgs(currRec, nextRec *account_record.AccountSubscription) (*validateAccountSubscriptionArgs, error) {
	args := &validateAccountSubscriptionArgs{
		currRec: currRec,
		nextRec: nextRec,
	}

	// Only gather existing subscriptions for create operations (when currRec is nil).
	// All subscription types require account_user_id; check for duplicates by user and status.
	if currRec == nil && nextRec.AccountUserID.Valid && nextRec.AccountUserID.String != "" {
		params := []coresql.Param{
			{Col: account_record.FieldAccountSubscriptionAccountUserID, Val: nextRec.AccountUserID.String},
			{Col: account_record.FieldAccountSubscriptionStatus, Val: account_record.AccountSubscriptionStatusActive},
		}
		if len(params) > 0 {
			existingSubs, err := m.GetManyAccountSubscriptionRecs(&coresql.Options{
				Params: params,
			})
			if err != nil {
				return nil, err
			}
			args.existingSubscriptions = existingSubs
		}
	}

	return args, nil
}

func (m *Domain) validateAccountSubscriptionRecForCreate(rec *account_record.AccountSubscription) error {
	args, err := m.populateAccountSubscriptionValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountSubscriptionRecForCreate(args)
}

func (m *Domain) validateAccountSubscriptionRecForUpdate(currRec, nextRec *account_record.AccountSubscription) error {
	args, err := m.populateAccountSubscriptionValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAccountSubscriptionRecForUpdate(args)
}

func (m *Domain) validateAccountSubscriptionRecForDelete(rec *account_record.AccountSubscription) error {
	args, err := m.populateAccountSubscriptionValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountSubscriptionRecForDelete(args)
}

func validateAccountSubscriptionRecForCreate(args *validateAccountSubscriptionArgs) error {
	rec := args.nextRec

	if rec.SubscriptionPeriod == "" {
		return InvalidField(account_record.FieldAccountSubscriptionSubscriptionPeriod, "", "subscription_period is required")
	}

	// Check if a subscription of the same type already exists
	for _, existingSub := range args.existingSubscriptions {
		if existingSub.SubscriptionType == rec.SubscriptionType {
			return InvalidField(account_record.FieldAccountSubscriptionSubscriptionType, rec.SubscriptionType, "account already has an active subscription of this type")
		}
	}

	return validateAccountSubscriptionRec(args)
}

func validateAccountSubscriptionRecForUpdate(args *validateAccountSubscriptionArgs) error {
	nextRec := args.nextRec
	currRec := args.currRec

	// Prevent modification of subscription type
	if currRec.SubscriptionType != nextRec.SubscriptionType {
		return InvalidField(account_record.FieldAccountSubscriptionSubscriptionType, nextRec.SubscriptionType, "subscription type cannot be modified")
	}

	return validateAccountSubscriptionRec(args)
}

func validateAccountSubscriptionRec(args *validateAccountSubscriptionArgs) error {
	rec := args.nextRec

	// All subscription types require an account user.
	if !rec.AccountUserID.Valid || rec.AccountUserID.String == "" {
		return InvalidField(account_record.FieldAccountSubscriptionAccountUserID, "", "account_user_id is required for all subscription types")
	}
	if err := domain.ValidateUUIDField(account_record.FieldAccountSubscriptionAccountUserID, rec.AccountUserID.String); err != nil {
		return err
	}

	// All subscription types require account_id and account_user_id.
	if !rec.AccountID.Valid || rec.AccountID.String == "" {
		return InvalidField(account_record.FieldAccountSubscriptionAccountID, "", "account_id is required for all subscription types")
	}
	if err := domain.ValidateUUIDField(account_record.FieldAccountSubscriptionAccountID, rec.AccountID.String); err != nil {
		return err
	}

	// Reject unknown subscription types
	subscriptionTypeSet := set.New(account_record.AccountSubscriptionTypeBasicGameDesigner, account_record.AccountSubscriptionTypeProfessionalGameDesigner, account_record.AccountSubscriptionTypeBasicManager, account_record.AccountSubscriptionTypeProfessionalManager, account_record.AccountSubscriptionTypeBasicPlayer, account_record.AccountSubscriptionTypeProfessionalPlayer)
	if !subscriptionTypeSet.Has(rec.SubscriptionType) {
		return InvalidField(account_record.FieldAccountSubscriptionSubscriptionType, rec.SubscriptionType, "subscription type is not valid")
	}

	subscriptionPeriodSet := set.New(account_record.AccountSubscriptionPeriodMonth, account_record.AccountSubscriptionPeriodYear, account_record.AccountSubscriptionPeriodEternal)
	if !subscriptionPeriodSet.Has(rec.SubscriptionPeriod) {
		return InvalidField(account_record.FieldAccountSubscriptionSubscriptionPeriod, rec.SubscriptionPeriod, "subscription period is not valid")
	}

	statusSet := set.New(account_record.AccountSubscriptionStatusActive, account_record.AccountSubscriptionStatusExpired)
	if !statusSet.Has(rec.Status) {
		return InvalidField(account_record.FieldAccountSubscriptionStatus, rec.Status, "status is not valid")
	}

	return nil
}

func validateAccountSubscriptionRecForDelete(args *validateAccountSubscriptionArgs) error {
	rec := args.nextRec

	if err := domain.ValidateUUIDField(account_record.FieldAccountSubscriptionID, rec.ID); err != nil {
		return err
	}

	return nil
}
