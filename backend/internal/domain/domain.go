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
	"gitlab.com/alienspaces/playbymail/internal/repository/game"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_character"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_character_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_creature"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_creature_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_item"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_item_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_location"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_location_instance"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_location_link"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_location_link_requirement"
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
				game_location.NewRepository,
				game_location_link.NewRepository,
				game_character.NewRepository,
				game_item.NewRepository,
				game_creature.NewRepository,
				game_item_instance.NewRepository,
				game_location_link_requirement.NewRepository,
				game_instance.NewRepository,
				game_location_instance.NewRepository,
				game_creature_instance.NewRepository,
				game_character_instance.NewRepository,
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
func (m *Domain) GameLocationLinkRepository() *repository.Generic[record.GameLocationLink, *record.GameLocationLink] {
	return m.Repositories[game_location_link.TableName].(*repository.Generic[record.GameLocationLink, *record.GameLocationLink])
}

// GameCharacterRepository -
func (m *Domain) GameCharacterRepository() *repository.Generic[record.GameCharacter, *record.GameCharacter] {
	return m.Repositories[game_character.TableName].(*repository.Generic[record.GameCharacter, *record.GameCharacter])
}

// GameLocationRepository -
func (m *Domain) GameLocationRepository() *repository.Generic[record.GameLocation, *record.GameLocation] {
	return m.Repositories[game_location.TableName].(*repository.Generic[record.GameLocation, *record.GameLocation])
}

// GameItemRepository -
func (m *Domain) GameItemRepository() *repository.Generic[record.GameItem, *record.GameItem] {
	return m.Repositories[game_item.TableName].(*repository.Generic[record.GameItem, *record.GameItem])
}

// GameCreatureRepository -
func (m *Domain) GameCreatureRepository() *repository.Generic[record.GameCreature, *record.GameCreature] {
	return m.Repositories[game_creature.TableName].(*repository.Generic[record.GameCreature, *record.GameCreature])
}

// GameItemInstanceRepository -
func (m *Domain) GameItemInstanceRepository() *repository.Generic[record.GameItemInstance, *record.GameItemInstance] {
	return m.Repositories[game_item_instance.TableName].(*repository.Generic[record.GameItemInstance, *record.GameItemInstance])
}

// GameLocationLinkRequirementRepository -
func (m *Domain) GameLocationLinkRequirementRepository() *repository.Generic[record.GameLocationLinkRequirement, *record.GameLocationLinkRequirement] {
	return m.Repositories[game_location_link_requirement.TableName].(*repository.Generic[record.GameLocationLinkRequirement, *record.GameLocationLinkRequirement])
}

// GameInstanceRepository -
func (m *Domain) GameInstanceRepository() *repository.Generic[record.GameInstance, *record.GameInstance] {
	return m.Repositories[game_instance.TableName].(*repository.Generic[record.GameInstance, *record.GameInstance])
}

// GameLocationInstanceRepository -
func (m *Domain) GameLocationInstanceRepository() *repository.Generic[record.GameLocationInstance, *record.GameLocationInstance] {
	return m.Repositories[game_location_instance.TableName].(*repository.Generic[record.GameLocationInstance, *record.GameLocationInstance])
}

// GameCreatureInstanceRepository -
func (m *Domain) GameCreatureInstanceRepository() *repository.Generic[record.GameCreatureInstance, *record.GameCreatureInstance] {
	return m.Repositories[record.TableGameCreatureInstance].(*repository.Generic[record.GameCreatureInstance, *record.GameCreatureInstance])
}

// GameCharacterInstanceRepository -
func (m *Domain) GameCharacterInstanceRepository() *repository.Generic[record.GameCharacterInstance, *record.GameCharacterInstance] {
	return m.Repositories[record.TableGameCharacterInstance].(*repository.Generic[record.GameCharacterInstance, *record.GameCharacterInstance])
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
