package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/repository/account"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_character"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_character_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_creature"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_creature_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_creature_placement"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_item"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_item_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_item_placement"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location_link"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location_link_requirement"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/repository/game"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_administration"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_configuration"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_instance_configuration"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_subscription"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// Domain -
type Domain struct {
	domain.Domain
	config config.Config
}

var _ domainer.Domainer = &Domain{}

func NewDomain(l logger.Logger, cfg config.Config) (*Domain, error) {

	l, err := l.NewInstance()
	if err != nil {
		return nil, err
	}

	repositoryConstructors := []domain.RepositoryConstructor{
		account.NewRepository,
		game.NewRepository,
		adventure_game_location.NewRepository,
		adventure_game_location_link.NewRepository,
		adventure_game_character.NewRepository,
		adventure_game_item.NewRepository,
		adventure_game_creature.NewRepository,
		adventure_game_item_instance.NewRepository,
		adventure_game_item_placement.NewRepository,
		adventure_game_creature_placement.NewRepository,
		adventure_game_location_link_requirement.NewRepository,
		adventure_game_location_instance.NewRepository,
		adventure_game_creature_instance.NewRepository,
		adventure_game_character_instance.NewRepository,
		adventure_game_turn_sheet.NewRepository,
		game_instance.NewRepository,
		game_configuration.NewRepository,
		game_instance_configuration.NewRepository,
		game_subscription.NewRepository,
		game_administration.NewRepository,
		game_turn_sheet.NewRepository,
	}

	cd, err := domain.NewDomain(l, repositoryConstructors)
	if err != nil {
		return nil, err
	}

	m := &Domain{
		Domain: *cd,
		config: cfg,
	}

	l.Info("returning domain %+v", m)

	return m, nil
}

// AccountRepository -
func (m *Domain) AccountRepository() *repository.Generic[account_record.Account, *account_record.Account] {
	return m.Repositories[account.TableName].(*repository.Generic[account_record.Account, *account_record.Account])
}

// GameRepository -
func (m *Domain) GameRepository() *repository.Generic[game_record.Game, *game_record.Game] {
	return m.Repositories[game.TableName].(*repository.Generic[game_record.Game, *game_record.Game])
}

// GameLocationLinkRepository -
func (m *Domain) AdventureGameLocationLinkRepository() *repository.Generic[adventure_game_record.AdventureGameLocationLink, *adventure_game_record.AdventureGameLocationLink] {
	return m.Repositories[adventure_game_location_link.TableName].(*repository.Generic[adventure_game_record.AdventureGameLocationLink, *adventure_game_record.AdventureGameLocationLink])
}

// AdventureGameCharacterRepository -
func (m *Domain) AdventureGameCharacterRepository() *repository.Generic[adventure_game_record.AdventureGameCharacter, *adventure_game_record.AdventureGameCharacter] {
	return m.Repositories[adventure_game_character.TableName].(*repository.Generic[adventure_game_record.AdventureGameCharacter, *adventure_game_record.AdventureGameCharacter])
}

// AdventureGameLocationRepository -
func (m *Domain) AdventureGameLocationRepository() *repository.Generic[adventure_game_record.AdventureGameLocation, *adventure_game_record.AdventureGameLocation] {
	return m.Repositories[adventure_game_location.TableName].(*repository.Generic[adventure_game_record.AdventureGameLocation, *adventure_game_record.AdventureGameLocation])
}

// AdventureGameItemRepository -
func (m *Domain) AdventureGameItemRepository() *repository.Generic[adventure_game_record.AdventureGameItem, *adventure_game_record.AdventureGameItem] {
	return m.Repositories[adventure_game_item.TableName].(*repository.Generic[adventure_game_record.AdventureGameItem, *adventure_game_record.AdventureGameItem])
}

// AdventureGameItemPlacementRepository -
func (m *Domain) AdventureGameItemPlacementRepository() *repository.Generic[adventure_game_record.AdventureGameItemPlacement, *adventure_game_record.AdventureGameItemPlacement] {
	return m.Repositories[adventure_game_item_placement.TableName].(*repository.Generic[adventure_game_record.AdventureGameItemPlacement, *adventure_game_record.AdventureGameItemPlacement])
}

// AdventureGameCreaturePlacementRepository -
func (m *Domain) AdventureGameCreaturePlacementRepository() *repository.Generic[adventure_game_record.AdventureGameCreaturePlacement, *adventure_game_record.AdventureGameCreaturePlacement] {
	return m.Repositories[adventure_game_creature_placement.TableName].(*repository.Generic[adventure_game_record.AdventureGameCreaturePlacement, *adventure_game_record.AdventureGameCreaturePlacement])
}

// AdventureGameCreatureRepository -
func (m *Domain) AdventureGameCreatureRepository() *repository.Generic[adventure_game_record.AdventureGameCreature, *adventure_game_record.AdventureGameCreature] {
	return m.Repositories[adventure_game_creature.TableName].(*repository.Generic[adventure_game_record.AdventureGameCreature, *adventure_game_record.AdventureGameCreature])
}

// AdventureGameItemInstanceRepository -
func (m *Domain) AdventureGameItemInstanceRepository() *repository.Generic[adventure_game_record.AdventureGameItemInstance, *adventure_game_record.AdventureGameItemInstance] {
	return m.Repositories[adventure_game_item_instance.TableName].(*repository.Generic[adventure_game_record.AdventureGameItemInstance, *adventure_game_record.AdventureGameItemInstance])
}

// AdventureGameLocationLinkRequirementRepository -
func (m *Domain) AdventureGameLocationLinkRequirementRepository() *repository.Generic[adventure_game_record.AdventureGameLocationLinkRequirement, *adventure_game_record.AdventureGameLocationLinkRequirement] {
	return m.Repositories[adventure_game_location_link_requirement.TableName].(*repository.Generic[adventure_game_record.AdventureGameLocationLinkRequirement, *adventure_game_record.AdventureGameLocationLinkRequirement])
}

// GameInstanceRepository -
func (m *Domain) GameInstanceRepository() *repository.Generic[game_record.GameInstance, *game_record.GameInstance] {
	return m.Repositories[game_instance.TableName].(*repository.Generic[game_record.GameInstance, *game_record.GameInstance])
}

// AdventureGameLocationInstanceRepository -
func (m *Domain) AdventureGameLocationInstanceRepository() *repository.Generic[adventure_game_record.AdventureGameLocationInstance, *adventure_game_record.AdventureGameLocationInstance] {
	return m.Repositories[adventure_game_location_instance.TableName].(*repository.Generic[adventure_game_record.AdventureGameLocationInstance, *adventure_game_record.AdventureGameLocationInstance])
}

// AdventureGameCreatureInstanceRepository -
func (m *Domain) AdventureGameCreatureInstanceRepository() *repository.Generic[adventure_game_record.AdventureGameCreatureInstance, *adventure_game_record.AdventureGameCreatureInstance] {
	return m.Repositories[adventure_game_record.TableAdventureGameCreatureInstance].(*repository.Generic[adventure_game_record.AdventureGameCreatureInstance, *adventure_game_record.AdventureGameCreatureInstance])
}

// AdventureGameCharacterInstanceRepository -
func (m *Domain) AdventureGameCharacterInstanceRepository() *repository.Generic[adventure_game_record.AdventureGameCharacterInstance, *adventure_game_record.AdventureGameCharacterInstance] {
	return m.Repositories[adventure_game_record.TableAdventureGameCharacterInstance].(*repository.Generic[adventure_game_record.AdventureGameCharacterInstance, *adventure_game_record.AdventureGameCharacterInstance])
}

// GameSubscriptionRepository -
func (m *Domain) GameSubscriptionRepository() *repository.Generic[game_record.GameSubscription, *game_record.GameSubscription] {
	return m.Repositories[game_subscription.TableName].(*repository.Generic[game_record.GameSubscription, *game_record.GameSubscription])
}

// GameAdministrationRepository -
func (m *Domain) GameAdministrationRepository() *repository.Generic[game_record.GameAdministration, *game_record.GameAdministration] {
	return m.Repositories[game_administration.TableName].(*repository.Generic[game_record.GameAdministration, *game_record.GameAdministration])
}

// GameConfigurationRepository -
func (m *Domain) GameConfigurationRepository() *repository.Generic[game_record.GameConfiguration, *game_record.GameConfiguration] {
	return m.Repositories[game_configuration.TableName].(*repository.Generic[game_record.GameConfiguration, *game_record.GameConfiguration])
}

// GameInstanceConfigurationRepository -
func (m *Domain) GameInstanceConfigurationRepository() *repository.Generic[game_record.GameInstanceConfiguration, *game_record.GameInstanceConfiguration] {
	return m.Repositories[game_instance_configuration.TableName].(*repository.Generic[game_record.GameInstanceConfiguration, *game_record.GameInstanceConfiguration])
}

// GameTurnSheetRepository -
func (m *Domain) GameTurnSheetRepository() *repository.Generic[game_record.GameTurnSheet, *game_record.GameTurnSheet] {
	return m.Repositories[game_turn_sheet.TableName].(*repository.Generic[game_record.GameTurnSheet, *game_record.GameTurnSheet])
}

// AdventureGameTurnSheetRepository -
func (m *Domain) AdventureGameTurnSheetRepository() *repository.Generic[adventure_game_record.AdventureGameTurnSheet, *adventure_game_record.AdventureGameTurnSheet] {
	return m.Repositories[adventure_game_turn_sheet.TableName].(*repository.Generic[adventure_game_record.AdventureGameTurnSheet, *adventure_game_record.AdventureGameTurnSheet])
}

// Logger - Returns a logger with package context and provided function context
func (m *Domain) Logger(functionName string) logger.Logger {
	return m.Log.WithFunctionContext(functionName)
}
