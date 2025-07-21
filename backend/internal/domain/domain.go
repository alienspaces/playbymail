package domain

import (
	"maps"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"

	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/repository/account"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_character"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_character_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_creature"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_creature_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_item"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_item_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location_link"
	"gitlab.com/alienspaces/playbymail/internal/repository/adventure_game_location_link_requirement"
	"gitlab.com/alienspaces/playbymail/internal/repository/game"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// Domain -
type Domain struct {
	domain.Domain
	config config.Config
}

var _ domainer.Domainer = &Domain{}

func NewDomain(l logger.Logger, j *river.Client[pgx.Tx], cfg config.Config) (*Domain, error) {

	l, err := l.NewInstance()
	if err != nil {
		return nil, err
	}

	m := &Domain{
		Domain: domain.Domain{
			Log: l.WithPackageContext("domain"),
			RepositoryConstructors: []domain.RepositoryConstructor{
				account.NewRepository,
				game.NewRepository,
				adventure_game_location.NewRepository,
				adventure_game_location_link.NewRepository,
				adventure_game_character.NewRepository,
				adventure_game_item.NewRepository,
				adventure_game_creature.NewRepository,
				adventure_game_item_instance.NewRepository,
				adventure_game_location_link_requirement.NewRepository,
				adventure_game_instance.NewRepository,
				adventure_game_location_instance.NewRepository,
				adventure_game_creature_instance.NewRepository,
				adventure_game_character_instance.NewRepository,
			},
		},
		config: cfg,
	}

	m.SetRLSFunc = m.SetRLS

	l.Info("returning domain %+v", m)

	return m, nil
}

// AccountRepository -
func (m *Domain) AccountRepository() *repository.Generic[record.Account, *record.Account] {
	return m.Repositories[account.TableName].(*repository.Generic[record.Account, *record.Account])
}

// GameRepository -
func (m *Domain) GameRepository() *repository.Generic[record.Game, *record.Game] {
	return m.Repositories[game.TableName].(*repository.Generic[record.Game, *record.Game])
}

// GameLocationLinkRepository -
func (m *Domain) AdventureGameLocationLinkRepository() *repository.Generic[record.AdventureGameLocationLink, *record.AdventureGameLocationLink] {
	return m.Repositories[adventure_game_location_link.TableName].(*repository.Generic[record.AdventureGameLocationLink, *record.AdventureGameLocationLink])
}

// AdventureGameCharacterRepository -
func (m *Domain) AdventureGameCharacterRepository() *repository.Generic[record.AdventureGameCharacter, *record.AdventureGameCharacter] {
	return m.Repositories[adventure_game_character.TableName].(*repository.Generic[record.AdventureGameCharacter, *record.AdventureGameCharacter])
}

// AdventureGameLocationRepository -
func (m *Domain) AdventureGameLocationRepository() *repository.Generic[record.AdventureGameLocation, *record.AdventureGameLocation] {
	return m.Repositories[adventure_game_location.TableName].(*repository.Generic[record.AdventureGameLocation, *record.AdventureGameLocation])
}

// AdventureGameItemRepository -
func (m *Domain) AdventureGameItemRepository() *repository.Generic[record.AdventureGameItem, *record.AdventureGameItem] {
	return m.Repositories[adventure_game_item.TableName].(*repository.Generic[record.AdventureGameItem, *record.AdventureGameItem])
}

// AdventureGameCreatureRepository -
func (m *Domain) AdventureGameCreatureRepository() *repository.Generic[record.AdventureGameCreature, *record.AdventureGameCreature] {
	return m.Repositories[adventure_game_creature.TableName].(*repository.Generic[record.AdventureGameCreature, *record.AdventureGameCreature])
}

// AdventureGameItemInstanceRepository -
func (m *Domain) AdventureGameItemInstanceRepository() *repository.Generic[record.AdventureGameItemInstance, *record.AdventureGameItemInstance] {
	return m.Repositories[adventure_game_item_instance.TableName].(*repository.Generic[record.AdventureGameItemInstance, *record.AdventureGameItemInstance])
}

// AdventureGameLocationLinkRequirementRepository -
func (m *Domain) AdventureGameLocationLinkRequirementRepository() *repository.Generic[record.AdventureGameLocationLinkRequirement, *record.AdventureGameLocationLinkRequirement] {
	return m.Repositories[adventure_game_location_link_requirement.TableName].(*repository.Generic[record.AdventureGameLocationLinkRequirement, *record.AdventureGameLocationLinkRequirement])
}

// AdventureGameInstanceRepository -
func (m *Domain) AdventureGameInstanceRepository() *repository.Generic[record.AdventureGameInstance, *record.AdventureGameInstance] {
	return m.Repositories[adventure_game_instance.TableName].(*repository.Generic[record.AdventureGameInstance, *record.AdventureGameInstance])
}

// AdventureGameLocationInstanceRepository -
func (m *Domain) AdventureGameLocationInstanceRepository() *repository.Generic[record.AdventureGameLocationInstance, *record.AdventureGameLocationInstance] {
	return m.Repositories[adventure_game_location_instance.TableName].(*repository.Generic[record.AdventureGameLocationInstance, *record.AdventureGameLocationInstance])
}

// AdventureGameCreatureInstanceRepository -
func (m *Domain) AdventureGameCreatureInstanceRepository() *repository.Generic[record.AdventureGameCreatureInstance, *record.AdventureGameCreatureInstance] {
	return m.Repositories[record.TableAdventureGameCreatureInstance].(*repository.Generic[record.AdventureGameCreatureInstance, *record.AdventureGameCreatureInstance])
}

// AdventureGameCharacterInstanceRepository -
func (m *Domain) AdventureGameCharacterInstanceRepository() *repository.Generic[record.AdventureGameCharacterInstance, *record.AdventureGameCharacterInstance] {
	return m.Repositories[record.TableAdventureGameCharacterInstance].(*repository.Generic[record.AdventureGameCharacterInstance, *record.AdventureGameCharacterInstance])
}

// SetRLS -
func (m *Domain) SetRLS(identifiers map[string][]string) {

	// We'll be resetting the "id" key when we use the map
	ri := maps.Clone(identifiers)

	for tableName := range m.Repositories {

		// When the repository table name matches an RLS identifier key, we apply the
		// RLS constraints to the "id" column to enforce any RLS constraints on itself!
		// Can this be done inside repository core code on itself? Absolutely... but it
		// would be making a naive assumption about conventions. This project's convention
		// is to name foreign key columns according to the table name it foreign keys to.
		// If that convention is not followed, then the following block would not work.
		if _, ok := ri[tableName+"_id"]; ok {
			ri["id"] = ri[tableName+"_id"]
			m.Repositories[tableName].SetRLS(ri)
			continue
		}
		m.Repositories[tableName].SetRLS(identifiers)
	}
}

// Logger - Returns a logger with package context and provided function context
func (m *Domain) Logger(functionName string) logger.Logger {
	return m.Log.WithFunctionContext(functionName)
}
