package harness

import "gitlab.com/alienspaces/playbymail/internal/record"

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

	GameInstanceOneRef = "game-instance-one"

	GameItemInstanceOneRef = "game-item-instance-one"

	GameLocationInstanceOneRef = "game-location-instance-one"
	GameLocationInstanceTwoRef = "game-location-instance-two"

	GameCreatureInstanceOneRef = "game-creature-instance-one"

	GameCharacterInstanceOneRef = "game-character-instance-one"
)

// DataConfig -
type DataConfig struct {
	GameConfigs    []GameConfig
	AccountConfigs []AccountConfig
}

type GameConfig struct {
	Reference               string // Reference to the game record
	Record                  *record.Game
	GameLocationConfigs     []GameLocationConfig     // Locations associated with this game
	GameLocationLinkConfigs []GameLocationLinkConfig // Links associated with this game
	GameItemConfigs         []GameItemConfig
	GameCreatureConfigs     []GameCreatureConfig
	GameCharacterConfigs    []GameCharacterConfig
	GameInstanceConfigs     []GameInstanceConfig
}

type GameCharacterConfig struct {
	Reference  string // Reference to the game_character record
	AccountRef string // Reference to the account
	Record     *record.GameCharacter
}

type GameItemConfig struct {
	Reference string // Reference to the game_item record
	Record    *record.GameItem
}

type GameCreatureConfig struct {
	Reference string // Reference to the game_creature record
	Record    *record.GameCreature
}

type AccountConfig struct {
	Reference string // Reference to the account record
	Record    *record.Account
}

type GameLocationConfig struct {
	Reference string // Reference to the game_location record
	Record    *record.GameLocation
}

type GameLocationLinkConfig struct {
	Reference                          string // Reference to the game_location_link record
	FromLocationRef                    string // Reference to the from location
	ToLocationRef                      string // Reference to the to location
	Record                             *record.GameLocationLink
	GameLocationLinkRequirementConfigs []GameLocationLinkRequirementConfig
}

type GameLocationLinkRequirementConfig struct {
	Reference   string // Reference to the game_location_link_requirement record
	GameItemRef string // Reference to the game_item
	Record      *record.GameLocationLinkRequirement
}

type GameInstanceConfig struct {
	Reference                    string // Reference to the game_instance record
	Record                       *record.GameInstance
	GameLocationInstanceConfigs  []GameLocationInstanceConfig
	GameItemInstanceConfigs      []GameItemInstanceConfig
	GameCreatureInstanceConfigs  []GameCreatureInstanceConfig
	GameCharacterInstanceConfigs []GameCharacterInstanceConfig
}

type GameLocationInstanceConfig struct {
	Reference       string // Reference to the game_location_instance record
	GameLocationRef string // Reference to the game_location (required)
	Record          *record.GameLocationInstance
}

type GameCreatureInstanceConfig struct {
	Reference       string // Reference to the game_creature_instance record
	GameCreatureRef string // Reference to the game_creature (required)
	GameLocationRef string // Reference to the game_location (required)
	Record          *record.GameCreatureInstance
}

type GameCharacterInstanceConfig struct {
	Reference        string // Reference to the game_character_instance record
	GameCharacterRef string // Reference to the game_character (required)
	GameLocationRef  string // Reference to the game_location (optional)
	Record           *record.GameCharacterInstance
}

type GameItemInstanceConfig struct {
	Reference        string // Reference to the game_item_instance record
	GameItemRef      string // Reference to the game_item (required)
	GameLocationRef  string // Reference to the game_location (optional)
	GameCharacterRef string // Reference to the game_character (optional)
	GameCreatureRef  string // Reference to the game_creature (optional)

	// TODO: Must be assigned to a location, a character, or a creature

	Record *record.GameItemInstance
}

// DefaultDataConfig -
func DefaultDataConfig() DataConfig {
	return DataConfig{
		GameConfigs: []GameConfig{
			{
				Reference: GameOneRef,
				Record: &record.Game{
					Name:     UniqueName("Default Game One"),
					GameType: record.GameTypeAdventure,
				},
				GameItemConfigs: []GameItemConfig{
					{
						Reference: GameItemOneRef,
						Record: &record.GameItem{
							Name:        UniqueName("Default Item One"),
							Description: "Default item one for handler tests",
						},
					},
					{
						Reference: GameItemTwoRef,
						Record: &record.GameItem{
							Name:        UniqueName("Default Item Two"),
							Description: "Default item two for handler tests",
						},
					},
				},
				GameLocationConfigs: []GameLocationConfig{
					{
						Reference: GameLocationOneRef,
						Record: &record.GameLocation{
							Name:        UniqueName("Default Location One"),
							Description: "Default location one for handler tests",
						},
					},
					{
						Reference: GameLocationTwoRef,
						Record: &record.GameLocation{
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
						Record: &record.GameLocationLink{
							Name:        UniqueName("The Red Door"),
							Description: "Travel by boat to the swamp of the long forgotten Frog God",
						},
						GameLocationLinkRequirementConfigs: []GameLocationLinkRequirementConfig{
							{
								Reference:   GameLocationLinkRequirementOneRef,
								GameItemRef: GameItemOneRef,
								Record: &record.GameLocationLinkRequirement{
									Quantity: 1,
								},
							},
						},
					},
				},
				GameCreatureConfigs: []GameCreatureConfig{
					{
						Reference: GameCreatureOneRef,
						Record: &record.GameCreature{
							Name: UniqueName("Default Creature One"),
						},
					},
					{
						Reference: GameCreatureTwoRef,
						Record: &record.GameCreature{
							Name: UniqueName("Default Creature Two"),
						},
					},
				},
				GameCharacterConfigs: []GameCharacterConfig{
					{
						Reference:  GameCharacterOneRef,
						AccountRef: AccountOneRef,
						Record: &record.GameCharacter{
							Name: UniqueName("Default Character One"),
						},
					},
					{
						Reference:  GameCharacterTwoRef,
						AccountRef: AccountTwoRef,
						Record: &record.GameCharacter{
							Name: UniqueName("Default Character Two"),
						},
					},
				},
				// Default game instance with a location and an item assigned to the location
				GameInstanceConfigs: []GameInstanceConfig{
					{
						Reference: GameInstanceOneRef,
						Record:    &record.GameInstance{},
						GameLocationInstanceConfigs: []GameLocationInstanceConfig{
							{
								Reference:       GameLocationInstanceOneRef,
								GameLocationRef: GameLocationOneRef,
								Record:          &record.GameLocationInstance{},
							},
							{
								Reference:       GameLocationInstanceTwoRef,
								GameLocationRef: GameLocationTwoRef,
								Record:          &record.GameLocationInstance{},
							},
						},
						GameItemInstanceConfigs: []GameItemInstanceConfig{
							{
								Reference:       GameItemInstanceOneRef,
								GameItemRef:     GameItemOneRef,
								GameLocationRef: GameLocationOneRef, // Assign to location
								Record:          &record.GameItemInstance{},
							},
						},
						GameCreatureInstanceConfigs: []GameCreatureInstanceConfig{
							{
								Reference:       GameCreatureInstanceOneRef,
								GameCreatureRef: GameCreatureOneRef,
								GameLocationRef: GameLocationOneRef, // Assign to location
								Record:          &record.GameCreatureInstance{},
							},
						},
						GameCharacterInstanceConfigs: []GameCharacterInstanceConfig{
							{
								Reference:        GameCharacterInstanceOneRef,
								GameCharacterRef: GameCharacterOneRef,
								GameLocationRef:  GameLocationOneRef,
								Record:           &record.GameCharacterInstance{},
							},
						},
					},
				},
			},
		},
		AccountConfigs: []AccountConfig{
			{
				Reference: AccountOneRef,
				Record: &record.Account{
					Email: UniqueEmail("default-account-one@example.com"),
				},
			},
			{
				Reference: AccountTwoRef,
				Record: &record.Account{
					Email: UniqueEmail("default-account-two@example.com"),
				},
			},
			{
				Reference: AccountThreeRef,
				Record: &record.Account{
					Email: UniqueEmail("default-account-three@example.com"),
				},
			},
		},
	}
}
