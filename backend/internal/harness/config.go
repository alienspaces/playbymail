// Package harness provides test data setup and teardown utilities.
package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

const (
	GameOneRef        = "game-one"
	GameTwoRef        = "game-two"
	GameDraftRef      = "game-draft"
	GameMechWargameRef = "game-mech-wargame"

	AccountStandardRef    = "account-standard"
	AccountProPlayerRef   = "account-pro-player"
	AccountProDesignerRef = "account-pro-designer"
	AccountProManagerRef  = "account-pro-manager"

	AccountUserStandardRef    = "account-user-standard"
	AccountUserProPlayerRef   = "account-user-pro-player"
	AccountUserProDesignerRef = "account-user-pro-designer"
	AccountUserProManagerRef  = "account-user-pro-manager"

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
	GameCharacterFourRef  = "game-character-four"

	GameInstanceOneRef   = "game-instance-one"
	GameInstanceTwoRef   = "game-instance-two"
	GameInstanceCleanRef = "game-instance-clean"

	GameInstanceParameterOneRef = "game-instance-parameter-one"

	GameCreaturePlacementOneRef = "game-creature-placement-one"
	GameItemPlacementOneRef     = "game-item-placement-one"

	GameLocationObjectOneRef = "game-location-object-one"
	GameLocationObjectTwoRef = "game-location-object-two"

	// States for the Ancient Shrine (ObjectOne)
	GameLocationObjectOneStateIntactRef    = "game-location-object-one-state-intact"
	GameLocationObjectOneStateActivatedRef = "game-location-object-one-state-activated"

	// States for the Hidden Passage (ObjectTwo)
	GameLocationObjectTwoStateSealedRef = "game-location-object-two-state-sealed"
	GameLocationObjectTwoStateOpenRef   = "game-location-object-two-state-open"

	GameLocationObjectEffectOneRef   = "game-location-object-effect-one"
	GameLocationObjectEffectTwoRef   = "game-location-object-effect-two"
	GameLocationObjectEffectThreeRef = "game-location-object-effect-three"
	GameLocationObjectEffectFourRef  = "game-location-object-effect-four"
	GameLocationObjectEffectFiveRef  = "game-location-object-effect-five"
	GameLocationObjectEffectSixRef   = "game-location-object-effect-six"
	GameLocationObjectEffectSevenRef = "game-location-object-effect-seven"

	GameItemInstanceOneRef = "game-item-instance-one"

	GameLocationInstanceOneRef = "game-location-instance-one"
	GameLocationInstanceTwoRef = "game-location-instance-two"

	GameCreatureInstanceOneRef = "game-creature-instance-one"

	GameCharacterInstanceOneRef = "game-character-instance-one"
	GameCharacterInstanceTwoRef = "game-character-instance-two"

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

	// Mech wargame specific references
	MechWargameChassisOneRef    = "mech-wargame-chassis-one"
	MechWargameWeaponOneRef     = "mech-wargame-weapon-one"
	MechWargameSectorOneRef     = "mech-wargame-sector-one"
	MechWargameSectorTwoRef     = "mech-wargame-sector-two"
	MechWargameSectorLinkOneRef = "mech-wargame-sector-link-one"
	MechWargameLanceOneRef      = "mech-wargame-lance-one"
	MechWargameLanceTwoRef      = "mech-wargame-lance-two"
	MechWargameLanceMechOneRef  = "mech-wargame-lance-mech-one"
)

// DataConfig -
type DataConfig struct {

	// Account configurations result in account, account_user, and account_user_contact records.
	AccountConfigs []AccountConfig

	// Game configurations result in game, game_instance, and game_image records.
	GameConfigs []GameConfig

	// Account user game subscription configurations may only be be resolved
	// once both accounts and games have been created.
	AccountUserGameSubscriptionConfigs []AccountUserGameSubscriptionConfig
}

type GameConfig struct {
	Reference        string // Reference to the game record
	Record           *game_record.Game
	GameImageConfigs []GameImageConfig // Game images for turn sheet backgrounds

	// Adventure game specific configurations
	AdventureGameLocationConfigs          []AdventureGameLocationConfig
	AdventureGameLocationLinkConfigs      []AdventureGameLocationLinkConfig
	AdventureGameItemConfigs              []AdventureGameItemConfig
	AdventureGameCreatureConfigs          []AdventureGameCreatureConfig
	AdventureGameCharacterConfigs         []AdventureGameCharacterConfig // Characters associated with this game
	AdventureGameLocationObjectConfigs    []AdventureGameLocationObjectConfig
	AdventureGameCreaturePlacementConfigs []AdventureGameCreaturePlacementConfig // Where creatures start in the world
	AdventureGameItemPlacementConfigs     []AdventureGameItemPlacementConfig     // Where items start in the world

	// Mech wargame specific configurations
	MechWargameChassisConfigs    []MechWargameChassisConfig
	MechWargameWeaponConfigs     []MechWargameWeaponConfig
	MechWargameSectorConfigs     []MechWargameSectorConfig
	MechWargameSectorLinkConfigs []MechWargameSectorLinkConfig
	MechWargameLanceConfigs      []MechWargameLanceConfig
}

type GameImageConfig struct {
	// Reference to the game_image record
	Reference string
	// Path to image file relative to test_data_images or testdata directories (loaded at runtime). Required.
	ImagePath string
	// The turn sheet type for this image
	TurnSheetType string
	// RecordID is the optional record ID (e.g. location ID) for the related record.
	RecordID string
}

type AccountConfig struct {
	Reference          string // Reference to the account record
	Record             *account_record.Account
	AccountUserConfigs []AccountUserConfig
}

type AccountUserConfig struct {
	Reference string // Reference to the account_user record
	Record    *account_record.AccountUser

	// Every account user by default will be assigned basic subscriptions.
	// Additional subscriptions may be explicitly configured here.
	AccountUserSubscriptionConfigs []AccountUserSubscriptionConfig
}

type AccountUserSubscriptionConfig struct {
	SubscriptionType string // e.g., AccountSubscriptionTypeProfessionalPlayer
	Record           *account_record.AccountSubscription
}

type GameInstanceConfig struct {
	Reference                    string // Reference to the game_instance record
	Record                       *game_record.GameInstance
	GameInstanceParameterConfigs []GameInstanceParameterConfig

	// PlayerSubscriptionRefs links player-type game subscriptions to this instance (for GetPlayerCountForGameInstance, etc.)
	PlayerSubscriptionRefs []string

	// ShouldStartGameInstance controls whether the harness calls StartGameInstance during setup.
	// When true, the instance is started (status "started", CurrentTurn 0) and initial turn sheets are created.
	ShouldStartGameInstance bool

	// TurnSheetRefConfigs maps turn sheet references to character instances, used to populate
	// Data.Refs.GameTurnSheetRefs so tests can look up turn sheets by reference.
	// Only used when ShouldStartGameInstance is true.
	TurnSheetRefConfigs []TurnSheetRefConfig
}

// TurnSheetRefConfig maps a turn sheet reference name to the character instance that owns it.
type TurnSheetRefConfig struct {
	Reference                string // e.g. GameTurnSheetOneRef
	GameCharacterInstanceRef string // e.g. GameCharacterInstanceOneRef
}

type AccountUserGameSubscriptionConfig struct {
	Reference        string                        // Reference to the game_subscription record
	AccountUserRef   string                        // Reference to account user (resolved from Data.Refs)
	SubscriptionType string                        // Type of subscription (Player, Manager, Designer)
	GameRef          string                        // Reference to the game
	InstanceLimit    *int32                        // Instance limit (nil = unlimited)
	Record           *game_record.GameSubscription // Optional record to use instead of default values

	// For Manager subscriptions: instances this manager creates. Each instance is created and linked
	// to this subscription via GameSubscriptionInstance.
	GameInstanceConfigs []GameInstanceConfig

	// For Player subscriptions: reference to the manager's game subscription (resolved from Data.Refs).
	// Semantically: this player joined via this manager's subscription.
	AccountUserManagerGameSubscriptionRef string

	// JoinGameScanData configures scan data for join game turn sheets
	// If provided, the harness will create a join game turn sheet and process it
	// For adventure games, use turn_sheet.AdventureGameJoinGameScanData
	// For other game types, use the appropriate game-specific scan data type
	JoinGameScanData any // Can be *turn_sheet.AdventureGameJoinGameScanData, etc.
}

// ------------------------------------------------------------
// Adventure game specific configuration
// ------------------------------------------------------------
type AdventureGameCharacterConfig struct {
	Reference  string // Reference to the game_character record
	AccountRef string // Reference to the account
	Record     *adventure_game_record.AdventureGameCharacter
}

type AdventureGameItemConfig struct {
	Reference string // Reference to the game_item record
	Record    *adventure_game_record.AdventureGameItem
	// Effects attached to this item
	AdventureGameItemEffectConfigs []AdventureGameItemEffectConfig
}

type AdventureGameItemEffectConfig struct {
	Reference string // Reference to the effect record
	Record    *adventure_game_record.AdventureGameItemEffect
	// ResultItemRef is an optional reference to an item whose ID should be set on result_adventure_game_item_id.
	ResultItemRef string
	// ResultLocationRef is an optional reference to a location whose ID should be set on result_adventure_game_location_id.
	ResultLocationRef string
	// ResultLinkRef is an optional reference to a location link whose ID should be set on result_adventure_game_location_link_id.
	ResultLinkRef string
	// ResultCreatureRef is an optional reference to a creature whose ID should be set on result_adventure_game_creature_id.
	ResultCreatureRef string
	// RequiredItemRef is an optional reference to an item whose ID should be set on required_adventure_game_item_id.
	RequiredItemRef string
	// RequiredLocationRef is an optional reference to a location whose ID should be set on required_adventure_game_location_id.
	RequiredLocationRef string
}

type AdventureGameCreatureConfig struct {
	Reference string // Reference to the game_creature record
	Record    *adventure_game_record.AdventureGameCreature
	// PortraitImage if set, creates a game_image (asset) for this creature's portrait.
	// Only ImagePath (and optionally Reference) need to be set; RecordID is set at runtime.
	PortraitImage *GameImageConfig
}

type AdventureGameLocationConfig struct {
	Reference string // Reference to the game_location record
	Record    *adventure_game_record.AdventureGameLocation
	// BackgroundImage if set, creates a game_image (turn_sheet_background) for this location.
	// Only ImagePath (and optionally Reference) are set in config; RecordID and TurnSheetType
	// are set at runtime from the created location ID and adventure_game_location_choice.
	BackgroundImage *GameImageConfig
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

	// Exactly one of GameItemRef or GameCreatureRef must be set
	GameItemRef     string // Reference to the game_item (for item-based conditions)
	GameCreatureRef string // Reference to the game_creature (for creature-based conditions)

	Record *adventure_game_record.AdventureGameLocationLinkRequirement
}

type AdventureGameCreaturePlacementConfig struct {
	Reference       string // Reference to the placement record
	GameCreatureRef string // Reference to the game_creature (required)
	GameLocationRef string // Reference to the game_location (required)
	InitialCount    int    // Number of creatures to place (defaults to 1)
	Record          *adventure_game_record.AdventureGameCreaturePlacement
}

type AdventureGameItemPlacementConfig struct {
	Reference       string // Reference to the placement record
	GameItemRef     string // Reference to the game_item (required)
	GameLocationRef string // Reference to the game_location (required)
	InitialCount    int    // Number of items to place (defaults to 1)
	Record          *adventure_game_record.AdventureGameItemPlacement
}

type AdventureGameLocationObjectStateConfig struct {
	Reference string // Reference to the state record (for resolving from other configs)
	Record    *adventure_game_record.AdventureGameLocationObjectState
}

type AdventureGameLocationObjectEffectConfig struct {
	Reference string // Reference to the effect record
	Record    *adventure_game_record.AdventureGameLocationObjectEffect
	// ResultObjectRef is an optional reference to another location object whose ID should
	// be resolved and set on result_adventure_game_location_object_id before creation.
	ResultObjectRef string
	// ResultItemRef is an optional reference to an item whose ID should be set on result_adventure_game_item_id.
	ResultItemRef string
	// ResultLocationRef is an optional reference to a location whose ID should be set on result_adventure_game_location_id.
	ResultLocationRef string
	// ResultLinkRef is an optional reference to a location link whose ID should be set on result_adventure_game_location_link_id.
	ResultLinkRef string
	// ResultCreatureRef is an optional reference to a creature whose ID should be set on result_adventure_game_creature_id.
	ResultCreatureRef string
	// RequiredItemRef is an optional reference to an item whose ID should be set on required_adventure_game_item_id.
	RequiredItemRef string
	// RequiredStateRef resolves to required_adventure_game_location_object_state_id
	RequiredStateRef string
	// ResultStateRef resolves to result_adventure_game_location_object_state_id
	ResultStateRef string
}

type AdventureGameLocationObjectConfig struct {
	Reference   string // Reference to the object record
	LocationRef string // Reference to the location this object belongs to
	Record      *adventure_game_record.AdventureGameLocationObject
	// InitialStateRef resolves to initial_adventure_game_location_object_state_id
	InitialStateRef string
	// States attached to this object - must be created before effects
	AdventureGameLocationObjectStateConfigs []AdventureGameLocationObjectStateConfig
	// Effects attached to this object
	AdventureGameLocationObjectEffectConfigs []AdventureGameLocationObjectEffectConfig
}

// ------------------------------------------------------------
// Mech wargame specific configuration
// ------------------------------------------------------------

type MechWargameChassisConfig struct {
	Reference string
	Record    *mech_wargame_record.MechWargameChassis
}

type MechWargameWeaponConfig struct {
	Reference string
	Record    *mech_wargame_record.MechWargameWeapon
}

type MechWargameSectorConfig struct {
	Reference string
	Record    *mech_wargame_record.MechWargameSector
}

type MechWargameSectorLinkConfig struct {
	Reference          string
	FromSectorRef      string
	ToSectorRef        string
	Record             *mech_wargame_record.MechWargameSectorLink
}

type MechWargameLanceMechConfig struct {
	Reference   string
	ChassisRef  string
	Record      *mech_wargame_record.MechWargameLanceMech
}

type MechWargameLanceConfig struct {
	Reference        string
	AccountRef       string
	Record           *mech_wargame_record.MechWargameLance
	LanceMechConfigs []MechWargameLanceMechConfig
}

// Helper methods for modifying DataConfig

// FindGameInstanceConfig finds a game instance config by reference.
// Instance configs live under manager subscription configs.
func (dc *DataConfig) FindGameInstanceConfig(gameInstanceRef string) (*GameInstanceConfig, error) {
	for i := range dc.AccountUserGameSubscriptionConfigs {
		sub := &dc.AccountUserGameSubscriptionConfigs[i]
		if sub.SubscriptionType != game_record.GameSubscriptionTypeManager {
			continue
		}
		for j := range sub.GameInstanceConfigs {
			if sub.GameInstanceConfigs[j].Reference == gameInstanceRef {
				return &sub.GameInstanceConfigs[j], nil
			}
		}
	}
	return nil, fmt.Errorf("game instance config with reference >%s< not found", gameInstanceRef)
}

// SetTurnSheetRefConfigs sets the turn sheet ref configs for a game instance.
// This allows tests to look up turn sheets by reference after the harness starts the instance.
func (dc *DataConfig) SetTurnSheetRefConfigs(gameInstanceRef string, configs []TurnSheetRefConfig) error {
	instanceConfig, err := dc.FindGameInstanceConfig(gameInstanceRef)
	if err != nil {
		return err
	}
	instanceConfig.TurnSheetRefConfigs = configs
	return nil
}

// AppendGameInstanceConfigs appends game instance configs to the first manager subscription
// config that has the given gameRef. Used by tests to add instances to a game's manager.
func (dc *DataConfig) AppendGameInstanceConfigs(gameRef string, configs []GameInstanceConfig) error {
	for i := range dc.AccountUserGameSubscriptionConfigs {
		sub := &dc.AccountUserGameSubscriptionConfigs[i]
		if sub.SubscriptionType == game_record.GameSubscriptionTypeManager && sub.GameRef == gameRef {
			sub.GameInstanceConfigs = append(sub.GameInstanceConfigs, configs...)
			return nil
		}
	}
	return fmt.Errorf("manager subscription config for game ref >%s< not found", gameRef)
}

// DefaultDataConfig returns default data config; currently adventure-game focused. Consider renaming to AdventureGameDataConfig when adding more game types.
func DefaultDataConfig() DataConfig {
	return DataConfig{
		AccountConfigs: []AccountConfig{
			{
				// Standard account with a single user withbasic subscriptions only
				Reference: AccountStandardRef,
				Record:    &account_record.Account{},
				AccountUserConfigs: []AccountUserConfig{
					{
						Reference: AccountUserStandardRef,
						Record:    &account_record.AccountUser{},
					},
				},
			},
			{
				// Pro player account with a single user with basic subscriptions + pro player subscription
				Reference: AccountProPlayerRef,
				Record:    &account_record.Account{},
				AccountUserConfigs: []AccountUserConfig{
					{
						Reference: AccountUserProPlayerRef,
						Record:    &account_record.AccountUser{},
						AccountUserSubscriptionConfigs: []AccountUserSubscriptionConfig{
							{
								SubscriptionType: account_record.AccountSubscriptionTypeProfessionalPlayer,
								Record:           &account_record.AccountSubscription{},
							},
						},
					},
				},
			},
			{
				// Pro designer account with a single user with basic subscriptions + pro designer subscription
				Reference: AccountProDesignerRef,
				Record:    &account_record.Account{},
				AccountUserConfigs: []AccountUserConfig{
					{
						Reference: AccountUserProDesignerRef,
						Record:    &account_record.AccountUser{},
						AccountUserSubscriptionConfigs: []AccountUserSubscriptionConfig{
							{
								SubscriptionType: account_record.AccountSubscriptionTypeProfessionalGameDesigner,
								Record:           &account_record.AccountSubscription{},
							},
						},
					},
				},
			},
			{
				// Pro manager account with a single user with basic subscriptions + pro manager subscription
				Reference: AccountProManagerRef,
				Record:    &account_record.Account{},
				AccountUserConfigs: []AccountUserConfig{
					{
						Reference: AccountUserProManagerRef,
						Record:    &account_record.AccountUser{},
						AccountUserSubscriptionConfigs: []AccountUserSubscriptionConfig{
							{
								SubscriptionType: account_record.AccountSubscriptionTypeProfessionalManager,
								Record:           &account_record.AccountSubscription{},
							},
						},
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
					{
						Reference: GameLocationThreeRef,
						Record: &adventure_game_record.AdventureGameLocation{
							Name:        UniqueName("Default Location Three"),
							Description: "Default location three for handler tests",
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
									Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
									Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory,
									Quantity:  1,
								},
							},
							{
								Reference:       GameLocationLinkRequirementTwoRef,
								GameCreatureRef: GameCreatureOneRef,
								Record: &adventure_game_record.AdventureGameLocationLinkRequirement{
									Purpose:   adventure_game_record.AdventureGameLocationLinkRequirementPurposeVisible,
									Condition: adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation,
									Quantity:  1,
								},
							},
						},
					},
					// Location Two -> Location One (return path, same link name)
					{
						Reference:       GameLocationLinkTwoRef,
						FromLocationRef: GameLocationTwoRef,
						ToLocationRef:   GameLocationOneRef,
						Record: &adventure_game_record.AdventureGameLocationLink{
							Name:        UniqueName("The Red Door"),
							Description: "Return through the red door from the swamp back to the start.",
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
						AccountRef: AccountUserStandardRef,
						Record: &adventure_game_record.AdventureGameCharacter{
							Name: UniqueName("Default Character One"),
						},
					},
					{
						Reference:  GameCharacterTwoRef,
						AccountRef: AccountUserProPlayerRef,
						Record: &adventure_game_record.AdventureGameCharacter{
							Name: UniqueName("Default Character Two"),
						},
					},
				},
				AdventureGameCreaturePlacementConfigs: []AdventureGameCreaturePlacementConfig{
					{
						Reference:       GameCreaturePlacementOneRef,
						GameCreatureRef: GameCreatureOneRef,
						GameLocationRef: GameLocationOneRef,
						InitialCount:    1,
						Record:          &adventure_game_record.AdventureGameCreaturePlacement{},
					},
				},
			AdventureGameItemPlacementConfigs: []AdventureGameItemPlacementConfig{
				{
					Reference:       GameItemPlacementOneRef,
					GameItemRef:     GameItemOneRef,
					GameLocationRef: GameLocationOneRef,
					InitialCount:    1,
					Record:          &adventure_game_record.AdventureGameItemPlacement{},
				},
			},
			AdventureGameLocationObjectConfigs: []AdventureGameLocationObjectConfig{
				// Object Two is created first so Object One's effects can reference it
				{
					Reference:       GameLocationObjectTwoRef,
					LocationRef:     GameLocationOneRef,
					InitialStateRef: GameLocationObjectTwoStateSealedRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        UniqueName("Hidden Passage"),
						Description: "A concealed passage behind the shrine.",
						IsHidden:    true,
					},
					AdventureGameLocationObjectStateConfigs: []AdventureGameLocationObjectStateConfig{
						{
							Reference: GameLocationObjectTwoStateSealedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "sealed",
								Description: "The passage is sealed shut.",
								SortOrder:   0,
							},
						},
						{
							Reference: GameLocationObjectTwoStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "open",
								Description: "The passage is open.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []AdventureGameLocationObjectEffectConfig{
						{
							Reference:        "game-location-object-two-effect-one",
							RequiredStateRef: GameLocationObjectTwoStateSealedRef,
							ResultStateRef:   GameLocationObjectTwoStateOpenRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeOpen,
								ResultDescription: "The passage opens, revealing a dark tunnel.",
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								IsRepeatable:      false,
							},
						},
					},
				},
				{
					Reference:       GameLocationObjectOneRef,
					LocationRef:     GameLocationOneRef,
					InitialStateRef: GameLocationObjectOneStateIntactRef,
					Record: &adventure_game_record.AdventureGameLocationObject{
						Name:        UniqueName("Ancient Shrine"),
						Description: "A weathered stone shrine covered in moss.",
						IsHidden:    false,
					},
					AdventureGameLocationObjectStateConfigs: []AdventureGameLocationObjectStateConfig{
						{
							Reference: GameLocationObjectOneStateIntactRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "intact",
								Description: "The shrine stands whole and undisturbed.",
								SortOrder:   0,
							},
						},
						{
							Reference: GameLocationObjectOneStateActivatedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectState{
								Name:        "activated",
								Description: "The shrine pulses with ancient power.",
								SortOrder:   1,
							},
						},
					},
					AdventureGameLocationObjectEffectConfigs: []AdventureGameLocationObjectEffectConfig{
						{
							Reference: GameLocationObjectEffectOneRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeInspect,
								ResultDescription: "The shrine glows faintly with ancient power.",
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
								IsRepeatable:      true,
							},
						},
						{
							Reference:        GameLocationObjectEffectTwoRef,
							RequiredStateRef: GameLocationObjectOneStateIntactRef,
							ResultStateRef:   GameLocationObjectOneStateActivatedRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeTouch,
								ResultDescription: "The shrine hums with power as you touch it.",
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
								IsRepeatable:      false,
							},
						},
						{
							Reference:        GameLocationObjectEffectThreeRef,
							RequiredStateRef: GameLocationObjectOneStateActivatedRef,
							ResultItemRef:    GameItemTwoRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeSearch,
								ResultDescription: "A hidden item appears from within the shrine.",
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
								IsRepeatable:      false,
							},
						},
						{
							Reference: GameLocationObjectEffectFourRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeBreak,
								ResultDescription: "Shards fly out cutting you.",
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage,
								ResultValueMin:    nullint32.FromInt32(5),
								ResultValueMax:    nullint32.FromInt32(5),
								IsRepeatable:      true,
							},
						},
						{
							Reference:       GameLocationObjectEffectFiveRef,
							ResultObjectRef: GameLocationObjectTwoRef,
							Record: &adventure_game_record.AdventureGameLocationObjectEffect{
								ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeUse,
								ResultDescription: "A hidden passage is revealed.",
								EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject,
								IsRepeatable:      false,
							},
						},
					{
						Reference:        GameLocationObjectEffectSixRef,
						RequiredStateRef: GameLocationObjectOneStateActivatedRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypeTouch,
							ResultDescription: "Warmth flows through you.",
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHeal,
							ResultValueMin:    nullint32.FromInt32(10),
							ResultValueMax:    nullint32.FromInt32(10),
							IsRepeatable:      true,
						},
					},
					{
						Reference:         GameLocationObjectEffectSevenRef,
						RequiredStateRef:  GameLocationObjectOneStateActivatedRef,
						ResultItemRef:     GameItemOneRef,
						ResultLocationRef: GameLocationTwoRef,
						Record: &adventure_game_record.AdventureGameLocationObjectEffect{
							ActionType:        adventure_game_record.AdventureGameLocationObjectEffectActionTypePull,
							ResultDescription: "An item materialises at the far location.",
							EffectType:        adventure_game_record.AdventureGameLocationObjectEffectEffectTypePlaceItem,
							IsRepeatable:      false,
						},
					},
				},
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
		// Mech wargame game for testing mech wargame specific resources
		{
			Reference: GameMechWargameRef,
			Record: &game_record.Game{
				Name:              UniqueName("Default Mech Wargame"),
				GameType:          game_record.GameTypeMechWargame,
				TurnDurationHours: 168,
			},
			MechWargameChassisConfigs: []MechWargameChassisConfig{
				{
					Reference: MechWargameChassisOneRef,
					Record: &mech_wargame_record.MechWargameChassis{
						Name:            UniqueName("Timber Wolf"),
						Description:     "A heavy assault mech.",
						ChassisClass:    mech_wargame_record.ChassisClassHeavy,
						ArmorPoints:     200,
						StructurePoints: 100,
						HeatCapacity:    40,
						Speed:           3,
					},
				},
			},
			MechWargameWeaponConfigs: []MechWargameWeaponConfig{
				{
					Reference: MechWargameWeaponOneRef,
					Record: &mech_wargame_record.MechWargameWeapon{
						Name:        UniqueName("Large Laser"),
						Description: "A powerful energy weapon.",
						Damage:      8,
						HeatCost:    8,
						RangeBand:   mech_wargame_record.WeaponRangeBandMedium,
						MountSize:   mech_wargame_record.WeaponMountSizeLarge,
					},
				},
			},
			MechWargameSectorConfigs: []MechWargameSectorConfig{
				{
					Reference: MechWargameSectorOneRef,
					Record: &mech_wargame_record.MechWargameSector{
						Name:             UniqueName("Ridge South"),
						Description:      "A rocky ridge offering good cover.",
						TerrainType:      mech_wargame_record.SectorTerrainTypeRough,
						Elevation:        2,
						IsStartingSector: true,
					},
				},
				{
					Reference: MechWargameSectorTwoRef,
					Record: &mech_wargame_record.MechWargameSector{
						Name:        UniqueName("Relay Station"),
						Description: "An abandoned communications relay station.",
						TerrainType: mech_wargame_record.SectorTerrainTypeUrban,
						Elevation:   0,
					},
				},
			},
			MechWargameSectorLinkConfigs: []MechWargameSectorLinkConfig{
				{
					Reference:     MechWargameSectorLinkOneRef,
					FromSectorRef: MechWargameSectorOneRef,
					ToSectorRef:   MechWargameSectorTwoRef,
					Record: &mech_wargame_record.MechWargameSectorLink{
						CoverModifier: 1,
					},
				},
			},
			MechWargameLanceConfigs: []MechWargameLanceConfig{
				{
				Reference:  MechWargameLanceOneRef,
				AccountRef: AccountUserStandardRef,
					Record: &mech_wargame_record.MechWargameLance{
						Name:        UniqueName("Alpha Lance"),
						Description: "First player lance.",
					},
					LanceMechConfigs: []MechWargameLanceMechConfig{
						{
							Reference:  MechWargameLanceMechOneRef,
							ChassisRef: MechWargameChassisOneRef,
							Record: &mech_wargame_record.MechWargameLanceMech{
								Callsign: "Wolf-1",
							},
						},
					},
				},
			},
		},
	},
		// Account user game subscription configurations may only be be resolved
		// once both accounts and games have been created.
		AccountUserGameSubscriptionConfigs: []AccountUserGameSubscriptionConfig{
			{
				Reference:        GameSubscriptionDesignerOneRef,
				AccountUserRef:   AccountUserStandardRef,
				GameRef:          GameOneRef,
				SubscriptionType: game_record.GameSubscriptionTypeDesigner,
				Record:           &game_record.GameSubscription{},
			},
			{
				Reference:                             GameSubscriptionPlayerThreeRef,
				AccountUserRef:                        AccountUserStandardRef,
				GameRef:                               GameOneRef,
				SubscriptionType:                      game_record.GameSubscriptionTypePlayer,
				AccountUserManagerGameSubscriptionRef: GameSubscriptionManagerOneRef,
				Record:                                &game_record.GameSubscription{},
			},
			{
				Reference:                             GameSubscriptionPlayerOneRef,
				AccountUserRef:                        AccountUserStandardRef,
				GameRef:                               GameOneRef,
				SubscriptionType:                      game_record.GameSubscriptionTypePlayer,
				AccountUserManagerGameSubscriptionRef: GameSubscriptionManagerOneRef,
				Record:                                &game_record.GameSubscription{},
			},
			{
				Reference:        GameSubscriptionDesignerOneRef,
				AccountUserRef:   AccountUserProDesignerRef,
				GameRef:          GameOneRef,
				SubscriptionType: game_record.GameSubscriptionTypeDesigner,
				Record:           &game_record.GameSubscription{},
			},
		{
			Reference:                             GameSubscriptionPlayerTwoRef,
			AccountUserRef:                        AccountUserProPlayerRef,
			GameRef:                               GameOneRef,
			SubscriptionType:                      game_record.GameSubscriptionTypePlayer,
			AccountUserManagerGameSubscriptionRef: GameSubscriptionManagerOneRef,
			Record:                                &game_record.GameSubscription{},
		},
		{
			Reference:        GameSubscriptionManagerOneRef,
			AccountUserRef:   AccountUserProManagerRef,
			GameRef:          GameOneRef,
			SubscriptionType: game_record.GameSubscriptionTypeManager,
			Record:           &game_record.GameSubscription{},
			GameInstanceConfigs: []GameInstanceConfig{
				{
					Reference:              GameInstanceOneRef,
					Record:                 &game_record.GameInstance{},
					PlayerSubscriptionRefs: []string{GameSubscriptionPlayerOneRef, GameSubscriptionPlayerThreeRef},
					GameInstanceParameterConfigs: []GameInstanceParameterConfig{
						{
							Reference: GameInstanceParameterOneRef,
							Record: &game_record.GameInstanceParameter{
								ParameterKey:   domain.AdventureGameParameterCharacterLives,
								ParameterValue: nullstring.FromString("3"),
							},
						},
					},
				ShouldStartGameInstance: true,
				},
				{
					Reference:              GameInstanceTwoRef,
					Record:                 &game_record.GameInstance{},
					PlayerSubscriptionRefs: []string{GameSubscriptionPlayerTwoRef},
				},
				{
					Reference: GameInstanceCleanRef,
					Record:    &game_record.GameInstance{},
				},
			},
		},
		},
	}
}
