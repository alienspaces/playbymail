package harness

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const (
	GameOneRef = "game-one"
	GameTwoRef = "game-two"

	AccountOneRef   = "account-one"
	AccountTwoRef   = "account-two"
	AccountThreeRef = "account-three"

	GameLocationOneRef   = "game-location-one"
	GameLocationTwoRef   = "game-location-two"
	GameLocationThreeRef = "game-location-three"
	GameLocationFourRef  = "game-location-four"

	GameLocationLinkOneRef   = "game-location-link-one"
	GameLocationLinkTwoRef   = "game-location-link-two"
	GameLocationLinkThreeRef = "game-location-link-three"
	GameLocationLinkFourRef  = "game-location-link-four"

	GameLocationLinkRequirementOneRef   = "game-location-link-requirement-one"
	GameLocationLinkRequirementTwoRef   = "game-location-link-requirement-two"
	GameLocationLinkRequirementThreeRef = "game-location-link-requirement-three"
	GameLocationLinkRequirementFourRef  = "game-location-link-requirement-four"

	GameItemOneRef   = "game-item-one"
	GameItemTwoRef   = "game-item-two"
	GameItemThreeRef = "game-item-three"
	GameItemFourRef  = "game-item-four"

	GameCreatureOneRef = "game-creature-one"
	GameCreatureTwoRef = "game-creature-two"

	GameCharacterOneRef   = "game-character-one"
	GameCharacterTwoRef   = "game-character-two"
	GameCharacterThreeRef = "game-character-three"

	GameInstanceOneRef   = "game-instance-one"
	GameInstanceTwoRef   = "game-instance-two"
	GameInstanceCleanRef = "game-instance-clean"

	GameInstanceParameterOneRef = "game-instance-parameter-one"

	GameItemInstanceOneRef = "game-item-instance-one"

	GameLocationInstanceOneRef = "game-location-instance-one"
	GameLocationInstanceTwoRef = "game-location-instance-two"

	GameCreatureInstanceOneRef = "game-creature-instance-one"

	GameCharacterInstanceOneRef = "game-character-instance-one"

	GameSubscriptionOneRef = "game-subscription-one"

	GameAdministrationOneRef = "game-administration-one"
)

// DataConfig -
type DataConfig struct {
	GameConfigs    []GameConfig
	AccountConfigs []AccountConfig
}

type GameConfig struct {
	Reference                 string // Reference to the game record
	Record                    *game_record.Game
	GameLocationConfigs       []GameLocationConfig     // Locations associated with this game
	GameLocationLinkConfigs   []GameLocationLinkConfig // Links associated with this game
	GameItemConfigs           []GameItemConfig
	GameCreatureConfigs       []GameCreatureConfig
	GameCharacterConfigs      []GameCharacterConfig
	GameInstanceConfigs       []GameInstanceConfig
	GameSubscriptionConfigs   []GameSubscriptionConfig
	GameAdministrationConfigs []GameAdministrationConfig
}

type GameCharacterConfig struct {
	Reference  string // Reference to the game_character record
	AccountRef string // Reference to the account
	Record     *adventure_game_record.AdventureGameCharacter
}

type GameItemConfig struct {
	Reference string // Reference to the game_item record
	Record    *adventure_game_record.AdventureGameItem
}

type GameCreatureConfig struct {
	Reference string // Reference to the game_creature record
	Record    *adventure_game_record.AdventureGameCreature
}

type AccountConfig struct {
	Reference string // Reference to the account record
	Record    *account_record.Account
}

type GameLocationConfig struct {
	Reference string // Reference to the game_location record
	Record    *adventure_game_record.AdventureGameLocation
}

type GameLocationLinkConfig struct {
	Reference                          string // Reference to the game_location_link record
	FromLocationRef                    string // Reference to the from location
	ToLocationRef                      string // Reference to the to location
	Record                             *adventure_game_record.AdventureGameLocationLink
	GameLocationLinkRequirementConfigs []GameLocationLinkRequirementConfig
}

type GameLocationLinkRequirementConfig struct {
	Reference   string // Reference to the game_location_link_requirement record
	GameItemRef string // Reference to the game_item
	Record      *adventure_game_record.AdventureGameLocationLinkRequirement
}

type GameInstanceConfig struct {
	Reference                    string // Reference to the game_instance record
	Record                       *game_record.GameInstance
	GameInstanceParameterConfigs []GameInstanceParameterConfig
	GameLocationInstanceConfigs  []GameLocationInstanceConfig
	GameItemInstanceConfigs      []GameItemInstanceConfig
	GameCreatureInstanceConfigs  []GameCreatureInstanceConfig
	GameCharacterInstanceConfigs []GameCharacterInstanceConfig
}

type GameLocationInstanceConfig struct {
	Reference       string // Reference to the game_location_instance record
	GameLocationRef string // Reference to the game_location (required)
	Record          *adventure_game_record.AdventureGameLocationInstance
}

type GameCreatureInstanceConfig struct {
	Reference       string // Reference to the game_creature_instance record
	GameCreatureRef string // Reference to the game_creature (required)
	GameLocationRef string // Reference to the game_location (required)
	Record          *adventure_game_record.AdventureGameCreatureInstance
}

type GameCharacterInstanceConfig struct {
	Reference        string // Reference to the game_character_instance record
	GameCharacterRef string // Reference to the game_character (required)
	GameLocationRef  string // Reference to the game_location (optional)
	Record           *adventure_game_record.AdventureGameCharacterInstance
}

type GameItemInstanceConfig struct {
	Reference        string // Reference to the game_item_instance record
	GameItemRef      string // Reference to the game_item (required)
	GameLocationRef  string // Reference to the game_location (optional)
	GameCharacterRef string // Reference to the game_character (optional)
	GameCreatureRef  string // Reference to the game_creature (optional)

	// TODO: Must be assigned to a location, a character, or a creature

	Record *adventure_game_record.AdventureGameItemInstance
}

type GameSubscriptionConfig struct {
	Reference        string // Reference to the game_subscription record
	AccountRef       string // Reference to the account
	SubscriptionType string // Type of subscription (Player, Manager, Collaborator)
	Record           *game_record.GameSubscription
}

type GameAdministrationConfig struct {
	Reference           string // Reference to the game_administration record
	AccountRef          string // Reference to the account
	GrantedByAccountRef string // Reference to the account that granted the administration rights
	Record              *game_record.GameAdministration
}

// DefaultDataConfig -
func DefaultDataConfig() DataConfig {
	return DataConfig{
		GameConfigs: []GameConfig{
			{
				Reference: GameOneRef,
				Record: &game_record.Game{
					Name:              UniqueName("Default Game One"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168, // 1 week
				},
				GameItemConfigs: []GameItemConfig{
					{
						Reference: GameItemOneRef,
						Record: &adventure_game_record.AdventureGameItem{
							Name:        UniqueName("Default Item One"),
							Description: "Default item one for handler tests",
						},
					},
					{
						Reference: GameItemTwoRef,
						Record: &adventure_game_record.AdventureGameItem{
							Name:        UniqueName("Default Item Two"),
							Description: "Default item two for handler tests",
						},
					},
				},
				GameLocationConfigs: []GameLocationConfig{
					{
						Reference: GameLocationOneRef,
						Record: &adventure_game_record.AdventureGameLocation{
							Name:        UniqueName("Default Location One"),
							Description: "Default location one for handler tests",
						},
					},
					{
						Reference: GameLocationTwoRef,
						Record: &adventure_game_record.AdventureGameLocation{
							Name:        UniqueName("Default Location Two"),
							Description: "Default location two for handler tests",
						},
					},
				},
				GameLocationLinkConfigs: []GameLocationLinkConfig{
					{
						Reference:       GameLocationLinkOneRef,
						FromLocationRef: GameLocationOneRef,
						ToLocationRef:   GameLocationTwoRef,
						Record: &adventure_game_record.AdventureGameLocationLink{
							Name:        UniqueName("The Red Door"),
							Description: "Travel by boat to the swamp of the long forgotten Frog God",
						},
						GameLocationLinkRequirementConfigs: []GameLocationLinkRequirementConfig{
							{
								Reference:   GameLocationLinkRequirementOneRef,
								GameItemRef: GameItemOneRef,
								Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
									Quantity: 1,
								},
							},
						},
					},
				},
				GameCreatureConfigs: []GameCreatureConfig{
					{
						Reference: GameCreatureOneRef,
						Record: &adventure_game_record.AdventureGameCreature{
							Name: UniqueName("Default Creature One"),
						},
					},
					{
						Reference: GameCreatureTwoRef,
						Record: &adventure_game_record.AdventureGameCreature{
							Name: UniqueName("Default Creature Two"),
						},
					},
				},
				GameCharacterConfigs: []GameCharacterConfig{
					{
						Reference:  GameCharacterOneRef,
						AccountRef: AccountOneRef,
						Record: &adventure_game_record.AdventureGameCharacter{
							Name: UniqueName("Default Character One"),
						},
					},
					{
						Reference:  GameCharacterTwoRef,
						AccountRef: AccountTwoRef,
						Record: &adventure_game_record.AdventureGameCharacter{
							Name: UniqueName("Default Character Two"),
						},
					},
				},
				GameSubscriptionConfigs: []GameSubscriptionConfig{
					{
						Reference:        GameSubscriptionOneRef,
						AccountRef:       AccountOneRef,
						SubscriptionType: "Player",
						Record:           &game_record.GameSubscription{},
					},
				},
				GameAdministrationConfigs: []GameAdministrationConfig{
					{
						Reference:           "game-administration-one",
						AccountRef:          AccountOneRef,
						GrantedByAccountRef: AccountOneRef,
						Record:              &game_record.GameAdministration{},
					},
				},
				// Default game instance with a location and an item assigned to the location
				GameInstanceConfigs: []GameInstanceConfig{
					{
						Reference: GameInstanceOneRef,
						Record:    &game_record.GameInstance{},
						GameInstanceParameterConfigs: []GameInstanceParameterConfig{
							{
								Reference: GameInstanceParameterOneRef,
								Record: &game_record.GameInstanceParameter{
									ParameterKey:   domain.AdventureGameParameterCharacterLives,
									ParameterValue: nullstring.FromString("3"),
								},
							},
						},
						GameLocationInstanceConfigs: []GameLocationInstanceConfig{
							{
								Reference:       GameLocationInstanceOneRef,
								GameLocationRef: GameLocationOneRef,
								Record:          &adventure_game_record.AdventureGameLocationInstance{},
							},
							{
								Reference:       GameLocationInstanceTwoRef,
								GameLocationRef: GameLocationTwoRef,
								Record:          &adventure_game_record.AdventureGameLocationInstance{},
							},
						},
						GameItemInstanceConfigs: []GameItemInstanceConfig{
							{
								Reference:       GameItemInstanceOneRef,
								GameItemRef:     GameItemOneRef,
								GameLocationRef: GameLocationOneRef, // Assign to location
								Record:          &adventure_game_record.AdventureGameItemInstance{},
							},
						},
						GameCreatureInstanceConfigs: []GameCreatureInstanceConfig{
							{
								Reference:       GameCreatureInstanceOneRef,
								GameCreatureRef: GameCreatureOneRef,
								GameLocationRef: GameLocationOneRef, // Assign to location
								Record:          &adventure_game_record.AdventureGameCreatureInstance{},
							},
						},
						GameCharacterInstanceConfigs: []GameCharacterInstanceConfig{
							{
								Reference:        GameCharacterInstanceOneRef,
								GameCharacterRef: GameCharacterOneRef,
								GameLocationRef:  GameLocationOneRef,
								Record:           &adventure_game_record.AdventureGameCharacterInstance{},
							},
						},
					},
					// Clean game instance with no parameters for testing
					{
						Reference: GameInstanceCleanRef,
						Record:    &game_record.GameInstance{},
						// No parameters, no instances - clean slate for testing
					},
				},
			},
		},
		AccountConfigs: []AccountConfig{
			{
				Reference: AccountOneRef,
				Record: &account_record.Account{
					Email: UniqueEmail("default-account-one@example.com"),
				},
			},
			{
				Reference: AccountTwoRef,
				Record: &account_record.Account{
					Email: UniqueEmail("default-account-two@example.com"),
				},
			},
			{
				Reference: AccountThreeRef,
				Record: &account_record.Account{
					Email: UniqueEmail("default-account-three@example.com"),
				},
			},
		},
	}
}
