package turn_sheet_processor

import (
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// lanceInstanceAccountUser resolves the account_user record for a player-owned
// lance instance by walking the subscription chain:
//
//	mecha_lance_instance.game_subscription_instance_id
//	  -> game_subscription_instance.game_subscription_id
//	  -> game_subscription.account_user_id
func lanceInstanceAccountUser(d *domain.Domain, lanceInstance *mecha_record.MechaLanceInstance) (*account_record.AccountUser, error) {
	if !lanceInstance.GameSubscriptionInstanceID.Valid {
		return nil, fmt.Errorf("lance instance >%s< has no game_subscription_instance_id", lanceInstance.ID)
	}

	subInstRec, err := d.GetGameSubscriptionInstanceRec(lanceInstance.GameSubscriptionInstanceID.String, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get game subscription instance: %w", err)
	}

	subRecs, err := d.GetManyGameSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionID, Val: subInstRec.GameSubscriptionID},
		},
		Limit: 1,
	})
	if err != nil || len(subRecs) == 0 {
		return nil, fmt.Errorf("failed to get game subscription for instance >%s<: %w", subInstRec.GameSubscriptionID, err)
	}

	accountUserRec, err := d.GetAccountUserRec(subRecs[0].AccountUserID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account user >%s<: %w", subRecs[0].AccountUserID, err)
	}

	return accountUserRec, nil
}
