package domain

import (
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

	// Only gather existing subscriptions for create operations (when currRec is nil)
	// We check for existing subscriptions based on AccountID (Tenant) OR AccountUserID (User)
	if currRec == nil {
		var params []coresql.Param
		if nextRec.AccountID.Valid && nextRec.AccountID.String != "" {
			params = []coresql.Param{
				{Col: account_record.FieldAccountSubscriptionAccountID, Val: nextRec.AccountID.String},
				{Col: account_record.FieldAccountSubscriptionStatus, Val: account_record.AccountSubscriptionStatusActive},
			}
		} else if nextRec.AccountUserID.Valid && nextRec.AccountUserID.String != "" {
			params = []coresql.Param{
				{Col: account_record.FieldAccountSubscriptionAccountUserID, Val: nextRec.AccountUserID.String},
				{Col: account_record.FieldAccountSubscriptionStatus, Val: account_record.AccountSubscriptionStatusActive},
			}
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

	// Set default subscription_period if not provided
	if rec.SubscriptionPeriod == "" {
		rec.SubscriptionPeriod = account_record.AccountSubscriptionPeriodEternal
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

	// Validate Access IDs based on Subscription Type
	isTenantSub := rec.SubscriptionType == account_record.AccountSubscriptionTypeBasicGameDesigner ||
		rec.SubscriptionType == account_record.AccountSubscriptionTypeProfessionalGameDesigner ||
		rec.SubscriptionType == account_record.AccountSubscriptionTypeBasicManager ||
		rec.SubscriptionType == account_record.AccountSubscriptionTypeProfessionalManager

	isUserSub := rec.SubscriptionType == account_record.AccountSubscriptionTypeBasicPlayer ||
		rec.SubscriptionType == account_record.AccountSubscriptionTypeProfessionalPlayer

	if isTenantSub {
		if !rec.AccountID.Valid || rec.AccountID.String == "" {
			return InvalidField(account_record.FieldAccountSubscriptionAccountID, rec.AccountID.String, "account_id is required for designer/manager subscriptions")
		}
		if err := domain.ValidateUUIDField(account_record.FieldAccountSubscriptionAccountID, rec.AccountID.String); err != nil {
			return err
		}
	} else if isUserSub {
		if !rec.AccountUserID.Valid || rec.AccountUserID.String == "" {
			return InvalidField(account_record.FieldAccountSubscriptionAccountUserID, rec.AccountUserID.String, "account_user_id is required for player subscriptions")
		}
		if err := domain.ValidateUUIDField(account_record.FieldAccountSubscriptionAccountUserID, rec.AccountUserID.String); err != nil {
			return err
		}
	} else {
		return InvalidField(account_record.FieldAccountSubscriptionSubscriptionType, rec.SubscriptionType, "subscription type is not valid")
	}

	if rec.SubscriptionPeriod != account_record.AccountSubscriptionPeriodMonth &&
		rec.SubscriptionPeriod != account_record.AccountSubscriptionPeriodYear &&
		rec.SubscriptionPeriod != account_record.AccountSubscriptionPeriodEternal {
		return InvalidField(account_record.FieldAccountSubscriptionSubscriptionPeriod, rec.SubscriptionPeriod, "subscription period is not valid")
	}

	if rec.Status != account_record.AccountSubscriptionStatusActive &&
		rec.Status != account_record.AccountSubscriptionStatusExpired {
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
