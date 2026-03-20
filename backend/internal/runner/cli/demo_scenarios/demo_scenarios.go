package demo_scenarios

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
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
// Accounts are managed by the CLI -- AccountConfigs is empty. The CLI (load_game_data.go) populates
// each subscription Record.AccountID and Record.AccountUserID from ensureDemoAccounts() in order:
// subscription[0] uses demoRecs.AccountUsers[0] (designer), subscription[1] uses [1] (manager).
func AdventureGameConfig() harness.DataConfig {
	return harness.DataConfig{
		GameConfigs: adventureGameConfigs(),
		AccountUserGameSubscriptionConfigs: []harness.AccountUserGameSubscriptionConfig{
			{
				Reference:        DemoSubscriptionDesignerOneRef,
				GameRef:          DemoAdventureGameRef,
				SubscriptionType: game_record.GameSubscriptionTypeDesigner,
				Record:           &game_record.GameSubscription{},
			},
			{
				Reference:        DemoSubscriptionManagerOneRef,
				GameRef:          DemoAdventureGameRef,
				SubscriptionType: game_record.GameSubscriptionTypeManager,
				Record:           &game_record.GameSubscription{},
				GameInstanceConfigs: []harness.GameInstanceConfig{
					{
						Reference: DemoInstanceOneRef,
						Record: &game_record.GameInstance{
							DeliveryEmail: true,
						},
					},
				},
			},
		},
	}
}
