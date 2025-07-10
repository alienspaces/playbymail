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
	"gitlab.com/alienspaces/playbymail/internal/repository/game_location"
	"gitlab.com/alienspaces/playbymail/internal/repository/game_location_link"
)

// Domain -
type Domain struct {
	jobClient *river.Client[pgx.Tx]
	domain.Domain
}

var _ domainer.Domainer = &Domain{}

func NewDomain(l logger.Logger, j *river.Client[pgx.Tx]) (*Domain, error) {

	l, err := l.NewInstance()
	if err != nil {
		return nil, err
	}

	m := &Domain{
		jobClient: j,
		Domain: domain.Domain{
			Log: l.WithPackageContext("domain"),
			RepositoryConstructors: []domain.RepositoryConstructor{
				account.NewRepository,
				game.NewRepository,
				game_location.NewRepository,
				game_location_link.NewRepository,
				game_character.NewRepository,
			},
		},
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
