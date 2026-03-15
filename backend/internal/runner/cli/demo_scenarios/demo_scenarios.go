package demo_scenarios

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
)

const (
	DemoAdventureGameName = "The Door Beneath the Staircase"
	DemoAdventureGameRef  = "demo-adventure-game"

	DemoAccountDesignerRef = "demo-account-designer"
	DemoAccountManagerRef  = "demo-account-manager"

	DemoAccountDesignerEmail = "demo-designer@example.com"
	DemoAccountManagerEmail  = "demo-manager@example.com"

	DemoSubscriptionDesignerOneRef = "demo-subscription-designer-one"
	DemoSubscriptionDesignerTwoRef = "demo-subscription-designer-two"
	DemoSubscriptionManagerOneRef  = "demo-subscription-manager-one"
	DemoSubscriptionManagerTwoRef  = "demo-subscription-manager-two"
)

// DemoAccountDefs defines the demo accounts needed by the CLI loader.
// Each entry maps a reference name to the account email. The CLI uses this
// to ensure accounts exist before the harness runs.
var DemoAccountDefs = []struct {
	Ref   string
	Email string
}{
	{Ref: DemoAccountDesignerRef, Email: DemoAccountDesignerEmail},
	{Ref: DemoAccountManagerRef, Email: DemoAccountManagerEmail},
}

// AdventureGameConfig returns a standalone demo scenario exercising all adventure game features.
// Accounts are managed by the CLI -- AccountConfigs is empty. Subscriptions are top-level
// GameSubscriptionConfigs; the caller must set each Record.AccountID and Record.AccountUserID
// from the demo account records (e.g. from ensureDemoAccounts()) in order: first subscription
// uses first demo account, second uses second.
func AdventureGameConfig() harness.DataConfig {
	return harness.DataConfig{
		GameConfigs: adventureGameConfigs(),
		// AccountUserGameSubscriptionConfigs: []harness.AccountUserGameSubscriptionConfig{
		// 	{
		// 		Reference:        DemoSubscriptionDesignerOneRef,
		// 		GameRef:          DemoAdventureGameRef,
		// 		SubscriptionType: game_record.GameSubscriptionTypeDesigner,
		// 		Record:           &game_record.GameSubscription{},
		// 	},
		// 	{
		// 		Reference:        DemoSubscriptionManagerOneRef,
		// 		GameRef:          DemoAdventureGameRef,
		// 		GameInstanceRefs: []string{DemoInstanceOneRef},
		// 		SubscriptionType: game_record.GameSubscriptionTypeManager,
		// 		Record:           &game_record.GameSubscription{},
		// 	},
		// },
	}
}
