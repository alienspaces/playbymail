package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const (
	GameOneRef   = "game-one"
	GameTwoRef   = "game-two"
	GameDraftRef = "game-draft"

	StandardAccountRef    = "account-standard"
	ProPlayerAccountRef   = "account-pro-player"
	ProDesignerAccountRef = "account-pro-designer"
	ProManagerAccountRef  = "account-pro-manager"

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

	// Designer subscriptions
	GameSubscriptionDesignerOneRef = "game-subscription-designer-one"
	GameSubscriptionDesignerTwoRef = "game-subscription-designer-two"

	// Manager subscriptions
	GameSubscriptionManagerOneRef = "game-subscription-manager-one"
	GameSubscriptionManagerTwoRef = "game-subscription-manager-two"

	// Player subscriptions
	GameSubscriptionPlayerOneRef   = "game-subscription-player-one"
	GameSubscriptionPlayerTwoRef   = "game-subscription-player-two"
	GameSubscriptionPlayerThreeRef = "game-subscription-player-three"
	GameSubscriptionPlayerFourRef  = "game-subscription-player-four"
	GameSubscriptionPlayerFiveRef  = "game-subscription-player-five"
	GameSubscriptionPlayerSixRef   = "game-subscription-player-six"
	GameSubscriptionPlayerSevenRef = "game-subscription-player-seven"
	GameSubscriptionPlayerEightRef = "game-subscription-player-eight"
	GameSubscriptionPlayerNineRef  = "game-subscription-player-nine"

	// Turn sheets
	GameTurnSheetOneRef   = "game-turn-sheet-one"
	GameTurnSheetTwoRef   = "game-turn-sheet-two"
	GameTurnSheetThreeRef = "game-turn-sheet-three"
	GameTurnSheetFourRef  = "game-turn-sheet-four"

	GameImageJoinGameRef  = "game-image-join-game"
	GameImageInventoryRef = "game-image-inventory"
)

// DataConfig -
type DataConfig struct {
	GameConfigs    []GameConfig
	AccountConfigs []AccountConfig
}

type GameConfig struct {
	Reference                        string // Reference to the game record
	Record                           *game_record.Game
	GameInstanceConfigs              []GameInstanceConfig
	AdventureGameLocationConfigs     []AdventureGameLocationConfig     // Locations associated with this game
	AdventureGameLocationLinkConfigs []AdventureGameLocationLinkConfig // Links associated with this game
	AdventureGameItemConfigs         []AdventureGameItemConfig
	AdventureGameCreatureConfigs     []AdventureGameCreatureConfig
	AdventureGameCharacterConfigs    []AdventureGameCharacterConfig
	GameImageConfigs []GameImageConfig // Game image configurations

	// Deprecated: Use GameImageConfigs instead
	//
	// Background image file path for game-level turn sheet backgrounds.
	// If set, creates a game_image record with type turn_sheet_background.
	BackgroundImagePath string // Path to image file relative to testdata directory
}

type AccountConfig struct {
	Reference                  string // Reference to the account record
	Record                     *account_record.AccountUser
	GameSubscriptionConfigs    []GameSubscriptionConfig    // Game subscriptions for this account
	AccountSubscriptionConfigs []AccountSubscriptionConfig // Account subscriptions for this account
}

type AccountSubscriptionConfig struct {
	SubscriptionType string // e.g., AccountSubscriptionTypeProfessionalPlayer
	Record           *account_record.AccountSubscription
}

type GameInstanceConfig struct {
	Reference                             string // Reference to the game_instance record
	Record                                *game_record.GameInstance
	GameInstanceParameterConfigs          []GameInstanceParameterConfig
	AdventureGameLocationInstanceConfigs  []AdventureGameLocationInstanceConfig
	AdventureGameItemInstanceConfigs      []AdventureGameItemInstanceConfig
	AdventureGameCreatureInstanceConfigs  []AdventureGameCreatureInstanceConfig
	AdventureGameCharacterInstanceConfigs []AdventureGameCharacterInstanceConfig
	GameTurnConfigs                       []GameTurnConfig
}

type GameSubscriptionConfig struct {
	Reference        string   // Reference to the game_subscription record
	GameRef          string   // Reference to the game (required)
	GameInstanceRefs []string // References to game_instances (optional, can link multiple instances)
	SubscriptionType string   // Type of subscription (Player, Manager, Designer)
	InstanceLimit    *int32   // Instance limit (nil = unlimited)
	Record           *game_record.GameSubscription
	// JoinGameScanData configures scan data for join game turn sheets
	// If provided, the harness will create a join game turn sheet and process it
	// For adventure games, use turn_sheet.AdventureGameJoinGameScanData
	// For other game types, use the appropriate game-specific scan data type
	JoinGameScanData any // Can be *turn_sheet.AdventureGameJoinGameScanData, etc.
}

type GameTurnSheetConfig struct {
	Reference        string // Reference to the game_turn_sheet record
	AccountRef       string // Reference to the account
	SheetType        string // Type of turn sheet (e.g., "location_choice", "join_game")
	ProcessingStatus string // Processing status (e.g., "pending", "processing", "completed")
	IsCompleted      bool   // Whether the turn sheet has been completed
	Record           *game_record.GameTurnSheet
	// ScanDataConfig configures scan data to apply to turn sheets created by job workers
	// If provided, the harness will apply this scan data after the turn sheet is created
	// For location choice turn sheets, use turn_sheet.LocationChoiceScanData
	// For other turn sheet types, add appropriate scan data types from turn_sheet package
	// Note: SheetData and SheetOrder are generated by job workers based on game state,
	// so they are not configurable here
	// Note: TurnNumber is determined by the parent GameTurnConfig, not specified here
	ScanDataConfig any // Can be *turn_sheet.LocationChoiceScanData, etc.
}

type GameTurnConfig struct {
	TurnNumber int

	// Adventure game specific turn sheet configurations
	AdventureGameTurnSheetConfigs []AdventureGameTurnSheetConfig

	// This is where additional turn sheet configurations can be added for other game types
}

// ------------------------------------------------------------
// Adventure game specific configuration
// ------------------------------------------------------------
type AdventureGameTurnSheetConfig struct {
	GameTurnSheetConfig      GameTurnSheetConfig
	GameCharacterInstanceRef string // Reference to the game_character_instance (required)
}

type AdventureGameCharacterConfig struct {
	Reference  string // Reference to the game_character record
	AccountRef string // Reference to the account
	Record     *adventure_game_record.AdventureGameCharacter
}

type AdventureGameItemConfig struct {
	Reference string // Reference to the game_item record
	Record    *adventure_game_record.AdventureGameItem
}

type AdventureGameCreatureConfig struct {
	Reference string // Reference to the game_creature record
	Record    *adventure_game_record.AdventureGameCreature
}

type AdventureGameLocationConfig struct {
	Reference string // Reference to the game_location record
	Record    *adventure_game_record.AdventureGameLocation
	// Background image file path for location-specific turn sheet backgrounds
	// If set, creates a game_image record with type turn_sheet_background
	// associated with this location
	// Path is relative to seed_images or testdata directories
	BackgroundImagePath string
}

type AdventureGameLocationLinkConfig struct {
	Reference                                   string // Reference to the game_location_link record
	FromLocationRef                             string // Reference to the from location
	ToLocationRef                               string // Reference to the to location
	Record                                      *adventure_game_record.AdventureGameLocationLink
	AdventureGameLocationLinkRequirementConfigs []AdventureGameLocationLinkRequirementConfig
}

type AdventureGameLocationLinkRequirementConfig struct {
	Reference string // Reference to the game_location_link_requirement record

	// Must be assigned to an item
	GameItemRef string // Reference to the game_item

	Record *adventure_game_record.AdventureGameLocationLinkRequirement
}

type AdventureGameLocationInstanceConfig struct {
	Reference string // Reference to the game_location_instance record

	// Must be assigned to a location
	GameLocationRef string // Reference to the game_location (required)

	Record *adventure_game_record.AdventureGameLocationInstance
}

type AdventureGameCreatureInstanceConfig struct {
	Reference string // Reference to the game_creature_instance record

	// Must be assigned to a creature and a location
	GameCreatureRef string // Reference to the game_creature (required)
	GameLocationRef string // Reference to the game_location (required)

	Record *adventure_game_record.AdventureGameCreatureInstance
}

type AdventureGameCharacterInstanceConfig struct {
	Reference string // Reference to the game_character_instance record

	// Must be assigned to a character and a location
	GameCharacterRef string // Reference to the game_character (required)
	GameLocationRef  string // Reference to the game_location (required)

	Record *adventure_game_record.AdventureGameCharacterInstance
}

type AdventureGameItemInstanceConfig struct {
	Reference string // Reference to the game_item_instance record

	// Must be assigned to an item
	GameItemRef string // Reference to the game_item (required)

	// Must be assigned to a location, a character, or a creature (one of these is required)
	GameLocationRef  string // Reference to the game_location
	GameCharacterRef string // Reference to the game_character
	GameCreatureRef  string // Reference to the game_creature

	Record *adventure_game_record.AdventureGameItemInstance
}

type GameImageConfig struct {
	Reference string // Reference to the game_image record
	// Path to image file relative to seed_images or testdata directories
	// (loaded at runtime). If set, the image is loaded from file.
	// If not set, the Record field must contain the image data.
	ImagePath     string
	TurnSheetType string // The turn sheet type for this image
	Record        *game_record.GameImage
}

// Helper methods for modifying DataConfig

// findGameInstanceConfig finds a game instance config by reference
func (dc *DataConfig) findGameInstanceConfig(gameInstanceRef string) (*GameInstanceConfig, error) {
	for i := range dc.GameConfigs {
		for j := range dc.GameConfigs[i].GameInstanceConfigs {
			if dc.GameConfigs[i].GameInstanceConfigs[j].Reference == gameInstanceRef {
				return &dc.GameConfigs[i].GameInstanceConfigs[j], nil
			}
		}
	}
	return nil, fmt.Errorf("game instance config with reference >%s< not found", gameInstanceRef)
}

// findGameTurnConfig finds a turn config for a given game instance and turn number
func (dc *DataConfig) findGameTurnConfig(gameInstanceRef string, turnNumber int) (*GameTurnConfig, error) {
	instanceConfig, err := dc.findGameInstanceConfig(gameInstanceRef)
	if err != nil {
		return nil, err
	}
	for i := range instanceConfig.GameTurnConfigs {
		if instanceConfig.GameTurnConfigs[i].TurnNumber == turnNumber {
			return &instanceConfig.GameTurnConfigs[i], nil
		}
	}
	return nil, nil
}

func normalizeTurnConfigs(configs []GameTurnConfig) []GameTurnConfig {
	norm := make([]GameTurnConfig, len(configs))
	prevTurn := 0
	for i := range configs {
		cfg := configs[i]
		if cfg.TurnNumber <= 0 {
			cfg.TurnNumber = prevTurn + 1
		}
		if cfg.TurnNumber <= prevTurn {
			cfg.TurnNumber = prevTurn + 1
		}
		norm[i] = cfg
		prevTurn = cfg.TurnNumber
	}
	return norm
}

// ReplaceGameTurnConfigs replaces the full set of turn configs for a game instance
func (dc *DataConfig) ReplaceGameTurnConfigs(gameInstanceRef string, configs []GameTurnConfig) error {
	instanceConfig, err := dc.findGameInstanceConfig(gameInstanceRef)
	if err != nil {
		return err
	}
	instanceConfig.GameTurnConfigs = normalizeTurnConfigs(configs)
	return nil
}

// AppendGameTurnConfigs appends additional turn configs to a game instance
func (dc *DataConfig) AppendGameTurnConfigs(gameInstanceRef string, configs []GameTurnConfig) error {
	instanceConfig, err := dc.findGameInstanceConfig(gameInstanceRef)
	if err != nil {
		return err
	}
	instanceConfig.GameTurnConfigs = append(instanceConfig.GameTurnConfigs, normalizeTurnConfigs(configs)...)
	return nil
}

// AppendAdventureGameTurnSheetConfigs appends configs to the first turn (for backwards compatibility)
func (dc *DataConfig) AppendAdventureGameTurnSheetConfigs(gameInstanceRef string, configs []AdventureGameTurnSheetConfig) error {
	instanceConfig, err := dc.findGameInstanceConfig(gameInstanceRef)
	if err != nil {
		return err
	}
	if len(instanceConfig.GameTurnConfigs) == 0 {
		instanceConfig.GameTurnConfigs = []GameTurnConfig{
			{
				TurnNumber:                    1,
				AdventureGameTurnSheetConfigs: configs,
			},
		}
		instanceConfig.GameTurnConfigs = normalizeTurnConfigs(instanceConfig.GameTurnConfigs)
		return nil
	}
	instanceConfig.GameTurnConfigs[0].AdventureGameTurnSheetConfigs = append(instanceConfig.GameTurnConfigs[0].AdventureGameTurnSheetConfigs, configs...)
	instanceConfig.GameTurnConfigs = normalizeTurnConfigs(instanceConfig.GameTurnConfigs)
	return nil
}

// ReplaceAdventureGameTurnSheetConfigs replaces the first turn config (backwards compatibility)
func (dc *DataConfig) ReplaceAdventureGameTurnSheetConfigs(gameInstanceRef string, configs []AdventureGameTurnSheetConfig) error {
	return dc.ReplaceGameTurnConfigs(gameInstanceRef, []GameTurnConfig{
		{
			TurnNumber:                    1,
			AdventureGameTurnSheetConfigs: configs,
		},
	})
}

// AppendGameTurnSheetConfigs appends configs to a specific turn (creating it if necessary)
func (dc *DataConfig) AppendGameTurnSheetConfigs(gameInstanceRef string, turnNumber int, configs []AdventureGameTurnSheetConfig) error {
	instanceConfig, err := dc.findGameInstanceConfig(gameInstanceRef)
	if err != nil {
		return err
	}
	var turnCfg *GameTurnConfig
	for i := range instanceConfig.GameTurnConfigs {
		if instanceConfig.GameTurnConfigs[i].TurnNumber == turnNumber {
			turnCfg = &instanceConfig.GameTurnConfigs[i]
			break
		}
	}
	if turnCfg == nil {
		instanceConfig.GameTurnConfigs = append(instanceConfig.GameTurnConfigs, GameTurnConfig{
			TurnNumber:                    turnNumber,
			AdventureGameTurnSheetConfigs: configs,
		})
		instanceConfig.GameTurnConfigs = normalizeTurnConfigs(instanceConfig.GameTurnConfigs)
		return nil
	}
	turnCfg.AdventureGameTurnSheetConfigs = append(turnCfg.AdventureGameTurnSheetConfigs, configs...)
	instanceConfig.GameTurnConfigs = normalizeTurnConfigs(instanceConfig.GameTurnConfigs)
	return nil
}

// AppendGameInstanceConfigs appends game instance configs to a game
func (dc *DataConfig) AppendGameInstanceConfigs(gameRef string, configs []GameInstanceConfig) error {
	for i := range dc.GameConfigs {
		if dc.GameConfigs[i].Reference == gameRef {
			dc.GameConfigs[i].GameInstanceConfigs = append(dc.GameConfigs[i].GameInstanceConfigs, configs...)
			return nil
		}
	}
	return fmt.Errorf("game config with reference >%s< not found", gameRef)
}

// TODO: Possibly rename the following to AdventureGameDataConfig when additional game types are added

// DefaultDataConfig
func DefaultDataConfig() DataConfig {
	return DataConfig{
		AccountConfigs: []AccountConfig{
			{
				// Standard account with basic subscriptions only
				Reference:                  StandardAccountRef,
				Record:                     &account_record.AccountUser{},
				AccountSubscriptionConfigs: []AccountSubscriptionConfig{},
				GameSubscriptionConfigs: []GameSubscriptionConfig{
					{
						Reference:        GameSubscriptionDesignerOneRef,
						GameRef:          GameOneRef,
						SubscriptionType: game_record.GameSubscriptionTypeDesigner,
						Record:           &game_record.GameSubscription{},
					},
					{
						Reference:        GameSubscriptionPlayerThreeRef, // Use a new ref or existing unused one
						GameRef:          GameOneRef,
						GameInstanceRefs: []string{GameInstanceOneRef},
						SubscriptionType: game_record.GameSubscriptionTypePlayer,
						Record:           &game_record.GameSubscription{},
					},
				},
			},
			{
				// Pro player account with basic subscriptions + pro player subscription
				Reference: ProPlayerAccountRef,
				Record:    &account_record.AccountUser{},
				AccountSubscriptionConfigs: []AccountSubscriptionConfig{
					{
						SubscriptionType: account_record.AccountSubscriptionTypeProfessionalPlayer,
						Record:           &account_record.AccountSubscription{},
					},
				},
				GameSubscriptionConfigs: []GameSubscriptionConfig{
					{
						Reference:        GameSubscriptionPlayerOneRef,
						GameRef:          GameOneRef,
						GameInstanceRefs: []string{GameInstanceOneRef},
						SubscriptionType: game_record.GameSubscriptionTypePlayer,
						Record:           &game_record.GameSubscription{},
					},
				},
			},
			{
				// Pro designer account with basic subscriptions + pro designer subscription
				Reference: ProDesignerAccountRef,
				Record:    &account_record.AccountUser{},
				AccountSubscriptionConfigs: []AccountSubscriptionConfig{
					{
						SubscriptionType: account_record.AccountSubscriptionTypeProfessionalGameDesigner,
						Record:           &account_record.AccountSubscription{},
					},
				},
				GameSubscriptionConfigs: []GameSubscriptionConfig{
					{
						Reference:        GameSubscriptionDesignerOneRef,
						GameRef:          GameOneRef,
						SubscriptionType: game_record.GameSubscriptionTypeDesigner,
						Record:           &game_record.GameSubscription{},
					},
				},
			},
			{
				// Pro manager account with basic subscriptions + pro manager subscription
				Reference: ProManagerAccountRef,
				Record:    &account_record.AccountUser{},
				AccountSubscriptionConfigs: []AccountSubscriptionConfig{
					{
						SubscriptionType: account_record.AccountSubscriptionTypeProfessionalManager,
						Record:           &account_record.AccountSubscription{},
					},
				},
				GameSubscriptionConfigs: []GameSubscriptionConfig{
					{
						Reference:        GameSubscriptionManagerOneRef,
						GameRef:          GameOneRef,
						GameInstanceRefs: []string{GameInstanceOneRef, GameInstanceCleanRef},
						SubscriptionType: game_record.GameSubscriptionTypeManager,
						Record:           &game_record.GameSubscription{},
					},
				},
			},
		},
		GameConfigs: []GameConfig{
			// Adventure game example
			{
				Reference: GameOneRef,
				Record: &game_record.Game{
					Name:              UniqueName("Default Game One"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168, // 1 week
				},
				// Game images for turn sheet backgrounds (loaded at runtime from testdata)
				GameImageConfigs: []GameImageConfig{
					{
						Reference:     GameImageJoinGameRef,
						ImagePath:     "background-darkforest.png",
						TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeJoinGame,
					},
					{
						Reference:     GameImageInventoryRef,
						ImagePath:     "background-dungeon.png",
						TurnSheetType: adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
					},
				},
				// Adventure game specific resources
				AdventureGameItemConfigs: []AdventureGameItemConfig{
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
				AdventureGameLocationConfigs: []AdventureGameLocationConfig{
					{
						Reference: GameLocationOneRef,
						Record: &adventure_game_record.AdventureGameLocation{
							Name:               UniqueName("Default Location One"),
							Description:        "Default location one for handler tests",
							IsStartingLocation: true,
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
				AdventureGameLocationLinkConfigs: []AdventureGameLocationLinkConfig{
					{
						Reference:       GameLocationLinkOneRef,
						FromLocationRef: GameLocationOneRef,
						ToLocationRef:   GameLocationTwoRef,
						Record: &adventure_game_record.AdventureGameLocationLink{
							Name:        UniqueName("The Red Door"),
							Description: "Travel by boat to the swamp of the long forgotten Frog God",
						},
						AdventureGameLocationLinkRequirementConfigs: []AdventureGameLocationLinkRequirementConfig{
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
				AdventureGameCreatureConfigs: []AdventureGameCreatureConfig{
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
				AdventureGameCharacterConfigs: []AdventureGameCharacterConfig{
					{
						Reference:  GameCharacterOneRef,
						AccountRef: StandardAccountRef,
						Record: &adventure_game_record.AdventureGameCharacter{
							Name: UniqueName("Default Character One"),
						},
					},
					{
						Reference:  GameCharacterTwoRef,
						AccountRef: ProPlayerAccountRef,
						Record: &adventure_game_record.AdventureGameCharacter{
							Name: UniqueName("Default Character Two"),
						},
					},
				},
				// Default adventure game instance with a location and an item assigned to the location
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
						// AdventureGameLocationInstanceConfigs: []AdventureGameLocationInstanceConfig{
						// 	{
						// 		Reference:       GameLocationInstanceOneRef,
						// 		GameLocationRef: GameLocationOneRef,
						// 		Record:          &adventure_game_record.AdventureGameLocationInstance{},
						// 	},
						// 	{
						// 		Reference:       GameLocationInstanceTwoRef,
						// 		GameLocationRef: GameLocationTwoRef,
						// 		Record:          &adventure_game_record.AdventureGameLocationInstance{},
						// 	},
						// },
						AdventureGameItemInstanceConfigs: []AdventureGameItemInstanceConfig{
							{
								Reference:       GameItemInstanceOneRef,
								GameItemRef:     GameItemOneRef,
								GameLocationRef: GameLocationOneRef, // Assign to location
								Record:          &adventure_game_record.AdventureGameItemInstance{},
							},
						},
						AdventureGameCreatureInstanceConfigs: []AdventureGameCreatureInstanceConfig{
							{
								Reference:       GameCreatureInstanceOneRef,
								GameCreatureRef: GameCreatureOneRef,
								GameLocationRef: GameLocationOneRef, // Assign to location
								Record:          &adventure_game_record.AdventureGameCreatureInstance{},
							},
						},
						AdventureGameCharacterInstanceConfigs: []AdventureGameCharacterInstanceConfig{
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
			// Minimal draft game for testing update operations
			{
				Reference: GameDraftRef,
				Record: &game_record.Game{
					Name:              UniqueName("Default Draft Game"),
					GameType:          game_record.GameTypeAdventure,
					TurnDurationHours: 168,
					Status:            game_record.GameStatusDraft,
				},
			},
		},
	}
}
